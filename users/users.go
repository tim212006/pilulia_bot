package users

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	Status    Status
	//Возможно оптимально было бы хранить данные о препаратах пользователя в мапе, но тогда надо пердусмотреть её
	//обнуление, что бы исключить нагрузку на память
	//Drugs map[int64]drugs.Drugs
}

type Status string

func (u *User) SetStatus(status Status) {
	u.Status = status
}

func (u *User) GetStatus() Status {
	return u.Status
}
