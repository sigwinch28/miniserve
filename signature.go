package main

import (
	"crypto"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"aead.dev/minisign"
	"golang.org/x/crypto/blake2b"
)

const digestSize = blake2b.Size

const commentType = "minisig"
const commentVersion = "1"

type digest [digestSize]byte

func decodeHexDigest(hexDigest []byte) (digest, error) {
	var digest digest

	if len(hexDigest) != hex.EncodedLen(digestSize) {
		return digest, errors.New("hex digest is the wrong length")
	}

	_, err := hex.Decode(digest[:], hexDigest)
	return digest, err
}

type _JSONTrustedComment struct {
	At      int64  `json:"at"`
	By      string `json:"by"`
	Type    string `json:"typ"`
	Version string `json:"v"`
}

type TrustedComment struct {
	At int64
	By string
}

func (tc *TrustedComment) Marshal() (string, error) {
	jtc := _JSONTrustedComment{
		At:      tc.At,
		By:      tc.By,
		Type:    commentType,
		Version: commentVersion,
	}
	bytes, err := json.Marshal(jtc)
	if err != nil {
		return "", err
	}
	return string(bytes), err

}

func (tc *TrustedComment) Unmarshal(data string) error {
	var jtc _JSONTrustedComment
	if err := json.Unmarshal([]byte(data), &jtc); err != nil {
		return err
	}

	if jtc.Type != commentType {
		return fmt.Errorf("incorrect 'typ' field. Expected %s, got %s", commentType, jtc.Type)
	}

	if jtc.Version != commentVersion {
		return fmt.Errorf("incorrect 'v' field. Expected %s, got %s", commentVersion, jtc.Version)
	}

	tc.At = jtc.At
	tc.By = jtc.By

	return nil
}

type Signer struct {
	By            string
	PrivateKey    minisign.PrivateKey
	PublicKey     minisign.PublicKey
	PublicKeyText []byte
}

func (signer *Signer) comments(at time.Time) (string, string, error) {
	tc := TrustedComment{
		At: at.UTC().Unix(),
		By: signer.By,
	}
	trustedComment, err := tc.Marshal()
	if err != nil {
		return "", "", err
	}

	untrustedComment := fmt.Sprintf("Signed by %s at %s", signer.By, at.UTC().Format(time.RFC1123Z))

	return trustedComment, untrustedComment, nil
}

func (signer *Signer) SignDigest(digest digest, at time.Time) (signature []byte, err error) {

	trustedComment, untrustedComment, err := signer.comments(at)
	if err != nil {
		return
	}

	return signer.PrivateKey.SignWithComments(digest[:], trustedComment, untrustedComment, crypto.BLAKE2b_512)
}

func (signer *Signer) VerifyDigest(digest digest, signature []byte) (verified bool, err error) {
	return signer.PublicKey.Verify(digest[:], signature, crypto.BLAKE2b_512)
}
