package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"web_chat/biz/db/mysql"
	"web_chat/biz/db/redis"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/po"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	accountInfoKeyPrefix = "web_chat_account_info_"
)

func accountInfoKey(accountID string) string {
	return fmt.Sprintf("%s%s", accountInfoKeyPrefix, accountID)
}

func CreateAccount(ctx context.Context, account *domain.Account) (string, error) {
	accountID := uuid.NewString()
	gormDB := mysql.GetGormDB().
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(
			&po.Account{
				AccountID:      accountID,
				Username:       account.Username,
				Password:       account.Password,
				Salt:           account.Salt,
				Status:         account.Status,
				ExpirationTime: account.ExpirationDate,
			},
		)
	if err := gormDB.Error; err != nil {
		hlog.CtxErrorf(ctx, "create account err: %v", err)
		return "", err
	}

	if gormDB.RowsAffected <= 0 {
		return "", nil
	}

	return accountID, nil
}

func UpdateAccountPassword(
	ctx context.Context,
	oldSalt, oldPassword, accountID string,
	newSalt, newPassword string,
) error {
	return mysql.GetGormDB().
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			// update firstly can block the login request
			if err := tx.Model(&po.Account{}).
				Update("salt", newSalt).
				Update("password", newPassword).
				Where("account_id", accountID).
				Where("salt", oldSalt).
				Where("password", oldPassword).
				Error; err != nil {
				hlog.CtxErrorf(ctx, "create account err: %v", err)
				return err
			}

			return deleteAccountInCache(ctx, accountID)
		})
}

func GetAccountByUsername(ctx context.Context, username string) (*domain.Account, error) {
	var result po.Account
	if err := mysql.GetGormDB().
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("username", username).
		Take(&result).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		hlog.CtxErrorf(ctx, "get account by username err: %v", err)
		return nil, err
	}

	return &domain.Account{
		AccountID:      result.AccountID,
		Username:       result.Username,
		Salt:           result.Salt,
		Password:       result.Password,
		Status:         result.Status,
		ExpirationDate: result.ExpirationTime,
	}, nil
}

func GetAccountByAccountID(ctx context.Context, accountID string) (*domain.Account, error) {
	account := getAccountFromCache(ctx, accountID)
	if account != nil {
		return account, nil
	}

	var result po.Account
	if err := mysql.GetGormDB().
		WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
		}).
		Where("account_id", accountID).
		Take(&result).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		hlog.CtxErrorf(ctx, "get account by account_id err: %v", err)
		return nil, err
	}

	account = &domain.Account{
		AccountID:      result.AccountID,
		Username:       result.Username,
		Salt:           result.Salt,
		Password:       result.Password,
		Status:         result.Status,
		ExpirationDate: result.ExpirationTime,
	}

	_ = setAccountInCache(ctx, account)

	return account, nil
}

func deleteAccountInCache(ctx context.Context, accountID string) error {
	if _, err := redis.GetRedisClient().
		HDel(ctx, accountInfoKey(accountID)).
		Result(); err != nil {
		hlog.CtxErrorf(ctx, "del account cache err: %v", err)
		return err
	}

	return nil
}

func getAccountFromCache(ctx context.Context, accountID string) *domain.Account {
	mapper, err := redis.GetRedisClient().HGetAll(ctx, accountInfoKey(accountID)).Result()
	if err != nil {
		hlog.CtxErrorf(ctx, "hgetall err: %v", err)
		return nil
	}

	expirationDate, err := time.Parse("2006-01-02", mapper["expiration_date"])
	if err != nil {
		hlog.CtxWarnf(ctx, "parse expiration date err: %v, err")
		expirationDate = time.Date(2099, 12, 31, 0, 0, 0, 0, time.Local)
	}

	return &domain.Account{
		AccountID:      mapper["account_id"],
		Username:       mapper["username"],
		Salt:           mapper["salt"],
		Password:       mapper["password"],
		Status:         mapper["status"],
		ExpirationDate: expirationDate,
	}
}

func setAccountInCache(ctx context.Context, account *domain.Account) error {
	_, err := redis.GetRedisClient().
		HSet(
			ctx, accountInfoKey(account.AccountID),
			"account_id", account.AccountID,
			"username", account.Username,
			"salt", account.Salt,
			"password", account.Password,
			"status", account.Status,
			"expiration_date", account.ExpirationDate.Format("2006-01-02"),
		).
		Result()
	if err != nil {
		hlog.CtxErrorf(ctx, "hset account err: %v", err)
		return err
	}

	return nil
}

const fieldSessionList = "session_list"

func AppendSessionInAccount(ctx context.Context, accountID, sessionID string) error {
	sessionList, err := GetSessionList(ctx, accountID)
	if err != nil {
		return err
	}

	sessionList = append(sessionList, sessionID)

	if _, err := redis.GetRedisClient().
		HSet(
			ctx, accountInfoKey(accountID),
			"session_list", sessionList,
		).
		Result(); err != nil {
		hlog.CtxErrorf(ctx, "set session list err: %v", err)
		return err
	}

	return nil
}

func GetSessionList(ctx context.Context, accountID string) ([]string, error) {
	byteSessionList, err := redis.GetRedisClient().
		HGet(ctx, accountInfoKey(accountID), fieldSessionList).
		Bytes()
	if err != nil {
		hlog.CtxErrorf(ctx, "get session list err: %v", err)
		return nil, err
	}

	var sessionList []string
	if err := json.Unmarshal(byteSessionList, &sessionList); err != nil {
		hlog.CtxWarnf(ctx, "session list invalid: %v", err)
		return nil, nil
	}

	return sessionList, nil
}

func RemoveSession(ctx context.Context, accountID, sessionID string) error {
	sessionList, err := GetSessionList(ctx, accountID)
	if err != nil {
		return err
	}

	var newSessionList []string
	for _, id := range sessionList {
		if id != sessionID {
			newSessionList = append(newSessionList, id)
		}
	}

	if _, err := redis.GetRedisClient().
		HSet(
			ctx, accountInfoKey(accountID),
			"session_list", newSessionList,
		).
		Result(); err != nil {
		hlog.CtxErrorf(ctx, "set session list err: %v", err)
		return err
	}

	return nil
}
