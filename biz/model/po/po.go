package po

import (
	"time"
)

type Account struct {
	ID             uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	AccountID      string    `gorm:"column:account_id;unique"`
	Username       string    `gorm:"column:username;unique"`
	Password       string    `gorm:"column:password"`
	Salt           string    `gorm:"column:salt"`
	Status         string    `gorm:"column:status"`
	ExpirationTime time.Time `gorm:"column:expiration_time"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (Account) TableName() string {
	return "account"
}
