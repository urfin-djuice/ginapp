package account

type SignInForm struct {
	Email    string `json:"email" form:"Email" binding:"required,email"`
	Password string `json:"password" form:"Password" binding:"required"`
}

type SignUpForm struct {
	Name            string `json:"name" form:"name"`
	Email           string `json:"email" form:"email" binding:"required,email,UniqueEmail"`
	Password        string `json:"password" form:"password" binding:"required,StrongPass"`
	PasswordConfirm string `json:"password_confirm" form:"password_confirm" binding:"required,eqfield=Password" `
}
