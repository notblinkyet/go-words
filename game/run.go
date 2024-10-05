package game

import (
	"log"
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
