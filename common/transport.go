package common

import (
	"encoding/json"
	"io"
)

func MustSerialize[T any](msg T) []byte {
	b, err := json.Marshal(msg)
	Must(err)
	return b
}

func MustDeserialize[T any](b []byte, t T) T {
	err := json.Unmarshal(b, &t)
	Must(err)
	return t
}

type MessageType string

const (
	StateUpdateMessage MessageType = "update"
	MoveMessage        MessageType = "move"
	ShootMessage       MessageType = "shoot"
)

type Message struct {
	Tag  MessageType
	Data json.RawMessage
}

type Position = int
type ID = string

type ClientGameState struct {
	You    Position
	Others map[ID]Position
}

func SerializeClientGameState(cgs ClientGameState) []byte {
	msg := Message{
		Tag:  StateUpdateMessage,
		Data: MustSerialize(cgs),
	}
	b := MustSerialize(msg)
	return b
}

func DeserializeClientGameState(b []byte) ClientGameState {
	msg := MustDeserialize(b, Message{})
	switch msg.Tag {
	case StateUpdateMessage:
		return MustDeserialize(msg.Data, ClientGameState{})
	default:
		panic("bad update")
	}
}

type MoveInput string

const (
	MoveLeftInput  MoveInput = "a"
	MoveRightInput MoveInput = "d"
)

func SerializeMoveInput(i MoveInput) []byte {
	msg := Message{
		Tag:  MoveMessage,
		Data: MustSerialize(i),
	}
	b := MustSerialize(msg)
	return b
}

func DeserializeMoveInput(b []byte) MoveInput {
	msg := MustDeserialize(b, Message{})
	switch msg.Tag {
	case MoveMessage:
		return MustDeserialize(msg.Data, MoveInput(""))
	default:
		panic("bad move msg")
	}
}

type ShootInput [2]int // (x, y).

func SerializeShootInput(i ShootInput) []byte {
	msg := Message{
		Tag:  ShootMessage,
		Data: MustSerialize(i),
	}
	b := MustSerialize(msg)
	return b
}

func DeserializeShootInput(b []byte) ShootInput {
	msg := MustDeserialize(b, Message{})
	switch msg.Tag {
	case ShootMessage:
		return MustDeserialize(msg.Data, ShootInput{})
	default:
		panic("bad shoot msg")
	}
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
