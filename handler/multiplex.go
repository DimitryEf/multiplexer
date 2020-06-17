package handler

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/DimitryEf/multiplexer/config"
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
func Multiplex(m *config.MultiplexerConfig) http.HandlerFunc {
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

		// проверяем валидность url
		for i := range urls {
			parsed, err := url.Parse(urls[i])
			if err != nil {
				m.Log.Errorf("invalid url: %q", urls[i])
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if parsed.Scheme == "" && parsed.Host == "" {
				m.Log.Errorf("invalid url: %q", urls[i])
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		mainCtx := r.Context()

		// Создаем канал для сбора ответов по url со сторонних серверов
		respChan := make(chan Resp, urlsLen)
		// Создаем канал для приема ошибки
		errorChan := make(chan error, 1)
		// Создаем канал для контроля количества исходящих подключений
		outputRequestLimit := make(chan struct{}, m.MaxOutputConnForOneInputConn)

		// Запускаем для каждого url отдельную горутину
		cancelableCtx, cancelRequests := context.WithCancel(mainCtx)
		defer cancelRequests()

		gotError := false
		for _, u := range urls {
			if gotError {
				break
			}

			go func(url string) {
				// Делаем запрос по url
				outputRequestLimit <- struct{}{}
				ctxTimeout, cancelTimeout := context.WithTimeout(cancelableCtx, m.UrlRequestTimeout)
				defer func() {
					cancelTimeout()
					<-outputRequestLimit
				}()

				body, err := DoRequest(ctxTimeout, url)
				if err != nil {
					gotError = true
					// Ошибку пишем в канал
					errorChan <- err
					return
				}
				// Ответ записываем в канал
				respChan <- Resp{url, body}
			}(u)
		}

		var result []Resp // Переменная для записи результата с ответами со сторонних серверов

		// В бесконечном цикле используем select по каналам
		run := true
		for run {
			select {
			case <-mainCtx.Done():
				// получен сигнал от родительского контекста
				cancelRequests()
				return
			case res := <-respChan:
				// Если получен ответ, то записываем его в слайс result
				result = append(result, res)
				// Если получены результаты по всем url, ты выходим из цикла
				if len(result) >= urlsLen {
					run = false
					break
				}
			case err := <-errorChan:
				cancelRequests()
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

// DoRequest делает запрос по переданному url
func DoRequest(ctx context.Context, url string) (string, error) {
	// Формируем новый запрос по url
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	// Делаем запрос по url
	resp, err := http.DefaultClient.Do(req)
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
