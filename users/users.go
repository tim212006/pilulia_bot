package users

import (
	"fmt"
	"pilulia_bot/drugs"
	"pilulia_bot/logger/consts"
)

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	Status    Status
	//Возможно оптимально было бы хранить данные о препаратах пользователя в мапе, но тогда надо пердусмотреть её
	//обнуление, что бы исключить нагрузку на память
	Drugs drugs.Drugs
}

type Status string

func (u *User) SetStatus(status Status) {
	u.Status = status
}

func (u *User) GetStatus() Status {
	return u.Status
}

func (u *User) SetDrugString(status Status, data string) {
	switch status {
	case consts.AddDrugName:
		u.Drugs.Drug_name = data
	case consts.AddDrugComment:
		u.Drugs.Comment = data
	default:
		fmt.Println("Статус пользователя: ", status, " не соответствует функции")
	}
}

func (u *User) SetDrugInt(status Status, data int64) {
	switch status {
	case consts.AddMorningDose:
		u.Drugs.M_dose = data
	case consts.AddAfternoonDose:
		u.Drugs.A_dose = data
	case consts.AddEvningDose:
		u.Drugs.E_dose = data
	case consts.AddNightDose:
		u.Drugs.N_dose = data
	case consts.AddDrugQuantity:
		u.Drugs.Quantity = data
	default:
		fmt.Println("Статус пользователя: ", status, " не соответствует функции")

	}
}

func (u *User) EraseDrug() {
	u.Drugs.Drug_name = ""
	u.Drugs.M_dose = 0
	u.Drugs.A_dose = 0
	u.Drugs.N_dose = 0
	u.Drugs.Quantity = 0
	u.Drugs.Comment = ""
}
