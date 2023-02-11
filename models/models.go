// fixadated is a daemon for a collaborative date finding tool.
//
//	Copyright (C) 2023 Daniel Steinhauer <d.steinhauer@mailbox.org>
//
//	This program is free software: you can redistribute it and/or modify
//	it under the terms of the GNU Affero General Public License as
//	published by the Free Software Foundation, either version 3 of the
//	License, or (at your option) any later version.
//
//	This program is distributed in the hope that it will be useful,
//	but WITHOUT ANY WARRANTY; without even the implied warranty of
//	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//	GNU Affero General Public License for more details.
//
//	You should have received a copy of the GNU Affero General Public License
//	along with this program.  If not, see <https://www.gnu.org/licenses/>.
package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"log"
	"bytes"

	"github.com/Pik-9/fixadated/util"
)

type Availability int8

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

	err := loadFromDisk()
	if err != nil {
		log.Printf("Could not load data from disk: %s\n", err)
	}
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

func SaveToDisk() error {
	path, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	fdir := path + "/fixadated"
	err = os.MkdirAll(fdir, 0755)
	if err != nil {
		return err
	}

	fout, err := os.OpenFile(fdir+"/dates.db", os.O_RDWR|os.O_CREATE, 0755)
	defer fout.Close()
	if err != nil {
		return err
	}

	evnts := make([]*Event, 0)
	for _, val := range eventById {
		evnts = append(evnts, val)
	}

	out, _ := json.Marshal(evnts)
	_, err = fout.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func loadFromDisk() error {
	path, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	fpath := path + "/fixadated/dates.db"

	rawContent, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(rawContent)
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	var events []*Event
	err = decoder.Decode(&events)

	if err != nil {
		return err
	}

	for _, evnt := range events {
		eventById[evnt.Id] = evnt
		eventByEdit[evnt.EditId] = evnt
		for _, part := range evnt.Participants {
			partByEdit[part.EditId] = part
		}
	}

	return nil
}
