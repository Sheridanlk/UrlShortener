package save

import (
	"UrlShortener/internal/http-server/handlers/url/save/mocks"
	"UrlShortener/internal/logger/sl"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		url          string
		expectedCode int
		respError    string
		mockError    error
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			url:          "https://google.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty alias",
			alias:        "",
			url:          "https://google.com",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Empty URL",
			url:          "",
			alias:        "some_alias",
			expectedCode: http.StatusBadRequest,
			respError:    "field URL is a required field",
		},
		{
			name:         "Invalid URL",
			url:          "some invalid URL",
			alias:        "some_alias",
			expectedCode: http.StatusBadRequest,
			respError:    "field URL is not a valid URL",
		},
		{
			name:         "SaveURL Error",
			alias:        "test_alias",
			url:          "https://google.com",
			expectedCode: http.StatusInternalServerError,
			respError:    "failed to add url",
			mockError:    errors.New("unexpected error"),
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlSaverMock := mocks.NewMockURLSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveURL", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := New(sl.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, tc.expectedCode)

			body := rr.Body.String()

			var resp Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}
}
