package dto

type RegisterRequest struct {
	Account  string `json:"account" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Nickname string `json:"nickname" binding:"required,min=1,max=100"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdateNicknameRequest struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=100"`
}

type UserResponse struct {
	ID       uint64 `json:"id"`
	Account  string `json:"account"`
	Nickname string `json:"nickname"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
