# Go IP Resolver

Универсальный HTTP сервер для резолвинга IP адресов доменов с конфигурируемыми группами доменов.

## Возможности

- Резолвинг IP адресов для групп доменов
- Конфигурируемые группы доменов через JSON файл
- HTTP API для получения IP адресов
- Настраиваемые таймауты DNS запросов
- Список доступных групп доменов
- **Кэширование с TTL для повышения производительности**
- **Периодическое обновление кэша каждые 10 минут**
- **Фильтрация по типу IP (IPv4/IPv6)**
- **Защита от утечек памяти**

## Структура проекта

```
├── cmd/ip-resolver/          # Основное приложение
├── internal/
│   ├── cache/               # Кэширование с TTL
│   ├── config/              # Конфигурация
│   ├── resolver/            # DNS резолвер
│   └── transport/           # HTTP транспорт
├── configs/                 # Конфигурационные файлы
│   └── domains.json         # Группы доменов
├── postman/                 # Postman коллекция
│   ├── Go-IP-Resolver.postman_collection.json
│   ├── Go-IP-Resolver-Environment.postman_environment.json
│   └── README.md            # Документация Postman
├── bin/                     # Собранные бинарники
└── Makefile                 # Команды сборки и запуска
```

## Установка и запуск

### Сборка

```bash
make build
```

### Запуск

```bash
make run
```

### Запуск в режиме разработки

```bash
make run-dev
```

### Другие команды

```bash
make test      # Запуск тестов
make clean     # Очистка
make deps      # Установка зависимостей
make fmt       # Форматирование кода
make lint      # Линтинг
make check     # Полная проверка
```

## Конфигурация

Домены настраиваются в файле `configs/domains.json`:

```json
{
  "youtube": [
    "youtube.com",
    "www.youtube.com",
    "m.youtube.com"
  ],
  "google": [
    "google.com",
    "www.google.com"
  ]
}
```

## API

### Получить IP адреса для группы доменов

```
GET /resolve?group=<group_name>&type=<ip_type>
```

**Параметры:**
- `group` (обязательный) - название группы доменов
- `type` (опциональный) - тип IP адресов: `ipv4`, `ipv6` или не указывать для всех

**Примеры:**
```bash
# Все IP адреса
curl "http://localhost:8080/resolve?group=youtube"

# Только IPv4
curl "http://localhost:8080/resolve?group=youtube&type=ipv4"

# Только IPv6
curl "http://localhost:8080/resolve?group=youtube&type=ipv6"
```

**Ответ:**
```json
[
  {
    "domain": "youtube.com",
    "ips": ["142.250.191.78", "2a00:1450:4010:c0a::65"]
  },
  {
    "domain": "www.youtube.com", 
    "ips": ["142.250.191.78", "2a00:1450:4010:c0a::65"]
  }
]
```

### Получить IP адреса для всех групп сразу

```
GET /resolve-all?type=<ip_type>
```

**Параметры:**
- `type` (опциональный) - тип IP адресов: `ipv4`, `ipv6` или не указывать для всех

**Примеры:**
```bash
# Все IP адреса
curl "http://localhost:8080/resolve-all"

# Только IPv4
curl "http://localhost:8080/resolve-all?type=ipv4"

# Только IPv6
curl "http://localhost:8080/resolve-all?type=ipv6"
```

**Ответ:**
```json
{
  "youtube": [
    {
      "domain": "youtube.com",
      "ips": ["142.250.191.78", "2a00:1450:4010:c0a::65"]
    }
  ],
  "google": [
    {
      "domain": "google.com", 
      "ips": ["142.250.191.78", "2a00:1450:4010:c0a::65"]
    }
  ]
}
```

### Получить список доступных групп

```
GET /groups
```

**Пример:**
```bash
curl "http://localhost:8080/groups"
```

**Ответ:**
```json
{
  "groups": ["youtube", "google", "github"]
}
```

### Получить статистику кэша

```
GET /cache/stats
```

**Пример:**
```bash
curl "http://localhost:8080/cache/stats"
```

**Ответ:**
```json
{
  "groups_count": 3,
  "total_entries": 9,
  "groups": {
    "youtube": 3,
    "google": 3,
    "github": 3
  }
}
```

## Параметры командной строки

- `-config` - путь к файлу конфигурации (по умолчанию: `configs/domains.json`)
- `-addr` - адрес сервера (по умолчанию: `:8080`)
- `-timeout` - таймаут DNS запросов (по умолчанию: `2s`)
- `-cache-ttl` - время жизни кэша (по умолчанию: `15m`)
- `-update-interval` - интервал обновления кэша (по умолчанию: `10m`)

## Примеры использования

```bash
# Запуск с кастомной конфигурацией
./bin/ip-resolver -config /path/to/config.json -addr :9090

# Запуск с увеличенным таймаутом
./bin/ip-resolver -timeout 10s
```

## Curl
```bash
$ sudo curl -L "https://github.com/darkjinnee/go-ip-resolver/releases/download/v1.0.0/ip-resolver-linux" -o /usr/local/bin/go-ip-resolver 
$ sudo chmod +x /usr/local/bin/go-ip-resolver
```

## Postman коллекция

Для удобного тестирования API включена готовая Postman коллекция:

### 📁 Файлы
- `postman/Go-IP-Resolver.postman_collection.json` - Коллекция запросов
- `postman/Go-IP-Resolver-Environment.postman_environment.json` - Переменные окружения
- `postman/README.md` - Подробная документация

### 🚀 Быстрый импорт
1. Откройте Postman
2. Импортируйте оба JSON файла из папки `postman/`
3. Выберите окружение "Go IP Resolver Environment"
4. Запустите сервер: `make run`
5. Готово! Все запросы настроены и готовы к использованию

### 📋 Что включено
- ✅ Все API endpoints с примерами
- ✅ Готовые запросы для всех групп доменов
- ✅ Примеры фильтрации IPv4/IPv6
- ✅ Тесты производительности
- ✅ Мониторинг кэша
- ✅ Переменные окружения для удобства

Подробная документация: [postman/README.md](postman/README.md)
Microservice that periodically resolves IP addresses for specified domains (e.g., YouTube, Telegram, or other blocked services).
