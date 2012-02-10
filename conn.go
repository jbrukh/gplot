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
	stdin  io.WriteCloser
}

// create a new conn
func newConn(persist bool) (*conn, os.Error) {
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
func (c *conn) closeConn() (err os.Error) {
	if c.handle != nil {
		c.stdin.Close()
		err = c.handle.Wait()
	}
	return err
}

func (c *conn) cmd(format string, a ...interface{}) os.Error {
	command := fmt.Sprintf(format, a...) + "\n"
	n, err := io.WriteString(c.stdin, command)
    fmt.Printf("cmd$ %v", command)
	fmt.Printf("cmd$ <%v\n", n)
	return err
}

// data will write binary data to the gp pipe
func (c *conn) data(data interface{}) (err os.Error) {
    err = binary.Write(c.stdin, binary.BigEndian, data)
    return
}
