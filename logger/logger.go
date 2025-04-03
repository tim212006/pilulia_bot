package logger

import (
	"log"
	"os"
)

// Объявляем структуру логера
type Logger struct {
	Err  *log.Logger
	Info *log.Logger
}

// Инициируем логера
func InitLogger() *Logger {
	Logger := Logger{}
	Logger.Info = log.New(os.Stdout, "[ИНФО]      ", log.Ldate|log.Ltime)
	Logger.Err = log.New(os.Stderr, "[ОШИБКА]      ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Info.Print("Инициализация логирования.")
	return &Logger
}
