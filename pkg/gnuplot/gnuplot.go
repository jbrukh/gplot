//
package gnuplot

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
    "exec"
)

var gnuplotCmd string
var gplotPrefix string = "gplot-"

func init() {
	var err os.Error
	gnuplotCmd, err = exec.LookPath("gnuplot")
	if err != nil {
		fmt.Printf("** could not find path to 'gnuplot':\n%v\n", err)
		panic("could not find 'gnuplot'")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type plotterProcess struct {
	handle *exec.Cmd
	stdin  io.WriteCloser
}

func newPlotterProcess(persist bool) (*plotterProcess, os.Error) {
	proc_args := []string{}
	if persist {
		proc_args = append(proc_args, "-persist")
	}
	fmt.Printf("--> [%v] %v\n", gnuplotCmd, proc_args)
	cmd := exec.Command(gnuplotCmd, proc_args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	return &plotterProcess{handle: cmd, stdin: stdin}, cmd.Start()
}

type tmpfiles map[string]*os.File

type Plotter struct {
	proc     *plotterProcess
	debug    bool
	plotcmd  string
	nplots   int    // number of currently active plots
	style    string // current plotting style
	tmpfiles tmpfiles
}

func (self *Plotter) Cmd(format string, a ...interface{}) os.Error {
	cmd := fmt.Sprintf(format, a...) + "\n"
	n, err := io.WriteString(self.proc.stdin, cmd)

	if self.debug {
		fmt.Printf("cmd> %v", cmd)
		fmt.Printf("res> %v\n", n)
	}

	return err
}

func (self *Plotter) CheckedCmd(format string, a ...interface{}) {
	err := self.Cmd(format, a...)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
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

func (self *Plotter) PlotX(data []float64, title string) os.Error {
	f, err := ioutil.TempFile(os.TempDir(), gplotPrefix)
	if err != nil {
		return err
	}
	fname := f.Name()
	self.tmpfiles[fname] = f
	for _, d := range data {
		f.WriteString(fmt.Sprintf("%v\n", d))
	}
	f.Close()
	cmd := self.plotcmd
	//if self.nplots > 0 {
	//	cmd = "replot"
	//}

	var line string
	if title == "" {
		line = fmt.Sprintf("%s \"%s\" with %s", cmd, fname, self.style)
	} else {
		line = fmt.Sprintf("%s \"%s\" title \"%s\" with %s lw .7",
			cmd, fname, title, self.style)
	}
	self.nplots += 1
	return self.Cmd(line)
}

func (self *Plotter) SetPlotCmd(cmd string) (err os.Error) {
	switch cmd {
	case "plot", "splot":
		self.plotcmd = cmd
	default:
		err = os.NewError("invalid plot cmd [" + cmd + "]")
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
	return self.Cmd(fmt.Sprintf("set xlabel '%s'", label))
}

func (self *Plotter) SetYLabel(label string) os.Error {
	return self.Cmd(fmt.Sprintf("set ylabel '%s'", label))
}

func (self *Plotter) SetZLabel(label string) os.Error {
	return self.Cmd(fmt.Sprintf("set zlabel '%s'", label))
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
	self.nplots = 0
	return err
}

func NewPlotter(fname string, persist, debug bool) (*Plotter, os.Error) {
	p := &Plotter{proc: nil, debug: debug, plotcmd: "plot",
		nplots: 0, style: "points"}
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
