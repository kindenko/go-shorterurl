package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"

	"strings"
	"testing"

	"github.com/kindenko/go-shorterurl/config"
	"github.com/stretchr/testify/assert"
)

var TestUrls = make(map[string]string)

func TestPostHandler(t *testing.T) {

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "First Post test",
			url:  "kfklr.com",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},

		{
			name: "Second Post test",
			url:  "",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	conf := &config.AppConfig{
		Host:      "localhost:8080",
		ResultURL: "http://localhost:8080",
		FilePATH:  "/tmp/short-url-db.json",
	}

	app := NewHandlers(conf)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "https://localhost:8080", strings.NewReader(tc.url))
			w := httptest.NewRecorder()

			app.PostHandler(w, r)
			res := w.Result()

			defer res.Body.Close()
			response, _ := io.ReadAll(res.Body)

			TestUrls[tc.url] = string(response)

			assert.Equal(t, tc.want.code, w.Code, "Код ответа Post не совпадает с ожидаемым")
			assert.Equal(t, tc.want.contentType, w.Header()["Content-Type"][0], "Заголовок Post ответа не совпадает с ожидаемым")
		})
	}
}

func TestPostJsonHandler(t *testing.T) {

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "First Post test",
			url: `{
				    "url": "https://practicum.yandex.ru"
			        } `,
			want: want{
				code:        http.StatusCreated,
				contentType: "application/json",
			},
		},

		{
			name: "Second Post test",
			url:  "",
			want: want{
				code:        http.StatusBadRequest,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	conf := &config.AppConfig{
		Host:      "localhost:8080",
		ResultURL: "http://localhost:8080",
		FilePATH:  "/tmp/short-url-db.json",
	}

	app := NewHandlers(conf)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "https://localhost:8080", strings.NewReader(tc.url))
			w := httptest.NewRecorder()

			app.PostJSONHandler(w, r)
			res := w.Result()

			defer res.Body.Close()
			response, _ := io.ReadAll(res.Body)

			TestUrls[tc.url] = string(response)

			assert.Equal(t, tc.want.code, w.Code, "Код ответа Post не совпадает с ожидаемым")
			assert.Equal(t, tc.want.contentType, w.Header()["Content-Type"][0], "Заголовок Post ответа не совпадает с ожидаемым")
		})
	}
}

// Разобраться че тут не так
// func TestGetHandler(t *testing.T) {

// 	type want struct {
// 		code        int
// 		contentType string
// 		location    string
// 		body        string
// 	}

// 	tests := []struct {
// 		name string
// 		url  string
// 		want want
// 	}{
// 		{
// 			name: "First Get Test",
// 			url:  "http://localhost:8080/XVlBzgba",
// 			want: want{
// 				code:        http.StatusTemporaryRedirect,
// 				contentType: "text/plain",
// 				location:    "kfklr.com",
// 				body:        "",
// 			},
// 		},

// 		{
// 			name: "Second Get test",
// 			url:  "http://localhost:8080/",
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "text/plain; charset=utf-8",
// 				location:    "",
// 				body:        "Bad URL\n",
// 			},
// 		},

// 		{
// 			name: "Third Get test",
// 			url:  "http://localhost:8080/utuyutusd",
// 			want: want{
// 				code:        http.StatusBadRequest,
// 				contentType: "text/plain; charset=utf-8",
// 				location:    "",
// 				body:        "Bad URL\n",
// 			},
// 		},
// 	}

// 	conf := &config.AppConfig{
// 		Host:      "localhost:8080",
// 		ResultURL: "http://localhost:8080",
// 		FilePATH:  "/tmp/short-url-db.json",
// 	}

// 	app := NewHandlers(conf)

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			r := httptest.NewRequest(http.MethodGet, tc.url, nil)
// 			w := httptest.NewRecorder()

// 			app.GetHandler(w, r)
// 			res := w.Result()

// 			defer res.Body.Close()
// 			resBody, _ := io.ReadAll(res.Body)

// 			assert.Equal(t, tc.want.code, w.Code, "Код ответа Get запроса не совпадает с ожидаемым")
// 			assert.Equal(t, tc.want.body, string(resBody), "Тело в Get запросе не совпадает с ожидаемым ответом")
// 			assert.Equal(t, tc.want.location, w.Header().Get("Location"), "Тело в Get запросе не совпадает с ожидаемым ответом")

// 			fmt.Println(TestUrls)
// 		})
// 	}
// }
