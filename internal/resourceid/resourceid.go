package resourceid

import (
	"encoding/base64"
	"math/rand"
	"strings"

	"github.com/google/uuid"
)

type ResourceType int

const (
	ResourceTypeDatabase ResourceType = iota
	ResourceTypeCollection
	ResourceTypeDocument
	ResourceTypeStoredProcedure
	ResourceTypeTrigger
	ResourceTypeUserDefinedFunction
	ResourceTypeConflict
	ResourceTypePartitionKeyRange
	ResourceTypeSchema
)

func New(resourceType ResourceType) string {
	var idBytes []byte
	switch resourceType {
	case ResourceTypeDatabase:
		idBytes = randomBytes(4)
	case ResourceTypeCollection:
		idBytes = randomBytes(4)
		// first byte should be bigger than 0x80 for collection ids
		// clients classify this id as "user" otherwise
		if (idBytes[0] & 0x80) <= 0 {
			idBytes[0] = byte(rand.Intn(0x80) + 0x80)
		}
	case ResourceTypeDocument:
		idBytes = randomBytes(8)
		idBytes[7] = byte(rand.Intn(0x10)) // Upper 4 bits = 0
	case ResourceTypeStoredProcedure:
		idBytes = randomBytes(8)
		idBytes[7] = byte(rand.Intn(0x10)) | 0x08 // Upper 4 bits = 0x08
	case ResourceTypeTrigger:
		idBytes = randomBytes(8)
		idBytes[7] = byte(rand.Intn(0x10)) | 0x07 // Upper 4 bits = 0x07
	case ResourceTypeUserDefinedFunction:
		idBytes = randomBytes(8)
		idBytes[7] = byte(rand.Intn(0x10)) | 0x06 // Upper 4 bits = 0x06
	case ResourceTypeConflict:
		idBytes = randomBytes(8)
		idBytes[7] = byte(rand.Intn(0x10)) | 0x04 // Upper 4 bits = 0x04
	case ResourceTypePartitionKeyRange:
		// we don't do partitions yet, so just use a fixed id
		idBytes = []byte{0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x69, 0x50}
	case ResourceTypeSchema:
		idBytes = randomBytes(8)
		idBytes[7] = byte(rand.Intn(0x10)) | 0x09 // Upper 4 bits = 0x09
	default:
		idBytes = randomBytes(4)
	}

	encoded := base64.StdEncoding.EncodeToString(idBytes)
	return strings.ReplaceAll(encoded, "/", "-")
}

func NewCombined(ids ...string) string {
	combinedIdBytes := make([]byte, 0)

	for _, id := range ids {
		idBytes, _ := base64.StdEncoding.DecodeString(strings.ReplaceAll(id, "-", "/"))
		combinedIdBytes = append(combinedIdBytes, idBytes...)
	}

	encoded := base64.StdEncoding.EncodeToString(combinedIdBytes)
	return strings.ReplaceAll(encoded, "/", "-")
}

func uintToBytes(id uint32) []byte {
	buf := make([]byte, 4)
	for i := 0; i < 4; i++ {
		buf[i] = byte(id >> (i * 8))
	}

	return buf
}

func randomBytes(count int) []byte {
	buf := make([]byte, count)
	for i := 0; i < count; i += 4 {
		id := uuid.New().ID()
		idBytes := uintToBytes(id)
		copy(buf[i:], idBytes)
	}
	return buf
}
