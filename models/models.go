package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Pik-9/fixadated/util"
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
	Id           util.Uuid          `json:"id"`
	EditId       util.Uuid          `json:"editID"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	Dates        []time.Time        `json:"dates"`
	Participants []*ParticipantPlus `json:"participants"`
}

var (
	eventById   map[util.Uuid]*EventPlus
	eventByEdit map[util.Uuid]*EventPlus
	partByEdit  map[util.Uuid]*ParticipantPlus
)

func init() {
	eventById = make(map[util.Uuid]*EventPlus)
	eventByEdit = make(map[util.Uuid]*EventPlus)
	partByEdit = make(map[util.Uuid]*ParticipantPlus)

	// TODO: Load events and participations from disk
}

func NewEvent(name string, description string, dates []time.Time) Event {
	ret := Event{
		Name:         name,
		Description:  description,
		Dates:        nil,
		Participants: make([]Participant, 0),
	}

	if dates == nil {
		ret.Dates = make([]time.Time, 0)
	} else {
		ret.Dates = dates
	}

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
		Id:           util.RandomUuid(),
		EditId:       util.RandomUuid(),
		Name:         event.Name,
		Description:  event.Description,
		Dates:        event.Dates,
		Participants: make([]*ParticipantPlus, 0),
	}

	eventById[ret.Id] = &ret
	eventByEdit[ret.EditId] = &ret

	return &ret
}

func EditEvent(editId util.Uuid, name string, description string, dates []time.Time) (*EventPlus, error) {
	evnt, found := eventByEdit[editId]
	if !found {
		return nil, errors.New("Not Found")
	}

	if len(evnt.Dates) != len(dates) {
		return nil, errors.New(fmt.Sprintf("Number of days mismatch: %d != %d.", len(evnt.Dates), len(dates)))
	}

	evnt.Name = name
	evnt.Description = description
	evnt.Dates = dates

	return evnt, nil
}

func EditParticipation(editId util.Uuid, name string, dates []Availability) (*ParticipantPlus, error) {
	part, found := partByEdit[editId]
	if !found {
		return nil, errors.New("Not Found")
	}

	if len(part.Days) != len(dates) {
		return nil, errors.New(fmt.Sprintf("Number of days mismatch: %d != %d.", len(part.Days), len(dates)))
	}

	part.Name = name
	part.Days = dates

	return part, nil
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

func (evnt EventPlus) ToClientJSON() []byte {
	ev := struct{
		Id           string          `json:"id"`
		EditId       string          `json:"editID"`
		Name         string             `json:"name"`
		Description  string             `json:"description"`
		Dates        []time.Time        `json:"dates"`
		Participants []*ParticipantPlus `json:"participants"`
	}{
		Id: evnt.Id.String(),
		EditId: evnt.EditId.String(),
		Name: evnt.Name,
		Description: evnt.Description,
		Dates: evnt.Dates,
		Participants: evnt.Participants,
	}

	ret, _ := json.Marshal(ev)
	return ret
}

func (evnt EventPlus) GetFlatEvent() Event {
	ret := Event{
		Name:         evnt.Name,
		Description:  evnt.Description,
		Dates:        evnt.Dates,
		Participants: make([]Participant, len(evnt.Participants)),
	}

	for index, part := range evnt.Participants {
		ret.Participants[index] = part.Participant
	}

	return ret
}

func (evnt *EventPlus) RegisterParticipant(participant Participant) (*ParticipantPlus, error) {
	if len(participant.Days) != len(evnt.Dates) {
		return nil, errors.New(fmt.Sprintf("Number of days mismatch: %d != %d.", len(participant.Days), len(evnt.Dates)))
	}

	npart := ParticipantPlus{
		participant,
		util.RandomUuid(),
	}

	evnt.Participants = append(evnt.Participants, &npart)
	partByEdit[npart.EditId] = &npart

	return &npart, nil
}
