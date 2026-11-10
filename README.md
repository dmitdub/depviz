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
* **Этап 4 — Дополнительные операции.** Обратные зависимости:
  кто зависит от заданного пакета. Обход с учётом `max_depth`
  и `filter_substring`, обработка циклов.
* **Этап 5 — Визуализация.** Генерация PlantUML-текста,
  SVG-файла через `plantuml.jar`, опциональное ASCII-дерево.

## Структура проекта

    depviz/
    ├── src/
    │   ├── main.go         точка входа
    │   ├── app/            вывод конфигурации и графа
    │   ├── apk/            скачивание и парсинг APKINDEX, обратные зависимости
    │   ├── config/         разбор и валидация TOML
    │   ├── graph/          построение и обход графа (прямой и обратный)
    │   └── render/         PlantUML, SVG, ASCII-дерево
    ├── tests/              тесты
    ├── testdata/           примеры конфигов и тестовых данных
    ├── scripts/            скрипты запуска и проверок
    ├── Makefile            сборка, тесты, запуск
    ├── go.mod              модуль Go
    └── README.md

## Требования

* Go 1.27 или новее.
* Java 11+ (для `plantuml.jar` на этапе 5).
* POSIX-совместимая оболочка (Linux, macOS, WSL, Git Bash) для
  скриптов в `scripts/`.

## Подготовка `plantuml.jar`

Для этапа 5 нужен `plantuml.jar`. Скачать:

    https://github.com/plantuml/plantuml/releases/latest/download/plantuml.jar

Положить в `bin/plantuml.jar`:

    mkdir -p bin
    # скачать plantuml.jar в bin/

Проверить:

    java -jar bin/plantuml.jar -version

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

Обратные зависимости:

    make run-reverse
    bash scripts/run-reverse.sh

Визуализация с PlantUML и ASCII-деревом:

    make run-render
    bash scripts/run-render.sh testdata/config.example.toml

Три пакета для демонстрации визуализации:

    bash scripts/run-three.sh

Создаёт `curl.svg`, `vim.svg`, `openssl.svg`.

## Тесты

    make test
    bash scripts/run-errors.sh

## Параметры конфигурации

Источник — TOML-файл. Примеры — `testdata/config.example.toml`,
`testdata/config.test.toml`, `testdata/config.reverse.toml`,
`testdata/config.curl.toml`, `testdata/config.vim.toml`,
`testdata/config.openssl.toml`.

| Ключ              | Тип    | Описание                                           |
|-------------------|--------|----------------------------------------------------|
| `package_name`    | string | Имя анализируемого пакета                          |
| `package_version` | string | Версия анализируемого пакета (для этапа 2)         |
| `repository_url`  | string | URL apk-репозитория                                |
| `test_repository` | bool   | Режим тестового репозитория                        |
| `test_file_path`  | string | Путь к файлу тестового репозитория                 |
| `output_image`    | string | Имя файла с изображением графа (SVG)               |
| `max_depth`       | int    | Максимальная глубина анализа (`>= 1`)              |
| `filter_substring`| string | Подстрока для фильтрации пакетов                   |
| `reverse`         | bool   | Режим обратных зависимостей                        |
| `ascii_tree`      | bool   | Печатать ASCII-дерево после PlantUML               |

## Тестовый репозиторий

Файл `testdata/graph_test.txt` описывает граф в формате:

    A: B C
    B: D
    C: D E
    D: E
    E: A

Строка `ИМЯ: ЗАВ1 ЗАВ2 ...` — зависимости через пробел.
Пустое поле после `:` — пакет без зависимостей.

## Обратные зависимости

При `reverse = true` приложение строит обратный граф: находит
все пакеты, которые зависят от заданного. Обход рекурсивный,
с учётом `max_depth` и `filter_substring`, циклы обрабатываются
через `seen`.

Пример конфига — `testdata/config.reverse.toml`.

## Визуализация (этап 5)

При запуске приложение:
1. Строит граф (прямой или обратный).
2. Печатает его в `stderr` в формате `root` / `nodes` / `graph`.
3. Генерирует PlantUML-текст и печатает его в `stderr`.
4. Рендерит SVG через `java -jar bin/plantuml.jar -tsvg`.
5. Если `ascii_tree = true` — печатает ASCII-дерево.

Пример конфига — `testdata/config.curl.toml`:

    package_name     = "curl"
    repository_url   = "https://dl-cdn.alpinelinux.org/alpine/v3.20/main/x86_64/"
    output_image     = "curl.svg"
    max_depth        = 3
    ascii_tree       = true

Запуск:

    bash scripts/run-render.sh testdata/config.curl.toml

Результат:
- `curl.svg` — изображение графа;
- в `stderr` — PlantUML-текст и ASCII-дерево.

Для трёх пакетов — `bash scripts/run-three.sh`.

## Сравнение с `apk`

На Windows `apk` недоступен, поэтому сравнение выполнено с веб-интерфейсом
Alpine: `https://pkgs.alpinelinux.org/package/v3.20/main/x86_64/curl`.

Для `curl` веб-интерфейс показывает зависимости:
`libcurl so:libc.musl-x86_64.so.1 so:libcurl.so.4 so:libz.so.1`.

Наш инструмент выводит: `libcurl` (без `so:`-зависимостей).

**Расхождение:** мы фильтруем `so:`-зависимости, потому что это библиотеки,
а не пакеты Alpine. Для построения графа пакетов они не нужны.
Прямая пакетная зависимость (`libcurl`) совпадает.

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

Визуализация:

    $ bash scripts/run-render.sh testdata/config.curl.toml
    @startuml
    "curl" --> "libcurl"
    "libcurl" --> "ca-certificates"
    @enduml
    SVG сохранён: curl.svg
    curl
    └── libcurl
        └── ca-certificates

## Лицензия

Учебный проект, свободное использование в образовательных целях.