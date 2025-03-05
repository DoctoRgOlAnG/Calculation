Оркестратор
Сервер, который имеет следующие endpoint-ы:

Добавление вычисления арифметического выражения
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": <строка с выражение>
}'
Коды ответа: 201 - выражение принято для вычисления, 422 - невалидные данные, 500 - что-то пошло не так

Тело ответа

{
    "id": <уникальный идентификатор выражения>
}
Получение списка выражений
curl --location 'localhost/api/v1/expressions'
Тело ответа

{
    "expressions": [
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        },
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        }
    ]
}
Коды ответа:

200 - успешно получен список выражений
500 - что-то пошло не так

Получение выражения по его идентификатору

curl --location 'localhost/api/v1/expressions/:id'
Коды ответа:

200 - успешно получено выражение
404 - нет такого выражения
500 - что-то пошло не так
Тело ответа

{
    "expression":
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        }
}
Получение задачи для выполнения
curl --location 'localhost/internal/task'
Коды ответа:

200 - успешно получена задача
404 - нет задачи
500 - что-то пошло не так
Тело ответа

{
    "task":
        {
            "id": <идентификатор задачи>,
            "arg1": <имя первого аргумента>,
            "arg2": <имя второго аргумента>,
            "operation": <операция>,
            "operation_time": <время выполнения операции>
        }
}
Прием результата обработки данных.
curl --location 'localhost/internal/task' \
--header 'Content-Type: application/json' \
--data '{
  "id": 1,
  "result": 2.5
}'
Коды ответа:

200 - успешно записан результат
404 - нет такой задачи
422 - невалидные данные
500 - что-то пошло не так
Время выполнения операций задается переменными среды в миллисекундах
TIME_ADDITION_MS - время выполнения операции сложения в миллисекундах
TIME_SUBTRACTION_MS - время выполнения операции вычитания в миллисекундах
TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в миллисекундах
TIME_DIVISIONS_MS - время выполнения операции деления в миллисекундах

Агент
Демон, который получает выражение для вычисления с сервера, вычисляет его и отправляет на сервер результат выражения.

При старте демон запускает несколько горутин, каждая из которых выступает в роли независимого вычислителя. Количество горутин регулируется переменной среды COMPUTING_POWER