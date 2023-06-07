package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	type want struct {
		code        int
		contentType string
	}

	testPostRequest := []struct {
		name     string
		url      string
		method   string
		wantPost want
		wantGet  want
	}{
		{
			name:   "Test post request first url",
			url:    "https://practicum.yandex.ru/",
			method: http.MethodPost,
			wantPost: want{
				code:        201,
				contentType: "text/plain",
			},
			wantGet: want{
				code:        307,
				contentType: "text/plain",
			},
		},

		{
			name:   "Test post request second url",
			url:    "https://www.youtube.com/",
			method: http.MethodPost,
			wantPost: want{
				code:        201,
				contentType: "text/plain",
			},
			wantGet: want{
				code:        307,
				contentType: "text/plain",
			},
		},
	}

	for _, tc := range testPostRequest {
		t.Run(tc.name, func(t *testing.T) {
			// Тест Post запроса
			rp := httptest.NewRequest(tc.method, "https://localhost:8080", strings.NewReader(tc.url))
			wp := httptest.NewRecorder()
			mainPost(wp, rp)
			postUtl := wp.Body.String()

			assert.Equal(t, tc.wantPost.code, wp.Code, "Код ответа Post не совпадает с ожидаемым")
			assert.Equal(t, tc.wantPost.contentType, wp.Header()["Content-Type"][0], "Заголовок Post ответа не совпадает с ожидаемым")

			// Тест Get запроса
			rg := httptest.NewRequest(http.MethodGet, postUtl, nil)
			wg := httptest.NewRecorder()
			mainPost(wg, rg)

			assert.Equal(t, tc.wantGet.code, wg.Code, "Код ответа Get запроса не совпадает с ожидаемым")
			assert.Equal(t, tc.url, wg.Header()["Location"][0], "Location в Get запросе не совпадает с ожидаемым ответом")

		})
	}
}
