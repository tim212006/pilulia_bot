package app

import (
	"errors"
	"pilulia_bot/config"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
)

type App struct {
	Bot     *Bot
	Config  *config.Config
	Logger  *logger.Logger
	Handler *Handler
	DB      *MySQLUserDb
}

func NewApp() (*App, error) {
	cfg, err := config.ConfigInit()
	if err != nil {
		return nil, err
	}
	lgr := logger.InitLogger()
	bot, err := NewBot(cfg, lgr)
	if err != nil {
		errors.New(consts.ErrorBot)
		return nil, err
	}
	db, dbErr := NewMySQLUserDb(cfg, lgr)
	if dbErr != nil {
		return nil, errors.New(consts.ErrorDBConnect)
	}
	handler := NewHandler(*bot, lgr, db)
	app := &App{
		Bot:     bot,
		Config:  cfg,
		Logger:  lgr,
		Handler: handler,
		DB:      db,
	}

	return app, nil
}

func (app *App) Start() {
	app.Logger.Info.Println(consts.InfoAppStart)
	app.Bot.Self.Handle("/start", app.Handler.HandleStart)
	app.Bot.Self.Handle("/exist", app.Handler.Exist)
	app.Bot.Start()
}
