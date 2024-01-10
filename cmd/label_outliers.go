package main

import (
	"github.com/jgbaldwinbrown/covplots/pkg"
)

func main() {
	err := covplots.FullLabelOutliers()
	if err != nil { panic(err) }
}
