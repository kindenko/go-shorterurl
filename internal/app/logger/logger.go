package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status int
		size   int
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
}

// WithLogging добавляет дополнительный код для регистрации сведений о запросе
// и возвращает новый http.Handler.

func WithLogging(h http.HandlerFunc) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {

		logger, err := zap.NewDevelopment()
		if err != nil {
			// вызываем панику, если ошибка
			panic(err)
		}

		defer func() {
			err = logger.Sync()
		}()

		// делаем регистратор SugaredLogger
		sugar := *logger.Sugar()

		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status, // получаем перехваченный код статуса ответа
			"duration", duration,
			"size", responseData.size, // получаем перехваченный размер ответа
		)
	}
	return logFn
}

//var sugar zap.SugaredLogger

// // WithLogging добавляет дополнительный код для регистрации сведений о запросе
// // и возвращает новый http.Handler.
// func WithLogging(h http.Handler) http.Handler {
// 	logFn := func(w http.ResponseWriter, r *http.Request) {
// 		// функция Now() возвращает текущее время
// 		start := time.Now()

// 		// эндпоинт /ping
// 		uri := r.RequestURI
// 		// метод запроса
// 		method := r.Method

// 		// точка, где выполняется хендлер pingHandler
// 		h.ServeHTTP(w, r) // обслуживание оригинального запроса

// 		// Since возвращает разницу во времени между start
// 		// и моментом вызова Since. Таким образом можно посчитать
// 		// время выполнения запроса.
// 		duration := time.Since(start)

// 		// отправляем сведения о запросе в zap
// 		sugar.Infoln(
// 			"uri", uri,
// 			"method", method,
// 			"duration", duration,
// 		)

// 	}
// 	// возвращаем функционально расширенный хендлер
// 	return http.HandlerFunc(logFn)
// }
