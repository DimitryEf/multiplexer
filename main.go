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
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/netutil"

	"github.com/DimitryEf/multiplexer/config"
	"github.com/DimitryEf/multiplexer/router"
)

//TODO Unit tests
//TODO integration tests

func main() {
	// В качестве логгера используется logrus
	log := logrus.New()
	log.SetOutput(os.Stdout) //Устанавливаем вывод логов в stdout

	m := &config.MultiplexerConfig{
		Log:                          log,
		Host:                         "",
		Port:                         ":8080",
		MaxUrls:                      20,
		MaxInputConn:                 100,
		MaxOutputConnForOneInputConn: 4,
		UrlRequestTimeout:            1 * time.Second,
		ShutdownTimeout:              10 * time.Second,
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		MaxHeaderBytes:               1 << 20,
	}

	m.Log.Info("Starting the app...")

	// Проверяем установку порта
	if len(m.Port) == 0 {
		m.Log.Fatal("Port is not set")
	}
	m.Log.Infof("Port is %v", m.Port)

	// Инициализируем сервер. Используется сервер из стандартной библиотеки
	server := &http.Server{
		Addr:           net.JoinHostPort(m.Host, m.Port),
		ReadTimeout:    m.ReadTimeout,
		WriteTimeout:   m.WriteTimeout,
		MaxHeaderBytes: m.MaxHeaderBytes,
		Handler:        router.Router(m),
	}

	// Запускаем сервер мультиплексора в горутине
	go RunMultiplexer(server, m)

	// Создаем канал для приема сигналов ОС
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // Отлавливаем в канал interrupt сигналы os.Interrupt и syscall.SIGTERM

	<-interrupt // Здесь исполнение кода блокируется, пока не не будет получен сигнал ОС

	m.Log.Info("Stopping server...")

	//Устанавливаем контекст с таймаутом для принудительного завершения работы сервера
	timeout, cancelFunc := context.WithTimeout(context.Background(), m.ShutdownTimeout)
	defer cancelFunc()
	err := server.Shutdown(timeout) // Функция Shutdown стандартного пакета http обеспечивает graceful shutdown
	if err != nil {
		log.Fatal(err)
	}

	log.Info("The server stopped.")
	os.Exit(0)
}

// RunMultiplexer запускает сервер с мультиплексором указанной конфигурации
func RunMultiplexer(server *http.Server, m *config.MultiplexerConfig) {
	m.Log.Info("Server is running...")

	// Инициализируем слушателя для протокола tcp на указанном порту
	l, err := net.Listen("tcp", m.Port)
	if err != nil {
		m.Log.Fatalf("error in net.Listen: %v", err)
	}

	defer func() {
		err := l.Close()
		if err != nil {
			m.Log.Errorf("error in close net.Listen: %v", err)
		}
	}()

	// Используем LimitListener для ограничения количества входящих соединений
	l = netutil.LimitListener(l, m.MaxInputConn)

	// Запускаем сервер
	m.Log.Fatal(server.Serve(l))
}
