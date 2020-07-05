package config

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// MultiplexerConfig - структура с настройками мультиплексора
type MultiplexerConfig struct {
	Log                          *logrus.Logger // Логгер
	Host                         string         // Хост
	Port                         string         // Порт, например ":8080"
	MaxUrls                      int            // Максимальное количество url в теле входящего запроса
	MaxInputConn                 int            // Максимальное количество одновременно входящих соединений
	MaxOutputConnForOneInputConn int            // Допустимое количество исходящих подключений на каждое входящее
	UrlRequestTimeout            time.Duration  // Таймаут на запрос одного url
	ShutdownTimeout              time.Duration  // Таймаут принудительной остановки сервера
	ReadTimeout                  time.Duration  // Таймаут для чтения запроса
	WriteTimeout                 time.Duration  // Таймаут для записи ответа
	MaxHeaderBytes               int            // Максимальный размер заголовков запроса
}

func NewMultiplexerConfig(logger *logrus.Logger) *MultiplexerConfig {
	logger.SetOutput(os.Stdout) //Устанавливаем вывод логов в stdout

	return &MultiplexerConfig{
		Log:                          logger,
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
}
