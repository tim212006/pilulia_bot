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
	"pilulia_bot/users"
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
		return fmt.Errorf("ошибка обновления последнего подключения для пользователя %d: %w", userId, err)
	}
	return nil
}

func (db *MySQLUserDb) InsertUser(userID int64, firstName, lastName, userName, status string) error {
	query := "INSERT INTO users (userid, firstname, lastname, username, userstatus) VALUES (?, ?, ?, ?, ?)"
	_, err := db.db.Exec(query, userID, firstName, lastName, userName, status)
	if err != nil {
		return err
	}
	return db.UpdateLastConnect(userID)
}

func (db *MySQLUserDb) UpdateUserStatus(userID int64, newStatus string) error {
	query := "UPDATE users SET userstatus = ? WHERE userid = ?"
	_, err := db.db.Exec(query, newStatus, userID)
	if err != nil {
		return err
	}
	return db.UpdateLastConnect(userID)
}

func (db *MySQLUserDb) GetUserStatus(userId int64) (users.Status, error) {
	var userStatus users.Status
	query := "SELECT userstatus FROM users WHERE userid = ?"
	err := db.db.Get(&userStatus, query, userId)
	if err != nil {
		return userStatus, fmt.Errorf("ошибка получения статуса пользователя")
	}
	return userStatus, nil
}

// GetUserDrugs Функции препаратов
func (db *MySQLUserDb) GetUserDrugs(userID int64) (map[string]drugs.Drugs, error) {
	query := `
		SELECT id, drug_name, m_dose, a_dose, e_dose, n_dose, quantity, comment 
		FROM drugs 
		WHERE userid = ?`
	rows, err := db.db.Queryx(query, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения препаратов пользователя %d: %w", userID, err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = fmt.Errorf("ошибка закрытия строк: %w", closeErr)
		}
	}()
	drugsMap := make(map[string]drugs.Drugs)
	for rows.Next() {
		var drug drugs.Drugs
		err = rows.StructScan(&drug)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения строки препарата %d: %w", userID, err)
		}
		drugsMap[drug.Drug_name] = drug
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return drugsMap, nil
}

func (db *MySQLUserDb) GetPeriodDrugs(userID int64, period string) (map[string]drugs.Drugs, error) {
	if period != "m_dose" && period != "a_dose" && period != "e_dose" && period != "n_dose" {
		return nil, fmt.Errorf("период: не соответствует")
	}
	query := `
		SELECT id, drug_name, m_dose, a_dose, e_dose, n_dose, quantity, comment 
		FROM drugs 
		WHERE userid = ? AND ? > 0`
	rows, err := db.db.Queryx(query, userID, period)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения препаратов периода пользователя %d: %w", userID, err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = fmt.Errorf("ошибка закрытия строк: %w", closeErr)
		}
	}()
	drugsMap := make(map[string]drugs.Drugs)
	for rows.Next() {
		var drug drugs.Drugs
		err = rows.StructScan(&drug)
		if err != nil {
			return nil, fmt.Errorf("ошибка получения строки (период) препарата %d: %w", userID, err)
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
		return drugs.Drugs{}, fmt.Errorf("ошибка получения препарата %d: %w", drugId, err)
	}
	return drug, nil
}

func (db *MySQLUserDb) DeleteDrug(drugId int64) error {
	query := "DELETE FROM drugs WHERE id = ?"
	result, err := db.db.Exec(query, drugId)
	if err != nil {
		return fmt.Errorf("ошибка удаления перпарата с ID %d: %s", drugId, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("препарат с ID %d не найден", drugId)
	}
	return nil
}

func (db *MySQLUserDb) InsertDrug(userId int64, drug drugs.Drugs) error {
	query := "INSERT INTO drugs (userid, drug_name, m_dose, a_dose, e_dose, n_dose, quantity, comment) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.db.Exec(query, userId, drug.Drug_name, drug.M_dose, drug.A_dose, drug.E_dose, drug.N_dose, drug.Quantity, drug.Comment)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLUserDb) UpdateDrugName(drugId int64, drugName string) error {
	query := "UPDATE drugs SET drug_name = ? WHERE id = ?"
	_, err := db.db.Queryx(query, drugName, drugId)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLUserDb) UpdateDrugDose(drugId int64, dose int64, period string) error {
	var query string
	switch period {
	case "m_dose":
		query = "UPDATE drugs SET m_dose = ? WHERE id = ?"
	case "a_dose":
		query = "UPDATE drugs SET a_dose = ? WHERE id = ?"
	case "e_dose":
		query = "UPDATE drugs SET e_dose = ? WHERE id = ?"
	case "n_dose":
		query = "UPDATE drugs SET n_dose = ? WHERE id = ?"
	default:
		fmt.Println("Неизвестный период")
	}
	_, err := db.db.Exec(query, dose, drugId)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLUserDb) UpdateDrugQuantity(drugId int64, quantity int64) error {
	query := "UPDATE drugs SET quantity = ? WHERE id = ?"
	_, err := db.db.Exec(query, quantity, drugId)
	if err != nil {
		return err
	}
	return nil
}

func (db *MySQLUserDb) UpdateDrugComment(drugId int64, drugComment string) error {
	query := "UPDATE drugs SET comment = ? WHERE id = ?"
	_, err := db.db.Queryx(query, drugComment, drugId)
	if err != nil {
		return err
	}
	return nil
}
