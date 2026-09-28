# TODO — что осталось добить

## Баги (срочно, вечер/завтра)
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