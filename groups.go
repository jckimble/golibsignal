package libsignal

import (
	"github.com/jckimble/golibsignal/signalservice"

	"errors"
	"fmt"
	"io"

	"bytes"
	"github.com/golang/protobuf/proto"

	"log"
)

// GroupStore is an interface for storing Signal Groups
type GroupStore interface {
	Get(string) (*Group, error)
	GetAll() ([]*Group, error)
	Save(*Group) error
}

// Group holds signal group information for sending group messages
type Group struct {
	ID      []byte
	Hexid   string
	Flags   uint32
	Name    string
	Members []string
	Avatar  *Attachment
}

// SendGroupMessage sends a text message to a given group
func (s *Signal) SendGroupMessage(hexid string, msg string) (uint64, error) {
	return s.sendGroup(hexid, msg, nil)
}

// SendGroupAttachment sends the contents of a reader, along with an optional message to a given group.
func (s *Signal) SendGroupAttachment(hexid, msg string, r io.Reader) (uint64, error) {
	ct, r := MIMETypeFromReader(r)
	a, err := s.uploadAttachment(r, ct)
	if err != nil {
		return 0, err
	}
	return s.sendGroup(hexid, msg, a)
}

func (s *Signal) sendGroup(hexid, msg string, a *signalservice.AttachmentPointer) (uint64, error) {
	group, err := s.GroupStore.Get(hexid)
	if err != nil {
		return 0, fmt.Errorf("Unknown Group: %s", err)
	}
	if group.Members == nil {
		return 0, errors.New("Unknown Group")
	}
	tel, err := s.Config.GetTel()
	if err != nil {
		return 0, err
	}
	var ts uint64
	for _, m := range group.Members {
		if m != tel {
			omsg := &outgoingMessage{
				tel: m,
				msg: msg,
				group: &groupMessage{
					id:  group.ID,
					typ: signalservice.GroupContext_DELIVER,
				},
			}
			if a != nil {
				omsg.attachments = []*signalservice.AttachmentPointer{a}
			}
			ts, err = s.sendMessage(omsg)
			if err != nil {
				return 0, err
			}
		}
	}
	return ts, nil
}

func (s *Signal) sendGroupUpdate() error {
	var buf bytes.Buffer
	groups, err := s.GroupStore.GetAll()
	if err != nil {
		return err
	}
	for _, g := range groups {
		gd := &signalservice.GroupDetails{
			Id:      g.ID,
			Name:    &g.Name,
			Members: g.Members,
		}
		b, err := proto.Marshal(gd)
		if err != nil {
			return err
		}
		buf.Write(varint32(len(b)))
		buf.Write(b)
	}
	a, err := s.uploadAttachment(&buf, "application/octet-stream")
	if err != nil {
		return err
	}
	sm := &signalservice.SyncMessage{
		Groups: &signalservice.SyncMessage_Groups{
			Blob: a,
		},
	}
	_, err = s.sendSyncMessage(sm)
	return err
}

func (s *Signal) handleGroups(src string, dm *signalservice.DataMessage) (*Group, error) {
	gr := dm.GetGroup()
	if gr == nil {
		return nil, nil
	}
	if s.GroupStore == nil {
		return nil, nil
	}
	group, err := s.GroupStore.Get(idToHex(gr.GetId()))
	if err != nil {
		log.Printf("Unknown Group: %s %s", idToHex(gr.GetId()), err)
		group = &Group{}
	}
	switch gr.GetType() {
	case signalservice.GroupContext_UPDATE:
		group.ID = gr.GetId()
		group.Hexid = idToHex(gr.GetId())
		group.Name = gr.GetName()
		group.Members = gr.GetMembers()
		if av := gr.GetAvatar(); av != nil {
			att, err := s.handleSingleAttachment(av)
			if err != nil {
				return nil, err
			}
			group.Avatar = att
		}
		if err := s.GroupStore.Save(group); err != nil {
			return nil, err
		}
		group.Flags = uint32(1)
	case signalservice.GroupContext_DELIVER:
		group.Flags = 0
	case signalservice.GroupContext_QUIT:
		for i, m := range group.Members {
			if m == src {
				group.Members = append(group.Members[:i], group.Members[i+1:]...)
				break
			}
		}
		if err := s.GroupStore.Save(group); err != nil {
			return nil, err
		}
		group.Flags = uint32(2)
	default:
		return nil, fmt.Errorf("Invalid Group Context: %s", gr.GetType().String())
	}
	return group, nil
}
