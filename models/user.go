package models

import "example/pkg"

type User struct {
	ID       pkg.PrimaryKey `json:"id"`
	Username string         `json:"username"`
}

func (u *User) GetID() pkg.PrimaryKey {
	return u.ID
}
func (u *User) SetID(id pkg.PrimaryKey) {
	u.ID = id
}

var _ pkg.Model = (*User)(nil)
