# TODO — что осталось добить

## Баги (срочно, вечер/завтра)

### 1. List — фильтр по completed не работает
- SQL: WHERE убран, но QueryContext всё ещё принимает `completed`
  → ошибка «ожидал 0 параметров, получил 1» → 500
- Решить: A2 (COALESCE) или A1 (динамический SQL)
  - A2: WHERE ($1::boolean IS NULL OR completed = $1)
  - A1: собирать SQL по условиям
- `?completed=true` сейчас игнорируется

### 2. h.Get — не читает query-параметр
- Сейчас: h.store.List(r.Context(), nil) — всегда nil
- Надо: читать `completed` из r.URL.Query(), парсить в *bool, передавать

### 3. Update в storge — два бага
- Scan: `tasks.Title` без `&` → надо `&tasks.Title`
- sql.ErrNoRows: `return nil, err` → надо `return nil, ErrTaskNotFound`

### 4. WriteError — опечатка
- "intnernail server error" → "internal server error"

### 5. TaskByID / Update в handler — TrimPrefix вместо PathValue
- Сейчас: strings.TrimPrefix(path, "/tasks/") + Atoi
- Надо: r.PathValue("id") + Atoi
- Убрать импорт strings, если больше не используется

### 6. TaskByID — сообщение в теле
- "400" → "invalid id"
- "405" → "method not allowed"

## Функциональность (следующие шаги)

### 7. DELETE /tasks/{id}
- store: Delete(ctx, id) error
- handler: Delete
- SQL: DELETE FROM tasks WHERE id = $1
- Проверка RowsAffected == 0 → ErrTaskNotFound
- Ответ: 204 No Content (без тела)

### 8. PATCH /tasks/{id}/complete — отдельная ручка
- Только completed, без title/description
- store: SetCompleted(ctx, id, completed bool) (*Task, error)
- Семантика: отдельная бизнес-операция, не смешивать с Update

### 9. Кириллица — проверить
- В БД через psql: SELECT * FROM tasks;
- Если в БД ?????? — проблема на входе (PowerShell/curl)
- Если в БД нормально, в консоли ?????? — только отображение

## Архитектура (на будущее — прокси/джобер)

### 10. NOTES.md — roadmap
- Прокси: отдельный сервис, использует TaskStore или REST API
- Джобер: pull тасков, смена статуса через SetCompleted
- Статусы: waited / in_progress / failed / completed
  (заменить bool → string, когда появятся требования)
- Retry: attempts, last_error, backoff
- Handoff: REST API → очередь → джобер
- Идемпотентность: как не запустить одну джобу дважды

## Принципы (чтобы не забывать)
- Один шаг — до зелёного curl — потом следующий
- Сверяться с работающим кодом (GetByID, Create)
- «Всё» ≠ работает. Пока curl не показал — не «всё»
- Активная модель: сначала слова (сигнатура, edge cases), потом тело