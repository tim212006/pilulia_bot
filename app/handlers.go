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

////////////////////////////////////////////////////////////Миддлварь

func (h *Handler) statusMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userId := c.Sender().ID

		statusData, err := h.DB.GetUserStatus(userId)
		if err != nil {
			return err
		}
		status := fmt.Sprintf("%.3s", statusData)
		if status != "Add" {
			return c.Send("Для использования этой функции начните процеду добавления перпарата сначала")
		}
		return next(c)
	}
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
	//return h.SendMenuWithInlineButtons(c)
	return h.SendMainMenu(c)
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
	inlineKeyboard := telebot.ReplyMarkup{}
	uniqueString := fmt.Sprintf("edit_%d", drug.Id)
	//fmt.Println(uniqueString)
	textString := fmt.Sprintf("Редактировать")
	btn := telebot.InlineButton{
		Unique: uniqueString,
		Text:   textString,
	}
	uniqueStringDel := fmt.Sprintf("delete_%d", drug.Id)
	btnDelete := telebot.InlineButton{
		Unique: uniqueStringDel,
		Text:   "Удалить",
	}
	//btnEdit := h.SendEditDrug(c, drugIdInt)

	inlineKeyboard.InlineKeyboard = [][]telebot.InlineButton{
		{btn, btnDelete},
	}

	// Отправка сообщения пользователю
	return c.Send(message, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: &inlineKeyboard,
	})
}

func (h *Handler) HandleDrugEdit(c telebot.Context) error {
	drugId := c.Callback().Data[6:]
	fmt.Println(drugId)
	drugIdInt, err := strconv.ParseInt(drugId, 10, 64)
	if err != nil {
		return err
	}
	fmt.Println(drugIdInt)
	drug, err := h.DB.GetDrug(drugIdInt)
	if err != nil {
		return err
	}
	RespString := fmt.Sprintf("Нажата кнопка - Редактировать препарат %s", drug.Drug_name)
	return c.Send(RespString)
}

func (h *Handler) HandleDrugDelete(c telebot.Context) error {
	drugId := c.Callback().Data[8:]
	fmt.Println(drugId)
	drugIdInt, err := strconv.ParseInt(drugId, 10, 64)
	if err != nil {
		return err
	}
	fmt.Println(drugIdInt)
	drug, err := h.DB.GetDrug(drugIdInt)
	if err != nil {
		return err
	}
	return h.SendDeleteConfirmation(c, drug.Id, drug.Drug_name)
}

/////////////////

// Обработчик для нажатия на кнопку "На главную"
func (h *Handler) handleMain(c telebot.Context) error {
	return c.Send("Добро пожаловать на главную страницу!")
}

// Обработчик для нажатия на кнопку "Препараты"
func (h *Handler) handleDrugs(c telebot.Context) error {
	return c.Send("Здесь будет информация о препаратах.")
}

// Обработчик для нажатия на кнопку "Помощь"
func (h *Handler) handleHelp(c telebot.Context) error {
	return c.Send("Здесь будет справочная информация.")
}

func (h *Handler) SendDeleteConfirmation(c telebot.Context, drugId int64, drugName string) error {
	confirmBtn := telebot.InlineButton{
		Unique: fmt.Sprintf("confirm_delete_%d", drugId),
		Text:   "Да, удалить",
	}
	cancelBtn := telebot.InlineButton{
		Unique: "c_del",
		Text:   "Отмена",
	}
	inlineKeyboard := &telebot.ReplyMarkup{}
	inlineKeyboard.InlineKeyboard = [][]telebot.InlineButton{
		{confirmBtn, cancelBtn}, // Кнопки в одной строке
	}
	respString := fmt.Sprintf("Вы уверены, что хотите удалить %s?", drugName)
	return c.Send(respString, inlineKeyboard)
}

func (h *Handler) AcceptedDeleteDrug(c telebot.Context) error {
	drugId := c.Callback().Data[16:]
	fmt.Println(drugId)
	drugIdInt, err := strconv.ParseInt(drugId, 10, 64)
	if err != nil {
		return err
	}
	errDel := h.DB.DeleteDrug(drugIdInt)
	if errDel != nil {
		return errDel
	}
	c.Send("Препарат удален из базы данных")
	return h.handleShowUserDrugs(c)
}

func (h *Handler) CancelDeleteDrug(c telebot.Context) error {
	c.Send("Удаление препарата отменено")
	return h.handleShowUserDrugs(c)
}

func (h *Handler) handleAddDrug(c telebot.Context) error {

	h.DB.UpdateUserStatus(c.Sender().ID, consts.AddDrugName)
	c.Send("Введите название препарата:")
	//var drug drugs.Drugs
	/*drug.Drug_name = "Андрогель 5мг"
	drug.M_dose = 1
	drug.A_dose = 0
	drug.E_dose = 0
	drug.N_dose = 0
	drug.Quantity = 50
	drug.Comment = ""
	err := h.DB.InsertDrug(c.Sender().ID, drug)
	if err != nil {
		return err
	}*/

	//return h.handleShowUserDrugs(c)
	return nil
}

func (h *Handler) SwitchStatus(c telebot.Context) error {
	// Извлечение статуса пользователя
	userId := c.Sender().ID
	userStatus, err := h.DB.GetUserStatus(userId)
	if err != nil {
		return c.Send("Ошибка получения статуса пользователя.")
	}
	// Проверка статуса пользователя и отправка соответствующего сообщения
	switch userStatus {
	case consts.AddDrugName:
		//22.43 - Необходимо создать структуру в которую будет отправляться текст от пользователя
		return c.Send("Ввод названия препарата")
	case consts.AddMorningDose:
		return c.Send("Ввод количества препарата утром")
	case consts.AddAfternoonDose:
		return c.Send("Ввод количества препарата днем")
	case consts.AddEvningDose:
		return c.Send("Ввод количества препарата вечером")
	case consts.AddNightDose:
		return c.Send("Ввод количества препарата ночью")
	case consts.AddDrugQuantity:
		return c.Send("Ввод количества оставшегося препарата")
	case consts.AddDrugComment:
		return c.Send("Ввод комментария к препарату")
	default:
		return c.Send("Неизвестный тип состояния, пожалуйста, повторите попытку.")
	}
}
