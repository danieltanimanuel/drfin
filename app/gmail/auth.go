package gmail

import (
	"drfin/app/shared/logdr"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
	TokenDir         = "tokens"
	ClientSecretFile = "resources/drfin-secret.json"
)

// tokenFileFor returns the per-owner token path, e.g. tokens/1.json.
func tokenFileFor(ownerID int) string {
	return filepath.Join(TokenDir, fmt.Sprintf("%d.json", ownerID))
}

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

func saveToken(file string, token *oauth2.Token) error {
	if err := os.MkdirAll(filepath.Dir(file), 0700); err != nil {
		return err
	}
	tokJSON, err := json.Marshal(token)
	if err != nil {
		return err
	}
	return os.WriteFile(file, tokJSON, 0600)
}

func loadOAuthConfig() (*oauth2.Config, error) {
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
	return config, nil
}

// AuthorizeOwner runs the interactive OAuth flow for one owner and stores the
// token. CLI only (see main.go); never called during a sync.
func AuthorizeOwner(ownerID int) error {
	config, err := loadOAuthConfig()
	if err != nil {
		return err
	}
	authURL := config.AuthCodeURL("", oauth2.AccessTypeOffline)
	fmt.Printf("Sign in with the Gmail account of owner %d and open:\n%v\n", ownerID, authURL)
	var code string
	fmt.Print("Enter the authorization code: ")
	if _, err := fmt.Scan(&code); err != nil {
		return err
	}
	tok, err := config.Exchange(context.Background(), code)
	if err != nil {
		return fmt.Errorf("unable to retrieve token: %w", err)
	}
	if err := saveToken(tokenFileFor(ownerID), tok); err != nil {
		return err
	}
	fmt.Printf("Saved token to %s\n", tokenFileFor(ownerID))
	return nil
}

// authO2 builds a Gmail client from the owner's stored token. It never prompts.
func authO2(ownerID int) (*gmail.Service, error) {
	config, err := loadOAuthConfig()
	if err != nil {
		return nil, err
	}
	tok, err := tokenFromFile(tokenFileFor(ownerID))
	if err != nil {
		return nil, fmt.Errorf("owner %d is not authorized (run: go run . -authorize-owner %d): %w", ownerID, ownerID, err)
	}
	client := config.Client(context.Background(), tok)
	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		logdr.Error("Unable to create Gmail service", err)
		return nil, err
	}
	return srv, nil
}