package libsignal

import (
	"crypto/tls"
	"net"
	"net/url"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/jckimble/golibsignal/axolotl"
	"github.com/jckimble/golibsignal/signalservice"

	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	writeWait     = 25 * time.Second
	pongWait      = 60 * time.Second
	pingPeriod    = (pongWait * 9) / 10
	websocketPath = "/v1/websocket/"
)

// Signal is golibsignal's Main Struct Config, HTTPClient, and All Stores required
type Signal struct {
	Config            Config
	IdentityStore     axolotl.IdentityStore
	PreKeyStore       axolotl.PreKeyStore
	SignedPreKeyStore axolotl.SignedPreKeyStore
	SessionStore      axolotl.SessionStore
	GroupStore        GroupStore
	ContactStore      ContactStore
	HTTPClient        *http.Client

	Handler interface{}

	ws *websocket.Conn

	send chan []byte
}

// Shutdown closes Websocket to stop receiving messages
func (s *Signal) Shutdown() error {
	if s.ws == nil {
		return errors.New("no listening connection to stop")
	}
	return s.ws.Close()
}

func (s *Signal) writeWorker() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.ws.Close()
	}()
	for {
		select {
		case message, ok := <-s.send:
			if !ok {
				s.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := s.write(websocket.BinaryMessage, message); err != nil {
				log.Printf("%s", err)
				return
			}
		case <-ticker.C:
			if err := s.write(websocket.PingMessage, nil); err != nil {
				log.Printf("%s", err)
				return
			}
		}
	}
}

func (s *Signal) write(mt int, payload []byte) error {
	s.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return s.ws.WriteMessage(mt, payload)
}

func (s *Signal) sendAck(id uint64) error {
	typ := signalservice.WebSocketMessage_RESPONSE
	message := "OK"
	status := uint32(200)
	wsm := &signalservice.WebSocketMessage{
		Type: &typ,
		Response: &signalservice.WebSocketResponseMessage{
			Id:      &id,
			Status:  &status,
			Message: &message,
		},
	}
	b, err := proto.Marshal(wsm)
	if err != nil {
		return err
	}
	s.send <- b
	return nil
}

// ListenAndServe connects to the server and handles incoming messages
func (s *Signal) ListenAndServe() error {
	if s.ws != nil {
		return errors.New("Already Listening")
	}
	if s.Config == nil {
		return errors.New("Config Not Configured")
	}
	s.send = make(chan []byte, 256)

	pass, err := s.Config.GetHTTPPassword()
	if err != nil {
		return err
	}
	tel, err := s.Config.GetTel()
	if err != nil {
		return err
	}
	server, err := s.Config.GetServer()
	if err != nil {
		return err
	}
	v := url.Values{}
	v.Set("login", tel)
	v.Set("password", pass)
	params := v.Encode()
	wsURL := strings.Replace(server, "http", "ws", 1) + "?" + params
	u, err := url.Parse(wsURL)
	u.Path = websocketPath
	if err != nil {
		return err
	}
	d := &websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	d.NetDial = func(network, addr string) (net.Conn, error) { return net.Dial(network, u.Host) }
	d.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	ws, _, err := d.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	s.ws = ws

	defer s.ws.Close()

	go s.writeWorker()
	s.ws.SetReadDeadline(time.Now().Add(pongWait))
	s.ws.SetPongHandler(func(string) error {
		s.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, bmsg, err := s.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				return err
			}
			return nil
		}
		wsm := &signalservice.WebSocketMessage{}
		if err := proto.Unmarshal(bmsg, wsm); err != nil {
			return err
		}
		m := wsm.GetRequest().GetBody()
		if len(m) > 0 {
			if err := s.handleRequestMessage(m); err != nil {
				log.Printf("%s", err)
			}
		}
		if err := s.sendAck(wsm.GetRequest().GetId()); err != nil {
			return err
		}
	}
}

