package covplots

import (
	"strings"
	"strconv"
	"bufio"
	"flag"
	"os"
	"io"
	"fmt"
)

type AllSingleFlags struct {
	Config string
	WinSize int
	WinStep int
	Threads int
	WholeGenome bool
	SelectWins string
	NoParent bool
}

func GetAllSingleFlags() AllSingleFlags {
	var f AllSingleFlags
	flag.StringVar(&f.Config, "i", "", "Input config file. Tab-separated columns containing input bed path, chromosome length bed path, and output prefix. Default stdin.")
	flag.IntVar(&f.WinSize, "w", 1000000, "Sliding window plot size (default = 1000000).")
	flag.IntVar(&f.WinStep, "s", 1000000, "Sliding window step distance (default = 1000000).")
	flag.Parse()

	return f
}

func SinglePlot(r io.Reader, outpre, chr string, start, end int) error {
	fr, err := Filter(r, chr, start, end)
	if err != nil {
		return err
	}

	if _, err = PlfmtSmall(fr, outpre, nil, false); err != nil {
		return err
	}

	if err = PlotSingle(outpre, false); err != nil {
		return err
	}
	return nil
}

func SinglePlotPath(path string, outpre, chr string, start, end int) error {
	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("SinglePlotPath: %w", err)
	}
	defer r.Close()

	err = SinglePlot(r, outpre, chr, start, end)
	if err != nil {
		return fmt.Errorf("SinglePlotPath: %w", err)
	}
	return nil
}

type Config struct {
	Inpath string
	Chrlenpath string
	Outpre string
	Inpath2 string
}

type ChrLenSet struct {
	Chr string
	Len int
}

func ParseChrLenSet(line string) (ChrLenSet, error) {
	var out ChrLenSet
	sl := strings.Split(line, "\t")
	if len(sl) < 3 {
		return out, fmt.Errorf("line %v has less than 3 fields", sl)
	}

	chrsplit := strings.Split(sl[0], "_")
	if len(chrsplit) < 2 {
		return out, fmt.Errorf("chr field %v less than 2 fields", chrsplit)
	}
	out.Chr = chrsplit[0]
	chrlen, err := strconv.ParseInt(sl[2], 0, 64)
	if err != nil {
		return out, fmt.Errorf("ParseChrLenSet: %w", err)
	}
	out.Len = int(chrlen)
	return out, nil
}

func GetChrLens(chrlenpath string) ([]ChrLenSet, error) {
	r, err := os.Open(chrlenpath)
	if err != nil {
		return nil, fmt.Errorf("GetChrLens: %w", err)
	}
	defer r.Close()
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	chrlens := make(map[string]int)
	for s.Scan() {
		set, err := ParseChrLenSet(s.Text())
		chrlen, ok := chrlens[set.Chr]
		if !ok {
			chrlen = 0
			chrlens[set.Chr] = 0
		}
		if set.Len > chrlen {
			chrlens[set.Chr] = set.Len
		}

		if err != nil {
			return nil, fmt.Errorf("GetChrLens: %w", err)
		}
	}

	var out []ChrLenSet
	for chr, chrlen := range chrlens {
		out = append(out, ChrLenSet{chr, chrlen})
	}
	return out, nil
}

func SinglePlotWins(inpath, chrlenpath, outpre string, winsize, winstep int) error {
	chrlens, err := GetChrLens(chrlenpath)
	if err != nil {
		return fmt.Errorf("SinglePlotWins: %w", err)
	}

	for _, chrlenset := range chrlens {
		chr, chrlen := chrlenset.Chr, chrlenset.Len
		for start := 0; start < chrlen; start += winstep {
			end := start + winsize
			outpre2 := fmt.Sprintf("%s_%v_%v_%v", outpre, chr, start, end)
			err = SinglePlotPath(inpath, outpre2, chr, start, end)
			if err != nil {
				return fmt.Errorf("SinglePlotWins: %w", err)
			}
		}
	}

	return nil
}

func ReadConfig(r io.Reader, subtract bool) ([]Config, error) {
	s := bufio.NewScanner(r)
	out := []Config{}
	for s.Scan() {
		line := strings.Split(s.Text(), "\t")
		if !subtract {
			if len(line) < 3 {
				return nil, fmt.Errorf("ReadConfig: line %v less than 3 fields", line)
			}
			out = append(out, Config{line[0], line[1], line[2], ""})
		} else {
			if len(line) < 4 {
				return nil, fmt.Errorf("ReadConfig: line %v less than 4 fields", line)
			}
			out = append(out, Config{line[0], line[1], line[2], line[3]})
		}
	}
	return out, nil
}

func GetConfig(cfgpath string, subtract bool) ([]Config, error) {
	errfunc := func(e error) error {
		return fmt.Errorf("GetConfig: %w", e)
	}
	var out []Config
	var err error
	if cfgpath == "" {
		out, err = ReadConfig(os.Stdin, subtract)
	} else {
		var r *os.File
		r, err = os.Open(cfgpath)
		if err != nil {
			return nil, errfunc(err)
		}
		defer r.Close()

		out, err = ReadConfig(r, subtract)
	}

	if err != nil {
		return nil, errfunc(err)
	}
	return out, nil
}

type Errors []error

func (e Errors) Error() string {
	var b strings.Builder
	for i, err := range e {
		fmt.Fprintf(&b, "Error %v:\n", i)
		fmt.Fprintln(&b, err)
	}
	return b.String()
}

func SinglePlotWinsParallel(cfgs []Config, winsize, winstep, threads int) error {
	// if threads == 1 {
	// 	for _, cfg := range cfgs {
	// 		err := SinglePlotWins(cfg.Inpath, cfg.Chrlenpath, cfg.Outpre, winsize, winstep)
	// 		if err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }

	jobs := make(chan Config, len(cfgs))
	for _, cfg := range cfgs {
		jobs <- cfg
	}
	close(jobs)

	errs := make(chan error, len(cfgs))

	for i:=0; i<threads; i++ {
		go func() {
			for cfg := range jobs {
				errs <- SinglePlotWins(cfg.Inpath, cfg.Chrlenpath, cfg.Outpre, winsize, winstep)
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

func RunAllSingle() {
	f := GetAllSingleFlags()
	cfgs, err := GetConfig(f.Config, false)
	if err != nil {
		panic(err)
	}

	err = SinglePlotWinsParallel(cfgs, f.WinSize, f.WinStep, 8)
	if err != nil {
		panic(err)
	}
}
