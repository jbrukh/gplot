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
	p.Dual([]float64{1, 2, 4, 8, 16, 32}, []float64{2, 3, 4, 5, 4, 5}, "jake1", "jake2")
}
