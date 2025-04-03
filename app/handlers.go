package app

import (
	"database/sql"
	"gopkg.in/telebot.v3"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
)

type Handler struct {
	Bot Bot
	Lgr *logger.Logger
	DB  *MySQLUserDb
}

func NewHandler(bot Bot, lgr *logger.Logger, db *MySQLUserDb) *Handler {
	return &Handler{Bot: bot, Lgr: lgr, DB: db}
}

func (h *Handler) HandleText(c telebot.Context) error {
	text := c.Text()
	return c.Send(text)
}

func (h *Handler) HandleStart(c telebot.Context) error {
	return c.Send("Приветствую, я бот-таблетница! Я могу учитывать, употребляемые вами препараты и напоминать об их приеме")
}

func (h *Handler) Exist(c telebot.Context) error {
	user := c.Sender()
	userId, err := h.DB.GetUserID(user)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				h.Lgr.Info.Println(consts.DBUserNotExist)
				return c.Send(consts.DBUserNotExist)
			}
		}
		h.Lgr.Err.Println(consts.DBErrorGetUser)
		return err
	}
	h.Lgr.Info.Println(consts.DBUserExist, userId)
	c.Send(userId)
	return nil
}
