package game

import (
	"database/sql"
	"fmt"
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
	Mod          int
}

var (
	games      = make(map[int64]*GameState)
	gamesMutex sync.Mutex
	BotMes     tgbotapi.MessageConfig
	words      = make(chan string, 9)
)

var keyBoard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("hard"),
		tgbotapi.NewKeyboardButton("meduim"),
		tgbotapi.NewKeyboardButton("easy"),
	))

func make_word(db *sql.DB, c chan<- string) {
	var cur string
	query := "SELECT word FROM nouns ORDER BY rand() LIMIT 1"
	err := db.QueryRow(query).Scan(&cur)
	if err != nil {
		panic(err)
	}
	c <- cur
}

func get_word(g *GameState, c <-chan string) {
	g.CorrectWord = <-c
	g.CorrectWordD = CreateDict(g.CorrectWord)
}

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
		game = &GameState{}
		gamesMutex.Lock()
		defer gamesMutex.Unlock()
		games[chatID] = game
	}

	if update.Message.Text == "/start" {
		username := update.Message.From.UserName
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Привет, %s. Это игра Go Word. Цель игры угадать слово которое я загадал. Это русское слово, состоящее из 5 букв. Хочешь сыграть? Отправь мне /play.", username))
	} else if update.Message.Text == "/play" {
		game.Mod = 1
		BotMes = tgbotapi.NewMessage(update.Message.Chat.ID, "Выбери уровень сложности:")
		BotMes.ReplyMarkup = keyBoard
		db, err := sql.Open("mysql", Dsn)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		go make_word(db, words)
	} else {
		if game.Mod == 1 {
			switch update.Message.Text {
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
		} else if game.Mod == 2 {

			NewMes := strings.TrimSpace(update.Message.Text)

			if IsValid(NewMes) {

				CurRes := Check(NewMes, game.CorrectWord, game.CorrectWordD)
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
