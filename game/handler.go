package game

import (
	"database/sql"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var keyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("hard"),
		tgbotapi.NewKeyboardButton("meduim"),
		tgbotapi.NewKeyboardButton("easy"),
	))

func handlePlay(update tgbotapi.Update, game *GameState) {
	game.Mod = 1
	BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери уровень сложности:")
	BotMes.ReplyMarkup = keyBoard
	db, err := sql.Open("mysql", Dsn)
	if err != nil {
		fmt.Println("\n\n\nWe have some problems.")
		panic(err)
	}
	go make_word(db, words)
}

func handleDiffucult(update tgbotapi.Update, diff string, game *GameState) {
	switch diff {
	case "hard":
		game.Mod = 2
		game.AttemptsLeft = 2
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбран hard уровень сложности, у вас есть 3 попытки угадать слово.")
		BotMes.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		get_word(game, words)
	case "meduim":
		game.Mod = 2
		game.AttemptsLeft = 3
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбран mediuk уровень сложности, у вас есть 4 попытки угадать слово.")
		BotMes.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		get_word(game, words)
	case "easy":
		game.Mod = 2
		game.AttemptsLeft = 4
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбран easy уровень сложности, у вас есть 5 попытки угадать слово.")
		BotMes.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		get_word(game, words)
	default:
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Не верно введен уровень сложности.\nВозможные варианты:\neasy/hard/medium.")
	}
}

func handleCheck(update tgbotapi.Update, mes string, game *GameState, chatID int64) {

	if IsValid(mes) {

		CurRes := Check(mes, game.CorrectWord, game.CorrectWordD)

		if Success(CurRes) {
			BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Ура! Вы угадали слово. Приходите поиграть ещё!")
			endGame(chatID)
		} else if game.AttemptsLeft == 0 {
			BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Игра окончена.Это было слово:\n"+game.CorrectWord+"\nНачните новую, отправив /play.")
			endGame(chatID)
		} else {
			game.AttemptsLeft--
			BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, CurRes)
		}
	} else {
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный ввод.")
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	gamesMutex.Lock()
	game, exists := games[chatID]
	gamesMutex.Unlock()

	if !exists {
		game = &GameState{}
		gamesMutex.Lock()
		defer gamesMutex.Unlock()
		games[chatID] = game
	}

	if update.Message.Text == "/start" {
		username := update.Message.From.UserName
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s. Это игра Go Word. Цель игры угадать слово которое я загадал. Это русское слово, состоящее из 5 букв. Хочешь сыграть? Отправь мне /play.", username))
	} else if update.Message.Text == "/play" {
		handlePlay(update, game)
	} else {
		if game.Mod == 1 {
			handleDiffucult(update, update.Message.Text, game)
		} else if game.Mod == 2 {
			handleCheck(update, strings.TrimSpace(update.Message.Text), game, chatID)
		} else {
			BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Для того что бы ознакомится с правилами отправьте /start, что бы играть отправьте /play")
		}
	}
	bot.Send(BotMes)
}

func endGame(chatID int64) {
	gamesMutex.Lock()
	defer gamesMutex.Unlock()
	delete(games, chatID)
}
