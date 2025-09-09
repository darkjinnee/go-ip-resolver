# Postman Collection для Go IP Resolver

Этот каталог содержит Postman коллекцию с примерами всех API запросов для Go IP Resolver.

## 📁 Файлы

- `Go-IP-Resolver.postman_collection.json` - Основная коллекция с запросами
- `Go-IP-Resolver-Environment.postman_environment.json` - Переменные окружения
- `README.md` - Документация (этот файл)

## 🚀 Быстрый старт

### 1. Импорт в Postman

1. Откройте Postman
2. Нажмите **Import**
3. Выберите файл `Go-IP-Resolver.postman_collection.json`
4. Выберите файл `Go-IP-Resolver-Environment.postman_environment.json`
5. Нажмите **Import**

### 2. Настройка окружения

1. В Postman выберите окружение **Go IP Resolver Environment**
2. Убедитесь, что `base_url` установлен в `http://localhost:8080`
3. При необходимости измените `group_name` для тестирования

### 3. Запуск сервера

```bash
# Запустите Go IP Resolver
make run

# Или с кастомными параметрами
./bin/ip-resolver -addr :8080
```

## 📋 Структура коллекции

### 🔍 Groups
- **Get All Groups** - Получить список всех доступных групп доменов

### 🌐 Resolve
- **Resolve Group - All IPs** - Получить все IP адреса для группы
- **Resolve Group - IPv4 Only** - Получить только IPv4 адреса
- **Resolve Group - IPv6 Only** - Получить только IPv6 адреса

### 🌍 Resolve All
- **Resolve All Groups - All IPs** - Получить все IP для всех групп
- **Resolve All Groups - IPv4 Only** - Получить только IPv4 для всех групп
- **Resolve All Groups - IPv6 Only** - Получить только IPv6 для всех групп

### 💾 Cache
- **Get Cache Statistics** - Получить статистику кэша

### 📝 Examples
- **YouTube - All IPs** - Пример для YouTube
- **OpenAI - IPv4 Only** - Пример для OpenAI (только IPv4)
- **GitHub Copilot - IPv6 Only** - Пример для Copilot (только IPv6)
- **Instagram - All IPs** - Пример для Instagram
- **Cursor AI - All IPs** - Пример для Cursor AI

## 🔧 Переменные окружения

| Переменная | Значение по умолчанию | Описание |
|------------|----------------------|----------|
| `base_url` | `http://localhost:8080` | Базовый URL сервера |
| `group_name` | `youtube` | Группа для тестирования |
| `youtube_group` | `youtube` | Группа YouTube |
| `google_group` | `google` | Группа Google |
| `github_group` | `github` | Группа GitHub |
| `instagram_group` | `instagram` | Группа Instagram |
| `openai_group` | `openai` | Группа OpenAI |
| `copilot_group` | `copilot` | Группа GitHub Copilot |
| `cursor_group` | `cursor` | Группа Cursor AI |
| `ipv4_type` | `ipv4` | Тип IP: IPv4 |
| `ipv6_type` | `ipv6` | Тип IP: IPv6 |

## 📊 Примеры ответов

### Список групп
```json
{
  "groups": [
    "youtube",
    "google", 
    "github",
    "instagram",
    "openai",
    "copilot",
    "cursor"
  ]
}
```

### Резолвинг группы
```json
[
  {
    "domain": "youtube.com",
    "ips": [
      "142.251.36.206",
      "2a00:1450:4016:809::200e"
    ]
  }
]
```

### Статистика кэша
```json
{
  "groups_count": 7,
  "total_entries": 21,
  "groups": {
    "youtube": 3,
    "google": 3,
    "github": 3,
    "instagram": 3,
    "openai": 3,
    "copilot": 3,
    "cursor": 3
  }
}
```

## 🎯 Рекомендации по использованию

### 1. Тестирование производительности
- Используйте **Resolve All Groups** для проверки общей производительности
- Сравните время ответа с кэшем и без кэша

### 2. Проверка фильтрации
- Тестируйте `type=ipv4` и `type=ipv6` для разных групп
- Убедитесь, что фильтрация работает корректно

### 3. Мониторинг кэша
- Регулярно проверяйте `/cache/stats` для мониторинга
- Отслеживайте количество записей в кэше

### 4. Тестирование групп
- Используйте готовые примеры в папке **Examples**
- Тестируйте все доступные группы доменов

## 🔍 Отладка

### Проблемы с подключением
1. Убедитесь, что сервер запущен на `http://localhost:8080`
2. Проверьте, что порт 8080 не занят другими приложениями
3. Попробуйте изменить `base_url` в переменных окружения

### Ошибки API
1. Проверьте логи сервера
2. Убедитесь, что группа доменов существует
3. Проверьте правильность параметров запроса

### Проблемы с кэшем
1. Проверьте статистику кэша через `/cache/stats`
2. Убедитесь, что кэш обновляется периодически
3. Проверьте TTL настройки сервера

## 📚 Дополнительные ресурсы

- [Документация API](../README.md) - Полная документация API
- [Конфигурация доменов](../configs/domains.json) - Настройка групп доменов
- [Makefile](../Makefile) - Команды сборки и запуска

## 🤝 Поддержка

Если у вас возникли проблемы с использованием Postman коллекции:

1. Проверьте, что сервер Go IP Resolver запущен
2. Убедитесь, что импортированы и коллекция, и окружение
3. Проверьте настройки переменных окружения
4. Обратитесь к документации API

---

**Примечание**: Эта коллекция создана для Go IP Resolver v1.0.0 и может потребовать обновления при изменении API.
