package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/lukejianu/gogame/common"
)

var port = 8080
var address = fmt.Sprintf(":%d", port)

func main() {
	l, err := net.Listen("tcp", address)
	common.Must(err)
	defer l.Close()

	for {
		conn, err := l.Accept()
		common.Must(err)
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	scan := bufio.NewScanner(conn)
	for scan.Scan() {
		fmt.Println(scan.Text())
	}
}

