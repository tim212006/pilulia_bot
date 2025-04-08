package app

import (
	"database/sql"
	"errors"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"pilulia_bot/config"
	"pilulia_bot/logger"
	"pilulia_bot/logger/consts"
	"pilulia_bot/users"
	"strconv"
	"unicode"
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

// Обновление статуса пользователя
func (h *Handler) UpdateUserStatus(userID int64, status users.Status) {
	h.Bot.UserMutex.Lock()
	defer h.Bot.UserMutex.Unlock()

	if u, exists := h.Bot.User[userID]; exists {
		u.Status = status
		h.Bot.User[userID] = u
	}
}

func (h *Handler) GetUserStatus(userID int64) (users.Status, error) {
	if u, exists := h.Bot.User[userID]; exists {
		status := u.Status
		return status, nil
	} else {
		fmt.Println("Информация о пользователе отсутствует в Map")
		status, err := h.DB.GetUserStatus(userID)
		return status, err
	}
}

func (h *Handler) SetDrugString(userID int64, status users.Status, data string) {
	h.Bot.UserMutex.Lock()
	defer h.Bot.UserMutex.Unlock()
	user, exists := h.Bot.User[userID]
	if !exists {
		fmt.Println("Пользователь не найден")
		return
	}
	switch status {
	case consts.AddDrugName:
		user.Drugs.Drug_name = data
	case consts.AddDrugComment:
		user.Drugs.Comment = data
	default:

		fmt.Println("Статус пользователя: ", status, " не соответствует функции")
		return
	}
	h.Bot.User[userID] = user
}

func (h *Handler) SetDrugInt(userID int64, status users.Status, data int64) {
	h.Bot.UserMutex.Lock()
	defer h.Bot.UserMutex.Unlock()
	user, exists := h.Bot.User[userID]
	if !exists {
		fmt.Println("Пользователь не найден")
		return
	}
	switch status {
	case consts.AddMorningDose:
		user.Drugs.M_dose = data
	case consts.AddAfternoonDose:
		user.Drugs.A_dose = data
	case consts.AddEvningDose:
		user.Drugs.E_dose = data
	case consts.AddNightDose:
		user.Drugs.N_dose = data
	case consts.AddDrugQuantity:
		user.Drugs.Quantity = data
	default:
		fmt.Println("Статус пользователя: ", status, " не соответствует функции")
	}
	h.Bot.User[userID] = user
}

// Проверяем наличие пользователя в БД, если нет, то добавляем в БД со статусом NewUser
func (h *Handler) HandleStart(c telebot.Context) error {
	user := c.Sender()
	h.Bot.UserMutex.Lock()

	if u, exists := h.Bot.User[user.ID]; !exists {
		h.Bot.User[user.ID] = users.User{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
		}
	} else {
		u.Status = consts.Default
		h.Bot.User[user.ID] = u
	}
	h.Bot.UserMutex.Unlock()
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
	h.DB.UpdateUserStatus(user.ID, consts.Default)
	h.UpdateUserStatus(user.ID, consts.Default)
	//Приветствие старому пользователю
	helloString := fmt.Sprintf("С возвращением, %s %s!", user.FirstName, user.LastName)
	c.Send(helloString)
	//h.menuCommand(c)
	//return h.menuCommand(c)
	//return h.SendMenuWithInlineButtons(c)

	fmt.Println(h.Bot.User[user.ID].Username, ": ", h.Bot.User[user.ID].Status)
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
		inlineKeyboard := &telebot.ReplyMarkup{}
		btn_add_drug := telebot.InlineButton{
			Unique: "add_d",
			Text:   "Добавить препарат",
		}
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn_add_drug})
		return c.Send("Нет доступных лекарств", inlineKeyboard)
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

	//fmt.Println(c.Sender().ID)
	//u.EraseDrug() //Очищаем структуру перед вводом препарата
	h.UpdateUserStatus(c.Sender().ID, consts.AddDrugName)
	fmt.Println(h.Bot.User[c.Sender().ID].Username, ": ", h.Bot.User[c.Sender().ID].Status)
	h.DB.UpdateUserStatus(c.Sender().ID, consts.AddDrugName)
	c.Send("Введите название препарата:")
	return nil
}

