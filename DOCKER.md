# Docker Deployment

## Описание

Этот документ описывает развертывание Go IP Resolver в Docker контейнере для продакшн среды.

## Файлы

- `Dockerfile` - многоэтапная сборка Go приложения
- `docker-compose.yml` - конфигурация для продакшн стенда
- `.dockerignore` - исключения для Docker контекста

## Быстрый старт

### Сборка и запуск

```bash
# Сборка образа
docker-compose build

# Запуск в фоне
docker-compose up -d

# Просмотр логов
docker-compose logs -f

# Остановка
docker-compose down
```

### Проверка работы

```bash
# Проверка health check
docker-compose ps

# Тест API
curl http://localhost:8080/groups
curl "http://localhost:8080/resolve?group=youtube&type=ipv4"
```

## Конфигурация

### Переменные окружения

- `TZ=UTC` - часовой пояс
- `GIN_MODE=release` - режим работы Gin (если используется)

### Порты

- `8080` - HTTP API сервер

### Volumes

- `./configs/domains.json:/app/configs/domains.json:ro` - конфигурация доменов (только чтение)
- `./logs:/app/logs` - директория для логов

## Мониторинг

### Health Check

Контейнер автоматически проверяет здоровье через endpoint `/groups` каждые 30 секунд.

### Логирование

Логи настроены с ротацией:
- Максимальный размер файла: 10MB
- Максимальное количество файлов: 3

### Ограничения ресурсов

- CPU: максимум 1.0, резерв 0.5
- Память: максимум 512MB, резерв 256MB

## Безопасность

- Приложение запускается под непривилегированным пользователем `appuser`
- Статическая сборка без CGO зависимостей
- Минимальный базовый образ Alpine Linux

## API Endpoints

- `GET /groups` - список доступных групп
- `GET /resolve?group=<name>&type=<ipv4|ipv6>` - резолвинг группы
- `GET /resolve-all?type=<ipv4|ipv6>` - резолвинг всех групп
- `GET /resolve-flat?group=<name>&type=<ipv4|ipv6>` - плоский список IP
- `GET /resolve-flat-all?type=<ipv4|ipv6>` - плоский список всех IP
- `GET /resolve-flat-all-ipv4` - все IPv4 адреса
- `GET /resolve-flat-all-ipv6` - все IPv6 адреса
- `GET /cache/stats` - статистика кэша
