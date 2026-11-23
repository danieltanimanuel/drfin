package gmail

import (
	"drfin/app/shared/logdr"
	"encoding/base64"
	"fmt"
	"html"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/gmail/v1"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

// fetchMessageByID retrieves a Gmail message by ID.
func fetchMessageByID(srv *gmail.Service, userID, messageID string) (*gmail.Message, error) {
	return srv.Users.Messages.Get(userID, messageID).Do()
}

// listRecentMessages lists recent messages for a user.
func listRecentMessages(srv *gmail.Service, userID string, maxResults int64) ([]*gmail.Message, error) {
	resp, err := srv.Users.Messages.List(userID).MaxResults(maxResults).Do()
	if err != nil {
		return nil, err
	}
	return resp.Messages, nil
}

func listMessages(
	srv *gmail.Service,
	userID string,
	query string,
) ([]*gmail.Message, error) {

	var messages []*gmail.Message
	pageToken := ""

	for {
		req := srv.Users.Messages.List(userID).Q(query)
		if pageToken != "" {
			req.PageToken(pageToken)
		}

		resp, err := req.Do()
		if err != nil {
			return nil, err
		}

		messages = append(messages, resp.Messages...)

		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken
	}

	return messages, nil
}

// decodeMessageBody decodes the body of a Gmail message part.
func decodeMessageBody(part *gmail.MessagePart) (string, error) {
	if part.Body == nil || part.Body.Data == "" {
		return "", nil
	}
	decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// parseHTMLTables parses and prints table data from HTML content.
func parseHTMLTables(htmlContent string, format EmailTitle, holderName string) ([]GmailDetailTransaction, error) {
	switch format {
	case EmailTitleCreditCardBCA:
		return parseCCBCAXML(htmlContent)
	case EmailTitleTransaksiBlu:
		return parseBLUXML(htmlContent, holderName)
	case EmailTitleTransaksiBCA:
		return parseBCAXML(htmlContent)
	default:
		return []GmailDetailTransaction{}, fmt.Errorf("email title not supported")
	}
}

// printMessageSummary prints a summary of a Gmail message and return the subject.
func extractMailTitle(msg *gmail.Message) string {
	// timestamp := msg.InternalDate / 1000
	var subject string
	for _, header := range msg.Payload.Headers {
		if header.Name == "Subject" {
			subject = header.Value
			break
		}
	}
	// fmt.Printf("%d - %s - %s\n", timestamp, msg.Id, subject)
	return subject
}

func isMostNestedDiv(s *goquery.Selection) bool {
	// Check if this div contains another div inside it
	return s.Find("div").Length() == 0
}

func parseBLUXML(htmlContent string, holderName string) ([]GmailDetailTransaction, error) {
	holderName = strings.TrimSpace(holderName)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logdr.Error("Error loading HTML", err)
		return []GmailDetailTransaction{}, err
	}
	// logdr.Info("Doc", htmlContent)
	var detailTransactions []GmailDetailTransaction
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		// logdr.Info("Text", "text", text)
		if isMostNestedDiv(s) {
			switch {
			case text == "Total Bayar":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Total Bayar",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("tr", "html").Next().Text()),
				})
			case text == "Tipe Transaksi":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Tipe Transaksi",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("div", "html").Next().Text()),
				})
			case text == "Total": // For QRIS, Total is on the header not as field, but for BI-FAST, total is on the field
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Total",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("div", "html").Next().Text()),
				})
			case text == "bluRewards":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "bluRewards",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("div", "html").Next().Text()),
				})
			case text == "Biaya Admin":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Biaya Admin",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("div", "html").Next().Text()),
				})
			case text == "Tgl & Jam Transaksi":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Tgl & Jam Transaksi",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("div", "html").Next().Text()),
				})
			case text == "Flazz":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Flazz",
					Value: strings.TrimSpace(s.Next().Text()),
				})
			case text == "Nominal Tagihan":
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Nominal Tagihan",
					Value: strings.TrimSpace(s.ParentsFilteredUntil("div", "html").Next().Text()),
				})
			case text == holderName:
				// Step 1: Get the parent <td>
				td := s.ParentsFiltered("td")

				// Step 2: Move to the next <td>
				nextTd := td.Next().Next()

				// Step 3: Find the <span> inside that td which has BCA + account
				info1 := nextTd.Find("span").First()
				info2 := nextTd.Find("span").Last()

				// Step 4: Extract text and clean it
				logdr.Info(info1.Text())
				logdr.Info(info2.Text())
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Beneficiary Full",
					Value: info1.Text() + " " + info2.Text(),
				})

				// Step 5 : extract beneficiary account number to be compared to owned accounts
				re := regexp.MustCompile(`\D`) // \D = non-digit
				beneficiary := re.ReplaceAllString(info2.Text(), "")
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   "Beneficiary",
					Value: beneficiary,
				})
			}
		}
	})

	logdr.Info("Detail Transactions", "detailTransactions", detailTransactions)
	return detailTransactions, nil
}

