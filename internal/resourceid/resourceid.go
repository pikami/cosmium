package resourceid

import (
	"encoding/base64"

	"github.com/google/uuid"
)

func New() string {
	id := uuid.New().ID()
	idBytes := uintToBytes(id)

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
