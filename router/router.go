package router

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router(m *MultiplexerConfig) *mux.Router {
	r := mux.NewRouter()

	// Устанавливаем единственную хэндлер-функцию
	r.HandleFunc("/", Multiplex(m)).Methods(http.MethodPost)

	// Используем middleware для логирования всех входящих запросов
	// и для добавления необходимого заголовка в ответ
	r.Use(LogMiddleware(m))
	r.Use(HeadersMiddleware())

	return r
}
