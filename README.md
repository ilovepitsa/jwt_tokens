# Test task Junior golang developer 

Тестовый микросервис для аутентификации. Сервис выдает пользователю access и refresh токены при аутентификации. Используемые технологии:
    - PostgreSQL (для хранения данных)
    - Docker (для запуска сервиса)
    - pgx (для работы с PostgreSQL)
    - jwt (для генерации токенов)

# Важная информация
    В условии было сказано следующее: 'Refresh токен хранится в базе исключительно в виде bcrypt хеша'. Т.к. bcrypt хэш не идемпотентный, то нельзя было просто хранить его в базе данных, а при надобности доставать его. Мой вариант с bcrypt хэшем представлен в ветке 'bcrypt_hash'. Прошу обратить на это внимание.


# Getting Started
    Для запуска сервиса нужно заполнить конфигурацию 'config/config.yaml'. Так же можно указать другой путь, использую флаг '-c'

# Usage
    Запустить сервис можно с помощью команды 'make up'

    Для запусков тестов необходимо выполнить команду 'make test', для запусков тестов с покрытием 'make cover-html'

## Example
    Некоторые примеры запросов
    - [Создание пользователя](#sing-up)
    - [Аутентификация](#sing-in)
    - [Обновление пары токенов](#refresh)



### Создание пользователя <a name="sing-up"></a>
    Создание пользователя:
    ```curl 
    curl --location --request POST htttp://localhost:8081/create
    ```
    Пример ответа:
    ```json
    {
        "id": 1
    }
    ```

### Аутентификация <a name="sing-in"></a>

    Аутентификация сервиса для получения токенов:
    ```curl
    curl --location --request GET http://localhost:8081/auth/sign-in?user_id=1
    ```
    Пример ответа:
    ```json
    {
        "accessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ3NDk0NDYsInN1YiI6IjEiLCJJcCI6IjE3Mi4xOS4wLjEifQ.nluHWn_mDy1dLOfWfYEpzYbajfNOqBLCMI4Ptt5PP-Sw1V_1R1vNuV6HSCSiyd-fvWIdGufRe_qh7LyjRQlI3A",
        
        "refreshToken":"PFCx6Hs5AoOr21ZnjI3zh72ytoy_ZMlBwSo_ZCpHsjE="
        }
    ```

### Обновление токенов <a name="refresh"></a>
     
     Refresh операция пары Access и Refresh токенов:
     ```curl
     curl --location --request POST http://localhost:8081/auth/refresh?token=PFCx6Hs5AoOr21ZnjI3zh72ytoy_ZMlBwSo_ZCpHsjE=

     ```

     Пример ответа:
     ```json
     {
        "accessToken":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ3NTA1NjEsInN1YiI6IjEiLCJJcCI6IjE3Mi4xOS4wLjEifQ.VLoGbcjGYlwgBSLpyahD5Dmf6ZaBR2Qxp2Y_m7mizZPB2rBGSV5V3hNBC1BiqeRciwFu-O0e8tWGmJcAR_dHCg",
        
        "refreshToken":"4Cn-AU5mb7bPMukGUE-HjJ9qnqpebCjhpzxa2neROYk="}
     ```