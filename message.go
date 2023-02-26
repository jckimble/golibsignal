package libsignal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jckimble/golibsignal/axolotl"
	"github.com/jckimble/golibsignal/signalservice"

	"log"
)

type messageHandler struct {
	f func(*Message)
}

func (mh messageHandler) Message(m *Message) {
	mh.f(m)
}

func MessageHandlerFunc(f func(*Message)) MessageHandler {
	return messageHandler{f}
}

// MessageHandler is an interface for recieving signal messages
type MessageHandler interface {
	Message(*Message)
}

// Message holds incoming message information
type Message struct {
	source      string
	message     string
	attachments []*Attachment
	group       *Group
	timestamp   uint64
	flags       uint32
	contact     *Contact
}

// Contact returns Contact from incoming Message
func (m Message) Contact() *Contact {
	return m.contact
}

// Source returns Message Source Phone Number
func (m Message) Source() string {
	return m.source
}

// Source returns Message Text
func (m Message) Message() string {
	return m.message
}

// Attachments returns Message Attachments
func (m Message) Attachments() []*Attachment {
	return m.attachments
}

// Group returns Group message was sent to, nil if direct
func (m Message) Group() *Group {
	return m.group
}

// Timestamp returns Message Timestamp
func (m Message) Timestamp() uint64 {
	return m.timestamp
}

// Flags returns Message Flags
func (m Message) Flags() uint32 {
	return m.flags
}

type sendResponse struct {
	NeedsSync bool `json:"needsSync"`
	Timestamp uint64
}
type mismatchedResponse struct {
	MissingDevices []uint32 `json:"missingDevices"`
	ExtraDevices   []uint32 `json:"extraDevices"`
}
type staleResponse struct {
	StaleDevices []uint32 `json:"staleDevices"`
}
type jsonMessage struct {
	Type               int32  `json:"type"`
	DestDeviceID       uint32 `json:"destinationDeviceId"`
	DestRegistrationID uint32 `json:"destinationRegistrationId"`
	Content            string `json:"content"`
	Relay              string `json:"relay,omitempty"`
}
type groupMessage struct {
	id      []byte
	name    string
	members []string
	typ     signalservice.GroupContext_Type
}
type outgoingMessage struct {
	tel         string
	msg         string
	group       *groupMessage
	attachments []*signalservice.AttachmentPointer
	flags       uint32
}

// SendMessage sends the given text message to the given number
func (s *Signal) SendMessage(tel, msg string) (uint64, error) {
	omsg := &outgoingMessage{
		tel: tel,
		msg: msg,
	}
	return s.sendMessage(omsg)
}

func (s *Signal) sendMessage(msg *outgoingMessage) (uint64, error) {
	dm := &signalservice.DataMessage{}
	if msg.msg != "" {
		dm.Body = &msg.msg
	}
	if msg.attachments != nil {
		dm.Attachments = msg.attachments
	}
	if msg.group != nil {
		dm.Group = &signalservice.GroupContext{
			Id:      msg.group.id,
			Type:    &msg.group.typ,
			Name:    &msg.group.name,
			Members: msg.group.members,
		}
	}
	dm.Flags = &msg.flags
	content := &signalservice.Content{
		DataMessage: dm,
	}
	b, err := proto.Marshal(content)
	if err != nil {
		return 0, err
	}
	resp, err := s.buildAndSendMessage(msg.tel, padMessage(b), false)
	if err != nil {
		return 0, err
	}
	if resp.NeedsSync {
		sm := &signalservice.SyncMessage{
			Sent: &signalservice.SyncMessage_Sent{
				Destination: &msg.tel,
				Timestamp:   &resp.Timestamp,
				Message:     dm,
			},
		}
		s.sendSyncMessage(sm)
	}
	return resp.Timestamp, err
}

func (s *Signal) buildMessage(tel string, paddedMessage []byte, isSync bool) ([]jsonMessage, error) {
	recid := tel[1:]
	messages := []jsonMessage{}

	contact, err := s.ContactStore.Get(tel)
	if err != nil {
		log.Printf("Unknown Contact: %s %s", tel, err)
		contact = &Contact{
			Tel:     tel,
			Devices: []uint32{1},
		}
	}
	for _, devid := range contact.Devices {
		if !s.SessionStore.ContainsSession(recid, devid) {
			pkb, err := s.makePreKeyBundle(tel, devid)
			if err != nil {
				return nil, err
			}
			sb := axolotl.NewSessionBuilder(s.IdentityStore, s.PreKeyStore, s.SignedPreKeyStore, s.SessionStore, recid, pkb.DeviceID)
			if err := sb.BuildSenderSession(pkb); err != nil {
				return nil, err
			}
		}
		sc := axolotl.NewSessionCipher(s.IdentityStore, s.PreKeyStore, s.SignedPreKeyStore, s.SessionStore, recid, devid)
		encryptedMessage, messageType, err := sc.SessionEncryptMessage(paddedMessage)
		if err != nil {
			return nil, err
		}
		rrID, err := sc.GetRemoteRegistrationID()
		if err != nil {
			return nil, err
		}
		jmsg := jsonMessage{
			Type:               messageType,
			DestDeviceID:       devid,
			DestRegistrationID: rrID,
			Content:            base64.StdEncoding.EncodeToString(encryptedMessage),
		}
		messages = append(messages, jmsg)
	}
	return messages, nil
}

func (s *Signal) buildAndSendMessage(tel string, paddedMessage []byte, isSync bool) (*sendResponse, error) {
	bm, err := s.buildMessage(tel, paddedMessage, isSync)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	m["messages"] = bm
	m["timestamp"] = uint64(time.Now().UnixNano() / 1000000)
	m["destination"] = tel
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		json.NewEncoder(w).Encode(m)
	}()
	var rsp sendResponse
	var stale staleResponse
	var mismatched mismatchedResponse
	if err := s.serverRequest("PUT", fmt.Sprintf("/v1/messages/%s", tel), r, "application/json", map[int]interface{}{
		200: &rsp,
		410: &stale,
		409: &mismatched,
	}); err != nil {
		return nil, err
	}
	if stale.StaleDevices != nil {
		for _, id := range stale.StaleDevices {
			s.SessionStore.DeleteSession(tel[1:], id)
		}
		return s.buildAndSendMessage(tel, paddedMessage, isSync)
	}
	if mismatched.MissingDevices != nil || mismatched.ExtraDevices != nil {
		contact, err := s.ContactStore.Get(tel)
		if err != nil {
			log.Printf("Unknown Contact: %s %s", tel, err)
			contact = &Contact{
				Tel:     tel,
				Devices: []uint32{1},
			}
		}
		devs := []uint32{}
		for _, id := range contact.Devices {
			in := true
			for _, eid := range mismatched.ExtraDevices {
				if id == eid {
					in = false
					break
				}
			}
			if in {
				devs = append(devs, id)
			}
		}
		contact.Devices = append(devs, mismatched.MissingDevices...)
		if err := s.ContactStore.Save(contact); err != nil {
			return nil, err
		}
		return s.buildAndSendMessage(tel, paddedMessage, isSync)
	}
	rsp.Timestamp = m["timestamp"].(uint64)
	return &rsp, nil
}
