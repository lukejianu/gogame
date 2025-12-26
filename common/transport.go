package common

import (
	"encoding/json"
)

type Position = int
type ID = string

type ClientGameState struct {
	You Position
	Others map[ID]Position
}

func SerializeClientGameState(cgs ClientGameState) []byte {
	b, err := json.Marshal(cgs)
	Must(err)
	return b
}

func DeserializeClientGameState(b []byte) ClientGameState {
	cgs := ClientGameState{}
	err := json.Unmarshal(b, &cgs)
	Must(err)
	return cgs
}

type MoveInput string

const (
	MoveLeftInput MoveInput = "a"
	MoveRightInput MoveInput = "d"
)

// TODO: See if you can de-duplicate by creating a generic SerializeT.
func SerializeMoveInput(i MoveInput) []byte {
	b, err := json.Marshal(i)
	Must(err)
	return b
}

func DeserializeMoveInput(b []byte) MoveInput {
	i := MoveInput("")
	err := json.Unmarshal(b, &i)
	Must(err)
	return i
}

