package gmail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Auth struct {
	service *Service
}

func NewAuth(service *Service) *Auth {
	return &Auth{
		service: service,
	}
}

const (
	TokenFile        = "token.json"
	ClientSecretFile = "resources/drfin-secret.json"
	SampleMessageID  = "196152a4ef35e1fa" // Replace with actual message ID
	TargetMessageID  = "197bfde226e8ec20" // Replace with actual message ID
)

// getClient retrieves a token, saves it if needed, and returns an HTTP client.
func getClient(config *oauth2.Config) *http.Client {
	tok, err := tokenFromFile(TokenFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(TokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// tokenFromFile retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	tokJSON, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	if err := json.Unmarshal(tokJSON, tok); err != nil {
		return nil, err
	}
	return tok, nil
}

// getTokenFromWeb requests a token from the web, prompting the user.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser:\n%v\n", authURL)

	var code string
	fmt.Print("Enter the authorization code: ")
	fmt.Scan(&code)

	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	return tok
}

// saveToken saves a token to a file path.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	tokJSON, _ := json.Marshal(token)
	os.WriteFile(file, tokJSON, 0600)
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

// decodeMessageBody2 decodes the body of a Gmail message part.
func decodeMessageBody2(part *gmail.MessagePart) (string, error) {
	if part.Body == nil || part.Body.Data == "" {
		return "", nil
	}
	decoded, err := base64.StdEncoding.DecodeString(part.Body.Data)
	log.Println(decoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// parseHTMLTables parses and prints table data from HTML content.
func parseHTMLTables(htmlContent string, emailTitle string) ([]GmailDetailTransaction, error) {
	switch emailTitle {
	case EmailTitleCreditCardBCA:
		return parseCCBCAXML(htmlContent)
	case EmailTitleTransaksiBlu:
		return parseBLUXML(htmlContent)
	default:
		return []GmailDetailTransaction{}, fmt.Errorf("email title not supported")
	}
}

func extractText(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
}

func parseBLUXML(htmlContent string) ([]GmailDetailTransaction, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		slog.Error("Error loading HTML", "error", err)
		return []GmailDetailTransaction{}, err
	}
	var detailTransactions []GmailDetailTransaction
	doc.Find("div").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		slog.Info("Text", "text", text)
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
		}
	})

	slog.Info("Detail Transactions", "detailTransactions", detailTransactions)
	return detailTransactions, nil
}

func parseCCBCAXML(htmlContent string) ([]GmailDetailTransaction, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		log.Println("Error loading HTML:", err)
		return []GmailDetailTransaction{}, err
	}
	var detailTransactions []GmailDetailTransaction
	doc.Find("table[style='width:100%']").Each(func(i int, table *goquery.Selection) {
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

// printMessageSummary prints a summary of a Gmail message and return the subject.
func printMessageSummary(msg *gmail.Message) string {
	timestamp := msg.InternalDate / 1000
	var subject string
	for _, header := range msg.Payload.Headers {
		if header.Name == "Subject" {
			subject = header.Value
			break
		}
	}
	fmt.Printf("%d - %s - %s\n", timestamp, msg.Id, subject)
	return subject
}

// authO2 demonstrates the OAuth2 flow and message processing.
func (a *Auth) authO2() {
	b, err := os.ReadFile(ClientSecretFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		slog.Error("Unable to create Gmail service", "error", err)
	}

	messages, err := listRecentMessages(srv, "me", 10)
	if err != nil {
		slog.Error("Unable to list messages", "error", err)
	}

	for _, m := range messages {
		log.Println(m.Id)
		// log.Println(printMessageSummary(m))
		log.Println("--------------------------------")
		// if m.Id == TargetMessageID {
		message, err := fetchMessageByID(srv, "me", m.Id)
		if err != nil {
			slog.Error("Unable to get message", "error", err)
		}
		emailTitle := printMessageSummary(message)
		log.Println("length ", len(message.Payload.Parts))
		for _, part := range message.Payload.Parts {
			log.Println("2- ", part.MimeType)
			if part.MimeType == "text/plain" || part.MimeType == "text/html" {
				body, err := decodeMessageBody(part)
				if err != nil {
					log.Printf("Error decoding body: %v", err)
					continue
				}
				if part.MimeType == "text/html" {
					detailTransactions, err := parseHTMLTables(body, emailTitle)
					if err != nil {
						log.Printf("Error parsing HTML tables: %v", err)
						continue
					}
					fmt.Println(detailTransactions)
					extractedData := &GmailExtractedData{
						Title:             emailTitle,
						DetailTransaction: detailTransactions,
					}
					fmt.Println(extractedData)
					err = a.service.createTransactionFromGmail(extractedData)
					if err != nil {
						slog.Error("Unable to create transaction from gmail", "error", err)
					}
				}
			} else if part.MimeType == "multipart/related" {
				for _, subPart := range part.Parts {
					log.Println("3 - ", subPart, " ", len(part.Parts))
					if subPart.MimeType == "text/html" {
						html, err := decodeMessageBody(subPart)
						if err != nil {
							slog.Error("Unable to decode message body", "error", err)
						}
						detailTransactions, err := parseHTMLTables(html, emailTitle)
						if err != nil {
							slog.Error("Unable to parse HTML tables", "error", err)
						}
						fmt.Println(detailTransactions)
						extractedData := &GmailExtractedData{
							Title:             emailTitle,
							DetailTransaction: detailTransactions,
						}
						fmt.Println(extractedData)
						// err = a.service.createTransactionFromGmail(extractedData)
						// if err != nil {
						// 	slog.Error("Unable to create transaction from gmail", "error", err)
						// }
					}
				}
			}
		}
		// }
	}
}
