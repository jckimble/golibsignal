package libsignal

import (
	"bytes"
	"gitlab.com/jckimble/golibsignal/signalservice"
	"io"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"crypto/sha256"
	"strconv"
)

// Attachment represents an attachment received from server
type Attachment struct {
	R    io.Reader
	Type string
}

// Reader returns io.Reader for attachment
func (a Attachment) Reader() io.Reader {
	return a.R
}

// MimeType returns Attachment MIME Type
func (a Attachment) MimeType() string {
	return a.Type
}

type jsonAttachment struct {
	ID       uint64 `json:"id"`
	Location string `json:"location"`
}

func (s *Signal) serverRequest(proto string, path string, data io.Reader, ct string, ret map[int]interface{}) error {
	server, err := s.Config.GetServer()
	if err != nil {
		return err
	}
	u, err := url.Parse(server)
	if err != nil {
		return err
	}
	u.Path = path
	req, err := http.NewRequest(proto, u.String(), data)
	if err != nil {
		return err
	}
	if ct != "" {
		req.Header.Add("Content-Type", ct)
	}
	/*	if config.UserAgent != "" {
		req.Header.Set("X-Signal-Agent", config.UserAgent)
	}*/
	pass, err := s.Config.GetHTTPPassword()
	if err != nil {
		return err
	}
	tel, err := s.Config.GetTel()
	if err != nil {
		return err
	}
	req.SetBasicAuth(tel, pass)
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if ret == nil {
		return nil
	}
	if resp.StatusCode == 413 {
		return errors.New("You have been ratelimited: Try Again Later")
	}
	if _, ok := ret[resp.StatusCode]; !ok {
		data, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(data))
		return fmt.Errorf("Invalid Code: %d", resp.StatusCode)
	}
	if ret[resp.StatusCode] == nil {
		return nil
	}
	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(ret[resp.StatusCode]); err != nil {
		return err
	}
	return nil
}
func (s *Signal) getAttachmentLocation(id uint64) (string, error) {
	var a jsonAttachment
	if err := s.serverRequest("GET", fmt.Sprintf("/v1/attachments/%d", id), nil, "", map[int]interface{}{
		200: &a,
	}); err != nil {
		return "", err
	}
	return a.Location, nil
}

func (s *Signal) getAttachment(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/octet-stream")
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (s *Signal) handleAttachments(dm *signalservice.DataMessage) ([]*Attachment, error) {
	atts := dm.GetAttachments()
	if atts == nil {
		return nil, nil
	}
	all := make([]*Attachment, len(atts))
	var err error
	for i, a := range atts {
		all[i], err = s.handleSingleAttachment(a)
		if err != nil {
			return nil, err
		}
	}
	return all, nil
}

func (s *Signal) handleSingleAttachment(a *signalservice.AttachmentPointer) (*Attachment, error) {
	loc, err := s.getAttachmentLocation(*a.Id)
	if err != nil {
		return nil, err
	}
	r, err := s.getAttachment(loc)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	l := len(b) - 32
	if !verifyMAC(a.Key[32:], b[:l], b[l:]) {
		return nil, errors.New("invalid MAC for attachment")
	}
	b, err = aesDecrypt(a.Key[:32], b[:l])
	if err != nil {
		return nil, err
	}
	return &Attachment{bytes.NewReader(b), a.GetContentType()}, nil
}

func (s *Signal) uploadAttachment(r io.Reader, ct string) (*signalservice.AttachmentPointer, error) {
	keys := make([]byte, 64)
	randBytes(keys)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	plaintextLength := uint32(len(b))
	e, err := aesEncrypt(keys[:32], b)
	if err != nil {
		return nil, err
	}
	m := appendMAC(keys[32:], e)
	id, location, err := s.allocateAttachment()
	if err != nil {
		return nil, err
	}
	digest, err := s.putAttachment(location, m)
	if err != nil {
		return nil, err
	}
	ap := &signalservice.AttachmentPointer{
		Id:          &id,
		ContentType: &ct,
		Key:         keys,
		Size:        &plaintextLength,
		Digest:      digest,
	}
	return ap, nil
}
func (s *Signal) putAttachment(url string, body []byte) ([]byte, error) {
	br := bytes.NewReader(body)
	req, err := http.NewRequest("PUT", url, br)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/octet-stream")
	req.Header.Add("Content-Length", strconv.Itoa(len(body)))

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP status %d\n", resp.StatusCode)
	}

	hasher := sha256.New()
	hasher.Write(body)

	return hasher.Sum(nil), nil
}
func (s *Signal) allocateAttachment() (uint64, string, error) {
	var a jsonAttachment
	if err := s.serverRequest("GET", "/v1/attachments/", nil, "", map[int]interface{}{
		200: &a,
	}); err != nil {
		return 0, "", err
	}
	return a.ID, a.Location, nil
}

// SendAttachment sends the contents of a reader, along with an optional message to a given number.
func (s *Signal) SendAttachment(tel, msg string, r io.Reader) (uint64, error) {
	ct, r := MIMETypeFromReader(r)
	a, err := s.uploadAttachment(r, ct)
	if err != nil {
		return 0, err
	}
	omsg := &outgoingMessage{
		tel:         tel,
		msg:         msg,
		attachments: []*signalservice.AttachmentPointer{a},
	}
	return s.sendMessage(omsg)
}
