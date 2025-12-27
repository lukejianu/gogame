package common

import (
	"encoding/json"
	"io"
)

type Position = int
type ID = string

type ClientGameState struct {
	You    Position
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
	MoveLeftInput  MoveInput = "a"
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

type lineWriter struct {
	io.Writer
}

func NewLineWriter(w io.Writer) io.Writer {
	return lineWriter{w}
}

func (w lineWriter) Write(p []byte) (int, error) {
	pCopy := make([]byte, len(p), len(p)+1)
	copy(pCopy, p)
	pCopy = append(pCopy, '\n')
	n, err := w.Writer.Write(pCopy)
	return min(len(p), n), err
}
