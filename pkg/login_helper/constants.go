/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package login_helper
import "time"


const (
	DefaultTimeout = time.Minute
	DefaultPath = "/tmp/"
	// TokenFileName with the name of the file we use to store the token.
	TokenFileName = "token"
	// RefreshTokenFileName with the name of the file that contains the refresh token
	RefreshTokenFileName = "refresh_token"
	AuthHeader = "Authorization"
)