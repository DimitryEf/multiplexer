package middleware

import "net/http"

// LogMiddleware для логирования всех входящих запросов
func LogMiddleware(m *MultiplexerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			m.Log.Infof("ip=%s, method=%s, route=%s", r.RemoteAddr, r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

// HeadersMiddleware для добавления необходимых заголовков
func HeadersMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}

}
