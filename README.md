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
* Этапы 3–5 — в разработке.

## Структура проекта

    depviz/
    ├── src/
    │   ├── main.go         точка входа
    │   ├── app/            вывод конфигурации и зависимостей
    │   ├── apk/            скачивание и парсинг APKINDEX
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

Прямые зависимости из реального репозитория:

    bash scripts/run-apk.sh

Все параметры конфигурации печатаются в **stderr** в формате
`key = value` (как в эталонном `shell-emulator`).

## Тесты

    make test
    bash scripts/run-errors.sh

## Параметры конфигурации

Источник — TOML-файл. Пример — `testdata/config.example.toml`.

| Ключ              | Тип    | Описание                                           |
|-------------------|--------|----------------------------------------------------|
| `package_name`    | string | Имя анализируемого пакета                          |
| `package_version` | string | Версия анализируемого пакета                       |
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
    package.version    = 8.5.0-r0
    repository.url     = https://dl-cdn.alpinelinux.org/alpine/v3.20/main/x86_64/
    test.repository    = false
    test.file.path     = 
    output.image       = deps.svg
    max.depth          = 5
    filter.substring   = test

Прямые зависимости:

    $ bash scripts/run-apk.sh
    package: curl 8.5.0-r0
    dependencies:
      libcurl
      so:libc.musl-x86_64.so.1
      ...

## Лицензия

Учебный проект, свободное использование в образовательных целях.