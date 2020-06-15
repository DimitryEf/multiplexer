package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// Структура, соответствующая телу запроса от клиента. Например:
// ["http://example.com/1", "http://example.com/2", "http://example.com/3"]
type Urls []string

// Структура для формирования ответа клиенту. Например:
// [{"url":"http://example.com/1","body":"api 1"},{"url":"http://example.com/2","body":"api 2"},{"url":"http://example.com/3","body":"api 3"}
type Resp struct {
	Url  string `json:"url"`  // Url для запроса к стороннему серверу
	Body string `json:"body"` // Ответ от стороннего сервера
}

// Multiplex осуществляет запросы к сторонним серверам, согласно списку переданному от клиента
// и возвращает клиенту json с массивом ответов
func Multiplex(m *MultiplexerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем тело запроса
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			m.Log.Errorf("error in read body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Преобразуем json тела запроса в структуру Urls
		urls := Urls{}
		err = json.Unmarshal(b, &urls)
		if err != nil {
			m.Log.Errorf("error in unmarshal urls: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		urlsLen := len(urls) // Количество url

		// Проверяем количество url
		if urlsLen < 1 || urlsLen > m.MaxUrls {
			m.Log.Errorf("there are more than %d urls: %d", m.MaxUrls, urlsLen)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Создаем канал для отслеживания отмены запроса от клиента или из-за ошибки
		cancelChan := make(chan struct{}, urlsLen)

		// Запускаем в отдельной горутине отслеживание отмены
		go func() {
			<-r.Context().Done()
			// Всем запросам сигнализируем об отмене
			for i := 0; i < urlsLen; i++ {
				cancelChan <- struct{}{}
			}
		}()

		// Создаем канал для сбора ответов по url со сторонних серверов
		respChan := make(chan Resp, urlsLen)
		// Создаем канал для приема ошибки
		errorChan := make(chan error, 1)
		// Создаем канал для контроля количества исходящих подключений
		outputConnChan := make(chan struct{}, m.MaxOutputConnForOneInputConn)

		// Запускаем для каждого url отдельную горутину
		for _, url := range urls {
			go func(url string) {
				// Делаем запрос по url
				body, err := MakeRequest(m.UrlRequestTimeout, url, &cancelChan, &outputConnChan)
				if err != nil {
					// Ошибку пишем в канал
					errorChan <- err
					return
				}
				// Ответ записываем в канал
				respChan <- Resp{url, body}
			}(url)
		}

		var result []Resp // Переменная для записи результата с ответами со сторонних серверов

		// В бесконечном цикле используем select по каналам
	LOOP:
		for {
			select {
			case res := <-respChan:
				// Если получен ответ, то записываем его в слайс result
				result = append(result, res)
				// Если получены результаты по всем url, ты выходим из цикла
				if len(result) >= urlsLen {
					break LOOP
				}
			case err := <-errorChan:
				// Если получена ошибка, то пишем в cancelChan для завершения остальных исходящих запросов
				for i := 0; i < urlsLen; i++ {
					cancelChan <- struct{}{}
				}
				w.WriteHeader(http.StatusInternalServerError)
				// Ошибку возвращаем клиенту
				_, errW := w.Write([]byte(err.Error()))
				if errW != nil {
					m.Log.Errorf("error in write to response: %v", errW)
				}
				return
			}
		}

		// Конвертируем результат в json
		resultBody, err := json.Marshal(result)
		if err != nil {
			m.Log.Errorf("error in write to response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			// Отправляем ошибку клиенту
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				m.Log.Errorf("error in write to response: %v", err)
			}
			return
		}

		// Отправляем результат клиенту
		_, err = w.Write(resultBody)
		if err != nil {
			m.Log.Errorf("error in write to response: %v", err)
		}
	}
}

// MakeRequest делает запрос по переданному url
func MakeRequest(urlRequestTimeout time.Duration, url string, cancelChan, outputConnChan *chan struct{}) (string, error) {
	// Блокируем выполнение, если превышен лимит исходящих подключений.
	// Конфигурируется в поле MaxOutputConnForOneInputConn у структуры MultiplexerConfig.
	*outputConnChan <- struct{}{}
	defer func() { <-*outputConnChan }()

	// Формируем новый запрос по url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Создаем контекст с таймаутом и функцией отмены
	ctx, cancel := context.WithTimeout(req.Context(), urlRequestTimeout)
	req = req.WithContext(ctx)

	// Запускаем горутину, в которой будем отслеживать отмену
	go func() {
		<-*cancelChan
		cancel()
	}()
	defer cancel()

	// Инициализируем клиента для запроса
	client := &http.Client{}

	// Делаем запрос по url
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	// Получаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Возвращаем результат
	return string(body), nil
}
