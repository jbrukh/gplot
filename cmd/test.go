//target:test
package main

import "gplot"

func main() {
    p, err := gplot.NewPlotter(true, true)
    if err != nil {
        println("problem")
        return
    }
    p.SetStyle("lines")
    p.PlotX([]float64{1, 2}, "jake")
}
