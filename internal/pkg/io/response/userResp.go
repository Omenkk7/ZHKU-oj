package response

type LoginData struct {
	Token string `json:"token"`
	Uname string `json:"uname"`
}

// 登录专用函数
func NewLoginResponse(token, uname string) Response[LoginData] {
	return NewResponse(0, "ok", LoginData{
		Token: token,
		Uname: uname,
	})
}

// 添加用户返回函数
func NewAddUserResponse(uid, uname, school, classes, major, adept, vjid, email, codeForceUser, headURL string, submited int64, solved uint32, rating int) Response[UserResp] {
	return NewResponse(0, "ok", UserResp{
		UID:     uid,
		Uname:   uname,
		School:  school,
		Classes: classes,
		Major:   major,
	})
}

// 获取用户信息返回函数
func NewGetUserResponse(uid, uname, school, classes, major, adept, vjid, email, codeForceUser, headURL string, submited int64, solved uint32, rating int) Response[UserResp] {
	return NewResponse(0, "ok", UserResp{
		UID:     uid,
		Uname:   uname,
		School:  school,
		Classes: classes,
		Major:   major,
	})
}

type UserResp struct {
	Response
	UID           string `json:"UID"`
	Uname         string `json:"UserName"`
	School        string `json:"School"`
	Classes       string `json:"Classes"`
	Major         string `json:"Major"`
	Adept         string `json:"Adept"`
	Vjid          string `json:"Vjid"`
	Email         string `json:"Email"`
	CodeForceUser string `json:"CodeForceUser"`
	HeadURL       string `json:"HeadURL"`
	Submited      int64  `json:"Submited"`
	Solved        uint32 `json:"Solved"`
	Rating        int    `json:"Rating"`
}
