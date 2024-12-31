package term

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

func ClientTtyNameToPath(clientTtyName string) string {
	return strings.Join(strings.Split(clientTtyName, "_"), "/")
}

func ServerTtyPathToName(serverTtyPath string) string {
	return strings.Join(strings.Split(serverTtyPath, "/"), "_")
}

type ConnectedTty struct {
	tty *os.File
	pty *os.File
	cmd *exec.Cmd
}

var ttys = make(map[string]*ConnectedTty)

func StartTty() (tty string) {
	conTty, err := getTty()
	if err != nil {
		println(err.Error())
	}
	name := conTty.tty.Name()
	ttys[name] = conTty
	return name
}

func getTty() (*ConnectedTty, error) {
	mypty, tty, err := pty.Open()
	if err != nil {
		println("Error could not open tty")
		return nil, err
	}

	println("PTY and TTY opened:", mypty.Name(), tty.Name())
	shell := os.Getenv("SHELL")
	cmd := exec.Command(shell)

	cmd.Stdin = tty
	cmd.Stdout = tty
	cmd.Stderr = tty
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setctty: true,
		Setsid:  true,
	}

	err = pty.Setsize(mypty, &pty.Winsize{
		Rows: 24,
		Cols: 78,
	})

	if err := cmd.Start(); err != nil {
		println("Error starting shell:", err)
		return nil, err
	}
	go func() {
		cmd.Wait()
		os.Exit(0)
	}()

	return &ConnectedTty{
		tty: tty,
		pty: mypty,
		cmd: cmd,
	}, nil
}

func ResizeTty(ttyPath string, col uint16, row uint16) {
	curTty := ttys[ttyPath].tty
	pty.Setsize(curTty, &pty.Winsize{
		Rows: row,
		Cols: col,
	})
}

func HandleConnection(connection *websocket.Conn, ttyName string) {
	connectedTty := ttys[ttyName]
	tty := connectedTty.tty
	mypty := connectedTty.pty

	defer func() {
		connection.Close()
		println("Closing connection")
	}()

	connection.SetCloseHandler(func(code int, text string) error {
		println("Closed TTY: " + tty.Name())
		tty.Close()
		return nil
	})

	// send output from the pty to the WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := mypty.Read(buf)
			if err != nil {
				if err != io.EOF {
					println("Error reading from PTY:", err)
					panic(err)
				}
				break
			}

			// Send the pty data to the WebSocket
			err = connection.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				println("Error writing to WebSocket:", err)
				break
			}
		}
	}()

	// handle websocket input and send it to the pty
	for {
		_, message, err := connection.ReadMessage()
		if err != nil {
			println("Read error:", err.Error())
			break
		}

		_, err = io.Copy(mypty, bytes.NewBuffer(message))
		if err != nil {
			println("Error copying message to pty:", err)
			break
		}
	}
}
