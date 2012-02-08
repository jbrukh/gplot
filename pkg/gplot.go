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

// temporary file prefix
var gplotPrefix string = "gplot-"

// resolve the gnuplot command on this
// system, or panic if it is not available
func init() {
	var err os.Error
	gnuplot, err = exec.LookPath("gnuplot")
	if err != nil {
		fmt.Printf("could not find path to 'gnuplot' (is it installed?):\n%v\n", err)
	    os.Exit(1)
    }
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// plotterProcess is the structure
// that represents the pipe to gnuplot
type plotterProcess struct {
	handle *exec.Cmd
	stdin  io.WriteCloser
}

// create a new plotterProcess
func newPlotterProcess(persist bool) (*plotterProcess, os.Error) {
	args := []string{}
	if persist {
		args = append(args, "-persist")
	}
	cmd := exec.Command(gnuplot, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	return &plotterProcess{cmd, stdin}, cmd.Start()
}

type tmpfiles map[string]*os.File

type Plotter struct {
	proc     *plotterProcess
	debug    bool
	style    string // current plotting style
	tmpfiles tmpfiles
}

func (self *Plotter) cmd(format string, a ...interface{}) os.Error {
	cmd := fmt.Sprintf(format, a...) + "\n"
	n, err := io.WriteString(self.proc.stdin, cmd)

	if self.debug {
		fmt.Printf("$ %v", cmd)
		fmt.Printf("%v\n", n)
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
	if self.proc != nil && self.proc.handle != nil {
		self.proc.stdin.Close()
		err = self.proc.handle.Wait()
	}
	self.ResetPlot()
	return err
}

const plotCmd = "plot \"-\" binary format=\"%%%%%%%%int%%%%float\" array=2x2 endian=big title \"%s\" with %s"

func (self *Plotter) PlotX(data []float64, title string) (err os.Error) {
    line := fmt.Sprintf(plotCmd, title, self.style)
    println("sending:", line)
    err = self.cmd(line)
    if err != nil {
        return
    }
    for i, v := range data {
        err = binary.Write(self.proc.stdin, binary.BigEndian, []float64{float64(i), v}) 
        if err != nil {
            return
        }
    }
    return
}

var allowed []string = []string{
		"lines",
        "points",
        "linepoints",
		"impulses",
        "dots",
		"steps",
		"errorbars",
		"boxes",
		"boxerrorbars",
		"pm3d",
}

func (self *Plotter) SetStyle(style string) (err os.Error) {

	for _, s := range allowed {
		if s == style {
			self.style = style
			err = nil
			return err
		}
	}

	fmt.Printf("** style '%v' not in allowed list %v\n", style, allowed)
	fmt.Printf("** default to 'points'\n")
	self.style = "points"
	err = os.NewError(fmt.Sprintf("invalid style '%s'", style))

	return err
}

func (self *Plotter) SetXLabel(label string) os.Error {
	return self.cmd(fmt.Sprintf("set xlabel '%s'", label))
}

func (self *Plotter) SetYLabel(label string) os.Error {
	return self.cmd(fmt.Sprintf("set ylabel '%s'", label))
}

func (self *Plotter) SetZLabel(label string) os.Error {
	return self.cmd(fmt.Sprintf("set zlabel '%s'", label))
}

func (self *Plotter) SetLabels(labels ...string) os.Error {
	ndims := len(labels)
	if ndims > 3 || ndims <= 0 {
		return os.NewError(fmt.Sprintf("invalid number of dims '%v'", ndims))
	}
	var err os.Error = nil

	for i, label := range labels {
		switch i {
		case 0:
			ierr := self.SetXLabel(label)
			if ierr != nil {
				err = ierr
				return err
			}
		case 1:
			ierr := self.SetYLabel(label)
			if ierr != nil {
				err = ierr
				return err
			}
		case 2:
			ierr := self.SetZLabel(label)
			if ierr != nil {
				err = ierr
				return err
			}
		}
	}
	return nil
}

func (self *Plotter) ResetPlot() (err os.Error) {
	for fname, fhandle := range self.tmpfiles {
		ferr := fhandle.Close()
		if ferr != nil {
			err = ferr
		}
		os.Remove(fname)
	}
	return err
}

func NewPlotter(fname string, persist, debug bool) (*Plotter, os.Error) {
	p := &Plotter{proc: nil, debug: debug, style: "points"}
	p.tmpfiles = make(tmpfiles)

	if fname != "" {
		panic("NewPlotter with fname is not yet supported")
	} else {
		proc, err := newPlotterProcess(persist)
		if err != nil {
			return nil, err
		}
		p.proc = proc
	}
	return p, nil
}

/* EOF */
