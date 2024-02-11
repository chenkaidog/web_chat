package dto

type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required,min=8,max=64alphanumunicode"`
	Password string `form:"password" json:"password" binding:"required,min=8,max=128,alphanumunicode"`
}

type LoginResp struct {
	CommonResp

	Username       string `json:"username"`
	ExpirationDate string `json:"expiration_date"`
}

type LogoutReq struct {

}

type LogoutResp struct {
	CommonResp
}


type PasswordUpdateReq struct {
	Password    string `form:"password" json:"password" binding:"required,min=8,max=128,alphanumunicode"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required,min=8,max=128,alphanumunicode"`
}

type PasswordUpdateResp struct {
	CommonResp
}
