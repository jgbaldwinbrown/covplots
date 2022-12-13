package covplots

import (
	"bufio"
	"math"
	"sort"
	"github.com/jgbaldwinbrown/shellout/pkg"
	"io"
	"os"
	"fmt"
)

func GetPlotFunc(fstr string) func(outpre string, ylim []float64, args any) error {
	switch fstr {
	case "plot_multi": return PlotMultiAny
	case "plot_multi_facet": return PlotMultiFacetAny
	case "": return PlotMultiAny
	case "plot_cov_vs_pair": return PlotCovVsPair
	default: return PlotPanic
	}
	return PlotPanic
}
func PlotPanic(outpre string, ylim []float64, args any) error {
	panic(fmt.Errorf("trying to use an unimplemented plot function"))
	return nil
}

func PlotMultiAny(outpre string, ylim []float64, args any) error {
	return PlotMulti(outpre, ylim)
}

func PlotMultiFacetAny(outpre string, ylim []float64, args any) error {
	return PlotMultiFacet(outpre, ylim)
}

func PlotCovVsPair(outpre string, ylim []float64, args any) error {
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

	fmt.Fprintf(os.Stderr, "%v\n", maps[1])

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

type Entry struct {
	Chr string
	Start int
	End int
	Val float64
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
