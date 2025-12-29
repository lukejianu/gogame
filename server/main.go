package main

import (
	"bufio"
	"fmt"
	"maps"
	"net"
	"sync"
	"time"

	"github.com/lukejianu/gogame/common"
)

var broadcastsPerSecond = 5
var timeBetweenBroadcastsMs = time.Duration(1000/broadcastsPerSecond) * time.Millisecond

var port = 8080
var address = fmt.Sprintf(":%d", port)

type serverGameState = map[serverID]common.Position
type strippedServerGameState = map[common.ID]common.Position

type serverID struct {
	id   common.ID
	conn net.Conn
}

func main() {
	l, err := net.Listen("tcp", address)
	common.Must(err)
	defer l.Close()

	idGen := mkIdGenerator()

	mu := sync.Mutex{}
	sgs := serverGameState{}

	go broadcastGameStateOnTicker(sgs, mu)

	for {
		conn, err := l.Accept()
		common.Must(err)
		newId := serverID{idGen(), conn}
		go handleConnection(newId, sgs, mu)
	}
}

func handleConnection(id serverID, sgs serverGameState, mu sync.Mutex) {
	defer id.conn.Close()
	defer func() {
		mu.Lock()
		defer mu.Unlock()
		delete(sgs, id)
	}()
	registerId(id, sgs, mu)
	scan := bufio.NewScanner(id.conn)
	for scan.Scan() {
		b := scan.Bytes()
		msg := common.MustDeserialize(b, common.Message{})
		switch msg.Tag {
		case common.KeyPressMessage:
			mi := common.DeserializeMoveInput(b)
			handleInput(mi, id, sgs, mu)
		case common.MouseClickMessage:
			// TODO: Handle mouse input.
		default:
			panic("bad msg")
		}
	}
}

func registerId(id serverID, sgs serverGameState, mu sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	startPosition := 0
	sgs[id] = startPosition
}

func handleInput(mi common.MoveInput, id serverID, sgs serverGameState, mu sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	switch mi {
	case common.MoveLeftInput:
		sgs[id] -= common.MoveStep
	case common.MoveRightInput:
		sgs[id] += common.MoveStep
	default:
		panic("bad move input")
	}
}

func broadcastGameStateOnTicker(sgs serverGameState, mu sync.Mutex) {
	ticker := time.NewTicker(timeBetweenBroadcastsMs)
	go func() {
		for {
			<-ticker.C
			broadcastGameState(sgs, mu)
		}
	}()
}

func broadcastGameState(sgs serverGameState, mu sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	smallSgs := stripSgs(sgs)
	for id, _ := range sgs {
		cgs := personalizeGs(id.id, smallSgs)
		msg := common.SerializeClientGameState(cgs)
		lineWriter := common.NewLineWriter(id.conn)
		_, err := lineWriter.Write(msg)
		common.Must(err)
	}
}

func stripSgs(sgs serverGameState) map[common.ID]common.Position {
	res := map[common.ID]common.Position{}
	for id, p := range sgs {
		res[id.id] = p
	}
	return res
}

func personalizeGs(id common.ID, smallSgs strippedServerGameState) common.ClientGameState {
	you := smallSgs[id]
	others := maps.Clone(smallSgs)
	delete(others, id)
	return common.ClientGameState{
		You:    you,
		Others: others,
	}
}

func mkIdGenerator() func() common.ID {
	c := 1
	return func() common.ID {
		c += 1
		return fmt.Sprintf("id%d", c)
	}
}
