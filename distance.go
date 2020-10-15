//go run distance.go -bins 10 -pl java

package main

import (
	"fmt"
	"math"
	"io/ioutil"
	"strings"
	"os"
	"flag"
	plot "./uniplot/histogram"
)


//main 
func main() {
	//load strings file

	bins := flag.Int("bins", 10, "Number of bins")
	name := flag.String("pl", "c#", "Programming language filename")
	flag.Parse()

	fmt.Println("filename:", *name, "\tbins:", *bins)

	data, err := ioutil.ReadFile(*name)
    if err != nil {
        fmt.Println("File reading error", err)
        return
    }
    
	//split strings
	x := string(data)

	strs := strings.Split(x, "\n")

	for i := 0; i < len(strs); i++{
		strs[i] = strings.TrimSpace(strs[i])
	}

	//create aux array histo
	histo := make([]float64, 10)

	t := (len(strs) * (len(strs) - 1)) / 2

	count := 0.0

	//calculate strings distance, and save it on histo
	for i := 0; i < len(strs); i++{
		for j := 0; j < len(strs); j++{
			if i < j{
				s := d(len(strs[i]) - 1, 
					   len(strs[j]) - 1, 
					   strs[i], 
					   strs[j])

				fmt.Print("\r")
				fmt.Print("Loading -> ",math.Ceil(count / float64(t) * 100),"%")

				histo = append(histo, float64(s))
				count++
			}
		}
	}
	
	fmt.Println("\nLoaded")
	fmt.Println("\nvalue - percent - hist - binCount")

	//uniplot library histogram
	hist := plot.Hist(*bins, histo)
	plot.Fprint(os.Stdout, hist, plot.Linear(5))

	mean := mean(histo)
	dev := dev(histo, mean)

	fmt.Println("\nMean:", mean, "\tstdr deviation:", dev)
}

//distance function, returns distance from string a, to string b
func d(i, j int, a, b string) int {
	min_v := 100

	cost := 1
	if a[i] == b[j] {
		cost = 0
	}

	if i == j && j == 0 {
		min_v = 0
	}
	if i > 0 {
		min_v = min(min_v, d(i - 1, j, a, b) + 1)
	}
	if j > 0 {
		min_v = min(min_v, d(i, j - 1, a, b) + 1)
	}
	if i > 0 && j > 0 {
		min_v = min(min_v, d(i - 1, j - 1, a, b) + cost)
	}
	if i > 1 && j > 1 && a[i] == b[j - 1] && a[i - 1] == b[j] {
		min_v = min(min_v, d(i - 2, j - 2, a, b) + 1)
	}

	return min_v

}

//min aux function
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//mean
func mean(dist []float64) float64{
	sum := 0.0
	for _, x := range dist{
		sum += x
	}

	return sum / float64(len(dist))
}

//standard deviation
func dev(dist []float64, mean float64) float64{
	sum := 0.0
	for _, x := range dist{
		sum += math.Pow(x - mean, 2.0)
	}

	return math.Sqrt(sum/float64(len(dist)))
}