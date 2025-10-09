package register

import (
	"UrlShortener/internal/clients/sso/grpc"
	"UrlShortener/internal/http-server/handlers/auth/register/mocks"
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
	passDefaultLen = 10
	email          = gofakeit.Email()
	password       = randomFakePassword()
)

func TestRegisterHandler(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		password     string
		expectedCode int
		respError    string
		mockError    error
	}{
		{
			name:         "Success",
			email:        email,
			password:     password,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty email",
			email:        "",
			password:     password,
			expectedCode: http.StatusBadRequest,
			respError:    "field Email is a required field",
		},
		{
			name:         "Empty password",
			email:        email,
			password:     "",
			expectedCode: http.StatusBadRequest,
			respError:    "field Password is a required field",
		},
		{
			name:         "Not valid email",
			email:        "aboba",
			password:     password,
			expectedCode: http.StatusBadRequest,
			respError:    "field Email is not a valid email",
		},
		{
			name:         "Registering an existing user",
			email:        email,
			password:     password,
			expectedCode: http.StatusConflict,
			respError:    "user is already exists",
			mockError:    grpc.ErrUserExists,
		},
		{
			name:         "Register Error",
			email:        email,
			password:     password,
			expectedCode: http.StatusInternalServerError,
			respError:    "failed to registrate",
			mockError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			registerServiceMock := mocks.NewMockRegisterService(t)

			if tc.respError == "" || tc.mockError != nil {
				registerServiceMock.On("Register", context.Background(), tc.email, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := New(sl.NewDiscardLogger(), registerServiceMock)

			input := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tc.email, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(input)))
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
