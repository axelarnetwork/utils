package proto_test

import (
	"testing"
	"time"

	"github.com/axelarnetwork/utils/proto"
	gogoprototypes "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	durationHour := gogoprototypes.DurationProto(time.Hour)
	durationMinute := gogoprototypes.DurationProto(time.Minute)

	actualHour := proto.Hash(durationHour)
	assert.Len(t, actualHour, 32)

	actualMinute := proto.Hash(durationMinute)
	assert.Len(t, actualMinute, 32)

	assert.NotEqual(t, actualHour, actualMinute)
}
