package resourceid

import (
	"encoding/base64"
	"math/rand"

	"github.com/google/uuid"
)

func New() string {
	id := uuid.New().ID()
	idBytes := uintToBytes(id)

	// first byte should be bigger than 0x80 for collection ids
	// clients classify this id as "user" otherwise
	if (idBytes[0] & 0x80) <= 0 {
		idBytes[0] = byte(rand.Intn(0x80) + 0x80)
	}

	return base64.StdEncoding.EncodeToString(idBytes)
}

func NewCombined(ids ...string) string {
	combinedIdBytes := make([]byte, 0)

	for _, id := range ids {
		idBytes, _ := base64.StdEncoding.DecodeString(id)
		combinedIdBytes = append(combinedIdBytes, idBytes...)
	}

	return base64.StdEncoding.EncodeToString(combinedIdBytes)
}

func uintToBytes(id uint32) []byte {
	buf := make([]byte, 4)
	for i := 0; i < 4; i++ {
		buf[i] = byte(id >> (i * 8))
	}

	return buf
}
