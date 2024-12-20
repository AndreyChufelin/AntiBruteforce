# Anti-Bruteforce
[![Go Report Card](https://goreportcard.com/badge/github.com/AndreyChufelin/AntiBruteforce)](https://goreportcard.com/report/github.com/AndreyChufelin/AntiBruteforce)

Микросервис для защиты систем авторизации от атак типа брутфорс. Сервис реализует ограничение частоты попыток входа на основе алгоритма **GCRA (Generic Cell Rate Algorithm)**, управляет списками разрешённых и заблокированных IP-адресов, предоставляет API для взаимодействия и CLI для администрирования.

## Основные функции:
- **Ограничение частоты запросов**:
    - До 10 попыток в минуту для конкретного логина.
    - До 100 попыток в минуту для конкретного пароля (защита от обратного брутфорса).
    - До 1000 попыток в минуту с одного IP.
    - Алгоритм ограничений реализован на основе **GCRA**, обеспечивающего эффективное управление частотой запросов при минимальных накладных расходах.
- **Управление списками IP**:
    - Белый список (whitelist): автоматическое разрешение авторизации.
    - Чёрный список (blacklist): автоматическое блокирование авторизации.
    - Поддержка работы с подсетями (например, 192.1.1.0/25).
- **API**:
    - gRPC методы для проверки попыток авторизации, сброса бакетов, управления списками IP.
    - Методы настройки и сброса:
        - Проверка попытки авторизации.
        - Сброс бакетов для заданных логина или IP.
        - Добавление и удаление подсетей в whitelist/blacklist.
- **CLI**:
    - Управление сервисом через командную строку:
        - Сброс бакетов.
        - Управление whitelist/blacklist.

**Технические особенности:**
- Использованы **Redis** для хранения бакетов и **PostgreSQL** для управления настройками и списками IP.
- Алгоритм **GCRA** применён для эффективного и точного контроля частоты запросов.
- Возможность запуска через `Docker Compose` с командами `make build`, `make run` и `make test`.
- Легкозаменяемая архитектура для работы с различными хранилищами (например, замена Redis или PostgreSQL).

**Тестирование:**
- Покрытие функциональности:
    - Юнит-тесты для проверки работы алгоритма **GCRA**.
    - Интеграционные тесты для проверки всех API-вызовов.
- Проверка сценариев работы с whitelist/blacklist и обработки ограничений.
