package gplot

import (
	"fmt"
	"exec"
	"os"
	"io"
	"encoding/binary"
)

// conn is the structure
// that represents the connection to gnuplot
type conn struct {
	handle *exec.Cmd
	in     io.WriteCloser
}

// create a new conn
func newConn(persist bool) (c *conn, err os.Error) {
	args := make([]string, 0, 2)
	if persist {
		args = append(args, "-persist")
	}
	cmd := exec.Command(gnuplot, args...)

	in, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	return &conn{cmd, in}, cmd.Start()
}

// close a connection
func (c *conn) closeConn() (err os.Error) {
	if c.handle != nil {
		c.in.Close()
		err = c.handle.Wait()
	}
	return err
}

func (c *conn) cmd(format string, a ...interface{}) os.Error {
	command := fmt.Sprintf(format, a...) + "\n"
	_, err := io.WriteString(c.in, command)
	return os.NewError(fmt.Sprintf("could not send command: ", err))
}

// data will write binary data to the gp pipe
func (c *conn) data(data interface{}) (err os.Error) {
	err = binary.Write(c.in, binary.BigEndian, data)
	return os.NewError(fmt.Sprintf("could not send data: ", err))
}
