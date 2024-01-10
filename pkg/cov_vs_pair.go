package covplots

import (
	"encoding/json"
	"bufio"
	"math"
	"sort"
	"github.com/jgbaldwinbrown/shellout/pkg"
	"io"
	"os"
	"fmt"
)

func GetPlotFunc(fstr string) func(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	switch fstr {
	case "plot_multi": return PlotMultiAny
	case "plot_multi_pretty": return PlotMultiPrettyAny
	case "plot_multi_pretty_blue": return PlotMultiPrettyBlueAny
	case "plot_multi_pretty_colorseries": return PlotMultiPrettyColorseriesAny
	case "plot_multi_facet": return PlotMultiFacetAny
	case "plot_multi_facet_scales": return PlotMultiFacetScalesAny
	case "plot_multi_facet_scales_boxed": return PlotMultiFacetScalesBoxedAny
	case "plot_multi_facetname_scales": return PlotMultiFacetnameScalesAny
	case "": return PlotMultiAny
	case "fixedorder": return PlotMultiFixedOrderAny
	case "plot_cov_vs_pair": return PlotCovVsPair
	case "plot_self_vs_pair": return PlotSelfVsPair
	case "plot_self_vs_pair_lim": return PlotSelfVsPairLim
	case "plot_self_vs_pair_pretty": return PlotSelfVsPairPretty
	case "plot_self_vs_pair_pretty_fixed": return PlotSelfVsPairPrettyFixed
	case "plot_boxwhisker": return PlotBoxwhisker
	case "plot_cov_hist": return PlotCovHist

	default: return PlotPanic
	}
	return PlotPanic
}
func PlotPanic(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	panic(fmt.Errorf("trying to use an unimplemented plot function"))
	return nil
}

func PlotMultiAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	return PlotMulti(outpre, ylim)
}

func PlotMultiFixedOrderAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	return PlotMultiFixedOrder(outpre, ylim)
}

func UnmarshalJsonOut(jsonOut any, dest any) error {
	buf, err := json.Marshal(jsonOut)
	if err != nil {
		return fmt.Errorf("UnmarshalJsonOut: during Marshal: %w", err)
	}

	err = json.Unmarshal(buf, dest)
	if err != nil {
		return fmt.Errorf("UnmarshalJsonOut: during Unmarshal: %w", err)
	}

	return nil
}

func PlotMultiPrettyAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	h := Handle("PlotMultiPrettyAny: %w")

	var cfg PrettyCfg
	err := UnmarshalJsonOut(args, &cfg)
	if err != nil {
		return h(err)
	}

	err = PlotMultiPretty(outpre, ylim, cfg)
	if err != nil {
		return h(err)
	}

	return nil
}

func PlotMultiPrettyBlueAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var cfg PrettyCfg
	err := UnmarshalJsonOut(args, &cfg)
	if err != nil {
		return err
	}
	return PlotMultiPrettyBlue(outpre, ylim, cfg)
}

func PlotMultiPrettyColorseriesAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var cfg PrettyCfg
	err := UnmarshalJsonOut(args, &cfg)
	if err != nil {
		return err
	}
	return PlotMultiPrettyColorseries(outpre, ylim, cfg)
}

func PlotMultiFacetAny(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	return PlotMultiFacet(outpre, ylim)
}