func (h *Handler) SwitchStatus(c telebot.Context) error {
	// Извлечение статуса пользователя
	//var userStatus users.Status
	//var err error
	/*u := h.Bot.User[c.Sender().ID]
	if u.Status == "" {
		userStatus, err = h.DB.GetUserStatus(c.Sender().ID)
		h.Bot.UserMutex.Lock()
		u.SetStatus(userStatus)
		h.Bot.UserMutex.Unlock()
		if err != nil {
			return c.Send("Ошибка получения статуса пользователя.")
		}
	} else {
		userStatus = u.Status
	}

	//fmt.Println(userStatus)
	if err != nil {
		return c.Send("Ошибка получения статуса пользователя.")
	}*/
	userStatus, err := h.GetUserStatus(c.Sender().ID)
	if err != nil {
		return err
	}
	//fmt.Println(u.Username, ": ", u.Status)
	// Проверка статуса пользователя и отправка соответствующего сообщения
	switch userStatus {
	case consts.AddDrugName:
		//fmt.Println(c.Text())
		h.SetDrugString(c.Sender().ID, userStatus, c.Text())
		//fmt.Println("Информация в структуре", u.Drugs.Drug_name)
		h.UpdateUserStatus(c.Sender().ID, consts.AddMorningDose)
		h.DB.UpdateUserStatus(c.Sender().ID, consts.AddMorningDose)
		//fmt.Println(h.Bot.User[c.Sender().ID].Username, ": ", h.Bot.User[c.Sender().ID].Status)
		fmt.Println("Название: ", h.Bot.User[c.Sender().ID].Drugs.Drug_name)
		return c.Send(" Введите количество препарата утром:")
	case consts.AddMorningDose:
		//fmt.Println(c.Text())
		dataInt, err := ParseInt64FromString(c.Text())
		if err != nil {
			return c.Send("Укажите число")
		}
		h.SetDrugInt(c.Sender().ID, userStatus, dataInt)
		//fmt.Println("Информация в структуре", u.Drugs.M_dose)
		h.UpdateUserStatus(c.Sender().ID, consts.AddAfternoonDose)
		h.DB.UpdateUserStatus(c.Sender().ID, consts.AddAfternoonDose)
		fmt.Println("Название: ", h.Bot.User[c.Sender().ID].Drugs.Drug_name, " У: ", h.Bot.User[c.Sender().ID].Drugs.M_dose)
		return c.Send("Введите количество препарата днем:")
	case consts.AddAfternoonDose:
		dataInt, err := ParseInt64FromString(c.Text())
		if err != nil {
			return c.Send("Укажите число")
		}
		h.SetDrugInt(c.Sender().ID, userStatus, dataInt)
		h.UpdateUserStatus(c.Sender().ID, consts.AddEvningDose)
		h.DB.UpdateUserStatus(c.Sender().ID, consts.AddEvningDose)
		fmt.Println("Название: ", h.Bot.User[c.Sender().ID].Drugs.Drug_name, " У: ", h.Bot.User[c.Sender().ID].Drugs.M_dose, " Д: ", h.Bot.User[c.Sender().ID].Drugs.A_dose)
		return c.Send("Введите количество препарата вечером:")
	case consts.AddEvningDose:
		dataInt, err := ParseInt64FromString(c.Text())
		if err != nil {
			return c.Send("Укажите число")
		}
		h.SetDrugInt(c.Sender().ID, userStatus, dataInt)
		h.UpdateUserStatus(c.Sender().ID, consts.AddNightDose)
		h.DB.UpdateUserStatus(c.Sender().ID, consts.AddNightDose)
		fmt.Println("Название: ", h.Bot.User[c.Sender().ID].Drugs.Drug_name, " У: ", h.Bot.User[c.Sender().ID].Drugs.M_dose, " Д: ", h.Bot.User[c.Sender().ID].Drugs.A_dose, " В: ", h.Bot.User[c.Sender().ID].Drugs.E_dose)
		return c.Send("Введите количество препарата на ночь:")
	case consts.AddNightDose:
		dataInt, err := ParseInt64FromString(c.Text())
		if err != nil {
			return c.Send("Укажите число")
		}
		h.SetDrugInt(c.Sender().ID, userStatus, dataInt)
		h.UpdateUserStatus(c.Sender().ID, consts.AddDrugQuantity)
		h.DB.UpdateUserStatus(c.Sender().ID, consts.AddDrugQuantity)
		return c.Send("Введите количества оставшегося препарата:")
	case consts.AddDrugQuantity:
		dataInt, err := ParseInt64FromString(c.Text())
		if err != nil {
			return c.Send("Укажите число")
		}

		h.SetDrugInt(c.Sender().ID, userStatus, dataInt)
		h.UpdateUserStatus(c.Sender().ID, consts.AddDrugComment)

		h.DB.UpdateUserStatus(c.Sender().ID, consts.AddDrugComment)
		return c.Send("Введите комментарий к препарату")
	case consts.AddDrugComment:
		h.SetDrugString(c.Sender().ID, userStatus, c.Text())
		return h.DrugApprove(c)
	default:
		return c.Send("Неизвестный тип состояния, пожалуйста, повторите попытку.")
	}
}

