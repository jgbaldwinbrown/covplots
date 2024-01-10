package covplots

import (
	"compress/gzip"
	"regexp"
	"errors"
	"os/exec"
	"math"
	"github.com/montanaflynn/stats"
	"strconv"
	"github.com/jgbaldwinbrown/shellout/pkg"
	"strings"
	"io"
	"bufio"
	"os"
	"fmt"
	"flag"
)

func GetAllMultiplotFlags() AllSingleFlags {
	var f AllSingleFlags
	flag.StringVar(&f.Config, "i", "", "Input config file. JSON, following the documented format.")
	flag.IntVar(&f.WinSize, "w", 1000000, "Sliding window plot size (default = 1000000).")
	flag.IntVar(&f.WinStep, "s", 1000000, "Sliding window step distance (default = 1000000).")
	flag.IntVar(&f.Threads, "t", 8, "Threads to run simultaneously")
	flag.BoolVar(&f.WholeGenome, "g", false, "Generate one plot for the whole genome, no windowing; this overrides all other options")
	// flag.BoolVar(&f.NoParent, "p", false, "Remove parent names from chromosomes")
	flag.StringVar(&f.SelectWins, "c", "", "Plot the windows specified in the provided .bed file path; this overrides sliding window options")
	flag.Parse()

	return f
}

func RunAllMultiplot() {
	fmt.Println("one")
	f := GetAllMultiplotFlags()
	fmt.Println(f)
	cfg, err := GetUltimateConfig(f.Config)
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	var selectWins []BedEntry
	if f.SelectWins != "" {
		selectWins, err = ReadBedPath(f.SelectWins)
		if err != nil {
			panic(err)
		}
	}

	err = AllMultiplotParallel(cfg, f.WinSize, f.WinStep, f.Threads, f.WholeGenome, selectWins)
	if err != nil {
		panic(fmt.Errorf("RunAllMultiplot: %w", err))
	}
}

func FilterMulti(chr string, start, end int, rs ...io.Reader) ([]io.Reader, error) {
	var out []io.Reader
	for _, r := range rs {
		fr, err := Filter(r, chr, start, end)
		if err != nil {
			return nil, fmt.Errorf("FilterMulti: %w", err)
		}
		out = append(out, fr)
	}
	return out, nil
}

func CombineSinglebpPlots(names []string, rs ...io.Reader) (*strings.Reader, error) {
	fmt.Printf("len(rs): %v; names: %v\n", len(rs), names)
	var out strings.Builder
	for i, r := range rs {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		nlines := 0
		for s.Scan() {
			fmt.Fprintf(&out, "%s\t%s\n", s.Text(), names[i])
			nlines++
		}
		fmt.Printf("rs[%v] nlines: %v\n", i, nlines)
	}
	return strings.NewReader(out.String()), nil
}

func PlotMulti(outpre string, ylim []float64) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFixedOrder(outpre string, ylim []float64) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_fixed_order %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

	// xlab = args[5]
	// ylab = args[6]

	// width = as.numeric(args[7])
	// height = as.numeric(args[8])
	// res = as.numeric(args[9])

type PrettyCfg struct {
	Xlab string
	Ylab string
	Width float64
	Height float64
	Res float64
	TextSize float64
}

