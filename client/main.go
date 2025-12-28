package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"log"
	"net"
	"sync"

	"github.com/lukejianu/gogame/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var port = 8080
var address = fmt.Sprintf(":%d", port)

var red = color.NRGBA{255, 0, 0, 255}
var blue = color.NRGBA{0, 0, 255, 255}

type Game struct {
	state common.ClientGameState

	updates []common.ClientGameState
	mu      sync.Mutex

	writer io.Writer
}

func (g *Game) Update() error {
	keyInputs := ebiten.InputChars()
	for _, key := range keyInputs {
		g.handleKeyInput(key)
	}
	updateState(g)
	return nil
}

func (g *Game) handleKeyInput(key rune) {
	mi, badKey := keyToMoveInput(key)
	if badKey {
		return
	}
	msg := common.SerializeMoveInput(mi)
	_, err := g.writer.Write(msg)
	common.Must(err)
	// TODO: Predict.
}

func keyToMoveInput(key rune) (common.MoveInput, bool) {
	switch key {
	case 'a':
		return common.MoveLeftInput, false
	case 'd':
		return common.MoveRightInput, false
	default:
		return common.MoveInput(0), true
	}
}

func updateState(g *Game) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if len(g.updates) > 0 {
		g.state = g.updates[0]
		g.updates = g.updates[1:]
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawCircle(screen, g.state.You, blue)
	for _, p := range g.state.Others {
		DrawCircle(screen, p, red)
	}
}

func DrawCircle(screen *ebiten.Image, x int, c color.Color) {
	vector.FillCircle(screen, float32(x), 250, 25, c, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	conn := connect()
	defer conn.Close()
	g := &Game{
		writer: common.NewLineWriter(conn),
	}
	go listenForUpdates(g, conn)
	if err := runGame(g, conn); err != nil {
		log.Fatal(err)
	}
}

func runGame(g *Game, conn net.Conn) error {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Hello, World!")
	return ebiten.RunGame(g)
}

func listenForUpdates(g *Game, conn net.Conn) {
	scan := bufio.NewScanner(conn)
	for scan.Scan() {
		line := scan.Bytes()
		cgs := common.DeserializeClientGameState(line)
		appendUpdate(g, cgs)
	}
}

func appendUpdate(g *Game, cgs common.ClientGameState) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.updates = append(g.updates, cgs)
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", address)
	common.Must(err)
	return conn
}
