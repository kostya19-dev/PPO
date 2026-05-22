# Лабораторная работа 9. Горизонтальный срез

## Тема среза

Добавлена поддержка MongoDB как нового способа работы с внешними данными.

До лабораторной работы №9 приложение работало через PostgreSQL-репозитории. В этой лабораторной работе добавлен отдельный слой MongoDB-репозиториев, который реализует те же интерфейсы бизнес-логики.

## Что изменилось

Новый storage-layer:

```text
VitaLineClinic/internal/storage/mongodb
```

Внутри реализованы репозитории:

```text
PatientRepository
AppointmentRepository
ScheduleRepository
VisitRepository
EmergencyRepository
```

Они реализуют те же usecase-интерфейсы, что и PostgreSQL-репозитории:

```text
usecase.PatientRepository
usecase.AppointmentRepository
usecase.ScheduleRepository
usecase.VisitRepository
usecase.EmergencyRepository
```

За счет этого бизнес-логика, HTTP API и GUI не зависят от конкретной СУБД.

## Конфигурация

Файл:

```text
VitaLineClinic/config/app.json
```

MongoDB включается через конфиг:

```json
{
  "database": {
    "driver": "mongodb",
    "url": "mongodb://localhost:27017",
    "name": "vitaline_clinic"
  }
}
```

Переключение на PostgreSQL возможно без изменения бизнес-логики:

```json
{
  "database": {
    "driver": "postgres",
    "url": "postgres://vitaline_app:vitaline_app@localhost:5432/vitaline_clinic?sslmode=disable",
    "name": "vitaline_clinic"
  }
}
```

Также поддержаны переменные окружения:

```text
VITALINE_DATABASE_DRIVER
DATABASE_URL
MONGODB_URI
MONGODB_DATABASE
```

## Где происходит выбор хранилища

Файл:

```text
VitaLineClinic/internal/bootstrap/bootstrap.go
```

В `bootstrap` выбирается конкретный набор репозиториев:

```text
database.driver = postgres -> internal/storage/postgres
database.driver = mongodb  -> internal/storage/mongodb
```

## MongoDB migration/seed

Файл:

```text
VitaLineClinic/migrations/mongo/001-init.js
```

Скрипт:

- создает индексы;
- добавляет демонстрационные данные батчами;
- использует `replaceOne` с `upsert: true`, поэтому повторный запуск не создает дубликаты;
- пишет состояние выполнения в коллекцию `schema_migrations`.

Коллекции:

```text
patients
appointments
schedule_slots
visits
emergency_requests
reschedule_offers
```

Если миграция оборвется на середине, уже обработанные документы останутся с теми же `_id`.
При повторном запуске скрипт не вставляет их второй раз, а заменяет по `_id`; необработанные документы будут добавлены.
Дополнительная защита от дублей находится в уникальных индексах, например `patients.email` и занятый слот врача.

## Docker Compose

Файл:

```text
docker-compose.yml
```

После лабораторной работы №9 приложение запускается на MongoDB:

```text
mongo       база данных
migrations одноразовый контейнер mongosh
api        backend
frontend   GUI
```

## C4 L4

```mermaid
flowchart LR
    UseCases["Usecase services"]
    Interfaces["Repository interfaces"]
    PgRepo["PostgreSQL repositories"]
    MongoRepo["MongoDB repositories"]
    PgDB["PostgreSQL"]
    MongoDB["MongoDB"]
    Config["database.driver"]

    UseCases --> Interfaces
    Config --> PgRepo
    Config --> MongoRepo
    Interfaces --> PgRepo
    Interfaces --> MongoRepo
    PgRepo --> PgDB
    MongoRepo --> MongoDB
```

## Что показать преподавателю

1. Diff merge request-а `lab9 -> master`.
2. Конфиг `database.driver = mongodb`.
3. Новый пакет `internal/storage/mongodb`.
4. `bootstrap.go`, где выбирается storage driver.
5. `docker-compose.yml`, где вместо PostgreSQL используется MongoDB.
6. `migrations/mongo/001-init.js`.
7. Работа GUI после запуска через Docker Compose.

## Запуск

```bash
cd ~/myhome/sem6/ppo_local
docker compose up -d --build
docker compose ps
```

GUI:

```text
http://localhost:8090/patients
```

Проверка MongoDB:

```bash
docker exec -it vitaline_mongo mongosh vitaline_clinic --eval 'db.patients.find().limit(5)'
```
