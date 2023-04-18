package eddilithium2jwt

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/cloudflare/circl/sign/eddilithium2"
	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEdDilithium2Verification = errors.New("eddilithium2: verification error")
)

type SigningMethodEdDilithium struct{}

var (
	SigningMethodEdDilithium2 *SigningMethodEdDilithium
)

func init() {
	SigningMethodEdDilithium2 = &SigningMethodEdDilithium{}
	jwt.RegisterSigningMethod(SigningMethodEdDilithium2.Alg(), func() jwt.SigningMethod {
		return SigningMethodEdDilithium2
	})
}

func (m *SigningMethodEdDilithium) Alg() string {
	return "EdDilithium2"
}

func (m *SigningMethodEdDilithium) Verify(signingString string, sig []byte, key interface{}) error {
	var pubKey eddilithium2.PublicKey
	var ok bool

	if pubKey, ok = key.(eddilithium2.PublicKey); !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(pubKey.Bytes()) != eddilithium2.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	if !eddilithium2.Verify(&pubKey, []byte(signingString), sig) {
		return ErrEdDilithium2Verification
	}

	return nil
}

func (m *SigningMethodEdDilithium) Sign(signingString string, key interface{}) ([]byte, error) {
	var privateKey eddilithium2.PrivateKey
	var ok bool

	if privateKey, ok = key.(eddilithium2.PrivateKey); !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	if _, ok := privateKey.Public().(*eddilithium2.PublicKey); !ok {
		return nil, jwt.ErrInvalidKey
	}

	sig := make([]byte, eddilithium2.SignatureSize)
	eddilithium2.SignTo(&privateKey, []byte(signingString), sig)
	return sig, nil
}

func ValidateIssuedAt(t *jwt.Token) error {
	when, err := t.Claims.GetIssuedAt()
	if err != nil {
		return err
	}
	now := time.Now()
	later := now.Add(30 * time.Second)
	recently := now.Add(-30 * time.Second)
	if when.After(later) || when.Before(recently) {
		return errors.New("token issued outside of acceptable window")
	}
	return nil
}

// signingKey is a base64url encoded string, passed, e.g, via env var
func NewSignedString(signingID string, signingKey string) (string, error) {
	token := jwt.NewWithClaims(SigningMethodEdDilithium2,
		jwt.MapClaims{"iss": signingID, "iat": time.Now().Unix()})
	skBytes, err := base64.RawURLEncoding.DecodeString(signingKey)
	if err != nil {
		return "", err
	}

	var sk eddilithium2.PrivateKey
	err = sk.UnmarshalBinary(skBytes)
	if err != nil {
		return "", err
	}

	return token.SignedString(sk)
}
