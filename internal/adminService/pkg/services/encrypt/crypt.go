package encrypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type (
	Hmac struct {
		params *HmacParams
	}
	HmacParams struct {
		SecretKey string
	}
)

func NewHmac(params *HmacParams) *Hmac {
	return &Hmac{params: params}
}

func (c *Hmac) EncryptString(msg string) string {
	mac := hmac.New(sha256.New, []byte(c.params.SecretKey))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func (c *Hmac) VerifySignature(msg string, hash string) (bool, error) {
	signature, err := hex.DecodeString(hash)
	if err != nil {
		return false, err
	}

	mac := hmac.New(sha256.New, []byte(c.params.SecretKey))
	mac.Write([]byte(msg))

	return hmac.Equal(signature, mac.Sum(nil)), nil
}

//pontus.poh.l@gmail.com1637050708
//pontus.poh.l@gmail.com1637048767799163000
