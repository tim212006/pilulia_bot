package users

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Username  string
	Status    Status
}

type Status string

func (u *User) SetStatus(status Status) {
	u.Status = status
}

func (u *User) GetStatus() Status {
	return u.Status
}