func CleanGmailHTML(raw string) string {
	// If raw comes directly from Gmail API JSON, escape is needed:
	unquoted, err := strconv.Unquote(`"` + raw + `"`)
	if err != nil {
		unquoted = raw // fallback
	}

	// Remove \r\n
	cleaned := strings.ReplaceAll(unquoted, "\r\n", "")

	// Fix HTML entities (&amp; -> &, &quot; -> ")
	return html.UnescapeString(cleaned)
}

func parseCCBCAXML(htmlContent string) ([]GmailDetailTransaction, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(CleanGmailHTML(htmlContent)))
	if err != nil {
		log.Println("Error loading HTML:", err)
		return []GmailDetailTransaction{}, err
	}
	// logdr.Info("Doc",htmlContent);
	var detailTransactions []GmailDetailTransaction
	doc.Find("table.es-table.cke_show_border.d.e").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(i int, tr *goquery.Selection) {
			cells := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
				return strings.TrimSpace(td.Text())
			})
			if len(cells) > 0 {
				detailTransactions = append(detailTransactions, GmailDetailTransaction{
					Key:   cells[0],
					Value: cells[2],
				})
			}
		})
	})
	return detailTransactions, nil
}

func parseBCAXML(htmlContent string) ([]GmailDetailTransaction, error) {
	htmlStr := html.UnescapeString(string(htmlContent))
	re := regexp.MustCompile(`(?s)<div class="table-box".*?>(.*?)</div>`)
	match := re.FindStringSubmatch(htmlStr)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(match[1]))
	if err != nil {
		logdr.Error("Error loading HTML:", err)
		return []GmailDetailTransaction{}, err
	}
	// logdr.Info("Doc",htmlContent);
	var detailTransactions []GmailDetailTransaction
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		tds := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
			return strings.TrimSpace(td.Text())
		})

		if len(tds) == 3 {
			detailTransactions = append(detailTransactions, GmailDetailTransaction{
				Key:   tds[0],
				Value: tds[2],
			})
		}
	})
	logdr.Info("Detail Transactions", "detailTransactions", detailTransactions)
	return detailTransactions, nil
}

func extractMail(message *gmail.Message, holderName string) (*GmailExtractedData, error) {
	emailTitle := extractMailTitle(message)
	for _, part := range message.Payload.Parts {
		// logdr.Info("","2- ", part.MimeType)
		if part.MimeType == "text/plain" || part.MimeType == "text/html" {
			body, err := decodeMessageBody(part)
			if err != nil {
				logdr.Error("error decoding body", err)
				continue
			}
			if part.MimeType == "text/html" {
				detailTransactions, err := parseHTMLTables(body, EmailTitle(emailTitle), holderName)
				if err != nil {
					logdr.Error("error parsing HTML tables", err)
					continue
				}
				// fmt.Println(detailTransactions)
				extractedData := &GmailExtractedData{
					Title:             EmailTitle(emailTitle),
					DetailTransaction: detailTransactions,
				}
				// fmt.Println(extractedData)
				return extractedData, nil
			}
		} else if part.MimeType == "multipart/related" {
			/* BLU */
			for _, subPart := range part.Parts {
				// log.Println("3 - ", subPart, " ", len(part.Parts))
				if subPart.MimeType == "text/html" {
					html, err := decodeMessageBody(subPart)
					if err != nil {
						logdr.Error("Unable to decode message body", err)
					}
					detailTransactions, err := parseHTMLTables(html, EmailTitle(emailTitle), holderName)
					if err != nil {
						logdr.Error("Unable to parse HTML tables", err)
					}
					// logdr.Info("log detail transactions","detailTransactions",detailTransactions)
					extractedData := &GmailExtractedData{
						Title:             EmailTitle(emailTitle),
						DetailTransaction: detailTransactions,
					}
					// logdr.Info("log extracted data", "extractedData",extractedData)
					return extractedData, nil
				}
			}
		} else if part.MimeType == "multipart/alternative" {
			/* CC BCA */
			for _, subPart := range part.Parts {
				// log.Println("3 - ", subPart, " ", len(part.Parts))
				body, err := decodeMessageBody(subPart)
				if err != nil {
					logdr.Error("error decoding body", err)
					continue
				}
				if subPart.MimeType == "text/html" {
					detailTransactions, err := parseHTMLTables(body, EmailTitle(emailTitle), holderName)
					if err != nil {
						logdr.Error("error parsing HTML tables", err)
						continue
					}
					// fmt.Println(detailTransactions)
					extractedData := &GmailExtractedData{
						Title:             EmailTitle(emailTitle),
						DetailTransaction: detailTransactions,
					}
					// fmt.Println(extractedData)
					return extractedData, nil
				}
			}
		}
	}

	return nil, nil
}
