/*
Тестовое задание HTTP-мультиплексор:

- приложение представляет собой http-сервер с одним хендлером
- хедлер на вход получает POST-запрос со списком url в json-формате
- сервер запрашивает данные по всем этим url и возвращает результат клиенту в json-формате
- если в процессе обработки хотя бы одного из url получена ошибка, обработка всего списка прекращается и клиенту возвращается ошибка в текстовом формате
Ограничения:
- сервер не принимает запрос если количество url  в нем больше 20
- сервер не обслуживает больше чем 100 одновременных входящих подключений
- для каждого входящего запроса должно быть не больше 4 одновременных исходящих
- таймаут на запрос одного url - секунда
- обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом
- сервис должен поддерживать 'graceful shutdown': при получении сигнала от OS перестать принимать входящие  запросы, завершить текущие запросы и остановиться
- для реализации задачи следует использовать Go 1.13 или выше
- результат должен быть выложен на github
*/
package main

import (
	"github.com/DimitryEf/multiplexer/config"
	"github.com/DimitryEf/multiplexer/router"
	"github.com/DimitryEf/multiplexer/server"
	"github.com/sirupsen/logrus"
)

//TODO Unit tests
//TODO integration tests

func main() {
	// В качестве логгера используется logrus
	logger := logrus.New()

	// Структура с настройками мультиплексора
	cfg := config.NewMultiplexerConfig(logger)

	// Роутер с handleFunc и Middleware
	r := router.NewRouter(cfg)

	// Инициализируем сервер. Используется сервер из стандартной библиотеки
	srv := server.NewMultiplexerServer(cfg, r)

	// Запускаем сервер мультиплексора в горутине
	go srv.Run()

	// Блокируемся на ожидании сигнала от ОС
	srv.WaitSignalOS()
}
