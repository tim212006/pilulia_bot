package app

import (
	"errors"
	"gopkg.in/telebot.v3"
	"pilulia_bot/config"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
	"pilulia_bot/logger/middleware"
	"pilulia_bot/users"
	"time"
)

type Bot struct {
	Self       *telebot.Bot
	Middleware *middleware.Middleware
	User       map[int64]users.User
	Cfg        *config.Config
	Lgr        *logger.Logger
}

func NewBot(cfg *config.Config, lgr *logger.Logger) (*Bot, error) {
	pref := telebot.Settings{
		Token:     cfg.BotSettings.TelegramToken,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeMarkdown,
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, errors.New(consts.ErrorBot)
	}
	middleware := middleware.NewMiddleware(lgr)
	user := make(map[int64]users.User)

	return &Bot{
		Self:       bot,
		Middleware: middleware,
		User:       user,
		Cfg:        cfg,
		Lgr:        lgr,
	}, nil
}

func (bot *Bot) Start() {
	bot.Self.Start()
}
