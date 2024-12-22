# Калькулятор

**Calc Go** — это сервис для выполнения математических вычислений на основе переданного выражения. Сервис предоставляет API для обработки запросов на вычисления, а также включает в себя набор тестов для проверки корректности работы.

Проект написан на языке Go, организован в модульной структуре и содержит примеры использования.

---

## Установка и запуск

Для запуска проекта выполните следующие шаги:

1. Склонируйте репозиторий:

```bash
git clone git@github.com:DoctoRgOlAnG/Calculation.git
cd Calculation
```

2. Убедитесь, что Go установлен и находится в `$PATH` (проверить версию можно командой `go version`).

3. Запустите API-сервер:

```bash
go run ./cmd/main.go
```

Сервер запустится на порту `8080` по умолчанию. 

## Использование API

### Endpoint

```
POST /api/v1/calculate
```

### Header

- `Content-Type: application/json`

### Body

Example:

```json
{
  "expression": "2+2*2"
}
```

### Answer

1. **Succese**

   **Статус-код:** `200 OK`  
   **Пример ответа:**

   ```json
   {
     "result": "6"
   }
   ```

2. **Ответ по ошибке при обработки выражения**

   **Статус-код:** `422 Unprocessable Entity`  
   **Пример ответа:**

   ```json
   {
     "error": "Unprocessable Entity"
   }
   ```

3. **Неподдерживаемый метод**

   **Статус-код:** `405 Method Not Allowed`  
   **Пример ответа:**

   ```json
   {
     "error": "Method Not Allowed"
   }
   ```

4. **Некорректное тело запроса**

   **Статус-код:** `400 Bad Request`  
   **Пример ответа:**

   ```json
   {
     "error": "Bad Request"
   }
   ```

---

## Примеры использования

1. **Успешный запрос**:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+10*10"
}'
```

Ответ:

```json
{
  "result": "102"
}
```

2. **Ошибка: некорректное выражение**:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "17/(13+9{)"
}'
```

Ответ:

```json
{
  "error": "Error calculate"
}
```

3. **Ошибка: неверный метод**:

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json'
```

Ответ:

```json
{
  "error": "Wrong Method"
}
