package request

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"code" validate:"required,len=6"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
}

type SendCodeReq struct {
	Email string `json:"email" validate:"required,email"`
}
