package entity

type UserInterface interface {
	GetId() uint32
	GetEmail() string
}

type User struct {
	Id uint32
}

func (u *User) GetID() uint32 {
	return u.Id
}
