package storage

import (
	"BusinessWallet/model"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u UserDB) Register(r *model.RegisterRequest) (*model.User, error) {
	hashedPassword, err := hashPassword(r.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:     r.Name,
		Surname:  r.Surname,
		Phone:    r.Phone,
		Password: hashedPassword,
		Email:    r.Email,
		Linkedin: r.Linkedin,
		Company:  r.Company,
		Contacts: make(pq.Int64Array, 0),
	}

	resp := u.db.Model(user).Create(&user)
	if resp.Error != nil {
		return nil, resp.Error
	}

	return user, err
}

func (u UserDB) Login(r *model.LoginRequest) (*model.User, error) {
	var user *model.User
	err := u.db.Model(user).Where("email = ?", r.Email).Find(&user).Error
	if err != nil {
		return nil, err
	}

	if !checkPassword(r.Password, user.Password) {
		return nil, errors.New("wrong email or pass")
	}

	return user, nil
}

func (u UserDB) GetUser(id int) (*model.User, error) {
	var user *model.User
	err := u.db.Model(user).Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u UserDB) AddContact(userId, contactId, eventId int) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.User{}).Where("id = ? and ? != ALL(contacts)", userId, contactId).
			Update("contacts", gorm.Expr("contacts || ?", pq.Int64Array([]int64{int64(contactId)}))).Error
		if err != nil {
			return err
		}

		err = tx.Model(model.User{}).Where("id = ? and ? != ALL(contacts)", contactId, userId).
			Update("contacts", gorm.Expr("contacts || ?", pq.Int64Array([]int64{int64(userId)}))).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (u UserDB) DeleteContact(userId, contactId, eventId int) error {
	return u.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(model.User{}).Where("id = ?", userId).
			Update("contacts", gorm.Expr("array_remove(contacts, ?)", contactId)).Error

		if err != nil {
			return err
		}

		err = tx.Model(model.User{}).Where("id = ?", contactId).
			Update("contacts", gorm.Expr("array_remove(contacts, ?)", userId)).Error

		if err != nil {
			return err
		}

		return nil
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
