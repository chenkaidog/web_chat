package handler

import (
	"context"
	"time"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/dto"
	"web_chat/biz/model/err"
	"web_chat/biz/repository"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

const (
	sessionAccountID = "account_id"
	sessionSessID    = "session_id"
)

func Login(ctx context.Context, c *app.RequestContext) {
	var stdErr error
	var req dto.LoginReq
	var resp dto.LoginResp
	if stdErr = c.BindAndValidate(&req); stdErr != nil {
		hlog.CtxErrorf(ctx, "BindAndValidate fail, %v", stdErr)
		dto.FailResp(c, &resp, err.ParamError)
		return
	}

	accountID, expirationDate, bizErr := accountLoginVerify(ctx, req.Username, req.Password)
	if bizErr != nil {
		// todo: security strategy
		dto.FailResp(c, &resp, bizErr)
		return
	}

	resp.Username = req.Username
	resp.ExpirationDate = expirationDate.Format("2006-01-02")

	sess := sessions.Default(c)
	sess.Set(sessionAccountID, accountID)
	if stdErr = sess.Save(); stdErr != nil {
		hlog.CtxErrorf(ctx, "save session err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	if stdErr := repository.AppendSessionInAccount(ctx, accountID, sess.ID()); stdErr != nil {
		hlog.CtxErrorf(ctx, "append session list err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	dto.SuccessResp(c, &resp)
	return
}

func accountLoginVerify(ctx context.Context, username, password string) (string, time.Time, err.Error) {
	account, bizErr := repository.GetAccountByUsername(ctx, username)
	if bizErr != nil {
		return "", time.Time{}, err.InternalServerError
	}
	if account == nil {
		hlog.CtxInfof(ctx, "username not exist: %s", username)
		return "", time.Time{}, err.AccountNotExistError
	}
	if account.IsInvalid() {
		hlog.CtxInfof(ctx, "account invalid: %s", account.Status)
		return "", time.Time{}, err.AccountStatusInvalidError
	}
	if !account.PasswordVerify(password) {
		hlog.CtxInfof(ctx, "password is incorrect")
		return "", time.Time{}, err.PasswordIncorrect
	}

	return account.AccountID, account.ExpirationDate, nil
}

func Logout(ctx context.Context, c *app.RequestContext) {
	var stdErr error
	var req dto.LogoutReq
	var resp dto.LogoutResp
	if stdErr = c.BindAndValidate(&req); stdErr != nil {
		hlog.CtxErrorf(ctx, "BindAndValidate fail, %v", stdErr)
		dto.FailResp(c, &resp, err.ParamError)
		return
	}

	if stdErr := repository.RemoveSession(
		ctx, c.GetString(sessionAccountID), c.GetString(sessionSessID)); stdErr != nil {
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	sess := sessions.Default(c)
	sess.Clear()
	if stdErr := sess.Save(); stdErr != nil {
		hlog.CtxErrorf(ctx, "save session err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	dto.SuccessResp(c, &resp)
	return
}

func UpdatePassword(ctx context.Context, c *app.RequestContext) {
	var stdErr error
	var req dto.PasswordUpdateReq
	var resp dto.PasswordUpdateResp
	if stdErr = c.BindAndValidate(&req); stdErr != nil {
		hlog.CtxErrorf(ctx, "BindAndValidate fail, %v", stdErr)
		dto.FailResp(c, &resp, err.ParamError)
		return
	}

	accountID := c.GetString(sessionAccountID)
	account, bizErr := passwordUpdateVerify(ctx, accountID, req.Password)
	if bizErr != nil {
		// todo: security strategy
		dto.FailResp(c, &resp, bizErr)
		return
	}

	salt, password := domain.EncodePassword(req.NewPassword)

	if stdErr = repository.UpdateAccountPassword(
		ctx, account.Salt, account.Password, accountID, salt, password); stdErr != nil {
		hlog.CtxErrorf(ctx, "update password err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	dto.SuccessResp(c, &resp)
	return
}

func passwordUpdateVerify(ctx context.Context, accountID, password string) (*domain.Account, err.Error) {
	account, bizErr := repository.GetAccountByAccountID(ctx, accountID)
	if bizErr != nil || account == nil {
		hlog.CtxErrorf(ctx, "get account err: %v", bizErr)
		return nil, err.InternalServerError
	}
	if !account.PasswordVerify(password) {
		hlog.CtxInfof(ctx, "password is incorrect")
		return nil, err.PasswordIncorrect
	}

	return &domain.Account{
		AccountID:      account.AccountID,
		Username:       account.Username,
		Salt:           account.Salt,
		Password:       account.Password,
		Status:         account.Status,
		ExpirationDate: account.ExpirationDate,
	}, nil
}