func PlotCovVsPair(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_cov_vs_pair %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotSelfVsPair(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_self_vs_pair %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotSelfVsPairLim(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var xlim []float64
	err := UnmarshalJsonOut(args, &xlim)
	if err != nil {
		return fmt.Errorf("PlotSelfVsPairLim: %w", err)
	}
	if len(xlim) != 2 {
		return fmt.Errorf("PlotSelfVsPairLim: len(xlim) %v != 2", len(xlim))
	}

	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_self_vs_pair_lim %v %v %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		xlim[0],
		xlim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

type PlotSelfVsPairArgs struct {
	Xmin float64
	Xmax float64
	Ylab string
	Xlab string
	Width float64
	Height float64
	ResScale float64
	TextSize float64
}

func PlotSelfVsPairPretty(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var a PlotSelfVsPairArgs
	err := UnmarshalJsonOut(args, &a)
	if err != nil { return fmt.Errorf("PlotSelfVsPairLim: %w", err) }

	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_self_vs_pair_pretty "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		a.Xmin,
		a.Xmax,
		a.Ylab,
		a.Xlab,
		a.Width,
		a.Height,
		a.ResScale,
		a.TextSize,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotSelfVsPairPrettyFixed(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var a PlotSelfVsPairArgs
	err := UnmarshalJsonOut(args, &a)
	if err != nil { return fmt.Errorf("PlotSelfVsPairLim: %w", err) }

	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_self_vs_pair_pretty_fixed "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		a.Xmin,
		a.Xmax,
		a.Ylab,
		a.Xlab,
		a.Width,
		a.Height,
		a.ResScale,
		a.TextSize,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}


func PosLess(p1, p2 Pos) bool {
	if p1.Chr < p2.Chr {
		return true
	}
	if p1.Chr > p2.Chr {
		return false
	}
	return p1.Bp < p2.Bp
}

func SortAndUniqPoses[T any](maps ...map[Pos]T) []Pos {
	uniques := map[Pos]struct{}{}
	for _, posmap := range maps {
		for pos, _ := range posmap {
			uniques[pos] = struct{}{}
		}
	}
	poses := make([]Pos, 0, len(uniques))
	for pos, _ := range uniques {
		poses = append(poses, pos)
	}
	sort.Slice(poses, func(i, j int) bool{return PosLess(poses[i], poses[j])})
	return poses
}

func CombineToOneLineOld(rs []io.Reader, args any) ([]io.Reader, error) {
	maps := []map[Pos]float64{}
	for _, r := range rs {
		posmap, err := CollectVals(r)
		if err != nil {
			return nil, fmt.Errorf("CombineToOneLine: error during CollectVals %w", err)
		}
		maps = append(maps, posmap)
	}

	// fmt.Fprintf(os.Stderr, "%v\n", maps[1])

	poses := SortAndUniqPoses(maps...)
	// fmt.Fprintf(os.Stderr, "%v\n", poses)

	out := PipeWrite(func(w io.Writer) {
		vals := []float64{}
		for _, pos := range poses {
			vals = vals[:0]
			for _, posmap := range maps {
				if val, ok := posmap[pos]; ok {
					vals = append(vals, val)
				} else {
					vals = append(vals, math.NaN())
				}
			}
			fmt.Fprintf(w, "%v\t%v\t%v", pos.Chr, pos.Bp, pos.Bp+1)
			for _, val := range vals {
				fmt.Fprintf(w, "\t%v", val)
			}
			fmt.Fprintf(w, "\n");
		}
	})
	return []io.Reader{out}, nil
}

type Span struct {
	Chr string
	Start int
	End int
}

type Entry struct {
	Span
	Val float64
}

type Sentry struct {
	Span
	Val string
}

type PosEntry struct {
	Pos
	Val float64
}

func CollectEntry(text string, sl *[]PosEntry) error {
	*sl = (*sl)[:0]
	var e Entry
	_, err := fmt.Sscanf(text, "%s	%d	%d	%f", &e.Chr, &e.Start, &e.End, &e.Val)
	if err != nil {
		return fmt.Errorf("CollectEntry: %w", err)
	}
	for i:=e.Start; i<e.End; i++ {
		*sl = append(*sl, PosEntry{Pos{e.Chr, i}, e.Val})
	}
	return nil
}

func InitNaNSl(length int) []float64 {
	out := make([]float64, length, length)
	for i:=0; i<length; i++ {
		out[i] = math.NaN()
	}
	return out
}

func InitEmptyStringSl(length int) []string {
	out := make([]string, length, length)
	for i:=0; i<length; i++ {
		out[i] = ""
	}
	return out
}

func CollectEntries(posmap map[Pos][]float64, idx, nidx int, ebuffer *[]PosEntry, r io.Reader) error {
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		err := CollectEntry(s.Text(), ebuffer)
		if err != nil {
			return err
		}
		for _, e := range *ebuffer {
			vals, ok := posmap[e.Pos]
			if !ok {
				vals = InitNaNSl(nidx)
				posmap[e.Pos] = vals
			}
			vals[idx] = e.Val
		}
	}
	return nil
}

func CombineToOneLine(rs []io.Reader, args any) ([]io.Reader, error) {
	posmap := map[Pos][]float64{}
	var es []PosEntry
	for i, r := range rs {
		err := CollectEntries(posmap, i, len(rs), &es, r)
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(posmap)
	out := PipeWrite(func(w io.Writer) {
		for pos, vals := range posmap {
			fmt.Fprintf(w, "%v\t%v\t%v", pos.Chr, pos.Bp, pos.Bp+1)
			for _, val := range vals {
				fmt.Fprintf(w, "\t%v", val)
			}
			fmt.Fprintf(w, "\n");
		}
	})
	return []io.Reader{out}, nil
}

func CollectEntriesDumb(posmap map[Span][]string, idx, nidx int, r io.Reader) error {
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		entry, err := CollectEntryDumb(s.Text())
		if err != nil {
			return err
		}
		vals, ok := posmap[entry.Span]
		if !ok {
			vals = InitEmptyStringSl(nidx)
			posmap[entry.Span] = vals
		}
		vals[idx] = entry.Val
	}
	return nil
}

func CollectEntryDumb(text string) (Sentry, error) {
	var e Sentry

	_, err := fmt.Sscanf(text, "%s	%d	%d	%s", &e.Chr, &e.Start, &e.End, &e.Val)
	if err != nil { return e, fmt.Errorf("CollectEntryDumb: line %v: %w", text, err) }

	return e, nil
}

func CombineToOneLineDumb(rs []io.Reader, args any) ([]io.Reader, error) {
	if len(rs) < 1 {
		return []io.Reader{}, nil
	}
	fmt.Println("CombineToOneLineDumb len(rs):", len(rs))

	posmap := map[Span][]string{}
	for i, r := range rs {
		err := CollectEntriesDumb(posmap, i, len(rs), r)
		if err != nil {
			return nil, err
		}
	}

	out := PipeWrite(func(w io.Writer) {
		for span, vals := range posmap {
			fmt.Println("printing span", span, "with vals", vals)
			fmt.Fprintf(w, "%v\t%v\t%v", span.Chr, span.Start, span.End)
			for _, val := range vals {
				fmt.Fprintf(w, "\t%v", val)
			}
			fmt.Fprintf(w, "\n");
		}
	})
	return []io.Reader{out}, nil
}

type PlotCovHistArgs struct {
	Xmin float64
	Xmax float64
	Ylab string
	Xlab string
	Binwidth float64
	Width float64
	Height float64
	ResScale float64
}

func PlotCovHist(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	h := handle("PlotCovHist: %w")

	var hargs PlotCovHistArgs
	err := UnmarshalJsonOut(args, &hargs)
	fmt.Println(hargs)
	if err != nil {
		return h(err)
	}

	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_cov_hist "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		hargs.Xmin,
		hargs.Xmax,
		hargs.Ylab,
		hargs.Xlab,
		hargs.Binwidth,
		hargs.Width,
		hargs.Height,
		hargs.ResScale,
	)
	fmt.Println(script)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

type PlotBoxwhiskerArgs struct {
	Xmin float64
	Xmax float64
	Ylab string
	Xlab string
	Width float64
	Height float64
	ResScale float64
	TextSize float64
	FillName string
}

func PlotBoxwhisker(outpre string, ylim []float64, args any, margs MultiplotPlotFuncArgs) error {
	var a PlotBoxwhiskerArgs
	err := UnmarshalJsonOut(args, &a)
	if err != nil { return fmt.Errorf("PlotSelfVsPairLim: %w", err) }

	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_boxwhisker "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.pdf", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		a.Xmin,
		a.Xmax,
		a.Ylab,
		a.Xlab,
		a.Width,
		a.Height,
		a.ResScale,
		a.TextSize,
		a.FillName,
	)

	fmt.Println(script)
	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

