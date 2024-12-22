# Калькулятор

**Calculation** — это сервис для выполнения математических вычислений на основе переданного выражения.

Проект написан на языке Go.

---

## Установка и запуск

Для запуска проекта выполните следующие шаги:

1. Склонируйте репозиторий:

```bash
git clone git@github.com:DoctoRgOlAnG/Calculation.git
cd Calculation
```

2. Убедитесь, что Go установлен и находится в `$PATH` (проверить версию можно командой `go version`).
3. Start testing
    ```
    go test ./internal
    Answer
    ok      doctor/internal 4.004s
    ------------------------------
    go test ./pkg
    ok      doctor/pkg      0.172s
    ```
4. Запустите API-сервер:

```bash
go run ./cmd/main.go
```

Сервер запустится на порту `8080` по умолчанию. 

## Usega

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

- **Successful**

   **Status code:** `200 OK`  
   **Example:**

   ```json
   {
     "result": "6"
   }
   ```

- **Ответ по ошибке при обработки выражения**

   **Статус-код:** `422 Unprocessable Entity`  
   **Example answer:**

   ```json
   {
     "error": "Unprocessable Entity"
   }
   ```

- **Неподдерживаемый метод**

   **Status code:** `405 Method Not Allowed`  
   **Example answer:**

   ```json
   {
     "error": "Method Not Allowed"
   }
   ```

- **Not correct body**

   **Statuc code:** `400 Bad Request`  
   **Example answer:**

   ```json
   {
     "error": "Bad Request"
   }
   ```

---

## Example

- **Successful**:

```bash
curl -H POST 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+10*10"
}'
```

Answer:

```json
{
  "result": "102"
}
```

- **Error: calculation**:

```bash
curl -H POST 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "17/(13+9{)"
}'
```

Answer:

```json
{
  "error": "Error calculate"
}
```

- **Error: not correct method**:

```bash
curl -H GET 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json'
```

Answer:

```json
{
  "error": "Wrong Method"
}
