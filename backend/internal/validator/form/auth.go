package form

type LoginAuth struct {
	UserName string `form:"username" json:"username" binding:"required,min=5"`
	PassWord string `form:"password" json:"password" binding:"required,min=6"`
}

func NewLoginForm() *LoginAuth {
	return &LoginAuth{}
}
