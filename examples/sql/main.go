package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"bufio"
	"database/sql"
	"github.com/jckimble/golibsignal"
	"github.com/jckimble/golibsignal/store/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	signal := libsignal.Signal{
		Config:            sqlstore.NewConfig(db),
		IdentityStore:     sqlstore.NewIdentityStore(db),
		PreKeyStore:       sqlstore.NewPreKeyStore(db),
		SignedPreKeyStore: sqlstore.NewSignedPreKeyStore(db),
		SessionStore:      sqlstore.NewSessionStore(db),
		GroupStore:        sqlstore.NewGroupStore(db),
		ContactStore:      sqlstore.NewContactStore(db),
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}

	signal.Handler = libsignal.MessageHandlerFunc(func(m *libsignal.Message) {
		fmt.Printf("%+v\n", m)
	})
	if signal.NeedsRegistration() {
		if err := signal.RequestCode("sms"); err != nil {
			panic(err)
		}
		var text string
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter Code: ")
			text, _ = reader.ReadString('\n')
			text = strings.TrimSpace(text)
			if text != "" && len(text) == 7 {
				break
			} else if text == "" {
				fmt.Println("Code can't be empty")
			} else if len(text) != 7 {
				fmt.Println("Code must be in format 000-000")
			}
		}
		if err := signal.VerifyCode(text); err != nil {
			panic(err)
		}
		if err := signal.RegisterKeys(); err != nil {
			panic(err)
		}
	}
	if err := signal.ListenAndServe(); err != nil {
		panic(err)
	}
}
