package util

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

type Uuid uint64

func (uuid Uuid) ToBase64() string {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(uuid))
	return base64.StdEncoding.EncodeToString(buf)
}

func Base64ToUuid(b64 string) (Uuid, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
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
