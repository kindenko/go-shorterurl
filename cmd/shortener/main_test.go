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
		method      string
	}

	testPostRequest := []struct {
		name string
		url  string
		Post want
		Get  want
	}{
		{
			name: "Test request first url",
			url:  "https://practicum.yandex.ru/",
			Post: want{
				method:      http.MethodPost,
				code:        201,
				contentType: "text/plain",
			},
			Get: want{
				method:      http.MethodGet,
				code:        307,
				contentType: "text/plain",
			},
		},

		{
			name: "Test request second url",
			url:  "https://www.youtube.com/",
			Post: want{
				method:      http.MethodPost,
				code:        201,
				contentType: "text/plain",
			},
			Get: want{
				method:      http.MethodGet,
				code:        307,
				contentType: "text/plain",
			},
		},
	}

	for _, tc := range testPostRequest {
		t.Run(tc.name, func(t *testing.T) {
			// Тест Post запроса
			rp := httptest.NewRequest(tc.Post.method, "https://localhost:8080", strings.NewReader(tc.url))
			wp := httptest.NewRecorder()
			newPost(wp, rp)
			postURL := wp.Body.String()

			assert.Equal(t, tc.Post.code, wp.Code, "Код ответа Post не совпадает с ожидаемым")
			assert.Equal(t, tc.Post.contentType, wp.Header()["Content-Type"][0], "Заголовок Post ответа не совпадает с ожидаемым")

			// Тест Get запроса
			rg := httptest.NewRequest(tc.Get.method, postURL, nil)
			wg := httptest.NewRecorder()

			newGet(wg, rg)

			assert.Equal(t, tc.Get.code, wg.Code, "Код ответа Get запроса не совпадает с ожидаемым")
			assert.Equal(t, tc.url, wg.Header()["Location"][0], "Location в Get запросе не совпадает с ожидаемым ответом")

		})
	}
}
