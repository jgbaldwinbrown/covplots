package main

import (
	"flag"
	"os"
	"github.com/jgbaldwinbrown/shellout/pkg"
	"io"
	"fmt"
)

func Plfmt(r io.Reader, outpre, chrbedpath string) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plfmt_flex -c 0 -b 1 -b2 2 -C %v -n > %v
`,
		chrbedpath,
		fmt.Sprintf("%v_plfmt.bed", outpre),
)

	return shellout.ShellPiped(script, r, os.Stdout, os.Stderr)
}

func Plot(outpre string) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_single_cov %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func main() {
	outpre := flag.String("o", "single_cov_plot", "Output prefix")
	chrbedpath := flag.String("C", "", "chromosome lengths bed path")
	flag.Parse()
	if *chrbedpath == "" {
		panic(fmt.Errorf("missing chrbedpath"))
	}
	err := Plfmt(os.Stdin, *outpre, *chrbedpath)
	if err != nil {
		panic(err)
	}
	err = Plot(*outpre)
	if err != nil {
		panic(err)
	}
}
