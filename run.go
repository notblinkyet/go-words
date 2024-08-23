package game

import (
	"database/sql"
	"log"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GameState struct {
	CorrectWord  string
	CorrectWordD map[rune]int
	AttemptsLeft int
	Playing      bool
}

var (
	games      = make(map[int64]*GameState)
	gamesMutex sync.Mutex
)

func Run() {
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	gamesMutex.Lock()
	game, exists := games[chatID]
	gamesMutex.Unlock()

	if !exists {
		game = &GameState{AttemptsLeft: 4}
		gamesMutex.Lock()
		games[chatID] = game
		gamesMutex.Unlock()
	}

	switch update.Message.Text {
	case "/start":
		BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, это игра Go Word. Цель игры угадать слово которое я загадал. Это русское слово, состоящее из 5 букв. Хочешь сыграть? Отправь мне /play.")
		bot.Send(BotMes)

	case "/play":
		db, err := sql.Open("mysql", Dsn)

		if err != nil {
			panic(err)
		}
		defer db.Close()

		query := "SELECT word FROM nouns ORDER BY rand() LIMIT 1"

		err = db.QueryRow(query).Scan(&game.CorrectWord)

		if err != nil {
			panic(err)
		}

		game.CorrectWordD = CreateDict(game.CorrectWord)
		BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, "Супер, я загадал слово! У тебя есть 5 попыток.")
		bot.Send(BotMes)
		game.Playing = true

	default:
		if !game.Playing {
			BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, "Для того что бы ознакомится с правилами отправьте /start, что бы играть отправьте /play")
			bot.Send(BotMes)
		} else if game.AttemptsLeft > 0 {
			NewMes := strings.TrimSpace(update.Message.Text)

			if IsValid(NewMes) {
				CurRes := Check(NewMes, game.CorrectWord, game.CorrectWordD)
				if Success(CurRes) {
					BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, "Ура! Вы угадали слово. Приходите поиграть ещё!")
					bot.Send(BotMes)
					endGame(chatID)
				} else {
					game.AttemptsLeft--
					BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, CurRes)
					bot.Send(BotMes)
				}
			} else {
				game.AttemptsLeft--
				BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный ввод.")
				bot.Send(BotMes)
			}
		} else {
			BotMes := tgbotapi.NewMessage(update.Message.Chat.ID, "Игра окончена.Это было слово:\n"+game.CorrectWord+"\nНачните новую, отправив /play.")
			bot.Send(BotMes)
			endGame(chatID)
		}
	}
}

func endGame(chatID int64) {
	gamesMutex.Lock()
	delete(games, chatID)
	gamesMutex.Unlock()
}
