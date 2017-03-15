package been

type UserFrom struct {
	Uname  string `form:"uname"`
	Mobile string `form:"mobile"`
	Pwd    string `form:"pwd"`
}
