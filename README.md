# depviz

Инструмент визуализации графа зависимостей для менеджера пакетов
Alpine Linux (apk). Учебный проект по дисциплине
«Конфигурационное управление», Вариант 7, ИКБО-10-25, РТУ МИРЭА.

Готовые менеджеры пакетов и сторонние библиотеки для извлечения
зависимостей **не используются**.

## Текущий статус

* **Этап 1 — Конфигурация.** TOML-файл, вывод всех параметров
  в формате `key = value`, валидация каждого параметра.
* **Этап 2 — Сбор данных.** Скачивание `APKINDEX.tar.gz`,
  поиск пакета по имени и версии, вывод прямых зависимостей.
* **Этап 3 — Основные операции.** Построение графа зависимостей
  BFS с рекурсией, учёт `max_depth`, фильтр `filter_substring`,
  обработка циклов, режим тестового репозитория.
* Этапы 4–5 — в разработке.

## Структура проекта

    depviz/
    ├── src/
    │   ├── main.go         точка входа
    │   ├── app/            вывод конфигурации и графа
    │   ├── apk/            скачивание и парсинг APKINDEX
    │   ├── config/         разбор и валидация TOML
    │   └── graph/          построение и обход графа зависимостей
    ├── tests/              тесты
    ├── testdata/           примеры конфигов и тестовых данных
    ├── scripts/            скрипты запуска и проверок
    ├── Makefile            сборка, тесты, запуск
    ├── go.mod              модуль Go
    └── README.md

## Требования

* Go 1.27 или новее.
* POSIX-совместимая оболочка (Linux, macOS, WSL, Git Bash) для
  скриптов в `scripts/`.

## Сборка

    make build

Бинарник появится в `bin/depviz` (в Windows — `bin/depviz.exe`).

## Запуск

Тестовый репозиторий (без сети):

    make run
    make run-graph
    bash scripts/run-graph.sh testdata/config.test.toml

Реальный APKINDEX:

    make run-apk
    bash scripts/run-apk.sh

С произвольным конфигом:

    ./bin/depviz path/to/config.toml
    bash scripts/run.sh path/to/config.toml

## Тесты

    make test
    bash scripts/run-errors.sh

## Параметры конфигурации

Источник — TOML-файл. Примеры — `testdata/config.example.toml`,
`testdata/config.test.toml`.

| Ключ              | Тип    | Описание                                           |
|-------------------|--------|----------------------------------------------------|
| `package_name`    | string | Имя анализируемого пакета                          |
| `package_version` | string | Версия анализируемого пакета (для этапа 2)         |
| `repository_url`  | string | URL apk-репозитория                                |
| `test_repository` | bool   | Режим тестового репозитория                        |
| `test_file_path`  | string | Путь к файлу тестового репозитория                 |
| `output_image`    | string | Имя файла с изображением графа                     |
| `max_depth`       | int    | Максимальная глубина анализа (`>= 1`)              |
| `filter_substring`| string | Подстрока для фильтрации пакетов                   |

## Тестовый репозиторий

Файл `testdata/graph_test.txt` описывает граф в формате:

    A: B C
    B: D
    C: D E
    D: E
    E: A

Строка `ИМЯ: ЗАВ1 ЗАВ2 ...` — зависимости через пробел.
Пустое поле после `:` — пакет без зависимостей.

## Обработка ошибок

Каждый параметр валидируется отдельной sentinel-ошибкой.
Проверяются: отсутствие файла, пустые обязательные поля,
`max_depth < 1`, `test_repository=true` без `test_file_path`.

Пример:

    $ ./bin/depviz nope.toml
    Ошибка конфигурации: config: файл не найден: nope.toml
    $ echo $?
    1

## Пример вывода

Граф в тестовом режиме:

    $ make run-graph
    root: A
    nodes: 5
    graph:
      A (depth=0): B, C
      B (depth=1): D
      D (depth=2): E
      E (depth=3): A
      C (depth=1): D, E

Граф по реальному APKINDEX:

    $ make run-apk
    root: curl
    nodes: 3
    graph:
      curl (depth=0): libcurl
      libcurl (depth=1): ca-certificates
      ca-certificates (depth=2): (none)

## Лицензия

Учебный проект, свободное использование в образовательных целях.