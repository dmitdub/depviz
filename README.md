# depviz

Инструмент визуализации графа зависимостей для менеджера пакетов
Alpine Linux (apk). Учебный проект по дисциплине
«Конфигурационное управление», Вариант 7, ИКБО-10-25, РТУ МИРЭА.

Готовые менеджеры пакетов и сторонние библиотеки для извлечения
зависимостей **не используются**.

## Текущий статус

* **Этап 1 — Конфигурация.** TOML-файл, вывод всех параметров
  в формате `key = value`, валидация каждого параметра.
* Этапы 2–5 — в разработке.

## Структура проекта

    depviz/
    ├── src/
    │   ├── main.go         точка входа
    │   ├── app/            вывод конфигурации
    │   └── config/         разбор и валидация TOML
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

    make run
    ./bin/depviz testdata/config.example.toml
    bash scripts/run.sh path/to/config.toml

Все параметры конфигурации печатаются в **stderr** в формате
`key = value` (как в эталонном `shell-emulator`).

## Тесты

    make test
    bash scripts/run-errors.sh

## Параметры конфигурации (этап 1)

Источник — TOML-файл. Пример — `testdata/config.example.toml`.

| Ключ              | Тип    | Описание                                           |
|-------------------|--------|----------------------------------------------------|
| `package_name`    | string | Имя анализируемого пакета                          |
| `repository_url`  | string | URL apk-репозитория                                |
| `test_repository` | bool   | Режим тестового репозитория                        |
| `test_file_path`  | string | Путь к файлу тестового репозитория                 |
| `output_image`    | string | Имя файла с изображением графа                     |
| `max_depth`       | int    | Максимальная глубина анализа (`>= 1`)              |
| `filter_substring`| string | Подстрока для фильтрации пакетов                   |

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

    $ make run
    configuration:
    package.name       = curl
    repository.url     = https://dl-cdn.alpinelinux.org/alpine/v3.20/main/x86_64/
    test.repository    = false
    test.file.path     = 
    output.image       = deps.svg
    max.depth          = 5
    filter.substring   = test

## Лицензия

Учебный проект, свободное использование в образовательных целях.