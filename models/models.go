package models

import (
	"encoding/json"
	"github.com/Pik-9/fixadated/util"
	"time"
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
	EditId Util.Uuid `json:"editID"`
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

func (prt Participant) ToJSON() string {
	ret, _ := json.Marshal(prt)
	return string(ret)
}

func (prt ParticipantPlus) ToJSON() string {
	ret, _ := json.Marshal(prt)
	return string(ret)
}

func (evnt Event) ToJSON() string {
	ret, _ := json.Marshal(evnt)
	return ret
}

func (evnt EventPlus) ToJSON() string {
	ret, _ := json.Marshal(evnt)
	return ret
}
