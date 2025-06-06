package app

import (
	"errors"
	"fmt"
	"gopkg.in/telebot.v3"
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
		return nil, errors.New(consts.ErrorBot)
	}
	existDb := СheckDatabaseAndTables(cfg, lgr)
	if existDb != nil {
		return nil, existDb
	}
	db, dbErr := NewMySQLUserDb(cfg, lgr)
	if dbErr != nil {
		return nil, errors.New(consts.ErrorDBConnect)
	}
	handler := NewHandler(bot, lgr, db)
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
	app.Bot.Self.Handle(&telebot.InlineButton{Unique: "showUserDrugs"}, app.Handler.handleShowUserDrugs)
	app.Bot.Self.Handle(&telebot.InlineButton{Unique: "getHelp"}, app.Handler.handleGetHelp)
	app.Bot.Self.Handle(telebot.OnCallback, func(c telebot.Context) error {
		callbackData := c.Callback().Data
		fmt.Println(callbackData[1:6])
		if len(callbackData) < 6 {
			return fmt.Errorf("некорректные данные callback: %s", callbackData)
		}
		switch callbackData[1:6] {
		case "edit_":
			return app.Handler.HandleDrugEdit(c)
		case "drug_":
			return app.Handler.HandleDrugInfo(c)
		case "add_d":
			return app.Handler.handleAddDrug(c)
		case "delet":
			return app.Handler.HandleDrugDelete(c)
		case "confi":
			return app.Handler.AcceptedDeleteDrug(c)
		case "c_del":
			return app.Handler.CancelDeleteDrug(c)
		case "d_sav":
			return app.Handler.SaveDrug(c)
		case "d_can":
			return app.Handler.EraseDrug(c)
		case "daily":
			return app.Handler.HandlePressDailyButton(c)
		case "dredi":
			return app.Handler.HandleDrugParametrEdit(c)
		case "cedit":
			return app.Handler.handleShowUserDrugs(c)
		default:
			return fmt.Errorf("неизвестный тип callback: %s", callbackData)
		}
	})

	app.Bot.Self.Handle(telebot.OnText, func(c telebot.Context) error {
		switch c.Text() {
		case "На главную":
			return app.Handler.HandleStart(c)
		case "Препараты":
			return app.Handler.handleShowUserDrugs(c)
		case "Помощь":
			return app.Handler.handleHelp(c)
		//case "Помощь": return app.Handler.handleShowDailyUserDrugs(c)
		default:
			return app.Handler.SwitchStatus(c)
		}
	})

	//
	/*app.Bot.Self.Handle(telebot.OnCallback, func(c telebot.Context) error {

		return nil
	})*/

	//
	//})

	app.Bot.Self.Handle("/start", app.Handler.HandleStart)
	app.Bot.Self.Handle("/exist", app.Handler.UpdateConnect)
	app.Bot.Start()

}
