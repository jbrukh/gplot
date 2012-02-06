package main

import "fmt"
import "bitbucket.org/binet/go-gnuplot/pkg/gnuplot"

func main() {
	fname := ""
	persist := false
	debug := true

	p,err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p.Close()

	p.CheckedCmd("set grid x")
	p.CheckedCmd("set grid y")
	p.CheckedCmd("set grid z")
	p.PlotXYZ(
		[]float64{0,1,2,3,4,5,6,7,8,9,10},
		[]float64{0,1,2,3,4,5,6,7,8,9,10},
		[]float64{0,1,2,3,4,5,6,7,8,9,10},
		"test 3d plot")
	p.SetLabels("x", "y", "z")
	p.CheckedCmd("set terminal pdf")
	p.CheckedCmd("set output 'plot005.pdf'")
	p.CheckedCmd("replot")


	p.CheckedCmd("q")
	return
}

/* EOF */
