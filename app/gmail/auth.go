package gmail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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

func getClient(config *oauth2.Config) *http.Client {
	// Check if we have a token in the file
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Retrieve a token from a local file
func tokenFromFile(file string) (*oauth2.Token, error) {
	tokJSON, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	err = json.Unmarshal(tokJSON, tok)
	return tok, err
}

// Get a token from the web
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

// Save the token to a file for later use
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	tokJSON, _ := json.Marshal(token)
	os.WriteFile(file, tokJSON, 0600)
}

func authO2() {
	// Set up OAuth2 config with credentials
	clientSecretFile := "resources/drfin-secret.json"
	b, err := os.ReadFile(clientSecretFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	// Create Gmail service
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Gmail service: %v", err)
	}

	// Assuming 'srv' is your Gmail service and you have the message ID of the older message.
	messageId := "196152a4ef35e1fa" // Replace with the message ID you want to filter by

	// Get the message to retrieve its date
	msg, err := srv.Users.Messages.Get("me", messageId).Do()
	if err != nil {
		log.Fatalf("Unable to get message: %v", err)
	}

	// Extract the date in a format suitable for the query, here we use the "InternalDate" field
	timestamp := msg.InternalDate // This is in milliseconds, convert to seconds

	// Use the "after" operator to get messages newer than the specified timestamp
	query := fmt.Sprintf("after:%d", (timestamp / 1000)) // divide by 1000 to convert to seconds
	log.Printf("%s", query)

	// Get the newer messages using the query
	messages, err := srv.Users.Messages.List("me").MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to list messages: %v", err)
	}

	for _, msg := range messages.Messages {
		message, err := srv.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			log.Fatalf("Unable to get message: %v", err)
		}
		tmp := message.InternalDate / 1000
		for _, header := range message.Payload.Headers {
			if header.Name == "Subject" {
				fmt.Printf("%d - %s - %s\n", tmp, msg.Id, header.Value)
			}
		}
		// Check for the body in the parts array
		var body string
		if msg.Id == "1965804f844ba06d" {
			for _, part := range message.Payload.Parts {
				// Check for the body if it's text or html part
				log.Printf("Testis : %s", part.MimeType)
				if part.MimeType == "text/plain" || part.MimeType == "text/html" {
					body += part.Body.Data
				}

				// If the part is base64url encoded, decode it
				if part.Body.Size > 0 {
					decoded, err := base64.URLEncoding.DecodeString(part.Body.Data)
					if err != nil {
						log.Fatalf("Error decoding body: %v", err)
					}
					body = string(decoded)
					// Load the HTML content into goquery
					doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
					if err != nil {
						log.Fatal("Error loading HTML: ", err)
					}

					doc.Find("table[style='width:100%']").Each(func(i int, table *goquery.Selection) {
						// Iterate over each row in the table
						table.Find("tr").Each(func(i int, tr *goquery.Selection) {
							// Extract the text of each <td> element in the row
							cells := tr.Find("td").Map(func(i int, td *goquery.Selection) string {
								fmt.Print(td.Each())
								return strings.TrimSpace(td.Text())
							})

							// Print the row if it has content
							if len(cells) > 0 {
								fmt.Println(cells)
							}
						})
					})
				}
			}
		}
	}
}
