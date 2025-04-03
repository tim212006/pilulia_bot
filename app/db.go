package app

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/telebot.v3"
	"pilulia_bot/config"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
)

type MySQLUserDb struct {
	db *sqlx.DB
}

func NewMySQLUserDb(cfg *config.Config, lgr *logger.Logger) (*MySQLUserDb, error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBSettings.Username,
		cfg.DBSettings.Password,
		cfg.DBSettings.Host,
		cfg.DBSettings.Port,
		cfg.DBSettings.Database)
	//lgr.Info.Println(connection)
	db, err := sqlx.Open("mysql", connection)
	if err != nil {
		lgr.Err.Println(consts.ErrorDBConnectSqlx, err)
		return nil, errors.New(consts.ErrorDBConnectSqlx)
	}
	err = db.Ping()
	if err != nil {
		lgr.Err.Println(consts.ErrorDBPing, err)
		return nil, errors.New(consts.ErrorDBPing)
	}
	lgr.Info.Println(consts.InfoDBConnected)
	return &MySQLUserDb{db: db}, err
}

func (db *MySQLUserDb) GetUserID(user *telebot.User) (int64, error) {
	var userId int64
	err := db.db.QueryRow("SELECT userid FROM users WHERE userid = ?", user.ID).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
