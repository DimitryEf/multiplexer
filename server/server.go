package server

import (
	"context"
	"github.com/DimitryEf/multiplexer/config"
	"github.com/gorilla/mux"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type MultiplexerServer struct {
	HttpServer *http.Server
	Cfg        *config.MultiplexerConfig
}

func NewMultiplexerServer(cfg *config.MultiplexerConfig, router *mux.Router) *MultiplexerServer {
	return &MultiplexerServer{
		HttpServer: &http.Server{
			Addr:           net.JoinHostPort(cfg.Host, cfg.Port),
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderBytes,
			Handler:        router,
		},
		Cfg: cfg,
	}
}

// RunMultiplexer запускает сервер с мультиплексором указанной конфигурации
func (srv *MultiplexerServer) Run() {
	srv.Cfg.Log.Info("Server is running...")

	// Инициализируем слушателя для протокола tcp на указанном порту
	l, err := net.Listen("tcp", srv.Cfg.Port)
	if err != nil {
		srv.Cfg.Log.Fatalf("error in net.Listen: %v", err)
	}

	defer func() {
		err := l.Close()
		if err != nil {
			srv.Cfg.Log.Errorf("error in close net.Listen: %v", err)
		}
	}()

	// Используем LimitListener для ограничения количества входящих соединений
	l = netutil.LimitListener(l, srv.Cfg.MaxInputConn)

	// Запускаем сервер
	srv.Cfg.Log.Fatal(srv.HttpServer.Serve(l))
}

func (srv *MultiplexerServer) WaitSignalOS() {
	// Создаем канал для приема сигналов ОС
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM) // Отлавливаем в канал interrupt сигналы os.Interrupt и syscall.SIGTERM

	<-interrupt // Здесь исполнение кода блокируется, пока не не будет получен сигнал ОС

	srv.Cfg.Log.Info("Stopping srv...")

	//Устанавливаем контекст с таймаутом для принудительного завершения работы сервера
	timeout, cancelFunc := context.WithTimeout(context.Background(), srv.Cfg.ShutdownTimeout)
	defer cancelFunc()

	if err := srv.HttpServer.Shutdown(timeout); err != nil && err != http.ErrServerClosed { // Функция Shutdown стандартного пакета http обеспечивает graceful shutdown
		srv.Cfg.Log.Fatal(err)
	}

	srv.Cfg.Log.Info("The srv stopped.")
	os.Exit(0)
}
