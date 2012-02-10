//target:gplot
package gplot

/*
gplot -- simple gnuplot interface for Go.

This is a fork, update, and rewrite of Binet's go-gnuplot; see the
original project at https://bitbucket.org/binet/go-gnuplot.

Currently, the main improvement that gplot offers over go-gnuplot
is that it pipes data to gnuplot in binary format, which makes it
slightly more appropriate for high-performance plotting.
*/

import (
	"fmt"
	"io"
	"os"
    "exec"
    "encoding/binary"
)

// the system command for gnuplot
var gnuplot string

// conn is the structure
// that represents the connection to gnuplot
type conn struct {
	handle *exec.Cmd
	stdin  io.WriteCloser
}

// create a new conn
func newConn(persist bool) (*conn, os.Error) {
	// TODO: make more efficient
    args := []string{}
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

// resolve the gnuplot command on this
// system, or panic if it is not available
func init() {
	var err os.Error
	gnuplot, err = exec.LookPath("gnuplot")
	if err != nil {
        str := fmt.Sprintf("error finding gnuplot (is it installed?):\n%v\n", err)
		panic(str)
    }
}

func (self *Plotter) PlotX(data []float64, title string) (err os.Error) {
    const plotCmd = "plot \"-\" binary array=%d title \"%s\" with %s"
    line := fmt.Sprintf(plotCmd, len(data), title, self.style)
    err = self.cmd(line)
    if err != nil {
        return
    }
    err = binary.Write(self.conn.stdin, binary.BigEndian, data)
    if err != nil {
          return
    }
    return
}

func (self *Plotter) SetStyle(style string) {
    // TODO: check validity
    self.style = style
}




type Plotter struct {
	conn     *conn
	debug    bool
	style    string // current plotting style
}

func (self *Plotter) cmd(format string, a ...interface{}) os.Error {
	cmd := fmt.Sprintf(format, a...) + "\n"
	n, err := io.WriteString(self.conn.stdin, cmd)

	if self.debug {
		fmt.Printf("cmd$ %v", cmd)
		fmt.Printf("cmd$ <%v\n", n)
	}

	return err
}

func (self *Plotter) CheckedCmd(format string, a ...interface{}) {
	err := self.cmd(format, a...)
	if err != nil {
		str := fmt.Sprintf("error: %v\n", err)
		panic(str)
	}
}

func (self *Plotter) Close() (err os.Error) {
	if self.conn != nil && self.conn.handle != nil {
		self.conn.stdin.Close()
		err = self.conn.handle.Wait()
	}
	return err
}

func (self *Plotter) SetXLabel(label string) os.Error {
	return self.cmd(fmt.Sprintf("set xlabel '%s'", label))
}

func (self *Plotter) SetYLabel(label string) os.Error {
	return self.cmd(fmt.Sprintf("set ylabel '%s'", label))
}

func NewPlotter(persist, debug bool) (plotter *Plotter, err os.Error) {
    const defaults = "set datafile binary format=\"%%float64\" endian=big"
	p := &Plotter{conn: nil, debug: debug, style: "points"}

	conn, err := newConn(persist)
	if err != nil {
		return nil, err
	}
	p.conn = conn
    p.cmd(defaults)
	return p, nil
}
