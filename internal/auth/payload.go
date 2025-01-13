package auth

type AuthRequest struct {
	AuthorizationCode string `json:"auth_code" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
