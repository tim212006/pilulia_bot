package app

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"pilulia_bot/drugs"
)

// Кнопка препараты

// Функция динамического создания кнопок в сообщении
//Возможно есть смысл переписать функцию динамического создания сообщения с кнопками, так что бы она принимала контекст,
//любую мапу и ключевое значение в любом виде, но пока так

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
	btn_add_drug := telebot.InlineButton{
		Unique: "add_d",
		Text:   "Добавить препарат",
	}
	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn_add_drug})
	return c.Send("Ваши препараты", inlineKeyboard)
}

//Ниже затупил - надо переписать. Функция вызова меню должна создавать кнопки, клавиатуру из кнопок и вызывать клаву.
//Сейчас функция вызова меню вызывает кнопку и уже кнопка создает клавиатуру.

func (h *Handler) InlineBtn(btnName, btnData string) telebot.Btn {
	return telebot.Btn{Text: btnName, Data: btnData}
}

// func (h *Handler) MarkupButton(c telebot.Context) error {}
// *******************************************
// Функция кнопки меню
func (h *Handler) menuButton(c telebot.Context, btnName, btnData string) error {
	btn := telebot.Btn{Text: btnName, Data: btnData}
	keyboard := telebot.ReplyMarkup{ResizeKeyboard: true}
	keyboard.Reply(keyboard.Row(btn))
	return c.Send("Здесь будет информация о препаратах на текущую дату", &keyboard)
}

// Функция вызова стартового меню
func (h *Handler) menuCommand(c telebot.Context) error {
	btnName := "Препараты"
	btnData := "BtnData"
	return h.menuButton(c, btnName, btnData)
}

func (h *Handler) SendMainMenu(c telebot.Context) error {
	h.handleShowDailyUserDrugs(c)
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
func (h *Handler) SendMenuWithInlineButtons(c telebot.Context) error {
	// Создаем первую inline-кнопку
	btn1 := telebot.InlineButton{
		Unique: "showUserDrugs",
		Text:   "Препараты",
	}

	// Создаем вторую inline-кнопку
	btn2 := telebot.InlineButton{
		Unique: "getHelp",
		Text:   "Помощь",
	}

	// Создаем объект для inline-клавиатуры
	keyboard := &telebot.ReplyMarkup{}
	keyboard.InlineKeyboard = [][]telebot.InlineButton{
		{btn1, btn2}, // Обе кнопки в одной строке
	}

	// Отправляем сообщение с inline-клавиатурой
	return c.Send("Выберите действие:", keyboard)
}

// MarkupMenu
func (h *Handler) SendMenuWithMarkupButtons(c telebot.Context) error {
	// Создаем кнопку "Препараты"
	btn1 := telebot.Btn{
		Text: "Препараты",
	}

	// Создаем кнопку "Помощь"
	btn2 := telebot.Btn{
		Text: "Помощь",
	}

	// Создаем объект для клавиатуры
	markup := &telebot.ReplyMarkup{ResizeKeyboard: true}
	markup.Reply(
		markup.Row(btn1, btn2), // Обе кнопки в одной строке
	)

	// Отправляем сообщение с клавиатурой
	return c.Send("Выберите действие:", markup)
}

// кнопка препарат в сообщении
func (h *Handler) InlineButtonMessege(c telebot.Context, drugName string) error {
	btn := telebot.InlineButton{
		Unique: "write_down_drugs",
		Text:   drugName,
	}
	inlineKeyboard := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btn},
		},
	}
	return c.Send("Список препаратов: ", &inlineKeyboard)
}

func (h *Handler) SendEditDrug(c telebot.Context, drugId int64) telebot.InlineButton {
	uniqueString := fmt.Sprintf("drug_edit_%d", drugId)
	textString := fmt.Sprintf("Редактировать")
	btn := telebot.InlineButton{
		Unique: uniqueString,
		Text:   textString,
	}
	return btn
}

/*
func (h *Handler) SendDailyDrugs(c telebot.Context, drugs map[string]drugs.Drugs) error {
	inlineKeyboard := &telebot.ReplyMarkup{}

	morningButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.M_dose > 0 {
			uniqueString := fmt.Sprintf("daily%d", drug.Id)
			textString := fmt.Sprintf("У, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			morningButtons = append(morningButtons, btn)
		}
	}
	afternoonButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.A_dose > 0 {
			uniqueString := fmt.Sprintf("daily%d", drug.Id)
			textString := fmt.Sprintf("Д, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			afternoonButtons = append(afternoonButtons, btn)
		}
	}
	eveningButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.E_dose > 0 {
			uniqueString := fmt.Sprintf("daily%d", drug.Id)
			textString := fmt.Sprintf("В, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			eveningButtons = append(eveningButtons, btn)
		}
	}
	nightButtons := []telebot.InlineButton{}
	for key, drug := range drugs {
		if drug.N_dose > 0 {
			uniqueString := fmt.Sprintf("daily%d", drug.Id)
			textString := fmt.Sprintf("Н, %s, ост.: %d", key, drug.Quantity)
			btn := telebot.InlineButton{
				Unique: uniqueString,
				Text:   textString,
			}
			nightButtons = append(nightButtons, btn)
		}
	}
	if len(morningButtons) > 0 {
		for _, btn := range morningButtons {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn})
		}
	} else {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{{Text: "Препаратов для приема утром нет"}})
	}

	if len(afternoonButtons) > 0 {
		for _, btn := range afternoonButtons {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn})
		}
	} else {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{{Text: "Препаратов для приема днем нет"}})
	}

	if len(eveningButtons) > 0 {
		for _, btn := range eveningButtons {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn})
		}
	} else {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{{Text: "Препаратов для приема вечером нет"}})
	}
	if len(nightButtons) > 0 {
		for _, btn := range nightButtons {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{btn})
		}
	} else {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, []telebot.InlineButton{{Text: "Препаратов для приема ночью нет"}})
	}
	return c.Send("Ваши препараты сегодня", inlineKeyboard)
}*/

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
