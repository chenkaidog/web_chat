package dto

type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required,max=64"`
	Password string `form:"password" json:"password" binding:"required,max=128"`
}

type LoginResp struct {
	CommonResp

	Username       string `json:"username,omitempty"`
	ExpirationDate string `json:"expiration_date,omitempty"`
}

type LogoutReq struct {

}

type LogoutResp struct {
	CommonResp
}

type PasswordUpdateReq struct {
	Password    string `form:"password" json:"password" binding:"required,max=128"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=8,max=128,alphanumunicode"`
}

type PasswordUpdateResp struct {
	CommonResp
}
