//target:test
package main

import "gplot"

func main() {
	p, err := gplot.NewPlotter(true)
	if err != nil {
		println("problem")
		return
	}
	defer p.Close()
	p.PlotX([]float64{1, 2, 4, 8, 16, 32}, "jake")
}
