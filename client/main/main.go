package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lukejianu/gogame/client"
	"github.com/lukejianu/gogame/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var port = 8080
var address = fmt.Sprintf(":%d", port)

var red = color.NRGBA{255, 0, 0, 255}
var blue = color.NRGBA{0, 0, 255, 255}
var gray = color.NRGBA{100, 100, 100, 255}

type Game struct {
	state common.ClientGameState

	updates []client.TimestampedClientGameState
	mu      sync.Mutex

	interp *client.Interpolater

	writer io.Writer
}

func (g *Game) Update() error {
	handleKeyInput(g)
	updateState(g)
	return nil
}

func handleKeyInput(g *Game) {
	mi, skip := keyToMoveInput()
	if skip {
		return
	}
	predictNewState(g, mi)
	msg := common.SerializeMoveInput(mi)
	_, err := g.writer.Write(msg)
	common.Must(err)
}

func keyToMoveInput() (common.MoveInput, bool) {
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyA) && ebiten.IsKeyPressed(ebiten.KeyD):
		// TODO: Check if ROTMG has this behavior.
		return common.MoveInput(0), true
	case ebiten.IsKeyPressed(ebiten.KeyA):
		return common.MoveLeftInput, false
	case ebiten.IsKeyPressed(ebiten.KeyD):
		return common.MoveRightInput, false
	default:
		return common.MoveInput(0), true
	}
}

func predictNewState(g *Game, mi common.MoveInput) {
	switch mi {
	case common.MoveLeftInput:
		g.state.You -= common.MoveStep
	case common.MoveRightInput:
		g.state.You += common.MoveStep
	default:
		panic("bad move input")
	}
}

func updateState(g *Game) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.updates) < 2 {
		return
	}

	if g.interp == nil {
		lastUpdate, currUpdate := g.updates[0], g.updates[1]
		interp := client.NewGameStateInterpolater(lastUpdate, currUpdate)
		g.interp = &interp
	}

	cgs, interpDone := (*g.interp).Interpolate(time.Now())
	if interpDone {
		g.updates = g.updates[1:]
		g.interp = nil
		return
	}

	g.state.Others = cgs.Others
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
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
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
	g.updates = append(g.updates, client.TimestampCgs(cgs))
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", address)
	common.Must(err)
	return conn
}
