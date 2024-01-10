package main

import (
	"github.com/jgbaldwinbrown/covplots/pkg"
)

func main() {
	err := covplots.FullFilterCov()
	if err != nil { panic(err) }
}
