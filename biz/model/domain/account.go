package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"web_chat/biz/util/random"
)

type Account struct {
	AccountID      string
	Username       string
	Salt           string
	Password       string
	Status         string
	ExpirationDate time.Time
}

func (account *Account) IsInvalid() bool {
	return account.Status != AccountStatusValid
}

func (account *Account) PasswordVerify(password string) bool {
	h := sha256.New()
	h.Write([]byte(account.Salt))
	h.Write([]byte(password))

	return hex.EncodeToString(h.Sum(nil)) == account.Password
}

func EncodePassword(password string) (string, string) {
	salt := random.RandStr(64)
	h := sha256.New()

	h.Write([]byte(salt))
	h.Write([]byte(password))

	return salt, hex.EncodeToString(h.Sum(nil))
}

const (
	AccountStatusValid   = "valid"
	AccountStatusInvalid = "invalid"
)
