package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveInputTransport(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(MoveLeftInput, DeserializeMoveInput(SerializeMoveInput(MoveLeftInput)))
	assert.Equal(MoveRightInput, DeserializeMoveInput(SerializeMoveInput(MoveRightInput)))
}

func TestShootInputTransport(t *testing.T) {
	assert := assert.New(t)
	i := ShootInput{50, 100}
	assert.Equal(i, DeserializeShootInput(SerializeShootInput(i)))
}

func TestClientGameStateTransport(t *testing.T) {
	assert := assert.New(t)
	cgs := ClientGameState{
		You: 50,
		Others: map[ID]Position{
			"id0": 65,
			"id1": 85,
		},
	}
	assert.Equal(cgs, DeserializeClientGameState(SerializeClientGameState(cgs)))
}
