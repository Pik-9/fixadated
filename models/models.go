package models

import (
	"encoding/json"
	"github.com/Pik-9/fixadated/util"
	"time"
	"errors"
	"fmt"
)

type Availability uint8

const (
	Available = Availability(iota)
	Maybe
	NotAvailble
)

type Participant struct {
	Name string         `json:"name"`
	Days []Availability `json:"declarations"`
}

type ParticipantPlus struct {
	Participant
	EditId util.Uuid `json:"editID"`
}

type Event struct {
	Name         string        `json:"name"`
	Description  string        `json:"description"`
	Dates        []time.Time   `json:"dates"`
	Participants []Participant `json:"participants"`
}

type EventPlus struct {
	Event
	Id     util.Uuid `json:"id"`
	EditId util.Uuid `json:"editID"`
}

var (
	eventById map[util.Uuid]*EventPlus
	eventByEdit map[util.Uuid]*EventPlus
	partById map[util.Uuid]*ParticipantPlus
	partByEdit map[util.Uuid]*ParticipantPlus
)

func init() {
	eventById = make(map[util.Uuid]*EventPlus)
	eventByEdit = make(map[util.Uuid]*EventPlus)
	partById = make(map[util.Uuid]*ParticipantPlus)
	partByEdit = make(map[util.Uuid]*ParticipantPlus)

	// TODO: Load events and participations from disk

	// Dummy code
	dummy := InsertNewEvent(Event{
		Name: "Test Event",
		Description: "Bla bla",
		Dates: make([]time.Time, 0),
		Participants: make([]Participant, 0),
	})

	fmt.Println("Created dummy event with uuid", dummy.Id.ToBase64())
}

func (prt Participant) ToJSON() []byte {
	ret, _ := json.Marshal(prt)
	return ret
}

func (prt ParticipantPlus) ToJSON() []byte {
	ret, _ := json.Marshal(prt)
	return ret
}

func (evnt Event) ToJSON() []byte {
	ret, _ := json.Marshal(evnt)
	return ret
}

func (evnt EventPlus) ToJSON() []byte {
	ret, _ := json.Marshal(evnt)
	return ret
}

func GetEventByUuid(uuid util.Uuid) (*EventPlus, error) {
	ret, ok := eventById[uuid]
	if ok {
		return ret, nil
	} else {
		return nil, errors.New("Event not found.")
	}
}

func InsertNewEvent(event Event) *EventPlus {
	ret := EventPlus{
		event,
		util.RandomUuid(),
		util.RandomUuid(),
	}

	eventById[ret.Id] = &ret
	eventByEdit[ret.EditId] = &ret

	return &ret
}
