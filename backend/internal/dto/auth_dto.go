package dto

// RegisterRequest 注册请求。
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Name     string `json:"name" binding:"required,min=1,max=64"`
	OwnerName string `json:"owner_name" binding:"required,min=1,max=128"`
	ContactName string `json:"contact_name" binding:"required,min=1,max=64"`
	Phone    string `json:"phone" binding:"required,max=32"`
}

// LoginRequest 登录请求。
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=1,max=64"`
}

// LoginResponse 登录响应。
type LoginResponse struct {
	Token string     `json:"token"`
	User  *UserView  `json:"user"`
}

// UserView 用户视图。
type UserView struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	RoleText  string `json:"role_text"`
	OwnerID   *uint  `json:"owner_id"`
	Status    string `json:"status"`
}

// MeResponse 当前用户响应。
type MeResponse struct {
	User *UserView `json:"user"`
}
