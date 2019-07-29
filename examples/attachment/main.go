package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"bufio"
	"gitlab.com/jckimble/golibsignal"
	"gitlab.com/jckimble/golibsignal/store/file"
	"os"
	"strings"
)

func main() {
	signal := libsignal.Signal{
		Config:            &filestore.Config{},
		IdentityStore:     &filestore.IdentityStore{},
		PreKeyStore:       &filestore.PreKeyStore{},
		SignedPreKeyStore: &filestore.SignedPreKeyStore{},
		SessionStore:      &filestore.SessionStore{},
		GroupStore:        &filestore.GroupStore{},
		ContactStore:      &filestore.ContactStore{},
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
		if m.Group() != nil {
			file, _ := os.Open("video.mp4")
			defer file.Close()
			signal.SendGroupAttachment(m.Group().Hexid, "Testing", file)
		}
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
