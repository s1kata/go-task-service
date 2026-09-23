# DONE — что уже работает

## Проект
`http-practive` — REST API на Go с слоями: Handler → TaskStore → PostgreSQL.

## Стек
- Go 1.22+ (метод-паттерны в ServeMux)
- PostgreSQL (через Docker)
- github.com/lib/pq (драйвер)

## Структура
- `internal/models` — Task, CreateTaskInput, UpdateTaskInput
- `internal/storge` — TaskStore (интерфейс), PostgresTaskStore, ErrTaskNotFound
- `internal/handler` — Handler, хендлеры, ResponseWithJSON / ResponseWithErrorJSON / WriteError
- `main/main.go` — точка входа, роуты, подключение БД

## Архитектура
- Handler зависит от интерфейса TaskStore, не от PostgresTaskStore
- Ошибки store (доменные) маппятся в HTTP через WriteError
- Ошибки клиента (валидация, JSON) обрабатываются в handler напрямую

## Роуты (метод-паттерны в mux)
- GET    /tasks            → список (List)
- POST   /tasks            → создание (Create) → 201
- GET    /tasks/{id}       → один таск (GetByID)
- PATCH  /tasks/{id}       → обновление (Update)
- GET    /ping             → ping
- GET    /hello            → hello
- /                        → 404 not-found

## Store-методы
- Create(ctx, CreateTaskInput) → *Task
- List(ctx, *bool) → []Task
- GetByID(ctx, id) → *Task
- Update(ctx, id, UpdateTaskInput) → *Task

## Обработка ошибок
- ErrTaskNotFound → 404
- sql.ErrNoRows → ErrTaskNotFound (в store)
- ошибки store → WriteError → 404 / 500
- ошибки клиента → ResponseWithErrorJSON → 400 / 405

## Инфраструктура
- docker-compose: только PostgreSQL (user/pass, db app, порт 5432)
- DATABASE_URL через env
- db.Ping() на старте — падаем сразу, если БД недоступна

## Что проверено curl
- POST /tasks       → 201 + task
- GET  /tasks       → 200 + [] или список
- GET  /tasks/{id}  → 200 / 404
- PATCH /tasks/{id} → 200 / 404
- GET  /tasks/abc   → 400
- GET  /ping        → 200
