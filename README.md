# demo-project

Учебный API на Go: регистрация, вход, обновление токенов, выход.

## Запуск

```bash
cp .env.example .env   # настройки для локального запуска
make infra-up          # Postgres (база demo-db) и Redis в Docker, миграции применятся сами
make run               # API на http://localhost:8080
```

Все команды для Docker, Postgres, Redis и проверки API — в [docs/commands.txt](docs/commands.txt).
Как добавить новую фичу или сущность — в [docs/adding-feature.md](docs/adding-feature.md).
Описание API в формате OpenAPI — [docs/openapi.yaml](docs/openapi.yaml) (открыть в GoLand, VS Code или на https://editor.swagger.io).
Описания фич: [посты (CRUD)](docs/features/posts.md) — реализовано создание, остальное — пошагово в [posts-implementation.md](docs/features/posts-implementation.md).

Логи: при `APP_ENV=local` — текстом в файл `LOG_FILE` (по умолчанию `app.log`), иначе — JSON в stdout.

## Слои

```
cmd/api/main.go        точка входа: только app.Run()
internal/
  app/                 сборка приложения: config → postgres/redis → repository → service → handler → server
  config/              настройки из .env с значениями по умолчанию
  handler/             HTTP-слой; handler.go — Handler и InitRoutes (все маршруты). Логирование — только здесь
    v1/                ручки API версии 1 (/api/v1/...): auth.go, user.go, dto.go
    middleware/        Logging, Recover, Auth
    response/          JSON-ответы и перевод ошибок в HTTP-статусы
  service/             бизнес-логика; service.go — интерфейсы и набор Services
    auth/              регистрация, вход, refresh, logout
    user/              профиль пользователя
  repository/          доступ к данным; repository.go — интерфейсы и набор Repositories
    user/ role/        Postgres
    session/           Redis: действующие сессии (пары токенов)
  models/              модели бизнес-логики и общие ошибки
pkg/                   подключение к Postgres и Redis, логгер
migrations/            SQL-схема и начальные роли
```

Свои структуры у каждого слоя:

| Слой | Пример | Теги |
|---|---|---|
| handler | `registerRequest`, `userResponse` | `json` |
| models | `User`, `RegisterInput` | нет |
| repository | `dbModel` | `db` |

## Токены

Оба токена — JWT, подписаны `JWT_SECRET`. Поле `type` не даёт выдать один токен за другой.

| | access | refresh |
|---|---|---|
| Поля | `sub` (id пользователя), `jti` (id сессии), `role_id`, `type=access`, `exp` | `sub`, `jti` (id сессии), `type=refresh`, `exp` |
| Живёт | `JWT_ACCESS_TTL` (15 минут) | `JWT_REFRESH_TTL` (30 дней) |
| Для чего | заголовок `Authorization: Bearer <token>` | получить новую пару (`/auth/refresh`) |

**Сессия** — одна пара токенов. При входе в Redis создаётся `session:<id>` → id пользователя, и этот id кладётся в `jti` обоих токенов. Токен действует, только пока его сессия есть в Redis:

- `/auth/refresh` удаляет старую сессию и создаёт новую: старые access и refresh сразу перестают работать.
- `/auth/logout` удаляет сессию: оба токена сразу перестают работать.
- Auth-middleware на каждый запрос проверяет подпись, срок и наличие сессии в Redis.

## Ручки

```bash
# Регистрация. role — user (служебные роли admin, moderator и т.д. так не получить)
curl -i -X POST localhost:8080/api/v1/auth/register \
  -d '{"email":"ivan@mail.com","password":"secret123","name":"Иван","role":"user"}'

# Вход → access_token и refresh_token
curl -i -X POST localhost:8080/api/v1/auth/login \
  -d '{"email":"ivan@mail.com","password":"secret123"}'

# Текущий пользователь
curl -i localhost:8080/api/v1/users/me -H "Authorization: Bearer <access_token>"

# Создать пост
curl -i -X POST localhost:8080/api/v1/posts -H "Authorization: Bearer <access_token>" \
  -d '{"title":"Мой первый пост","content":"Привет, мир!"}'

# Новая пара токенов
curl -i -X POST localhost:8080/api/v1/auth/refresh -d '{"refresh_token":"<refresh_token>"}'

# Выход
curl -i -X POST localhost:8080/api/v1/auth/logout -d '{"refresh_token":"<refresh_token>"}'
```

| Ошибка | Статус |
|---|---|
| битый JSON, неверные данные | 400 |
| неверный email/пароль, недействительный токен | 401 |
| пользователь не найден | 404 |
| email уже занят | 409 |
| всё остальное | 500 |
