package libsignal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/jckimble/golibsignal/axolotl"
	"io"
	"strings"
)

// NeedsRegistration checks if libsignal is initalized
func (s *Signal) NeedsRegistration() bool {
	return !s.PreKeyStore.ContainsPreKey(0xFFFFFF)
}

// RequestCode sends code to phone for verification, allowed methods are sms and voice
func (s *Signal) RequestCode(method string) error {
	tel, err := s.Config.GetTel()
	if err != nil {
		return err
	}
	b := make([]byte, 16)
	randBytes(b[:])
	if err := s.Config.SetHTTPPassword(base64EncWithoutPadding(b)); err != nil {
		return err
	}
	return s.serverRequest("GET", fmt.Sprintf("/v1/accounts/%s/code/%s", method, tel), nil, "", map[int]interface{}{200: nil})
}

// VerifyCode verifies code sent to phone
func (s *Signal) VerifyCode(code string) error {
	code = strings.Replace(code, "-", "", -1)
	regid := randUint32() & 0x3fff
	if err := s.IdentityStore.SetLocalRegistrationID(regid); err != nil {
		return err
	}
	sk := make([]byte, 52)
	randBytes(sk[:])
	if err := s.Config.SetHTTPSignalingKey(sk); err != nil {
		return err
	}
	vd := struct {
		SignalingKey    string `json:"signalingKey"`
		RegistrationID  uint32 `json:"registrationId"`
		FetchesMessages bool   `json:"fetchesMessages"`
	}{
		SignalingKey:    base64.StdEncoding.EncodeToString(sk),
		FetchesMessages: true,
		RegistrationID:  regid,
	}
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		json.NewEncoder(w).Encode(vd)
	}()
	return s.serverRequest("PUT", fmt.Sprintf("/v1/accounts/code/%s", code), r, "application/json", map[int]interface{}{204: nil})
}

// RegisterKeys creates and uploads keys to signal server, must be ran after verification
func (s *Signal) RegisterKeys() error {
	identityKey := axolotl.GenerateIdentityKeyPair()
	if err := s.IdentityStore.SetIdentityKeyPair(identityKey); err != nil {
		return err
	}
	prekeys := generatePreKeys()
	pks := struct {
		IdentityKey   string          `json:"identityKey"`
		PreKeys       []*preKeyEntity `json:"preKeys"`
		LastResortKey *preKeyEntity   `json:"lastResortKey"`
		SignedPreKey  *preKeyEntity   `json:"signedPreKey"`
	}{
		PreKeys:       []*preKeyEntity{},
		LastResortKey: &preKeyEntity{},
		SignedPreKey:  &preKeyEntity{},
	}
	signedKey := generateSignedPreKey(identityKey)
	for i, key := range prekeys {
		pks.PreKeys = append(pks.PreKeys, &preKeyEntity{
			ID:        *key.Pkrs.Id,
			PublicKey: encodeKey(key.Pkrs.PublicKey),
		})
		if err := s.PreKeyStore.StorePreKey(i, key); err != nil {
			return err
		}
	}
	if err := s.SignedPreKeyStore.StoreSignedPreKey(*signedKey.Spkrs.Id, signedKey); err != nil {
		return err
	}
	pks.SignedPreKey.ID = *signedKey.Spkrs.Id
	pks.SignedPreKey.PublicKey = encodeKey(signedKey.Spkrs.PublicKey)
	pks.SignedPreKey.Signature = base64EncWithoutPadding(signedKey.Spkrs.Signature)
	pks.LastResortKey = pks.PreKeys[len(pks.PreKeys)-1]
	pks.IdentityKey = base64EncWithoutPadding(identityKey.PublicKey.Serialize())
	r, w := io.Pipe()
	go func() {
		defer w.Close()
		json.NewEncoder(w).Encode(pks)
	}()
	return s.serverRequest("PUT", "/v2/keys/", r, "application/json", map[int]interface{}{204: nil})
}

/*

   const getSocket = this.server.getProvisioningSocket.bind(this.server);
   const queueTask = this.queueTask.bind(this);
   const provisioningCipher = new libsignal.ProvisioningCipher();
*/
