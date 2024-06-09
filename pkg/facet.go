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

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
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

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
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

	err = shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
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

func PlotMultiTissue(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiTissue\n")
	fmt.Fprintf(os.Stderr, "PlotMultiTissue scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_tissues %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiTissueAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiTissueAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiTissueAny: args %v == \"\"", args)
	}
	return PlotMultiTissue(outpre, scalespath)
}

func PlotMultiRescue(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiRescue\n")
	fmt.Fprintf(os.Stderr, "PlotMultiRescue scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_rescue %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiRescueAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiRescueAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiRescueAny: args %v == \"\"", args)
	}
	return PlotMultiRescue(outpre, scalespath)
}

type PlotMultiSawamuraArgs struct {
	Scales string
	Boxes string
}

func PlotMultiSawamura(outpre string, args PlotMultiSawamuraArgs, margs MultiplotPlotFuncArgs) error {
	h := func(e error) error {
		return fmt.Errorf("PlotMultiSawamura: %w", e)
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

plot_sawamura %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		args.Scales,
		boxpath,
	)

	err = shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return h(err)
	}
	return nil
}

func PlotMultiSawamuraAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var args2 PlotMultiSawamuraArgs
	err := UnmarshalJsonOut(args, &args2)
	if err != nil {
		return fmt.Errorf("PlotMultiSawamuraAny: %w", err)
	}
	return PlotMultiSawamura(outpre, args2, margs)
}

func PlotMultiSawamuraSdist(outpre string, args PlotMultiSawamuraArgs, margs MultiplotPlotFuncArgs) error {
	h := func(e error) error {
		return fmt.Errorf("PlotMultiSawamura: %w", e)
	}

	boxpre := fmt.Sprintf("%v_boxes", outpre)
	boxpath := fmt.Sprintf("%v_boxes_plfmt.bed", outpre)
	err := PlfmtPath(args.Boxes, boxpre, margs)
	if err != nil {
		return h(err)
	}

	fmt.Fprintf(os.Stderr, "running PlotMultiSawamuraSdist\n")
	fmt.Fprintf(os.Stderr, "PlotMultiSawamuraSdist scalespath: %v\n", args.Scales);
	fmt.Fprintf(os.Stderr, "PlotMultiSawamuraSdist boxes path: %v\n", boxpath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_sawamura_sdist %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		args.Scales,
		boxpath,
	)

	err = shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return h(err)
	}
	return nil
}

func PlotMultiSawamuraSdistAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var args2 PlotMultiSawamuraArgs
	err := UnmarshalJsonOut(args, &args2)
	if err != nil {
		return fmt.Errorf("PlotMultiSawamuraSdistAny: %w", err)
	}
	return PlotMultiSawamuraSdist(outpre, args2, margs)
}

func PlotMultiVsill(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiVsill\n")
	fmt.Fprintf(os.Stderr, "PlotMultiVsill scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_vsill %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiVsillAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiVsillAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiVsillAny: args %v == \"\"", args)
	}
	return PlotMultiVsill(outpre, scalespath)
}

func PlotMultiHybrid(outpre string, scalespath string) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiHybrid\n")
	fmt.Fprintf(os.Stderr, "PlotMultiHybrid scalespath: %v\n", scalespath);
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_hybrids %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		scalespath,
	)

	return shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiHybridAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	scalespath, ok := args.(string)
	if !ok {
		return fmt.Errorf("PlotMultiHybridAny: args %v not a string", args)
	}
	if scalespath == "" {
		return fmt.Errorf("PlotMultiHybridAny: args %v == \"\"", args)
	}
	return PlotMultiHybrid(outpre, scalespath)
}


func PlotMultiSawamuraMelcolor(outpre string, args PlotMultiSawamuraArgs, margs MultiplotPlotFuncArgs) error {
	h := func(e error) error {
		return fmt.Errorf("PlotMultiSawamuraMelcolor: %w", e)
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

plot_sawamura_melcolor %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		args.Scales,
		boxpath,
	)

	err = shellout.ShellPiped(script, os.Stdin, os.Stdout, os.Stderr)
	if err != nil {
		return h(err)
	}
	return nil
}

func PlotMultiSawamuraMelcolorAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var args2 PlotMultiSawamuraArgs
	err := UnmarshalJsonOut(args, &args2)
	if err != nil {
		return fmt.Errorf("PlotMultiSawamuraMelcolorAny: %w", err)
	}
	return PlotMultiSawamuraMelcolor(outpre, args2, margs)
}
