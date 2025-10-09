package login

import (
	"UrlShortener/internal/clients/sso/grpc"
	"UrlShortener/internal/http-server/handlers/auth/login/mocks"
	"UrlShortener/internal/logger/sl"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	email          = gofakeit.Email()
	password       = randomFakePassword()
	appID          = gofakeit.Int32()
	passDefaultLen = 10
)

func TestLoginHandler(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		appID        int32
		expectedCode int
		respError    string
		mockError    error
	}{
		{
			name:         "Success",
			email:        email,
			password:     password,
			appID:        appID,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty email",
			email:        "",
			password:     password,
			appID:        appID,
			expectedCode: http.StatusBadRequest,
			respError:    "field Email is a required field",
		},
		{
			name:         "Empty password",
			email:        email,
			password:     "",
			appID:        appID,
			expectedCode: http.StatusBadRequest,
			respError:    "field Password is a required field",
		},
		{
			name:         "Empty app_id",
			email:        email,
			password:     password,
			appID:        0,
			expectedCode: http.StatusBadRequest,
			respError:    "field AppID is a required field",
		},
		{
			name:         "Not valid email",
			email:        "aboba",
			password:     password,
			appID:        appID,
			expectedCode: http.StatusBadRequest,
			respError:    "field Email is not a valid email",
		},
		{
			name:         "Login of a non-existent user",
			email:        email,
			password:     password,
			appID:        appID,
			expectedCode: http.StatusUnauthorized,
			respError:    "the user does not exist",
			mockError:    grpc.ErrUserNotFound,
		},
		{
			name:         "Login Error",
			email:        email,
			password:     password,
			appID:        appID,
			expectedCode: http.StatusInternalServerError,
			respError:    "failed to login user",
			mockError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			loginServiceMock := mocks.NewMockLoginService(t)

			if tc.respError == "" || tc.mockError != nil {
				loginServiceMock.On("Login", context.Background(), tc.email, mock.AnythingOfType("string"), mock.AnythingOfType("int32")).
					Return("jwt", tc.mockError).
					Once()
			}

			handler := New(sl.NewDiscardLogger(), loginServiceMock)

			input := fmt.Sprintf(`{"email": "%s", "password": "%s", "app_id": %d}`, tc.email, tc.password, tc.appID)

			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, rr.Code, tc.expectedCode)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, resp.Error, tc.respError)
		})
	}
}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefaultLen)
}
