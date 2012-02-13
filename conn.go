package gplot

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
)

// conn is the structure
// that represents the connection to gnuplot
type conn struct {
	handle *exec.Cmd
	stdin  io.WriteCloser
}

// create a new conn
func newConn(persist bool) (*conn, error) {
	args := make([]string, 0, 10)
	if persist {
		args = append(args, "-persist")
	}
	cmd := exec.Command(gnuplot, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	return &conn{cmd, stdin}, cmd.Start()
}

// close a connection
func (c *conn) closeConn() (err error) {
	if c.handle != nil {
		c.stdin.Close()
		err = c.handle.Wait()
	}
	return err
}

func (c *conn) cmd(format string, a ...interface{}) error {
	command := fmt.Sprintf(format, a...) + "\n"
	_, err := io.WriteString(c.stdin, command)
	return err
}

// data will write binary data to the gp pipe
func (c *conn) data(data interface{}) (err error) {
	err = binary.Write(c.stdin, binary.BigEndian, data)
	return
}
