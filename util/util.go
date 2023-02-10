// fixadated is a daemon for a collaborative date finding tool.
//    Copyright (C) 2023 Daniel Steinhauer <d.steinhauer@mailbox.org>
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as
//    published by the Free Software Foundation, either version 3 of the
//    License, or (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.
package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net/http"
)

type Uuid uint64

func (uuid Uuid) ToBase64() string {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(uuid))
	return base64.URLEncoding.EncodeToString(buf)
}

func Base64ToUuid(b64 string) (Uuid, error) {
	data, err := base64.URLEncoding.DecodeString(b64)
	if err != nil {
		return 0, err
	}

	if len(data) != 8 {
		return 0, errors.New("UUID Size Mismatch.")
	}

	ret := Uuid(binary.LittleEndian.Uint64(data))
	return ret, nil
}

func RandomUuid() Uuid {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return Uuid(binary.LittleEndian.Uint64(buf))
}

func DecodeJsonBody[T interface{}](w http.ResponseWriter, r *http.Request, obj *T) error {
	ctype := r.Header.Get("Content-Type")
	if ctype != "" {
		if ctype != "application/json" {
			return errors.New("Wrong Content-Type: " + ctype)
		}
	}

	bodyStream := http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(bodyStream)
	decoder.DisallowUnknownFields()

	return decoder.Decode(obj)
}
