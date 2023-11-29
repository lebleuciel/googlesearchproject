package gateway

import "github.com/lebleuciel/maani/models"

// swagger:route POST /api/auth/login Auth login
// Signs in user.
// responses:
//   200: Token

// swagger:response Token
type SignInUserResponse struct {
	// in:body
	Body models.UserTokenResponse
}

// swagger:parameters login
type SignInUserRequest struct {
	// in:body
	Body models.UserLoginCredentials
}

// swagger:route POST /api/auth/logout Auth logout
// Logs out the user.
// Security:
//    bearerAuth: []
// responses:
//   204: logout

// swagger:response logout
type LogOutUserResponse struct {
	Code int
}

// swagger:route POST /api/auth/refresh Auth getRefreshToken
// Generate JWT RefreshToken for current user.
// Security:
//    bearerAuth: []
// responses:
//   200: refreshToken

// swagger:response refreshToken
type RefreshTokenResponse struct {
	models.UserTokenResponse
}

// swagger:route POST /api/auth/register Auth register
// Register new customer user.
// Security:
//    bearerAuth: []
// responses:
//   200: registerUser

// swagger:response registerUser
type RegisterUserResponse struct {
	models.User
}

// swagger:parameters register
type RegisterUserRequest struct {
	// in:body
	Body models.UserRegisterParameters
}

// swagger:route GET /api/user/list File list
// Its only for admin user.
// Its only for admin user
// Security:
//    bearerAuth: []
// responses:
//   200:
