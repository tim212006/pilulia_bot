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
	//app.Bot.Self.Handle(&telebot.InlineButton{Unique: "write_down_drugs"}, app.Handler.HandleWriteDownDrug)
	//app.Bot.Self.Handle("/btn", app.Handler.HandleWriteDownDrug)
	//app.Bot.Self.Handle(&telebot.InlineButton{Unique: "showUserDrugs"}, app.Handler.HandleButtonPress)
	//app.Bot.Self.Handle(&telebot.InlineButton{Unique: "getHelp"}, app.Handler.HandleButtonPress)
	/*app.Bot.Self.Handle(telebot.OnText, func(c telebot.Context) error {
		switch c.Text() {
		case "Препараты":
			return app.Handler.handleShowUserDrugs(c)
		case "Помощь":
			return app.Handler.handleGetHelp(c)
		default:
			return nil
		}
	})*/
	//
	//app.Bot.Self.Handle(telebot.OnCallback, app.Handler.handleButtonPress)
	//
	app.Bot.Self.Handle(&telebot.InlineButton{Unique: "showUserDrugs"}, app.Handler.handleShowUserDrugs)
	app.Bot.Self.Handle(&telebot.InlineButton{Unique: "getHelp"}, app.Handler.handleGetHelp)
	//

	app.Bot.Self.Handle(telebot.OnCallback, func(c telebot.Context) error {
		callbackData := c.Callback().Data
		fmt.Println(callbackData)
		if len(callbackData) < 6 {
			return fmt.Errorf("Некорректные данные callback: %s", callbackData)
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
		default:
			return fmt.Errorf("Неизвестный тип callback: %s", callbackData)
		}
		/*fmt.Println(c.Callback().Data)

		if c.Callback().Data[1:6] == "edit_" {
			return app.Handler.HandleDrugEdit(c)
		}
		if c.Callback().Data[1:6] == "drug_" {
			return app.Handler.HandleDrugInfo(c)
		}
		return nil*/
	})

	app.Bot.Self.Handle(telebot.OnText, func(c telebot.Context) error {
		switch c.Text() {
		case "На главную":
			return app.Handler.HandleStart(c)
		case "Препараты":
			return app.Handler.handleShowUserDrugs(c)
		case "Помощь":
			return app.Handler.handleHelp(c)
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
