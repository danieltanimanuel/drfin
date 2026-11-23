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
	"google.golang.org/api/gmail/v1"
	"errors"
	"drfin/app/owners"
	"strings"
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
	ownerRepository 		 *owners.Repository
}

// ErrSkipEmail is returned by a handler to say "intentionally not imported".
// It is not a failure: the sync loop records the email and moves on.
var ErrSkipEmail = errors.New("email intentionally skipped")

func NewService(transactionService *transactions.Service,
	accountService *accounts.Service,
	lastEmailRepository *lastemail.Repository,
	processedGmailRepository *processedgmail.Repository,
	ownerRepository *owners.Repository) *Service {
	return &Service{transactionService: transactionService, accountService: accountService, lastEmailRepository: lastEmailRepository, processedGmailRepository: processedGmailRepository, ownerRepository: ownerRepository}
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

func (s *Service) syncMail(ctx context.Context, owner owners.Owner) error {
	srv, err := authO2(owner.ID)
	if err != nil {
		return err
	}

	// 1. Load checkpoint
	lastInternalDate, err := s.lastEmailRepository.GetLastInternalDate(owner.ID)
	if err != nil {
		return err
	}
	logdr.Info("checkpoint internal_date:", "lastinternaldate", lastInternalDate)

	// 2. Coarse Gmail query (date-based)
	query := fmt.Sprintf(`after:%d {subject:"%s" subject:"%s" subject:"%s"}`,
		lastInternalDate/1000, EmailTitleCreditCardBCA, EmailTitleTransaksiBlu, EmailTitleTransaksiBCA)
	logdr.Info("Fetch email with query", "query", query)

	msgs, err := listMessages(srv, "me", query)
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
		exists, err := s.processedGmailRepository.Exists(owner.ID, message.Id)
		if err != nil {
			return err
		}

		if exists {
			logdr.Info("Skip this email because its already been processed", "Email Title", title, "Email ID", message.Id)
			continue
		}

		// Start processing email
		logdr.Info("Processing email", "email_id", message.Id, "subject", title)
		data, err := s.processMessage(ctx, owner, message)
		if errors.Is(err, ErrSkipEmail) {
			logdr.Warn("Skipping email: unsupported currency", "email_id", message.Id, "subject", title)
			// Record as processed so it is not retried, and let the checkpoint advance.
			_ = s.processedGmailRepository.RecordSuccess(owner.ID, message.Id, message.InternalDate, 0)
			safeInternalDate = message.InternalDate
			continue
		}
		if err != nil {
			s.failMessage(owner.ID, message, title, data, err)
			break
		}

		// ✅ SUCCESS
		_ = s.processedGmailRepository.RecordSuccess(
			owner.ID,
			message.Id,
			message.InternalDate,
			0, // optional transaction ID
		)

		safeInternalDate = message.InternalDate

		totalProcessed++
	}

	logdr.Info("Sync finished")
	logdr.Info("Total checked:" + strconv.Itoa(totalChecked))
	logdr.Info("Total processed:" + strconv.Itoa(totalProcessed))
	logdr.Info("New checkpoint internal_date:" + strconv.Itoa(int(safeInternalDate)))

	// 5. Advance checkpoint ONLY to last safe success
	if safeInternalDate > lastInternalDate {
		if err := s.lastEmailRepository.UpdateLastInternalDate(owner.ID, safeInternalDate); err != nil {
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

// processMessage extracts one email and creates its transaction.
// A panic (bad HTML, short-string slicing) becomes a recorded error instead of a lost request.
func (s *Service) processMessage(ctx context.Context, owner owners.Owner, message *gmail.Message) (data *GmailExtractedData, err error) {
	defer func() { /* unchanged recover block */ }()

	// Without the holder name the blu parser would silently import a transfer as an Expense.
	if EmailTitle(extractMailTitle(message)) == EmailTitleTransaksiBlu && strings.TrimSpace(owner.BankHolderName) == "" {
		return nil, fmt.Errorf("bank_holder_name is not set for owner %d", owner.ID)
	}
	if data, err = extractMail(message, owner.BankHolderName); err != nil {
		return data, err
	}
	if data == nil {
		return nil, fmt.Errorf("no parsable html part found")
	}
	data.OwnerID = owner.ID
	return data, s.createTransactionFromGmail(ctx, data)
}

// failMessage logs and stores everything needed to debug one email.
func (s *Service) failMessage(ownerID int, msg *gmail.Message, title string, data *GmailExtractedData, err error) {
	details := map[string]string{}
	if data != nil {
		details = toMap(data.DetailTransaction)
	}
	logdr.Info("EMAIL FAILED", "email_id", msg.Id, "subject", title,
		"internal_date", msg.InternalDate, "error", err.Error(), "details", details)
	if e := s.processedGmailRepository.RecordFailure(ownerID, msg.Id, msg.InternalDate,
		fmt.Sprintf("%v | subject=%q | details=%v", err, title, details)); e != nil {
		logdr.Error("record failure", e)
	}
}