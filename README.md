# URL Shortener REST API

REST-сервис для сокращения ссылок.  
Регистрация и вход пользователей выполняются через отдельный **[Auth-сервис (gRPC)](https://github.com/Sheridanlk/gRPCAuth)**,  
который отвечает за проверку логина и выдачу JWT-токенов.  
REST-сервис использует эти токены для авторизации запросов.

---

## 1. Технологический стек
| Категория | Используемые технологии |
|------------|-------------------------|
| HTTP | net/http + chi |
| Конфигурация | cleanenv + YAML |
| БД | PostgreSQL |
| Миграции | golang-migrate/migrate |
| Авторизация | JWT, внешнее gRPC |
| gRPC-клиент | [google.golang.org/grpc](https://pkg.go.dev/google.golang.org/grpc) |
| Логирование | slog |
| Тестирование | testify, mockery |
| Контейнеризация | Docker + Docker Compose |

---

## 2. Основной функционал
- Регистрация и вход пользователя (`POST /auth/register`, `POST /auth/login`)
- Создание коротких ссылок (`POST /url/create`)
- Переход по короткой ссылке (`GET /{alias}`)
- Аутентификация и авторизация через gRPC-SSO
- JWT-middleware для защиты эндпоинтов
- Ведение логов и базовая аналитика обращений

---

## 3. Структура

```
UrlShortener/
├── cmd/
│ └── url-shortener/     # Точка входа (main.go)
├── internal/
│ ├── config/            # Загрузка/валидация конфигурации
│ ├── logger/            # Обёртка и настройка slog
│ ├── clients/
│ │ └── ssogrpc/         # gRPC-клиент к SSO: Register, Login
│ ├── storage/
│ │ └── postgresql/      # Работа с PostgreSQL
│ ├── http-server/
│ │ ├── router.go        # Инициализация chi, регистрация маршрутов и middleware
│ │ ├── middleware/
│ │ │ ├── jwt/           # JWTAuth: парсинг Bearer, проверки, user_id в context
│ │ │ └── logger/        # HTTP access-логгер, requestID, recover и пр.
│ │ └── handlers/        # HTTP-обработчики (auth, url)
│ └── lib/               # Вспомогательные утилиты: генератор alias, валидация, ошибки...
├── migrations/          # SQL-миграции для PostgreSQL
├── config/
│ └── config.yaml        # Файл конфигурации сервиса
├── Dockerfile           # Мультистейдж сборка Go-бинарника
├── docker-compose.yml   # REST + Postgres + migrate
├── Makefile
├── go.mod / go.sum
└── README.md
```

---

## 4. Запуск
### 4.1. Конфигурация 
Создать файл конфигурации `config/config.yaml`:
``` yaml
env: local #dev, prod
http_server:
  address: "0.0.0.0:8080"
  timeout: 4s
  idle_timeout: 30s
postgres:
  host: shortener-db # при Docker: имя сервиса БД в compose
  port: 5432
  user: postgres
  password: postgres
  name: shortener
clients:
  sso:
    address: "auth:34443" # при Docker: имя gRPC-сервиса в сети
    timeout: 10s
    retries_count: 5
```
(для локального запуска без Docker: поставь host: localhost и реальный порт/БД)

### 4.2. Старт
A) Локально(Go):
``` 
make run-local
```
B) В Docker контейнере(рекомендуется)
```
make run-docker
```

### 4.3. Миграции(для новых)
#### Создать новую:
```
make create-migration NAME=<имя_миграции>
```
#### Поднять через конейнер 
```
migrate-up-docker    # применить все
migrate-down-docker  # откатить одну
```
#### Локально(поменять путь для подключения к бд на свой)
```
migrate-up-local    # применить все
migrate-down-local  # откатить одну
```







