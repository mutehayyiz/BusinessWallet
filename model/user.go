package model

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"html"
	"strings"
)

type User struct {
	gorm.Model
	Name     string        `json:"name"`
	Surname  string        `json:"surname"`
	Phone    string        `json:"phone" gorm:"index:idx_users_phone,unique"`
	Password string        `json:"password"`
	Email    string        `json:"email" gorm:"index:idx_users_phone,unique"`
	Linkedin string        `json:"linkedin"`
	Company  string        `json:"company"`
	Contacts pq.Int64Array `json:"contacts" gorm:"type:integer[]"`
}

type Contact struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Linkedin string `json:"linkedin"`
	Company  string `json:"company"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Linkedin string `json:"linkedin"`
	Company  string `json:"company"`
}

func (r RegisterRequest) Validate() error {
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.Surname = html.EscapeString(strings.TrimSpace(r.Surname))
	r.Phone = html.EscapeString(strings.TrimSpace(r.Surname))
	r.Email = html.EscapeString(strings.TrimSpace(r.Email))
	r.Linkedin = html.EscapeString(strings.TrimSpace(r.Linkedin))
	r.Company = html.EscapeString(strings.TrimSpace(r.Company))

	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Surname == "" {
		return errors.New("surname is required")
	}
	if r.Password == "" {
		return errors.New("password is required")
	}
	if r.Phone == "" {
		return errors.New("phone is required")
	}

	if err := checkmail.ValidateFormat(r.Email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginRequest) Validate() error {
	r.Email = html.EscapeString(strings.TrimSpace(r.Email))

	if r.Password == "" {
		return errors.New("password is required")
	}

	if err := checkmail.ValidateFormat(r.Email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserData User   `json:"user_data"`
}
