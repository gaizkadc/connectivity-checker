/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package login_helper

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-login-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpc_status "google.golang.org/grpc/status"
	"sync"
)

const (
	// Maximum number of retries for authentication
	MaxAuthRetries = 10
)

type LoginHelper struct {
	Connection
	useTLS      bool
	email       string
	password    string
	Credentials *Credentials
	mu          sync.RWMutex
}

// NewLogin creates a new LoginHelper structure.
func NewLogin(hostname string, port int, useTLS bool, email string, password string, caCertPath string, clientCertPath string, skipCAValidation bool) *LoginHelper {
	return &LoginHelper{
		Connection: *NewConnection(hostname, port, useTLS, "", "", true),
		email:      email,
		password:   password,
	}
}

func (l *LoginHelper) Login() derrors.Error {
	// Lock incoming
	l.mu.Lock()
	defer l.mu.Unlock()
	c, err := l.GetConnection()
	if err != nil {
		return err
	}
	defer c.Close()
	loginClient := grpc_login_api_go.NewLoginClient(c)
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	loginRequest := &grpc_authx_go.LoginWithBasicCredentialsRequest{
		Username: l.email,
		Password: l.password,
	}
	response, lErr := loginClient.LoginWithBasicCredentials(ctx, loginRequest)
	if lErr != nil {
		return conversions.ToDerror(lErr)
	}
	l.Credentials = NewCredentials(DefaultPath, response.Token, response.RefreshToken)
	sErr := l.Credentials.Store()
	if sErr != nil {
		return sErr
	}

	return nil
}

func (l *LoginHelper) GetContext() (context.Context, context.CancelFunc) {
	return l.Credentials.GetContext()
}

type GenericGRPCCall func(context.Context, interface{}, ...grpc.CallOption) (interface{}, error)

// Generic function to wrap GRPC calls inside a logged-in context.
//  params:
//   request	The request to be sent
//   call 		The GRPC function to be called
//  return:
//   interface of the object to be returned
func (l *LoginHelper) AuthenticatedGrpcCall(
	request interface{},
	call GenericGRPCCall,
) (interface{}, error) {
	// Get the logged-in context
	ctx, cancel := l.GetContext()
	defer cancel()
	// execute the GRPC call
	answer, err := call(ctx, &request)
	if err != nil {
		st := grpc_status.Convert(err).Code()
		if codes.Unauthenticated == st {
			authError := l.RerunAuthentication()
			if authError != nil {
				log.Error().Interface("call", call).Msg("impossible to run reauthentication... skip call")
				return nil, authError
			}
			ctx2, cancel2 := l.GetContext()
			defer cancel2()
			answer, err = call(ctx2, &request)
		}
	}
	return answer, err
}

// Internal function that runs the authentication process if and only if
func (l *LoginHelper) RerunAuthentication() derrors.Error {
	log.Info().Msg("reauthentication launched...")
	authenticated := false
	retries := 0
	for !authenticated && retries < MaxAuthRetries {
		loginError := l.Login()
		if loginError != nil {
			if grpc_status.Convert(loginError).Code() == codes.Unauthenticated {
				log.Error().Err(loginError).Int("retries", retries).Msg("unanthenticated when retrying login")
			}
			log.Error().Err(loginError).Int("retries", retries).Msg("retrying login...")
		} else {
			log.Info().Msg("login renegotiation successful")
			authenticated = true
		}
		retries = retries + 1
	}
	if authenticated {
		return nil
	}

	return derrors.NewUnauthenticatedError("authenticated, failed after reaching max retries")
}
