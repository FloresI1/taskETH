# TaskETH

TaskETH - это приложение на Go для анализа изменений балансов Ethereum кошельков по блокам. Приложение использует API GetBlock для получения данных о блоках и балансах.

## Оглавление

- [Установка](#установка)
- [Настройка](#настройка)
- [Использование](#использование)
- [Структура проекта](#структура-проекта)

## Установка

1. Склонируйте репозиторий:

    ```sh
    git clone https://github.com/FloresI1/taskETH.git
    ```

2. Перейдите в директорию проекта:

    ```sh
    cd taskETH
    ```

3. Установите зависимости:

    ```sh
    go mod tidy
    ```

## Настройка

1. Создайте файл `.env` в корневой директории проекта и добавьте следующие переменные окружения:

    ```env
    GETBLOCK_API_KEY=your_api_key_here
    GETBLOCK_BASE_URL=https://eth.getblock.io/mainnet/
    NUM_WORKERS=5
    ```

    - `GETBLOCK_API_KEY`: Ваш API ключ для GetBlock.
    - `GETBLOCK_BASE_URL`: Базовый URL для GetBlock API.
    - `NUM_WORKERS`: Количество параллельных потоков для обработки адресов.

## Использование

1. Запустите приложение:

    ```sh
    go run main.go
    ```

    Пример вывода:

    ```
    Start
    Адрес: 0x123..., Баланс на блоке 0xabc... (с учётом 100 блоков назад): 1000000000000000000 Wei
    Адрес: 0x123..., Баланс на начальном блоке: 1000000000000000000 Wei, Баланс на конечном блоке: 900000000000000000 Wei, Изменение: 100000000000000000 Wei
    Кошелек с максимальным изменением: 0x123..., изменение: 100000000000000000 Wei
    ```

## Структура проекта

- `main.go`: Основной файл, запускающий приложение.
- `request/`: Пакет для выполнения запросов к API.
- `utils/`: Пакет утилит для работы с данными.
- `workers/`: Пакет для параллельной обработки адресов.
