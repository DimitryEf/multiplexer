package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/DimitryEf/multiplexer/config"
	"github.com/DimitryEf/multiplexer/handler"
	"github.com/DimitryEf/multiplexer/middleware"
)

func NewRouter(m *config.MultiplexerConfig) *mux.Router {
	r := mux.NewRouter()

	// Устанавливаем единственную хэндлер-функцию
	r.HandleFunc("/", handler.Multiplex(m)).Methods(http.MethodPost)

	// Используем middleware для логирования всех входящих запросов
	// и для добавления необходимого заголовка в ответ
	r.Use(middleware.LogMiddleware(m))
	r.Use(middleware.HeadersMiddleware())

	return r
}
