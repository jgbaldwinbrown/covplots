package covplots

import (
	"strings"
	"io"
	"bufio"
	"os"
	"fmt"
	"flag"
)

func GetAllSubtractSingleFlags() AllSingleFlags {
	var f AllSingleFlags
	flag.StringVar(&f.Config, "i", "", "Input config file. Tab-separated columns containing input bed path 1, input bed path 2, chromosome length bed path, and output prefix. Default stdin.")
	flag.IntVar(&f.WinSize, "w", 1000000, "Sliding window plot size (default = 1000000).")
	flag.IntVar(&f.WinStep, "s", 1000000, "Sliding window step distance (default = 1000000).")
	flag.Parse()

	return f
}

func SubtractSinglePlotWinsParallel(cfgs []Config, winsize, winstep, threads int) error {
	jobs := make(chan Config, len(cfgs))
	for _, cfg := range cfgs {
		jobs <- cfg
	}
	close(jobs)

	errs := make(chan error, len(cfgs))

	for i:=0; i<threads; i++ {
		go func() {
			for cfg := range jobs {
				errs <- SubtractSinglePlotWins(cfg.Inpath, cfg.Inpath2, cfg.Chrlenpath, cfg.Outpre, winsize, winstep)
			}
		}()
	}

	var out Errors
	for i:=0; i<len(cfgs); i++ {
		err := <-errs
		if err != nil {
			out = append(out, err)
		}
	}
	if len(out) < 0 {
		return out
	}
	return nil
}


func RunAllSubtractSinglePlots() {
	f := GetAllSubtractSingleFlags()
	cfgs, err := GetConfig(f.Config, true)
	if err != nil {
		panic(err)
	}

	err = SubtractSinglePlotWinsParallel(cfgs, f.WinSize, f.WinStep, 8)
	if err != nil {
		panic(err)
	}
}

func SubtractSinglePlotWins(inpath1, inpath2, chrlenpath, outpre string, winsize, winstep int) error {
	chrlens, err := GetChrLens(chrlenpath)
	if err != nil {
		return fmt.Errorf("SubtactSinglePlotWins: %w", err)
	}

	for _, chrlenset := range chrlens {
		chr, chrlen := chrlenset.Chr, chrlenset.Len
		for start := 0; start < chrlen; start += winstep {
			end := start + winsize
			outpre2 := fmt.Sprintf("%s_%v_%v_%v", outpre, chr, start, end)
			err = SubtractSinglePlotPath(inpath1, inpath2, outpre2, chr, start, end)
			if err != nil {
				return fmt.Errorf("SubtractSinglePlotWins: %w", err)
			}
		}
	}

	return nil
}


func SubtractSinglePlotPath(path1, path2 string, outpre, chr string, start, end int) error {
	r1, err := os.Open(path1)
	if err != nil {
		return fmt.Errorf("SubtractSinglePlotPath: %w", err)
	}
	defer r1.Close()

	r2, err := os.Open(path2)
	if err != nil {
		return fmt.Errorf("SubtractSinglePlotPath: %w", err)
	}
	defer r2.Close()

	err = SubtractSinglePlot(r1, r2, outpre, chr, start, end)
	if err != nil {
		return fmt.Errorf("SubtractSinglePlotPath: %w", err)
	}
	return nil
}

type Pos struct {
	Chr string
	Bp int
}

type SubVal struct {
	Val float64
	Subtracted bool
}

func ParsePosVal(line string, outbuf []PosEntry) ([]PosEntry, error) {
	outbuf = outbuf[:0]
	var chr string
	var start int
	var end int
	var v float64
	_, err := fmt.Sscanf(line, "%s	%d	%d	%f", &chr, &start, &end, &v)
	if err != nil {
		return nil, fmt.Errorf("ParsePosVal: %w", err)
	}
	for i:=start; i<end; i++ {
		outbuf = append(outbuf, PosEntry{Pos{chr, i}, v})
	}
	return outbuf, nil
}

func CollectVals(r io.Reader) (map[Pos]float64, error) {
	out := make(map[Pos]float64)
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		var chr string
		var start int
		var end int
		var v float64
		_, err := fmt.Sscanf(s.Text(), "%s	%d	%d	%f", &chr, &start, &end, &v)
		if err != nil {
			return nil, fmt.Errorf("CollectVals: %w", err)
		}
		for i:=start; i<end; i++ {
			out[Pos{chr, i}] = v
		}
	}
	return out, nil
}

func SubtractInternal(r1, r2 io.Reader) (map[Pos]SubVal, error) {
	out := map[Pos]SubVal{}
	s1 := bufio.NewScanner(r1)
	s1.Buffer([]byte{}, 1e12)
	var posvals []PosEntry
	for s1.Scan() {
		posvals, err := ParsePosVal(s1.Text(), posvals)
		if err != nil {
			return nil, fmt.Errorf("SubtractInternal: %w", err)
		}
		for _, pv := range posvals {
			out[pv.Pos] = SubVal{pv.Val, false}
		}
	}

	s2 := bufio.NewScanner(r2)
	s2.Buffer([]byte{}, 1e12)
	for s2.Scan() {
		posvals, err := ParsePosVal(s2.Text(), posvals)
		if err != nil {
			return nil, fmt.Errorf("SubtractInternal: %w", err)
		}
		for _, pv2 := range posvals {
			if sv1, ok := out[pv2.Pos]; ok {
				out[pv2.Pos] = SubVal{sv1.Val - pv2.Val, true}
			}
		}
	}
	return out, nil
}

func Subtract(r1, r2 io.Reader) (*strings.Reader, error) {
	sub, err := SubtractInternal(r1, r2)
	if err != nil {
		return nil, fmt.Errorf("Subtract: %w", err)
	}

	var out strings.Builder
	for pos, sval := range sub {
		if sval.Subtracted {
			fmt.Fprintf(&out, "%s\t%d\t%d\t%f\n", pos.Chr, pos.Bp, pos.Bp+1, sval.Val)
		}
	}
	return strings.NewReader(out.String()), nil
}

func SubtractOld(r1, r2 io.Reader) (*strings.Reader, error) {
	vals1, err := CollectVals(r1)
	if err != nil {
		return nil, fmt.Errorf("Subtract: %w", err)
	}
	vals2, err := CollectVals(r2)
	if err != nil {
		return nil, fmt.Errorf("Subtract: %w", err)
	}
	sub := make(map[Pos]float64)
	for pos, val := range vals1 {
		if val2, ok := vals2[pos]; ok {
			sub[pos] = val - val2
		}
	}

	var out strings.Builder
	for pos, val := range sub {
		fmt.Fprintf(&out, "%s\t%d\t%d\t%f\n", pos.Chr, pos.Bp, pos.Bp+1, val)
	}
	return strings.NewReader(out.String()), nil
}

func SubtractSinglePlot(r1, r2 io.Reader, outpre, chr string, start, end int) error {
	fr1, err := Filter(r1, chr, start, end)
	if err != nil {
		return err
	}

	fr2, err := Filter(r2, chr, start, end)
	if err != nil {
		return err
	}

	fr3, err := Subtract(fr1, fr2)
	if err != nil {
		return err
	}

	_, err = PlfmtSmall(fr3, outpre, nil, false)
	if err != nil {
		return err
	}

	err = PlotSingle(outpre, true)
	if err != nil {
		return err
	}
	return nil
}