func PlotMultiPretty(outpre string, ylim []float64, cfg PrettyCfg) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_pretty "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		cfg.Xlab,
		cfg.Ylab,
		cfg.Width,
		cfg.Height,
		cfg.Res,
		cfg.TextSize,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiPrettyBlue(outpre string, ylim []float64, cfg PrettyCfg) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_multi_pretty_blue "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		cfg.Xlab,
		cfg.Ylab,
		cfg.Width,
		cfg.Height,
		cfg.Res,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiPrettyColorseries(outpre string, ylim []float64, cfg PrettyCfg) error {
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_multi_pretty_colorseries "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v" "%v"
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
		cfg.Xlab,
		cfg.Ylab,
		cfg.Width,
		cfg.Height,
		cfg.Res,
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func PlotMultiFacet(outpre string, ylim []float64) error {
	fmt.Fprintf(os.Stderr, "running PlotMultiFacet\n")
	script := fmt.Sprintf(
		`#!/bin/bash
set -e

plot_singlebp_multiline_cov_facet %v %v %v %v
`,
		fmt.Sprintf("%v_plfmt.bed", outpre),
		fmt.Sprintf("%v_plotted.png", outpre),
		ylim[0],
		ylim[1],
	)

	return shellout.ShellOutPiped(script, os.Stdin, os.Stdout, os.Stderr)
}

func Nop([]io.Reader, any) ([]io.Reader, error) {return nil, nil}

func Panic([]io.Reader, any) ([]io.Reader, error) {
	panic(fmt.Errorf("trying to use an unimplemented function"))
	return nil, nil
}

func GetFunc(fstr string) func(rs []io.Reader, args any) ([]io.Reader, error) {
	switch fstr {
	case "add_facet": return AddFacet
	case "subtract_two": return SubtractTwo
	case "dumb_subtract_two": return DumbSubtractTwo
	case "unchanged": return Unchanged
	case "normalize": return Normalize
	case "fourcolumns": return FourColumns
	case "fourcolumns_some": return FourColumnsSome
	case "columns": return Columns
	case "columns_some": return ColumnsSome
	case "hic_self_cols": return HicSelfColumns
	case "hic_self_cols_some": return HicSelfColumnsSome
	case "hic_pair_cols": return HicPairColumns
	case "hic_pair_cols_some": return HicPairColumnsSome
	case "hic_pair_prop_cols": return HicPairPropColumns
	case "hic_pair_prop_cols_some": return HicPairPropColumnsSome
	case "hic_pair_prop_fpkm_cols": return HicPairPropFpkmColumns
	case "hic_pair_prop_fpkm_cols_some": return HicPairPropFpkmColumnsSome
	case "rechr": return ReChr
	case "cov_win_cols": return WindowCovColumns
	case "cov_win_cols_some": return WindowCovColumnsSome
	case "per_bp": return MultiplePerBpNormalize
	case "combine_to_one_line": return CombineToOneLine
	case "combine_to_one_line_dumb": return CombineToOneLineDumb
	case "log10": return Log10
	case "abs": return Abs
	case "add": return Add
	case "gunzip": return Gunzip
	case "chrgrep": return ChrGrep
	case "colgrep": return ColGrep
	case "colgrep_some": return ColGrepSome
	case "colsed": return ColSed
	case "colsed_some": return ColSedSome
	case "sliding_mean": return SlidingMean
	case "strip_header": return StripHeader
	case "strip_header_some": return StripHeaderSome
	case "subset_dumb": return SubsetDumb
	case "subset_dumb_some": return SubsetDumbSome
	case "shell": return Shell
	case "shell_some": return ShellSome
	case "hic_ovl_cols": return HicOvlColumns
	case "hic_ovl_cols_some": return HicOvlColumnsSome
	case "hic_nonovl_cols": return HicNonovlColumns
	case "hic_nonovl_cols_some": return HicNonovlColumnsSome
	case "hic_nonovl_prop_cols": return HicNonovlPropColumns
	case "hic_nonovl_prop_cols_some": return HicNonovlPropColumnsSome
	case "hic_nonovl_prop_fpkm_cols": return HicNonovlPropFpkmColumns
	case "hic_nonovl_prop_fpkm_cols_some": return HicNonovlPropFpkmColumnsSome
	case "hic_ovl_prop_cols": return HicOvlPropColumns
	case "hic_ovl_prop_cols_some": return HicOvlPropColumnsSome
	case "hic_ovl_prop_fpkm_cols": return HicOvlPropFpkmColumns
	case "hic_ovl_prop_fpkm_cols_some": return HicOvlPropFpkmColumnsSome
	default: return Panic
	}
	return Panic
}

type GzReader struct {
	f *os.File
	*gzip.Reader
}

func (r *GzReader) Close() error {
	e1 := r.Reader.Close()
	e2 := r.f.Close()
	if e1 != nil {
		return e1
	}
	if e2 != nil {
		return e2
	}
	return nil
}

func OpenMaybeGz(path string) (io.ReadCloser, error) {
	re := regexp.MustCompile(`\.gz$`)

	r, e := os.Open(path)
	if e != nil {
		return nil, e
	}

	if !re.MatchString(path) {
		return r, nil
	}

	gr, e := gzip.NewReader(r)
	if e != nil {
		r.Close()
		return nil, e
	}

	return &GzReader{f: r, Reader: gr}, nil
}

func OpenPaths(paths ...string) ([]io.Reader, error) {
	fmt.Printf("opening paths %v\n", paths)
	var out []io.Reader
	for _, path := range paths {
		r, err := OpenMaybeGz(path)
		if err != nil {
			CloseAny(out...)
			return nil, err
		}
		out = append(out, r)
	}
	return out, nil
}

func CloseAny[T any](ts ...T) {
	for _, t := range ts {
		a := any(t)
		if c, ok := a.(io.Closer); ok {
			c.Close()
		}
	}
}

func MultiplotInputSet(cfg InputSet, chr string, start, end int, fullchr bool) (io.Reader, []io.Closer, error) {
	rs, err := OpenPaths(cfg.Paths...)
	if err != nil {
		return nil, nil, fmt.Errorf("MultiplotInputSet: during OpenPaths: %w", err)
	}
	var closers []io.Closer
	for _, r := range rs {
		closers = append(closers, r.(io.Closer))
	}

	var frs []io.Reader
	if !fullchr {
		frs, err = FilterMulti(chr, start, end, rs...)
		if err != nil {
			CloseAny(closers...)
			return nil, nil, fmt.Errorf("MultiplotInputSet: during FilterMulti: %w", err)
		}
	} else {
		frs = rs
	}

	for i, funcstr := range cfg.Functions {
		fmt.Println("running", funcstr)
		f := GetFunc(funcstr)
		if len(cfg.FunctionArgs) > i {
			frs, err = f(frs, cfg.FunctionArgs[i])
		} else {
			frs, err = f(frs, nil)
		}
		if err != nil {
			CloseAny(closers...)
			return nil, nil, fmt.Errorf("error when running %v: %w; paths: %v", funcstr, err, cfg.Paths)
		}
	}
	if len(frs) != 1 {
		CloseAny(closers...)
		return nil, nil, fmt.Errorf("Need exactly one reader")
	}


	var out io.Reader = frs[0]
	if !fullchr {
		outs, err := FilterMulti(chr, start, end, frs[0])
		if err != nil {
			CloseAny(closers...)
			return nil, nil, fmt.Errorf("MultiplotInputSet: during FilterMulti 2: %w", err)
		}
		if len(outs) != 1 {
			CloseAny(closers...)
			return nil, nil, fmt.Errorf("Need exactly one reader")
		}
		out = outs[0]
	}


	return out, closers, err
}

func CheckPathExists(path string) bool {
	_, err := os.Stat("/path/to/whatever")
	return !errors.Is(err, os.ErrNotExist)
}

func GzPath(path string, threads int) error {
	cmd := exec.Command("pigz", "-f", "-p", fmt.Sprintf("%d", threads), path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StripParent(r io.Reader) (io.Reader, error) {
	newr := PipeWrite(func(w io.Writer) {
		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		re := regexp.MustCompile(`^([^_	]*)_([^	])*`)
		for s.Scan() {
			line := re.ReplaceAllString(s.Text(), "$1")
			fmt.Fprintln(w, line)
		}
	})
	return newr, nil
}

func GetManualChrs(path string) (chrs []string, err error) {
	h := Handle("GetManualChrs: %w")

	bed, e := ReadBedPath(path)
	if e != nil { return nil, h(e) }

	for _, bede := range bed {
		chrs = append(chrs, bede.Chr)
	}

	return chrs, nil
}

type MultiplotPlotFuncArgs struct {
	Plformatter *Plformatter
	Cfg UltimateConfig
	Chr string
	Start int
	End int
	Fullchr bool
}

func Multiplot(cfg UltimateConfig, chr string, start, end int) error {
	outpre := fmt.Sprintf("%s_%v_%v_%v", cfg.Outpre, chr, start, end)
	var rs []io.Reader
	fullchr := cfg.Fullchr || chr == "full_genome"
	for _, set := range cfg.InputSets {
		r, closers, err := MultiplotInputSet(set, chr, start, end, fullchr)
		if err != nil {
			return fmt.Errorf("Multiplot: during MultiplotInputSet: %w", err)
		}
		defer CloseAny(closers...)
		rs = append(rs, r)
	}

	var names []string
	for _, set := range cfg.InputSets {
		names = append(names, set.Name)
	}

	var combined io.Reader
	combined, err := CombineSinglebpPlots(names, rs...)
	if err != nil {
		return fmt.Errorf("Multiplot: during CombineSinglebpPlots: %w", err)
	}

	if cfg.NoParent {
		combined, err = StripParent(combined)
		if err != nil {
			return fmt.Errorf("Multiplot: during StripParent: %w", err)
		}
	}

	var pf *Plformatter

	if len(cfg.ManualChrs) > 0 {
		pf, err = PlfmtSmall(combined, outpre, cfg.ManualChrs, true)
	} else if cfg.ManualChrsBedPath != "" {
		manualChrs, err := GetManualChrs(cfg.ManualChrsBedPath)
		if err != nil {
			return fmt.Errorf("Multiplot: during GetManualChrs: %w", err)
		}
		pf, err = PlfmtSmall(combined, outpre, manualChrs, true)
	} else {
		pf, err = PlfmtSmall(combined, outpre, nil, false)
	}
	if err != nil {
		return fmt.Errorf("Multiplot: during PlfmtSmall: %w", err)
	}

	ylim := []float64{-300,300}
	if cfg.Ylim != nil {
		ylim = cfg.Ylim
	}

	mPlotFuncArgs := MultiplotPlotFuncArgs{
		Plformatter: pf,
		Cfg: cfg,
		Chr: chr,
		Start: start,
		End: end,
		Fullchr: fullchr,
	}

	plotfunc := GetPlotFunc(cfg.Plotfunc)

	err = plotfunc(outpre, ylim, cfg.PlotfuncArgs, mPlotFuncArgs)
	if err != nil {
		return fmt.Errorf("Multiplot: during plotfunc: %w", err)
	}

	err = GzPath(outpre + "_plfmt.bed", 8)
	if err != nil {
		return fmt.Errorf("Multiplot: during GzPath: %w", err)
	}

	return nil
}

func ParseSubArgs(args any) ([][]int, error) {
	var out [][]int
	anysl, ok := args.([]any)
	if !ok {
		return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
	}
	for _, anypair := range anysl {
		pair, ok := anypair.([]any)
		if !ok {
			return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
		}
		if len(pair) != 2 {
			return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
		}
		entry := make([]int, 2)
		for i, val := range pair {
			ival, ok := val.(int)
			if !ok {
				return nil, fmt.Errorf("ParseSubArgs: parsing of %v failed", args)
			}
			entry[i] = ival
		}
		out = append(out, entry)
	}
	return out, nil
}

func SubtractTwo(rs []io.Reader, args any) ([]io.Reader, error) {
	newreader, err := Subtract(rs[0], rs[1])
	if err != nil {
		return nil, fmt.Errorf("SubtractSome: %w", err)
	}
	return []io.Reader{newreader}, nil
}

func Unchanged(rs []io.Reader, args any) ([]io.Reader, error) {
	if len(rs) != 1 {
		return nil, fmt.Errorf("Unchanged: wrong number of paths (%v)", len(rs))
	}
	return rs, nil
}

func NormalizeFloats(in []float64) []float64 {
	var nanfree []float64
	for _, f := range in {
		if !math.IsNaN(f) {
			nanfree = append(nanfree, f)
		}
	}
	m, err := stats.Mean(nanfree)
	if err != nil {
		m = 0
	}
	s, err := stats.StdDevP(nanfree)
	if err != nil {
		s = 1
	}
	out := make([]float64, len(in))
	for i, f := range in {
		out[i] = (f-m) / s
	}
	return out
}

func Normalize(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Println("normalizing now")
	if len(rs) != 1 {
		return nil, fmt.Errorf("Normalize: wrong number of paths (%v)", len(rs))
	}
	s := bufio.NewScanner(rs[0])
	s.Buffer([]byte{}, 1e12)
	var lines [][]string
	var vals []float64
	for s.Scan() {
		line := strings.Split(s.Text(), "\t")
		if len(line) < 4 {
			return nil, fmt.Errorf("Normalize: line %v has length %v < 4", line, len(line))
		}
		lines = append(lines, line)
		f, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			f = math.NaN()
		}
		vals = append(vals, f)
	}
	vals = NormalizeFloats(vals)
	if len(vals) != len(lines) {
		return nil, fmt.Errorf("Normalize: len(vals) %v != len(lines) %v", len(vals), len(lines))
	}

	var out strings.Builder
	for i, line := range lines {
		line[3] = fmt.Sprintf("%f", vals[i])
		fmt.Fprintln(&out, strings.Join(line, "\t"))
	}
	return []io.Reader{strings.NewReader(out.String())}, nil
}

func MultiplotFullchr(cfg UltimateConfig) error {
	err := Multiplot(cfg, "full_genome", 0, 0)
	if err != nil {
		return fmt.Errorf("MultiplotFullchr: %w", err)
	}

	return nil
}

func MultiplotSelectWins(cfg UltimateConfig, wins []BedEntry) error {
	h := Handle("MultiplotSelectWins: %w")
	fmt.Printf("MultiplotSelectWins: input: %v\n", wins)

	for _, entry := range wins {
		e := Multiplot(cfg, entry.Chr, int(entry.Start), int(entry.End))
		if E(e) { return h(e) }
	}

	return nil
}

func MultiplotSlide(cfg UltimateConfig, winsize, winstep int) error {
	chrlens, err := GetChrLens(cfg.Chrlens)
	if err != nil {
		return fmt.Errorf("MultiplotSlide: %w", err)
	}

	for _, chrlenset := range chrlens {
		chr, chrlen := chrlenset.Chr, chrlenset.Len
		for start := 0; start < chrlen; start += winstep {
			end := start + winsize
			err := Multiplot(cfg, chr, start, end)
			if err != nil {
				return fmt.Errorf("MultiplotSlide loop: %w", err)
			}
		}
	}

	return nil
}

func AllMultiplotParallel(cfgs []UltimateConfig, winsize, winstep, threads int, fullgenome bool, selectWins []BedEntry) error {
	jobs := make(chan UltimateConfig, len(cfgs))
	for _, cfg := range cfgs {
		jobs <- cfg
	}
	close(jobs)

	errs := make(chan error, len(cfgs))

	for i:=0; i<threads; i++ {
		go func() {
			for cfg := range jobs {
				if cfg.Fullchr || fullgenome {
					errs <- MultiplotFullchr(cfg)
				} else if selectWins != nil {
					errs <- MultiplotSelectWins(cfg, selectWins)
				} else {
					errs <- MultiplotSlide(cfg, winsize, winstep)
				}
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
	if len(out) > 0 {
		return fmt.Errorf("AllMultiplotParallel: %w", out)
	}

	fmt.Println("done with parallel")
	return nil
}


