package gmail

import (
	"drfin/app/shared/logdr"
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

type Auth struct {
}

func NewAuth() *Auth {
	return &Auth{}
}

const (
	TokenFile        = "token.json"
	ClientSecretFile = "resources/drfin-secret.json"
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

// authO2 demonstrates the OAuth2 flow and message processing.
func authO2() (*gmail.Service, error) {
	b, err := os.ReadFile(ClientSecretFile)
	if err != nil {
		logdr.Error("Unable to read client secret file: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		logdr.Error("Unable to parse client secret file to config: %v", err)
		return nil, err
	}

	client := getClient(config)

	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logdr.Error("Unable to create Gmail service", err)
		return nil, err
	}
	return srv, err
}
