package gmailclient

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type TokenSource struct {
	secretFile string
	tokenFile  string
	scopes     []string
	cfg        *oauth2.Config
}

func NewTokenSource(secretFile string, tokenFile string, scopes ...string) *TokenSource {
	return &TokenSource{
		secretFile: secretFile,
		tokenFile:  tokenFile,
		scopes:     scopes,
	}
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	oauth2cfg, err := t.oauthConfig()
	if err != nil {
		return nil, err
	}

	needSave := false
	tok, err := t.tokenFromFile()
	if err != nil {
		tok, err = t.tokenFromWeb(oauth2cfg)
		if err != nil {
			return nil, err
		}
		needSave = true
	} else if tok.Expiry.Before(time.Now()) {
		rs := oauth2cfg.TokenSource(context.TODO(), tok)
		tok, err = rs.Token()
		if err != nil {
			return nil, err
		}
		needSave = true
	}

	if needSave {
		err := t.saveToken(tok)
		if err != nil {
			log.Printf("Error saving token: %s", err.Error())
		}
	}
	return tok, nil
}

func (t *TokenSource) oauthConfig() (*oauth2.Config, error) {
	if t.cfg == nil {
		secret, err := ioutil.ReadFile(t.secretFile)
		if err != nil {
			return nil, err
		}

		cfg, err := google.ConfigFromJSON(secret, t.scopes...)
		if err != nil {
			return nil, err
		}
		t.cfg = cfg
	}
	return t.cfg, nil
}

// Taken from example here: https://developers.google.com/people/quickstart/go
func (t *TokenSource) tokenFromFile() (*oauth2.Token, error) {
	f, err := os.Open(t.tokenFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func (t *TokenSource) tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

func (t *TokenSource) saveToken(token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", t.tokenFile)
	f, err := os.OpenFile(t.tokenFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(token)
}
