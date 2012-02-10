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
	"os"
    "exec"
)

// the system command for gnuplot
var gnuplot string

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

type Plotter struct {
	conn     *conn
	style    string // current plotting style
}

func NewPlotter(persist bool) (plotter *Plotter, err os.Error) {
    const defaults = "set datafile binary format=\"%%float64\" endian=big"

    p := &Plotter{style: "lines"}
	conn, err := newConn(persist)
	if err != nil {
		return nil, err
	}
	p.conn = conn
    
    err = p.conn.cmd(defaults)
	if err != nil {
        println("could not set binary mode")
        return nil, p.conn.closeConn()
    }
    return p, nil
}

// plot a basic line graph
func (p *Plotter) PlotX(data []float64, title string) (err os.Error) {
    // the default plot command
    const plotCmd = "plot \"-\" binary array=%d title \"%s\" with %s"
    line := fmt.Sprintf(plotCmd, len(data), title, p.style)
    err = p.conn.cmd(line)
    if err != nil {
        return
    }
    p.conn.data(data)
    if err != nil {
        return
    }
    return
}

func (p *Plotter) Close() os.Error {
    return p.conn.closeConn()
}

func (p *Plotter) CheckedCmd(format string, a ...interface{}) {
	err := p.conn.cmd(format, a...)
	if err != nil {
		str := fmt.Sprintf("error: %v\n", err)
		panic(str)
	}
}
