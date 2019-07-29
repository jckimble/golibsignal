package libsignal

import (
	"errors"
	"fmt"
	"gitlab.com/jckimble/golibsignal/signalservice"

	"github.com/golang/protobuf/proto"
)

// SyncSentHandler is required to recieve SyncSent Messages
type SyncSentHandler interface {
	SyncSent(*Message, uint64)
}

// SyncReadHandler is required to recieve SyncRead Messages
type SyncReadHandler interface {
	SyncRead(string, uint64)
}

func (s *Signal) sendSyncMessage(sm *signalservice.SyncMessage) (uint64, error) {
	tel, err := s.Config.GetTel()
	if err != nil {
		return 0, err
	}
	content := &signalservice.Content{
		SyncMessage: sm,
	}
	b, err := proto.Marshal(content)
	if err != nil {
		return 0, err
	}
	resp, err := s.buildAndSendMessage(tel, padMessage(b), true)
	return resp.Timestamp, err
}

func (s *Signal) handleSyncMessage(src string, timestamp uint64, sm *signalservice.SyncMessage) error {
	fmt.Printf("SyncMesssage: %+v\n", sm)
	if sm.GetSent() != nil {
		return s.handleSyncSent(sm.GetSent(), timestamp)
	} else if sm.GetRequest() != nil {
		return s.handleSyncRequest(sm.GetRequest())
	} else if sm.GetRead() != nil {
		return s.handleSyncRead(sm.GetRead())
	} else {
		return errors.New("SyncMessage contains no known sync types")
	}
}

func (s *Signal) handleSyncSent(sm *signalservice.SyncMessage_Sent, ts uint64) error {
	dm := sm.GetMessage()
	dest := sm.GetDestination()
	timestamp := sm.GetTimestamp()
	if dm == nil {
		return errors.New("DataMessage was nil for SyncMessage_Sent")
	}
	flags, err := s.handleFlags(dest, dm)
	if err != nil {
		return err
	}
	atts, err := s.handleAttachments(dm)
	if err != nil {
		return err
	}
	gr, err := s.handleGroups(dest, dm)
	if err != nil {
		return err
	}
	msg := &Message{
		source:      dest,
		message:     dm.GetBody(),
		attachments: atts,
		group:       gr,
		timestamp:   timestamp,
		flags:       flags,
	}
	if s.Handler != nil {
		if sh, ok := s.Handler.(SyncSentHandler); ok {
			sh.SyncSent(msg, ts)
		}
	}
	return nil
}

func (s *Signal) handleSyncRequest(request *signalservice.SyncMessage_Request) error {
	if request.GetType() == signalservice.SyncMessage_Request_GROUPS {
		return s.sendGroupUpdate()
	}
	return fmt.Errorf("Unimplemented: %s", request.GetType().String())
}

func (s *Signal) handleSyncRead(readMessages []*signalservice.SyncMessage_Read) error {
	if s.Handler != nil {
		if sh, ok := s.Handler.(SyncReadHandler); ok {
			for _, r := range readMessages {
				sh.SyncRead(r.GetSender(), r.GetTimestamp())
			}
		}
	}
	return nil
}
