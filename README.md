# Notify Hub

Централизованный сервис для отправки уведомлений через различные каналы (Telegram, Email, SMS).

## Особенности

- ✅ gRPC API
- ✅ Поддержка Telegram и Email
- ✅ Безопасное хранение конфигураций
- ✅ Асинхронная очередь
- ✅ Clean Architecture
- ✅ Логирование
- ✅ Docker поддержка

## Быстрый старт

1. Клонируйте репозиторий
2. Настройте `configs/integrations.yaml`
3. Запустите сервер:

```bash
make deps
make proto
make run
```