func (s *Signal) handleRequestMessage(msg []byte) error {
	macpos := len(msg) - 10
	tmac := msg[macpos:]
	signalingKey, err := s.Config.GetHTTPSignalingKey()
	if err != nil {
		return err
	}
	aesKey := signalingKey[:32]
	macKey := signalingKey[32:]
	if !axolotl.ValidTruncMAC(msg[:macpos], tmac, macKey) {
		return errors.New("Invalid MAC")
	}
	ciphertext := msg[1:macpos]
	plaintext, err := axolotl.Decrypt(aesKey, ciphertext)
	if err != nil {
		return err
	}
	env := &signalservice.Envelope{}
	if err := proto.Unmarshal(plaintext, env); err != nil {
		return err
	}
	recid := env.GetSource()[1:]
	sc := axolotl.NewSessionCipher(s.IdentityStore, s.PreKeyStore, s.SignedPreKeyStore, s.SessionStore, recid, env.GetSourceDevice())
	switch *env.Type {
	case signalservice.Envelope_RECEIPT:
		return nil
	case signalservice.Envelope_CIPHERTEXT:
		msg := env.GetContent()
		if msg == nil {
			return errors.New("Legacy messages unsupported")
		}
		wm, err := axolotl.LoadWhisperMessage(msg)
		if err != nil {
			return err
		}
		b, err := sc.SessionDecryptWhisperMessage(wm)
		if _, ok := err.(axolotl.DuplicateMessageError); ok {
			return nil
		}
		if _, ok := err.(axolotl.InvalidMessageError); ok {
			return nil
		}
		if err != nil {
			return err
		}
		if err := s.handleMessage(env, b); err != nil {
			return err
		}
	case signalservice.Envelope_PREKEY_BUNDLE:
		msg := env.GetContent()
		pkwm, err := axolotl.LoadPreKeyWhisperMessage(msg)
		if err != nil {
			return err
		}
		b, err := sc.SessionDecryptPreKeyWhisperMessage(pkwm)
		if _, ok := err.(axolotl.DuplicateMessageError); ok {
			return nil
		}
		if _, ok := err.(axolotl.PreKeyNotFoundError); ok {
			return nil
		}
		if _, ok := err.(axolotl.InvalidMessageError); ok {
			return nil
		}
		if err != nil {
			return err
		}
		if err := s.handleMessage(env, b); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Non-Implemented Envelope Message: %s", (*env.Type).String())
	}
	return nil
}
func (s *Signal) handleMessage(env *signalservice.Envelope, b []byte) error {
	tel, err := s.Config.GetTel()
	if err != nil {
		return err
	}
	b = stripPadding(b)
	content := &signalservice.Content{}
	if err := proto.Unmarshal(b, content); err != nil {
		return err
	}
	if dm := content.GetDataMessage(); dm != nil {
		return s.handleDataMessage(env.GetSource(), env.GetTimestamp(), dm)
	} else if sm := content.GetSyncMessage(); sm != nil && tel == env.GetSource() {
		return s.handleSyncMessage(env.GetSource(), env.GetTimestamp(), sm)
	}
	return fmt.Errorf("Non-Implemented Content Message: %s", content.String())
}

func (s *Signal) handleDataMessage(src string, timestamp uint64, dm *signalservice.DataMessage) error {
	flags, err := s.handleFlags(src, dm)
	if err != nil {
		return err
	}
	atts, err := s.handleAttachments(dm)
	if err != nil {
		return err
	}
	gr, err := s.handleGroups(src, dm)
	if err != nil {
		return err
	}
	contact, err := s.handleContact(src, dm)
	if err != nil {
		return err
	}
	msg := &Message{
		source:      src,
		message:     dm.GetBody(),
		attachments: atts,
		group:       gr,
		timestamp:   timestamp,
		flags:       flags,
		contact:     contact,
	}
	if s.Handler != nil {
		if mh, ok := s.Handler.(MessageHandler); ok {
			mh.Message(msg)
		}
	}
	return nil
}

func (s *Signal) handleFlags(src string, dm *signalservice.DataMessage) (uint32, error) {
	flags := uint32(0)
	if dm.GetFlags() == uint32(signalservice.DataMessage_END_SESSION) {
		flags = uint32(1)
		s.SessionStore.DeleteAllSessions(src[1:])
	}
	return flags, nil
}
