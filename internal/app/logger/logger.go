package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {

	logger := func(w http.ResponseWriter, r *http.Request) {

		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}

		defer func() {
			err = logger.Sync()
		}()

		sugar := *logger.Sugar()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logger)
}

// func WithLogging(h http.HandlerFunc) http.HandlerFunc {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {

// 		logger, err := zap.NewDevelopment()
// 		if err != nil {
// 			// вызываем панику, если ошибка
// 			panic(err)
// 		}

// 		defer func() {
// 			err = logger.Sync()
// 		}()

// 		// делаем регистратор SugaredLogger
// 		sugar := *logger.Sugar()

// 		start := time.Now()

// 		responseData := &responseData{
// 			status: 0,
// 			size:   0,
// 		}
// 		lw := loggingResponseWriter{
// 			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
// 			responseData:   responseData,
// 		}
// 		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

// 		duration := time.Since(start)

// 		sugar.Infoln(
// 			"uri", r.RequestURI,
// 			"method", r.Method,
// 			"status", responseData.status, // получаем перехваченный код статуса ответа
// 			"duration", duration,
// 			"size", responseData.size, // получаем перехваченный размер ответа
// 		)
// 	}
// 	return logFn
// }
