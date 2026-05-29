package gmail

import (
	"context"
	"drfin/app/accounts"
	"drfin/app/lastemail"
	"drfin/app/processedgmail"
	"drfin/app/shared/logdr"
	"drfin/app/transactions"
	"fmt"
	"sort"
	"strconv"
	"time"

	"google.golang.org/api/gmail/v1"
)

type AmountFormat string
type LayoutDateTime string
type EmailTitle string

const (
	EmailTitleCreditCardBCA       EmailTitle     = "Credit Card Transaction Notification"
	EmailTitleTransaksiBlu        EmailTitle     = "Transaksimu Pakai blu Berhasil"
	EmailTitleTransaksiBCA        EmailTitle     = "Internet Transaction Journal"
	LayoutDateTimeCCBCA           LayoutDateTime = "02-01-2006 15:04:05"
	LayoutDateTimeDDMMMYYYYHHMMSS LayoutDateTime = "02 Jan 2006 15:04:05"
	AmountFormatRp                AmountFormat   = "Rp"
	AmountFormatBCA               AmountFormat   = "IDR"
	TargetMessageID                              = "198d6af7ea2020da" // Replace with actual message ID
)

type Service struct {
	transactionService       *transactions.Service
	accountService           *accounts.Service
	lastEmailRepository      *lastemail.Repository
	processedGmailRepository *processedgmail.Repository
}

func NewService(transactionService *transactions.Service,
	accountService *accounts.Service,
	lastEmailRepository *lastemail.Repository,
	processedGmailRepository *processedgmail.Repository) *Service {
	return &Service{transactionService: transactionService, accountService: accountService, lastEmailRepository: lastEmailRepository, processedGmailRepository: processedGmailRepository}
}

func shouldProcessTitle(format EmailTitle) bool {
	switch format {
	case EmailTitleCreditCardBCA, EmailTitleTransaksiBlu, EmailTitleTransaksiBCA:
		// case EmailTitleTransaksiBlu:
		return true
	default:
		return false
	}
}

var allowedTitles = map[EmailTitle]struct{}{
	EmailTitleCreditCardBCA: {},
	EmailTitleTransaksiBlu:  {},
	EmailTitleTransaksiBCA:  {},
}

func IsAllowedEmailTitle(subject string) bool {
	_, ok := allowedTitles[EmailTitle(subject)]
	return ok
}

func (s *Service) syncMail(ctx context.Context) error {
	srv, err := authO2()
	if err != nil {
		return err
	}

	// 1. Load checkpoint
	lastInternalDate, err := s.lastEmailRepository.GetLastInternalDate()
	if err != nil {
		return err
	}
	logdr.Info("checkpoint internal_date:", "lastinternaldate", lastInternalDate)

	// 2. Coarse Gmail query (date-based)
	after := time.UnixMilli(lastInternalDate).UTC()
	before := time.Now().AddDate(0,0,1)
	query := fmt.Sprintf(
		"after:%s before:%s",
		after.Format("2006/01/02"),
		before.Format("2006/01/02"),
	)
	logdr.Info("Fetch email with query : ", "Query", query)

	msgs, _ := listMessages(srv, "me", query)
	if err != nil {
		logdr.Info("Unable to get last email", err)
		return err
	}
	logdr.Info("msg list length : ", "length", len(msgs))

	var fullMsgs []*gmail.Message

	logdr.Info("Fetching each messages...")
	for _, m := range msgs {
		msg, err := fetchMessageByID(srv, "me", m.Id)
		if err != nil {
			return err
		}
		fullMsgs = append(fullMsgs, msg)
		logdr.Info("Progress bar " + fmt.Sprint(len(fullMsgs)) + "/" + fmt.Sprint(len(msgs)))
	}

	logdr.Info("Sorting...")
	// 3. Sort by internalDate (Gmail does NOT guarantee order)
	sort.Slice(fullMsgs, func(i, j int) bool {
		return fullMsgs[i].InternalDate < fullMsgs[j].InternalDate
	})

	totalChecked := 0
	totalProcessed := 0

	// 🔐 Strategy-B variable
	safeInternalDate := lastInternalDate

	for _, message := range fullMsgs {
		totalChecked++

		title := extractMailTitle(message)

		// Millisecond-precision cutoff
		if message.InternalDate <= lastInternalDate {
			logdr.Info("Skip this email because of internal date", "Email Title", title, "Email ID", message.Id)
			continue
		}

		// Skip if not in the list of should process title
		if !shouldProcessTitle(EmailTitle(title)) {
			logdr.Info("Skip this email because its not in the listed title", "Email Title", title, "Email ID", message.Id)
			continue
		}

		// Skip if already successfully processed
		exists, err := s.processedGmailRepository.Exists(message.Id)
		if err != nil {
			return err
		}

		if exists {
			logdr.Info("Skip this email because its already been processed", "Email Title", title, "Email ID", message.Id)
			continue
		}

		// Start processing email
		logdr.Info("Processing this email", message.Id, title)
		logdr.Debug(strconv.Itoa(len(message.Payload.Parts)))
		extractedData, err := extractMail(message)
		if err != nil {
			_ = s.processedGmailRepository.RecordFailure(
				message.Id,
				message.InternalDate,
				err.Error(),
			)
			break
		}

		logdr.Info("Creating transaction..  ")
		if err := s.createTransactionFromGmail(ctx, extractedData); err != nil {
			_ = s.processedGmailRepository.RecordFailure(
				message.Id,
				message.InternalDate,
				err.Error(),
			)
			break
		}

		// ✅ SUCCESS
		_ = s.processedGmailRepository.RecordSuccess(
			message.Id,
			message.InternalDate,
			0, // optional transaction ID
		)

		safeInternalDate = message.InternalDate

		totalProcessed++
	}

	logdr.Info("Sync finished")
	logdr.Info("Total checked:", totalChecked)
	logdr.Info("Total processed:", totalProcessed)
	logdr.Info("New checkpoint internal_date:", safeInternalDate)

	// 5. Advance checkpoint ONLY to last safe success
	if safeInternalDate > lastInternalDate {
		if err := s.lastEmailRepository.UpdateLastInternalDate(safeInternalDate); err != nil {
			return err
		}
	}
	return nil
}

// --------------------
// Dispatcher
// --------------------
func (s *Service) createTransactionFromGmail(ctx context.Context, data *GmailExtractedData) error {
	handlers := map[EmailTitle]func(context.Context, *GmailExtractedData) error{
		EmailTitleCreditCardBCA: s.handleCreditCardBCA,
		EmailTitleTransaksiBlu:  s.handleTransaksiBlu,
		EmailTitleTransaksiBCA:  s.handleTransaksiBCA,
	}

	logdr.Info(string(data.Title))
	if handler, ok := handlers[data.Title]; ok {
		return handler(ctx, data)
	}
	return fmt.Errorf("no handler for email title: %s", data.Title)
}
