package libsignal

import (
	"fmt"
	"strconv"

	"github.com/aebruno/textsecure/curve25519sign"
	"github.com/jckimble/golibsignal/axolotl"

	"time"
)

type preKeyEntity struct {
	ID        uint32 `json:"keyId"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature,omitempty"`
}

type preKeyResponseItem struct {
	DeviceID       uint32        `json:"deviceId"`
	RegistrationID uint32        `json:"registrationId"`
	SignedPreKey   *preKeyEntity `json:"signedPreKey"`
	PreKey         *preKeyEntity `json:"preKey"`
}

type preKeyResponse struct {
	IdentityKey string `json:"identityKey"`
	Devices     []preKeyResponseItem
}

func (s *Signal) makePreKeyBundle(tel string, deviceID uint32) (*axolotl.PreKeyBundle, error) {
	var pkr preKeyResponse
	if err := s.serverRequest("GET", fmt.Sprintf("/v2/keys/%s/%s", tel, strconv.Itoa(int(deviceID))), nil, "", map[int]interface{}{
		200: &pkr,
	}); err != nil {
		return nil, err
	}
	if len(pkr.Devices) != 1 {
		return nil, fmt.Errorf("no prekeys for contact %s, device %d\n", tel, deviceID)
	}
	d := pkr.Devices[0]
	if d.PreKey == nil {
		return nil, fmt.Errorf("no prekey for contact %s, device %d\n", tel, deviceID)
	}
	decPK, err := decodeKey(d.PreKey.PublicKey)
	if err != nil {
		return nil, err
	}
	if d.SignedPreKey == nil {
		return nil, fmt.Errorf("no signed prekey for contact %s, device %d\n", tel, deviceID)
	}
	decSPK, err := decodeKey(d.SignedPreKey.PublicKey)
	if err != nil {
		return nil, err
	}
	decSig, err := decodeSignature(d.SignedPreKey.Signature)
	if err != nil {
		return nil, err
	}
	decIK, err := decodeKey(pkr.IdentityKey)
	if err != nil {
		return nil, err
	}
	pkb, err := axolotl.NewPreKeyBundle(d.RegistrationID, d.DeviceID, d.PreKey.ID, axolotl.NewECPublicKey(decPK), int32(d.SignedPreKey.ID), axolotl.NewECPublicKey(decSPK), decSig, axolotl.NewIdentityKey(decIK))
	if err != nil {
		return nil, err
	}
	return pkb, nil
}
func generateSignedPreKey(identityKey *axolotl.IdentityKeyPair) *axolotl.SignedPreKeyRecord {
	kp := axolotl.NewECKeyPair()
	id := randUint32() & 0xffffff
	var random [64]byte
	randBytes(random[:])
	priv := identityKey.PrivateKey.Key()
	signature := curve25519sign.Sign(priv, kp.PublicKey.Serialize(), random)
	record := axolotl.NewSignedPreKeyRecord(id, uint64(time.Now().UnixNano()*1000), kp, signature[:])
	return record
}

func generatePreKeys() map[uint32]*axolotl.PreKeyRecord {
	m := map[uint32]*axolotl.PreKeyRecord{}
	startID := randUint32() & 0xffffff
	for i := 0; i < 100; i++ {
		kp := axolotl.NewECKeyPair()
		m[startID+uint32(i)] = axolotl.NewPreKeyRecord(startID+uint32(i), kp)
	}
	kp := axolotl.NewECKeyPair()
	m[0xffffff] = axolotl.NewPreKeyRecord(0xffffff, kp)
	return m
}
