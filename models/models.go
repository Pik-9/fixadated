package models

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Pik-9/fixadated/util"
)

type Availability uint8

const (
	NotAvailable = Availability(iota)
	Availble
	Maybe
)

type Participant struct {
	EditId util.Uuid      `json:"editID"`
	Name   string         `json:"name"`
	Days   []Availability `json:"declarations"`
}

type Event struct {
	Id           util.Uuid      `json:"id"`
	EditId       util.Uuid      `json:"editID"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Dates        []int64        `json:"dates"`
	Participants []*Participant `json:"participants"`
}

var (
	eventById   map[util.Uuid]*Event
	eventByEdit map[util.Uuid]*Event
	partByEdit  map[util.Uuid]*Participant
)

func init() {
	eventById = make(map[util.Uuid]*Event)
	eventByEdit = make(map[util.Uuid]*Event)
	partByEdit = make(map[util.Uuid]*Participant)

	// TODO: Load events and participations from disk
}

func (prt Participant) ToJSON() []byte {
	ret, _ := json.Marshal(prt)
	return ret
}

func (part Participant) ToClientJSON(flat bool) []byte {
	pt := make(map[string]interface{})
	pt["name"] = part.Name
	pt["dates"] = part.Days
	if !flat {
		pt["editID"] = part.EditId.ToBase64()
	}

	ret, _ := json.Marshal(pt)
	return ret
}

func (part *Participant) EditSafely(npart Participant) {
	if len(part.Days) == len(npart.Days) {
		part.Days = npart.Days
	}

	if npart.Name != "" {
		part.Name = npart.Name
	}
}

func (evnt Event) ToJSON() []byte {
	ret, _ := json.Marshal(evnt)
	return ret
}

func (evnt Event) ToClientJSON(flat bool) []byte {
	type strippedDownPart struct {
		Name string         `json:"name"`
		Days []Availability `json:"declarations"`
	}

	ev := make(map[string]interface{})
	if !flat {
		ev["id"] = evnt.Id.ToBase64()
		ev["editID"] = evnt.EditId.ToBase64()
	}
	ev["name"] = evnt.Name
	ev["description"] = evnt.Description
	ev["dates"] = evnt.Dates
	parts := make([]strippedDownPart, len(evnt.Participants))

	for index, item := range evnt.Participants {
		parts[index].Name = item.Name
		parts[index].Days = item.Days
	}

	ev["participants"] = parts

	ret, _ := json.Marshal(ev)
	return ret
}

func (evnt *Event) EditSafely(nevnt Event) {
	if nevnt.Name != "" {
		evnt.Name = nevnt.Name
	}

	if nevnt.Description != "" {
		evnt.Description = nevnt.Description
	}
}

func (evnt *Event) RegisterParticipant(participant Participant) (*Participant, error) {
	if len(participant.Days) != len(evnt.Dates) {
		return nil, errors.New(fmt.Sprintf("Number of days mismatch: %d != %d.", len(participant.Days), len(evnt.Dates)))
	}

	npart := new(Participant)
	npart.EditId = util.RandomUuid()
	npart.Name = participant.Name
	npart.Days = participant.Days

	evnt.Participants = append(evnt.Participants, npart)
	partByEdit[npart.EditId] = npart

	return npart, nil
}

func GetEventByUuid(uuid util.Uuid) (*Event, error) {
	ret, ok := eventById[uuid]
	if ok {
		return ret, nil
	} else {
		return nil, errors.New("Event not found.")
	}
}

func GetEventByEditUuid(uuid util.Uuid) (*Event, error) {
	ret, ok := eventByEdit[uuid]
	if ok {
		return ret, nil
	} else {
		return nil, errors.New("Event not found.")
	}
}

func GetParticipantByUuid(uuid util.Uuid) (*Participant, error) {
	ret, ok := partByEdit[uuid]
	if ok {
		return ret, nil
	} else {
		return nil, errors.New("Participant not found.")
	}
}

func InsertNewEvent(event Event) *Event {
	ret := new(Event)
	ret.Name = event.Name
	ret.Description = event.Description
	ret.Dates = event.Dates

	ret.Id = util.RandomUuid()
	ret.EditId = util.RandomUuid()

	eventById[ret.Id] = ret
	eventByEdit[ret.EditId] = ret

	return ret
}
