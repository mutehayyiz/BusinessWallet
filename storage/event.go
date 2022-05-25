package storage

import (
	"BusinessWallet/model"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"html"
	"strconv"
	"time"
)

func (e EventDB) Create(c *model.CreateEventRequest, organizer int) (*model.Event, error) {
	event := &model.Event{
		Name:        c.Name,
		Description: c.Description,
		Organizer:   organizer,
		Start:       c.Start,
		Finish:      c.Finish,
	}

	event.Attendees = append(event.Attendees, int64(organizer))

	err := e.db.Model(event).Create(event).Error
	if err != nil {
		return nil, err
	}

	return event, err
}

func (e EventDB) Get(eventId string) (*model.Event, error) {
	var event *model.Event
	err := e.db.Model(event).Where("id = ?", eventId).First(&event).Error
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (e EventDB) Attend(eventId int, userID int) error {
	return e.db.Model(model.Event{}).Where("id = ? and ? != ALL(attendees)", eventId, userID).
		Update("attendees", gorm.Expr("attendees || ?", pq.Int64Array([]int64{int64(userID)}))).Error
}

func (e EventDB) Leave(eventId string, userID int) error {
	return e.db.Model(model.Event{}).Where("id = ?", eventId).
		Update("attendees", gorm.Expr("array_remove(attendees, ?)", userID)).Error
}

func (e EventDB) Past(id int) (events []*model.Event, err error) {
	err = e.db.Model(model.Event{}).Where("? = ANY(attendees) and finish <= ?", id, time.Now()).Find(&events).Error
	return
}

func (e EventDB) Active() (events []*model.Event, err error) {
	err = e.db.Model(model.Event{}).Where("start >= ?", time.Now()).Find(&events).Error
	return
}

func (e EventDB) Now(id int) (events []*model.Event, err error) {
	err = e.db.Model(model.Event{}).Where("? = ANY(attendees) and finish >= ? and start <= ? ", id, time.Now(), time.Now()).Find(&events).Error
	return
}

func (e EventDB) Delete(eventId string, id int) error {
	eventId = html.EscapeString(eventId)
	eId, err := strconv.ParseUint(eventId, 10, 64)
	if err != nil {
		return err
	}

	ret := e.db.Model(model.Event{}).Where("id = ? and organizer = ?", eId, id).Delete(&model.Event{})
	if ret.Error != nil {
		logrus.Error(ret.Error)
		return ret.Error
	}

	return nil
}

func (e EventDB) Together(userId, contactId int) (events []*model.Event, err error) {
	var sel = "SELECT * FROM events WHERE ? <@ attendees"
	ids := []int{userId, contactId}
	err = e.db.Raw(sel, pq.Array(ids)).Take(&events).Error
	return
}
