package covplots

import (
	"os"
	"strings"
	"flag"
	"fmt"
	"strconv"
	"github.com/montanaflynn/stats"
)

func ReadPathBed(path string) ([]BedEntry, error) {
	h := Handle("ReadPathBed: %w")

	r, e := os.Open(path)
	if e != nil { return nil, h(e) }
	defer r.Close()

	b, e := ReadBed(r)
	if e != nil { return nil, h(e) }

	return b, nil
}

func GetBedEntryLabels(labeller func(BedEntry) string, entries ...BedEntry) []string {
	var out []string
	for _, e := range entries {
		out = append(out, labeller(e))
	}
	return out
}

func LabelByDistOutliers(cols, distcols []int, outlierPerc, lowOutlierPerc float64, entries []BedEntry, dist []BedEntry) ([]string, error) {
	h := Handle("LbelByDistOutliers: %w")

	getval := MustGetColSums(cols...)
	getdistval := MustGetColSums(distcols...)

	vals, err := GetBedVals(getdistval, dist...)
	if err != nil { return nil, h(err) }

	lowthresh, err := stats.Percentile(vals, lowOutlierPerc)
	if err != nil { return nil, h(err) }

	hithresh, err := stats.Percentile(vals, 100.0 - outlierPerc)
	if err != nil { return nil, h(err) }

	labeller := func(b BedEntry) string {
		val, err := getval(b)
		if err != nil { return "NA" }
		if val > hithresh { return "high" }
		if val < lowthresh { return "low" }
		return "mid"
	}

	return GetBedEntryLabels(labeller, entries...), nil
	// func Percentile(input Float64Data, percent float64) (percentile float64, err error) {
}

func AppendLabels(entries []BedEntry, labels []string) error {
	if len(labels) != len(entries) {
		return fmt.Errorf("AppendLabels: len(labels) %v != len(entries) %v", len(labels), len(entries))
	}

	for i, _ := range entries {
		ep := &entries[i]
		ep.Fields = append(ep.Fields, labels[i])
	}

	return nil
}

type LabellerArgs struct {
	TestCols []int
	TestPath string
	ThreshPerc float64
	LowThreshPerc float64
	DistCols []int
	DistPath string
}

func SplitCommas(in string) []int {
	colstrs := strings.Split(in, ",")
	var cols []int
	for _, colstr := range colstrs {
		col, err := strconv.ParseInt(colstr, 0, 64)
		if err != nil {
			panic(fmt.Errorf("Could not parse col %v", colstr))
		}
		cols = append(cols, int(col))
	}
	return cols
}

func ReadLabellerFlags() LabellerArgs {
	var a LabellerArgs
	var tcstring string
	var dcstring string
	var percstring string
	var lowpercstring string
	flag.StringVar(&tcstring, "tc", "", "comma-separated columns to filter by")
	flag.StringVar(&a.TestPath, "tp", "", "test input path")
	flag.StringVar(&dcstring, "dc", "", "comma-separated columns to establish distribution")
	flag.StringVar(&a.DistPath, "dp", "", "distribution input path")
	flag.StringVar(&percstring, "p", "", "Percent to keep")
	flag.StringVar(&lowpercstring, "lp", "", "Lower percent to keep (if different from high)")
	flag.Parse()

	if a.TestPath == "" { panic(fmt.Errorf("missing -tp")) }
	if a.DistPath == "" { panic(fmt.Errorf("missing -dp")) }

	a.TestCols = SplitCommas(tcstring)
	a.DistCols = SplitCommas(dcstring)

	var err error
	a.ThreshPerc, err = strconv.ParseFloat(percstring, 64)
	if err != nil {
		panic(fmt.Errorf("Could not parse stdev %v", percstring))
	}

	a.LowThreshPerc = a.ThreshPerc
	if lowpercstring != "" {
		a.LowThreshPerc, err = strconv.ParseFloat(lowpercstring, 64)
		if err != nil {
			panic(fmt.Errorf("Could not parse stdev %v", percstring))
		}
	}

	return a
}

func FullLabelOutliers() error {
	args := ReadLabellerFlags()

	testEntries, err := ReadPathBed(args.TestPath)
	if err != nil { return err }

	distEntries, err := ReadPathBed(args.DistPath)
	if err != nil { return err }

	labels, err := LabelByDistOutliers(args.TestCols, args.DistCols, args.ThreshPerc, args.LowThreshPerc, testEntries, distEntries)
	if err != nil { return err }

	AppendLabels(testEntries, labels)

	err = WriteBedEntries(os.Stdout, testEntries...)
	if err != nil { return err }

	return nil
}
