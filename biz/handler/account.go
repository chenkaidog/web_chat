package handler

import (
	"context"
	"web_chat/biz/model/domain"
	"web_chat/biz/model/dto"
	"web_chat/biz/model/err"
	"web_chat/biz/repository"
	"web_chat/biz/util/origin"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/hertz-contrib/sessions"
)

const (
	sessionAccountID       = "account_id"
	sessionAccountUsername = "username"
	sessionAccountStatus   = "status"
	sessionSessID          = "session_id"
	maxFailTime            = 10
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

	account, bizErr := accountLoginVerify(ctx, req.Username, req.Password)
	if bizErr != nil {
		if !err.ErrorEqual(bizErr, err.InternalServerError) {
			recordFailOperation(ctx, c)
		}
		dto.FailResp(c, &resp, bizErr)
		return
	}

	resp.Username = req.Username
	resp.ExpirationDate = account.ExpirationDate.Format("2006-01-02")

	sess := sessions.Default(c)
	sess.Set(sessionAccountID, account.AccountID)
	sess.Set(sessionAccountUsername, account.Username)
	sess.Set(sessionAccountStatus, account.Status)
	if stdErr = sess.Save(); stdErr != nil {
		hlog.CtxErrorf(ctx, "save session err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	if stdErr := repository.AppendSessionInAccount(ctx, account.AccountID, sess.ID()); stdErr != nil {
		hlog.CtxErrorf(ctx, "append session list err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	dto.SuccessResp(c, &resp)
}

func accountLoginVerify(ctx context.Context, username, password string) (*domain.Account, err.Error) {
	account, bizErr := repository.GetAccountByUsername(ctx, username)
	if bizErr != nil {
		return nil, err.InternalServerError
	}
	if account == nil {
		hlog.CtxInfof(ctx, "username not exist: %s", username)
		return nil, err.AccountNotExistError
	}
	if account.IsInvalid() {
		hlog.CtxInfof(ctx, "account invalid: %s", account.Status)
		return nil, err.AccountStatusInvalidError
	}
	if !account.PasswordVerify(password) {
		hlog.CtxInfof(ctx, "password is incorrect")
		return nil, err.PasswordIncorrect
	}

	return account, nil
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
		if !err.ErrorEqual(bizErr, err.InternalServerError) {
			recordFailOperation(ctx, c)
		}
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

	sess := sessions.Default(c)
	sess.Clear()
	if stdErr := sess.Save(); stdErr != nil {
		hlog.CtxErrorf(ctx, "save session err: %v", stdErr)
		dto.FailResp(c, &resp, err.InternalServerError)
		return
	}

	dto.SuccessResp(c, &resp)
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

func recordFailOperation(ctx context.Context, c *app.RequestContext) {
	remoteAddr := origin.GetIp(c)

	failTime, _ := repository.RecordFailTime(ctx, remoteAddr)

	if failTime > maxFailTime {
		_ = repository.RecordRemoteAddrBlackList(ctx, remoteAddr)
	}
}