/*
	func (h *Handler) CheckInt(c telebot.Context) (int64, error) {
		for i := range c.Text() {
			if !unicode.IsDigit(rune(c.Text()[i])) {
				return 0, c.Send("Введите число")
			}
		}
		return strconv.ParseInt(c.Text(), 10, 64)
	}
*/
func ParseInt64FromString(input string) (int64, error) {
	// Проверяем, что строка не пустая
	if len(input) == 0 {
		return 0, errors.New("строка пуста")
	}

	// Проверяем, что каждый символ строки является цифрой
	for _, r := range input {
		if !unicode.IsDigit(r) {
			return 0, errors.New("строка содержит недопустимые символы")
		}
	}

	// Преобразуем строку в int64
	result, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("ошибка преобразования строки в int64: %w", err)
	}

	return result, nil
}

func (h *Handler) DrugApprove(c telebot.Context) error {
	user, exists := h.Bot.User[c.Sender().ID]
	if !exists {
		fmt.Println("Пользователь не найден")
	}

	message := fmt.Sprintf("*Информация о созданном препарате:*\n- Название: %s\n- Доза утром: %d\n- Доза днем: %d\n- Доза вечером: %d\n- Доза ночью: %d\n- Количество: %d\n- Комментарий: %s",
		user.Drugs.Drug_name, user.Drugs.M_dose, user.Drugs.A_dose, user.Drugs.E_dose, user.Drugs.N_dose, user.Drugs.Quantity, config.EscapeMarkdown(user.Drugs.Comment))
	inlineKeyboard := telebot.ReplyMarkup{}
	uniqueString := fmt.Sprintf("d_sav")
	textString := fmt.Sprintf("Сохранить")
	btnApprove := telebot.InlineButton{
		Unique: uniqueString,
		Text:   textString,
	}
	uniqueStringCancel := fmt.Sprintf("d_can", user.Drugs.Id)
	btnCancel := telebot.InlineButton{
		Unique: uniqueStringCancel,
		Text:   "Отмена",
	}
	//btnEdit := h.SendEditDrug(c, drugIdInt)

	inlineKeyboard.InlineKeyboard = [][]telebot.InlineButton{
		{btnApprove, btnCancel},
	}

	// Отправка сообщения пользователю
	return c.Send(message, &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: &inlineKeyboard,
	})

}

func (h *Handler) EraseDrug(c telebot.Context) error {
	user, exists := h.Bot.User[c.Sender().ID]
	if !exists {
		fmt.Println("Пользователь не найден")
	}
	user.Drugs.Drug_name = ""
	user.Drugs.M_dose = 0
	user.Drugs.A_dose = 0
	user.Drugs.N_dose = 0
	user.Drugs.Quantity = 0
	user.Drugs.Comment = ""
	h.Bot.User[c.Sender().ID] = user
	return nil
}

func (h *Handler) SaveDrug(c telebot.Context) error {
	user, exists := h.Bot.User[c.Sender().ID]
	if !exists {
		fmt.Println("Пользователь не найден")
	}
	h.DB.InsertDrug(c.Sender().ID, user.Drugs)
	return h.handleShowUserDrugs(c)
}
