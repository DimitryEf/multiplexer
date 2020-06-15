# HTTP-мультиплексор

## Описание:
- [приложение представляет собой http-сервер с одним хендлером](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/router/router.go#L15)
- [хедлер на вход получает POST-запрос со списком url в json-формате](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L14)
- [сервер запрашивает данные по всем этим url и возвращает результат клиенту в json-формате](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L137)
- [если в процессе обработки хотя бы одного из url получена ошибка, обработка всего списка прекращается и клиенту возвращается ошибка в текстовом формате](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L102)
## Ограничения:
- [сервер не принимает запрос если количество url  в нем больше 20](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L47)
- [сервер не обслуживает больше чем 100 одновременных входящих подключений](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/main.go#L118)
- [для каждого входящего запроса должно быть не больше 4 одновременных исходящих](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L140)
- [таймаут на запрос одного url - секунда](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L150)
- [обработка запроса может быть отменена клиентом в любой момент, это должно повлечь за собой остановку всех операций связанных с этим запросом](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/handler/multiplex.go#L57)
- [сервис должен поддерживать 'graceful shutdown': при получении сигнала от OS перестать принимать входящие  запросы, завершить текущие запросы и остановиться](https://github.com/DimitryEf/multiplexer/blob/cb76c47928ada4312943daf3cfb15938acacd188/main.go#L81)