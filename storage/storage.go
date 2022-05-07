package storage

import (
	"BusinessWallet/config"
	"BusinessWallet/model"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserDB struct {
	db *gorm.DB
}

type EventDB struct {
	db *gorm.DB
}

var User UserDB
var Event EventDB

func Connect(options config.StorageOptions) error {
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		options.Username,
		options.Password,
		options.Host,
		options.Port,
		options.DB,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&model.User{}, &model.Event{})
	if err != nil {
		return err
	}

	User = UserDB{db: db}
	Event = EventDB{db: db}

	return nil
}

func Seed() {
	registerRequests := []model.RegisterRequest{
		{Name: "ahmet", Surname: "Şenharputlu", Email: "ahmet@test.com", Phone: "91864198274", Password: "testpsw"},
		{Name: "gizem", Surname: "bulut", Email: "gizem@test.com", Phone: "91864198244", Password: "testpsw"},
		{Name: "damla", Surname: "damla", Email: "damla@test.com", Phone: "91864298274", Password: "testpsw"},
		{Name: "zeynep", Surname: "emre", Email: "zeynep@test.com", Phone: "941864198274", Password: "testpsw"},
		{Name: "emre", Surname: "iuc", Email: "emre@test.com", Phone: "91864118274", Password: "testpsw"},
		{Name: "mira", Surname: "hello", Email: "mira@test.com", Phone: "91864198234", Password: "testpsw"},
	}

	userList := make([]*model.User, 0)

	for _, u := range registerRequests {
		user, err := User.Register(&u)

		if user == nil || err != nil {
			continue
		}

		userList = append(userList, user)
	}

	createEventRequests := []model.CreateEventRequest{
		{Name: "Meet", Description: "lets meet, it is passed, and you were attending !", Start: time.Now().Add(time.Hour * -12), Finish: time.Now().Add(time.Hour * -12)},
		{Name: "Aws Conference", Description: "it is current (now), and you are attending!", Start: time.Now().Add(time.Hour * -6), Finish: time.Now().Add(time.Hour * 12)},
		{Name: "Security Event", Description: "its active and ghonna start after a while, and you are attending!", Start: time.Now().Add(time.Hour * 12), Finish: time.Now().Add(time.Hour * 24)},
		{Name: "Fair Event", Description: "its ghonna start after a while, and you are not attending!", Start: time.Now().Add(time.Hour * 12), Finish: time.Now().Add(time.Hour * 24)},
	}

	eventList := make([]*model.Event, 0)

	for _, e := range createEventRequests {
		x, _ := json.Marshal(e.Start)
		fmt.Println(string(x))
		event, err := Event.Create(&e, int(userList[0].ID))
		if event == nil || err != nil {
			logrus.WithError(err).Error("event create: ")
			continue
		}

		eventList = append(eventList, event)
	}

	for _, user := range userList[2:] {
		for _, event := range eventList {
			err := Event.Attend(int(event.ID), int(user.ID))
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	Event.Attend(int(eventList[0].ID), int(userList[1].ID))
	Event.Attend(int(eventList[1].ID), int(userList[1].ID))
	Event.Attend(int(eventList[2].ID), int(userList[1].ID))

	err := User.AddContact(int(userList[0].ID), int(userList[1].ID), 0)
	if err != nil {
		logrus.Error(err)
	}
}

func Delete() error {
	err := User.db.Unscoped().Select(clause.Associations).Where("1 = 1").Delete(model.User{}).Error
	if err != nil {
		return err
	}

	err = Event.db.Unscoped().Select(clause.Associations).Where("1 = 1").Delete(model.Event{}).Error
	if err != nil {
		return err
	}

	return nil
}
