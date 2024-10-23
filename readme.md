# Go-words Telegram Bot

Go-words Telegram Bot — это реализация популярной игры в угадывание слов, написанная на языке Go (Golang). Игроку загадывается случайное слово, и его задача — угадать это слово за несколько попыток, получая подсказки по каждой букве. 

## Описание проекта

Бот использует [Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api/v5) и взаимодействует с базой данных MySQL, в которой хранятся загаданные слова. Проект построен с использованием пакета database/sql для работы с базой данных и запускается в контейнере Docker.

## Функциональность

- Выбор уровня сложности. Пользователь может выбрать уровень сложности игры.
- Угадывание слова. Игрок вводит слова, пытаясь угадать загаданное.
- Проверка корректности. Реализована логика проверки существования слова и сравнения с загаданным словом.
- Подсказки. Игроку предоставляются подсказки, помогающие угадать слово.

## Структура проекта

- [`database.go`](https://github.com/notblinkyet/go-words/blob/master/game/database.go): Запрос к базе данных для получения случайного слова с использованием пакета database/sql.
- [`game.go`](https://github.com/notblinkyet/go-words/blob/master/game/game.go): Реализация основных функций игры, включая проверку правильности введённого слова, сравнение с загаданным словом и другие вспомогательные функции.
- [`game_test.go`](https://github.com/notblinkyet/go-words/blob/master/game/game_test.go): Тесты, проверяющие корректность работы игровых функций, с использованием пакета testing.
- [`handler.go`](https://github.com/notblinkyet/go-words/blob/master/game/handler.go): Хэндлеры для обработки входящих сообщений, выбор уровня сложности и проверка введённого слова. Используется github.com/go-telegram-bot-api/telegram-bot-api/v5.
- [`run.go`](https://github.com/notblinkyet/go-words/blob/master/game/run.go): Основной цикл обработки входящих сообщений от Telegram и передача их в handleUpdate. Также импортируются библиотеки для работы с ботом и github.com/go-sql-driver/mysql для подключения к MySQL.
- [`main.go`](https://github.com/notblinkyet/go-words/blob/master/cmd/main.go): Точка входа в приложение.

## Идеи для улучшения

- Добавление статистики. Реализовать хранение и отображение статистики игроков.
- Проверка существования слова. Добавить валидацию слов по базе данных, чтобы избежать ввода несуществующих слов.
- Расширение словаря. Добавить больше слов в базу данных, включая английские слова.

## Используемые библиотеки

- [Golang](https://golang.org/)
- [Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api/v5)
- [MySQL Driver](https://github.com/go-sql-driver/mysql)