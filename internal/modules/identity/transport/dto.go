package transport

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email,max=255"`
	Username    string `json:"username" validate:"required,min=3,max=50"`
	DisplayName string `json:"displayName" validate:"required,min=2,max=100"`
	Password    string `json:"password" validate:"required,min=8,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type AuthResponse struct {
	User         AuthUser `json:"user"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	TokenType    string   `json:"token_type"`
	ExpiresIn    int64    `json:"expires_in"`
}

type MeResponse struct {
	User AuthUser `json:"user"`
}
