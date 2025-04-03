package middleware

import (
	tele "gopkg.in/telebot.v3"
	myLogger "pilulia_bot/logger"
)

type Middleware struct {
	logger *myLogger.Logger
}

func NewMiddleware(logger *myLogger.Logger) *Middleware {
	return &Middleware{logger}
}

func (m *Middleware) LoggingMidleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		username := "unknown"
		if c.Sender().Username != "" {
			username = c.Sender().Username
		}
		m.logger.Info.Printf("%s | %s", username, c.Message().Text)
		return next(c)
	}
}
