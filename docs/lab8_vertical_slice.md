# Лабораторная работа 8. Вертикальный срез

## Тема среза

Добавлен вертикальный срез **"Срочная запись пациента"**.

Сценарий: пациенту требуется срочный прием у врача. Если выбранное время свободно, система сразу создает запись. Если время занято, система создает предложение переноса для пациента, чья запись занимает слот.

## Как проходит сценарий

1. Администратор открывает вкладку **Срочные** в GUI.
2. Выбирает пациента, врача, желаемое время и указывает причину срочности.
3. GUI отправляет запрос в backend API.
4. Backend вызывает usecase срочной записи.
5. Usecase проверяет, есть ли активная запись у выбранного врача на это время.
6. Если слот свободен, создается обычная запись на прием.
7. Если слот занят, создаются:
   - `emergency_request` — срочная заявка;
   - `reschedule_offer` — предложение переноса существующей записи.
8. В GUI предложение можно принять или отклонить.
9. При принятии старая запись переносится на предложенное время, а срочный пациент занимает освободившийся слот.
10. При отклонении заявка получает статус `NEEDS_REVIEW`, и администратор может повторить сценарий с другим временем.

## Какие слои затронуты

Срез проходит через все основные уровни приложения:

```text
GUI
  -> backend HTTP API
  -> usecase
  -> repository interface
  -> PostgreSQL repository
  -> PostgreSQL tables
```

## Новые таблицы БД

Файл:

```text
VitaLineClinic/migrations/schema.sql
```

Добавлены таблицы:

```sql
emergency_requests
reschedule_offers
```

`emergency_requests` хранит срочные заявки:

```text
id
patient_id
doctor_id
requested_date
reason
status
created_at
```

`reschedule_offers` хранит уведомления/предложения переноса:

```text
id
emergency_request_id
appointment_id
patient_id
doctor_id
old_date
proposed_date
status
created_at
```

## Основные файлы реализации

Доменная модель:

```text
VitaLineClinic/internal/domain/emergency.go
```

Бизнес-логика:

```text
VitaLineClinic/internal/usecase/emergency.go
```

PostgreSQL-репозиторий:

```text
VitaLineClinic/internal/storage/postgres/emergency.go
```

Backend API:

```text
VitaLineClinic/internal/delivery/http/api_handlers.go
```

GUI-клиент к backend:

```text
VitaLineClinic/internal/delivery/gui/backend_client.go
```

GUI-обработчики:

```text
VitaLineClinic/internal/delivery/gui/server.go
```

GUI-шаблон:

```text
VitaLineClinic/internal/delivery/gui/templates.go
```

## API

Добавлены endpoints:

```text
GET  /api/emergencies
POST /api/emergencies

GET  /api/reschedule-offers
POST /api/reschedule-offers/{id}/accept
POST /api/reschedule-offers/{id}/decline
```

## Статусы

Срочная заявка:

```text
WAITING_RESPONSE
SCHEDULED
NEEDS_REVIEW
CANCELLED
```

Предложение переноса:

```text
PENDING
ACCEPTED
DECLINED
```

## Что показать преподавателю

1. Открыть GUI:

```text
http://localhost:8090/emergencies
```

2. Показать вкладку **Срочные**.
3. Создать срочную запись на занятое время.
4. Показать, что появилась заявка и предложение переноса.
5. Нажать **Принять**.
6. Показать, что предложение стало `ACCEPTED`, а заявка стала `SCHEDULED`.
7. Объяснить, что это не отдельная кнопка, а вертикальный срез через БД, backend, бизнес-логику и GUI.

## Проверка

Команды:

```bash
cd ~/myhome/sem6/ppo_local/VitaLineClinic
go test ./cmd/clinic ./internal/appconfig ./internal/delivery/http ./internal/delivery/gui ./internal/storage/memory ./internal/usecase
make build
```

Для запуска всего приложения:

```bash
cd ~/myhome/sem6/ppo_local
docker compose up -d --build
```

GUI:

```text
http://localhost:8090/emergencies
```
