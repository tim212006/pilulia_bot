package app

import (
	"database/sql"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"pilulia_bot/config"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
	"strconv"
)

type Handler struct {
	Bot Bot
	Lgr *logger.Logger
	DB  *MySQLUserDb
}

func NewHandler(bot Bot, lgr *logger.Logger, db *MySQLUserDb) *Handler {
	return &Handler{Bot: bot, Lgr: lgr, DB: db}
}

//////////////////////////////////////////////////////////Обработчики

func (h *Handler) HandleText(c telebot.Context) error {
	text := c.Text()
	return c.Send(text)
}

func (h *Handler) UpdateConnect(c telebot.Context) error {
	h.DB.UpdateLastConnect(c.Sender().ID)
	string := fmt.Sprintf("Обновление пользователя %d", c.Sender().ID)
	return c.Send(string)
}

// Проверяем наличие пользователя в БД, если нет, то добавляем в БД со статусом NewUser
func (h *Handler) HandleStart(c telebot.Context) error {
	user := c.Sender()
	_, err := h.DB.GetUserID(user)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				status := consts.NewUser
				h.DB.InsertUser(user.ID, user.FirstName, user.LastName, user.Username, status)
				//Приветствие новому пользователю
				return c.Send("Приветствую, я бот-таблетница! Я могу учитывать, употребляемые вами препараты и напоминать об их приеме. Давайте добавим препараты для приема и учета")
			}
		}
		h.Lgr.Err.Println(consts.DBErrorGetUser)
		return err
	}
	status := consts.Default
	h.DB.UpdateUserStatus(user.ID, status)
	//Приветствие старому пользователю
	helloString := fmt.Sprintf("С возвращением, %s %s!", user.FirstName, user.LastName)
	c.Send(helloString)
	//h.menuCommand(c)
	//return h.menuCommand(c)
	return h.SendMenuWithInlineButtons(c)
}

// Обработчик текстового сообщения для кнопки "Препараты"
func (h *Handler) handleShowUserDrugs(c telebot.Context) error {
	// Логика для обработки сообщения "Препараты"
	userDrugs, err := h.DB.GetUserDrugs(c.Sender().ID)
	if err != nil {
		return err
	}
	if len(userDrugs) == 0 {
		return c.Send("Нет доступных лекарств")
	}
	//Отправляем пользователю кнопки с названием препаратов
	return h.SendDynamicButtonMessage(c, userDrugs)
	//Элемент кода для отображения препаратов списком
	/*var message string
	for name := range userDrugs {
		//Экранируем специальные символы Markdown
		escapedName := config.EscapeMarkdown(name)
		message += fmt.Sprintf("- %s\n", escapedName)
	}
	return c.Send(fmt.Sprintf("Ваши препараты: \n%s", message), &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})*/
}

//Функция отправки пользователю информации о выбранном препарате
//func (h *Handlers)

// Обработчик текстового сообщения для кнопки "Помощь"
func (h *Handler) handleGetHelp(c telebot.Context) error {
	// Логика для обработки сообщения "Помощь"
	return c.Send("Здесь будет справочная информация...")
}

func (h *Handler) HandleWriteDownDrug(c telebot.Context) error {
	return c.Respond(&telebot.CallbackResponse{
		Text: "Текст сообщения",
	})
}

func (h *Handler) handleButtonPress(c telebot.Context) error {
	// Логируем идентификатор и данные callback
	callbackData := c.Callback().Data
	log.Printf("Получен callback: ID = %s, Data = %s", c.Callback().ID, callbackData)

	// Пример обработки callback
	return c.Respond(&telebot.CallbackResponse{
		Text: "Вы нажали на кнопку!",
	})
}

func (h *Handler) HandleDrugInfo(c telebot.Context) error {
	drugId := c.Callback().Data[6:]
	fmt.Println(c.Callback().Data)
	fmt.Println(drugId)
	drugIdInt, err := strconv.ParseInt(drugId, 10, 64)
	if err != nil {
		return err
	}
	drug, err := h.DB.GetDrug(drugIdInt)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("*Информация о препарате %s:*\n- Доза утром: %d\n- Доза днем: %d\n- Доза вечером: %d\n- Доза ночью: %d\n- Количество: %d\n- Комментарий: %s",
		drug.Drug_name, drug.M_dose, drug.A_dose, drug.E_dose, drug.N_dose, drug.Quantity, config.EscapeMarkdown(drug.Comment))

	// Отправка сообщения пользователю
	return c.Send(message, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
}
