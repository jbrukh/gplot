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
	debug    bool
	style    string // current plotting style
}

func NewPlotter(persist, debug bool) (plotter *Plotter, err os.Error) {
    const defaults = "set datafile binary format=\"%%float64\" endian=big"
	p := &Plotter{conn: nil, debug: debug, style: "points"}

	conn, err := newConn(persist)
	if err != nil {
		return nil, err
	}
	p.conn = conn
    p.conn.cmd(defaults)
	return p, nil
}

func (p *Plotter) PlotX(data []float64, title string) (err os.Error) {
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

func (p *Plotter) SetStyle(style string) {
    // TODO: check validity
    p.style = style
}






func (p *Plotter) CheckedCmd(format string, a ...interface{}) {
	err := p.conn.cmd(format, a...)
	if err != nil {
		str := fmt.Sprintf("error: %v\n", err)
		panic(str)
	}
}


func (p *Plotter) SetXLabel(label string) os.Error {
	return p.conn.cmd(fmt.Sprintf("set xlabel '%s'", label))
}

func (p *Plotter) SetYLabel(label string) os.Error {
	return p.conn.cmd(fmt.Sprintf("set ylabel '%s'", label))
}

