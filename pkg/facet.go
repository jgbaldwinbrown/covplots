package covplots

import (
	"github.com/jgbaldwinbrown/shellout/pkg"
	"os"
	"fmt"
	"bufio"
	"io"
)

func AddFacetToOneReader(r io.Reader, facetname string) (io.Reader) {
	out := PipeWrite(func(w io.Writer) {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)

		for s.Scan() {
			fmt.Fprintf(w, "%v\t%v\n", s.Text(), facetname)
		}
	})
	return out
}

func AddFacet(rs []io.Reader, args any) ([]io.Reader, error) {
	var out []io.Reader
	anysl, ok := args.([]any)
	if !ok {
		return nil, fmt.Errorf("AddFacet: input args %v not []any", args)
	}

	var facetnames []string
	for _, arg := range anysl {
		name, ok := arg.(string)
		if !ok {
			return nil, fmt.Errorf("AddFacet: input arg %v not string", arg)
		}
		facetnames = append(facetnames, name)
	}

	for i, r := range rs {
		out = append(out, AddFacetToOneReader(r, facetnames[i]))
	}
	return out, nil
}

func PlotMultiFacetScales(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiFacetScales\n")
	fmt.Fprintf(os.Stderr, "PlotMultiFacetScales scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facetscales %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacetScalesAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiFacetScalesAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiFacetScalesAny: args %v == \"\"", args)
	}
	return PlotMultiFacetScales(outpre, scalespath)
}

func PlotMultiFacetnameScales(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiFacetnameScales\n")
	fmt.Fprintf(os.Stderr, "PlotMultiFacetnameScales scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facetname_scales %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacetnameScalesAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiFacetnameScalesAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiFacetnameScalesAny: args %v == \"\"", args)
	}
	return PlotMultiFacetnameScales(outpre, scalespath)
}

type PlotMultiFacetScalesBoxedArgs struct {
	Scales string
	Boxes string
}

func PlfmtPath(inpath, outpre string, margs MultiplotPlotFuncArgs) error {
	h := func(e error) error {
		return fmt.Errorf("PlfmtPath: %w", e)
	}

	fp, err := os.Open(inpath)
	if err != nil {
		return h(err)
	}
	defer fp.Close()
	var r io.Reader = fp

	if !margs.Fullchr {
		r2, err := FilterMulti(margs.Chr, margs.Start, margs.End, r)
		if err != nil {
			return h(err)
		}
		if len(r2) != 1 {
			return h(err)
		}
		defer CloseAny(r2[0])
		r = r2[0]
	}


	data, _, err := PlfmtSmallRead(r, nil, false)
	if err != nil {
		return h(err)
	}

	if err = PlfmtSmallWrite(outpre, data, margs.Plformatter); err != nil {
		return h(err)
	}
	return nil
}

func PlotMultiFacetScalesBoxed(outpre string, args PlotMultiFacetScalesBoxedArgs, margs MultiplotPlotFuncArgs) error {
	h := func(e error) error {
		return fmt.Errorf("PlotMultiFacetScalesBoxed: %w", e)
	}

	boxpre := fmt.Sprintf("%v_boxes", outpre)
	boxpath := fmt.Sprintf("%v_boxes_plfmt.bed", outpre)
	err := PlfmtPath(args.Boxes, boxpre, margs)
	if err != nil {
		return h(err)
	}

	fmt.Fprintf(os.Stderr, "running PlotMultiFacetScalesBoxed\n")
	fmt.Fprintf(os.Stderr, "PlotMultiFacetScalesBoxed scalespath: %v\n", args.Scales);
	fmt.Fprintf(os.Stderr, "PlotMultiFacetScalesBoxed boxes path: %v\n", boxpath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facetscales_boxed %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		args.Scales,
		boxpath,
	)

	err = shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return h(err)
	}
	return nil
}

func PlotMultiFacetScalesBoxedAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var args2 PlotMultiFacetScalesBoxedArgs
	err := UnmarshalJsonOut(args, &args2)
	if err != nil {
		return fmt.Errorf("PlotMultiFacetScalesBoxedAny: %w", err)
	}
	return PlotMultiFacetScalesBoxed(outpre, args2, margs)
}

