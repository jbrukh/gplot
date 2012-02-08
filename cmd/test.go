//target:test
package main

import "gplot"

func main() {
    p, err := gplot.NewPlotter("", true, true)
    if err != nil {
        println("problem")
        return
    }
    p.PlotX([]float64{0,1}, "jake")
}
