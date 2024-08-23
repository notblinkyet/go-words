# Go Words

Это мой проект, который я написал для изучения golang.

Go words - это телеграмм бот, в котором реализована игра на подобии Wordly.

Когда игрок начинает игру, из базы данных достается случайное слово из 5 букв, которое нужно угадать. Всего есть 5 попыток, а так же при каждой из попыток ввыводится результат на правильном ли месте конкретная буква, есть ли это буква в слове или буква вовсе осутствует.

Для начала я написал логику [проверки](https://github.com/notblinkyet/go-words/blob/master/game.go) слов на корректность: состоит ли оно из русских букв и подходит по длине. После написал проверку заданного коректно слова, какие буквы на своем месте, какие не на своем.

Потом написал [тестики](https://github.com/notblinkyet/go-words/blob/master/game_test.go) для этих проверок.

Написал [несложного бота](https://github.com/notblinkyet/go-words/blob/master/run.go) добавив в него проверку слова.

Потом решил увеличить словарный запас бота и нашел MySql запросы для создание таблицы русских слов. Создал базу данных в ней таблицу, удалил не подходящие по длине слова и написал запрос к ней в Golang.

Так же есть идеи как улучшить проект, например добавить уровни сложность и/или добавить английский язык и английские слова и много других.