package app

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"pilulia_bot/drugs"
)

// Кнопка препараты

// Функция динамического создания кнопок в сообщении

func (h *Handler) SendDynamicButtonMessage(c telebot.Context, drugs map[string]drugs.Drugs) error {
	inlineKeyboard := &telebot.ReplyMarkup{}
	for key, drug := range drugs {
		uniqueString := fmt.Sprintf("drug_%d", drug.Id)
		textString := fmt.Sprintf("%s, Ост.:%d", key, drug.Quantity)
		btn := telebot.InlineButton{
			Unique: uniqueString,
			Text:   textString,
		}
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn})
	}
	btnAddDrug := telebot.InlineButton{
		Unique: "add_d",
		Text:   "Добавить препарат",
	}
	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btnAddDrug})
	return c.Send("Ваши препараты", inlineKeyboard)
}

// *******************************************
// Функция кнопки меню
func (h *Handler) menuButton(c telebot.Context, btnName, btnData string) error {
	btn := telebot.Btn{Text: btnName, Data: btnData}
	keyboard := telebot.ReplyMarkup{ResizeKeyboard: true}
	keyboard.Reply(keyboard.Row(btn))
	return c.Send("Здесь будет информация о препаратах на текущую дату", &keyboard)
}

// Функция вызова стартового меню

func (h *Handler) SendMainMenu(c telebot.Context) error {
	err := h.handleShowDailyUserDrugs(c)
	if err != nil {
		return err
	}
	btnMain := telebot.Btn{Text: "На главную"}
	btnDrugs := telebot.Btn{Text: "Препараты"}
	btnHelp := telebot.Btn{Text: "Помощь"}
	markup := &telebot.ReplyMarkup{ResizeKeyboard: true}
	markup.Reply(markup.Row(btnMain, btnDrugs, btnHelp))
	return c.Send("Выберите действие:", &telebot.SendOptions{
		ReplyMarkup: markup,
	})
}

// ******************************************
// InlineMenu

func (h *Handler) SendDailyDrugs(c telebot.Context, drugs map[string]drugs.Drugs) error {
	// Функция для отправки сообщений с кнопками для конкретного времени суток
	sendDrugPeriod := func(period string, buttons []telebot.InlineButton) error {
		if len(buttons) > 0 {
			// Создаем и отправляем клавиатуру с кнопками
			inlineKeyboard := &telebot.ReplyMarkup{}
			for _, btn := range buttons {
				inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn})
			}
			return c.Send(fmt.Sprintf("Препараты для приема %s:", period), inlineKeyboard)
		} else {
			// Сообщение, если препаратов для данного времени суток нет
			return c.Send(fmt.Sprintf("Препаратов для приема %s нет", period))
		}
	}

	// Создание кнопок для утренних препаратов
	morningButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.M_dose > 0 {
			uniqueString := fmt.Sprintf("daily_m_%d", drug.Id)
			textString := fmt.Sprintf("У, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			morningButtons = append(morningButtons, btn)
		}
	}

	// Создание кнопок для дневных препаратов
	afternoonButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.A_dose > 0 {
			uniqueString := fmt.Sprintf("daily_a_%d", drug.Id)
			textString := fmt.Sprintf("Д, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			afternoonButtons = append(afternoonButtons, btn)
		}
	}

	// Создание кнопок для вечерних препаратов
	eveningButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.E_dose > 0 {
			uniqueString := fmt.Sprintf("daily_e_%d", drug.Id)
			textString := fmt.Sprintf("В, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			eveningButtons = append(eveningButtons, btn)
		}
	}

	// Создание кнопок для ночных препаратов
	nightButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.N_dose > 0 {
			uniqueString := fmt.Sprintf("daily_n_%d", drug.Id)
			textString := fmt.Sprintf("Н, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			nightButtons = append(nightButtons, btn)
		}
	}

	// Отправка сообщений для каждого времени суток
	if err := sendDrugPeriod("утром", morningButtons); err != nil {
		return err
	}
	if err := sendDrugPeriod("днем", afternoonButtons); err != nil {
		return err
	}
	if err := sendDrugPeriod("вечером", eveningButtons); err != nil {
		return err
	}
	if err := sendDrugPeriod(" на ночь", nightButtons); err != nil {
		return err
	}

	return nil
}
