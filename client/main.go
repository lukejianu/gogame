package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"log"
	"net"

	"github.com/lukejianu/gogame/common"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var port = 8080
var address = fmt.Sprintf(":%d", port)

var red = color.NRGBA{255, 0, 0, 255}

type Game struct {
	positions []common.Position
	writer    io.Writer
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		msg := common.SerializeMoveInput(common.MoveLeftInput)
		_, err := g.writer.Write(msg)
		common.Must(err)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		msg := common.SerializeMoveInput(common.MoveRightInput)
		_, err := g.writer.Write(msg)
		common.Must(err)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, p := range g.positions {
		DrawCircle(screen, p)
	}
}

func DrawCircle(screen *ebiten.Image, x int) {
	vector.FillCircle(screen, float32(x), 250, 25, red, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	conn := connect()
	defer conn.Close()
	g := &Game{
		positions: []common.Position{},
		writer:    common.NewLineWriter(conn),
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
		positions := []common.Position{}
		positions = append(positions, cgs.You)
		for _, p := range cgs.Others {
			positions = append(positions, p)
		}
		g.positions = positions // TODO: Lock.
	}
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", address)
	common.Must(err)
	return conn
}
