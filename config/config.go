package config

import (
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
