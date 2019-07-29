package libsignal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	"encoding/base64"
	"encoding/binary"
	"encoding/hex"

	"bitbucket.org/taruti/mimemagic"
)

func decodeSignature(sig string) ([]byte, error) {
	b, err := base64DecodeNonPadded(sig)
	if err != nil {
		return nil, err
	}
	if len(b) != 64 {
		return nil, fmt.Errorf("signature is %d, not 64 bytes", len(b))
	}
	return b, nil
}

func encodeKey(key []byte) string {
	return base64EncWithoutPadding(append([]byte{5}, key[:]...))
}

func decodeKey(k string) ([]byte, error) {
	b, err := base64DecodeNonPadded(k)
	if err != nil {
		return nil, err
	}
	if len(b) != 33 || b[0] != 5 {
		return nil, errors.New("public key not formatted correctly")
	}
	return b[1:], nil
}

func base64EncWithoutPadding(b []byte) string {
	str := base64.StdEncoding.EncodeToString(b)
	return strings.TrimRight(str, "=")
}

func base64DecodeNonPadded(str string) ([]byte, error) {
	if len(str)%4 != 0 {
		str = str + strings.Repeat("=", 4-len(str)%4)
	}
	return base64.StdEncoding.DecodeString(str)
}

func idToHex(id []byte) string {
	return hex.EncodeToString(id)
}

func varint32(value int) []byte {
	buf := make([]byte, binary.MaxVarintLen32)
	n := binary.PutUvarint(buf, uint64(value))
	return buf[:n]
}

func stripPadding(msg []byte) []byte {
	for i := len(msg) - 1; i >= 0; i-- {
		if msg[i] == 0x80 {
			return msg[:i]
		}
	}
	return msg
}

func padMessage(msg []byte) []byte {
	l := (len(msg) + 160)
	l = l - l%160
	n := make([]byte, l)
	copy(n, msg)
	n[len(msg)] = 0x80
	return n
}

// MIMETypeFromReader guesses MIME Type From Reader
func MIMETypeFromReader(r io.Reader) (mime string, reader io.Reader) {
	var buf bytes.Buffer
	io.CopyN(&buf, r, 1024)
	mime = mimemagic.Match("", buf.Bytes())
	return mime, io.MultiReader(&buf, r)
}
