package gmail

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"net/http"

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

	// List the messages in the user's inbox
	r, err := srv.Users.Messages.List("me").MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	// Print the subject of the last 10 emails
	if len(r.Messages) == 0 {
		fmt.Println("No messages found.")
	} else {
		fmt.Println("Messages:")
		for _, msg := range r.Messages {
			message, err := srv.Users.Messages.Get("me", msg.Id).Do()
			if err != nil {
				log.Fatalf("Unable to get message: %v", err)
			}
			for _, header := range message.Payload.Headers {
				if header.Name == "Subject" {
					fmt.Printf("- %s\n", header.Value)
				}
			}
		}
	}
}
