package app

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"gopkg.in/telebot.v3"
	"pilulia_bot/config"
	"pilulia_bot/drugs"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
	"time"
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

func (db *MySQLUserDb) UpdateLastConnect(userId int64) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := "UPDATE users SET last_connect = ? WHERE userid = ?"
	_, err := db.db.Exec(query, currentTime, userId)
	if err != nil {
		return fmt.Errorf("Ошибка обновления последнего подключения для пользователя %d: %w", userId, err)
	}
	return nil
}

func (db *MySQLUserDb) InsertUser(userID int64, firstName, lastName, userName, status string) error {
	query := "INSERT INTO users (userid, firstname, lastname, username, userstatus) VALUES (?, ?, ?, ?, ?)"
	_, err := db.db.Exec(query, userID, firstName, lastName, userName, status)
	db.UpdateLastConnect(userID)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLUserDb) UpdateUserStatus(userID int64, newStatus string) error {
	query := "UPDATE users SET userstatus = ? WHERE userid = ?"
	_, err := db.db.Exec(query, newStatus, userID)
	db.UpdateLastConnect(userID)
	if err != nil {
		return err
	}
	return nil
}

// Функции препаратов
func (db *MySQLUserDb) GetUserDrugs(userID int64) (map[string]drugs.Drugs, error) {
	query := `
		SELECT id, drug_name, m_dose, a_dose, e_dose, n_dose, quantity, comment 
		FROM drugs 
		WHERE userid = ?`
	rows, err := db.db.Queryx(query, userID)
	if err != nil {
		return nil, fmt.Errorf("Ошибка получения препаратов пользователя %d: %w", userID, err)
	}
	defer rows.Close()
	drugsMap := make(map[string]drugs.Drugs)
	for rows.Next() {
		var drug drugs.Drugs
		err = rows.StructScan(&drug)
		if err != nil {
			return nil, fmt.Errorf("Ошибка получения строки препарата %d: %w", userID, err)
		}
		drugsMap[drug.Drug_name] = drug
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return drugsMap, nil
}

func (db *MySQLUserDb) GetDrug(drugId int64) (drugs.Drugs, error) {
	var drug drugs.Drugs
	query := `
		SELECT id, drug_name, m_dose, a_dose, e_dose, n_dose, quantity, comment 
		FROM drugs 
		WHERE id = ?`
	err := db.db.Get(&drug, query, drugId)
	if err != nil {
		return drugs.Drugs{}, fmt.Errorf("Ошибка получения препарата %d: %w", drugId, err)
	}
	return drug, nil
}
