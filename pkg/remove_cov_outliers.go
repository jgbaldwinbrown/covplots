package covplots

import (
	"math"
	"os"
	"io"
	"strings"
	"flag"
	"fmt"
	"strconv"
	"github.com/montanaflynn/stats"
)

func FilterBedEntries(filt func(BedEntry) bool, entries ...BedEntry) []BedEntry {
	var out []BedEntry
	for _, e := range entries {
		if filt(e) {
			out = append(out, e)
		}
	}
	return out
}

func GetBedVals(extract func(BedEntry) (float64, error), entries ...BedEntry) ([]float64, error) {
	var out []float64
	for _, entry := range entries {
		val, e := extract(entry)
		if e != nil {
			return nil, fmt.Errorf("GetBedVals: %w", e)
		}
		out = append(out, val)
	}
	return out, nil
}

func GetColFloat(col int) func(BedEntry) (float64, error) {
	return func(b BedEntry) (float64, error) {
		if len(b.Fields) <= col {
			return 0, fmt.Errorf("GetColFloat: b.Fields %v <= col %v", b.Fields, col)
		}
		f, err := strconv.ParseFloat(b.Fields[col], 64)
		if err != nil {
			return 0, fmt.Errorf("GetColFloat: %w", err)
		}
		return f, nil
	}
}

func GetColSums(cols ...int) func(BedEntry) (float64, error) {
	return func(b BedEntry) (float64, error) {
		sum := 0.0
		for _, col := range cols {
			if len(b.Fields) <= col {
				return 0, fmt.Errorf("GetColSums: b.Fields %v <= col %v", b.Fields, col)
			}
			f, err := strconv.ParseFloat(b.Fields[col], 64)
			if err != nil {
				return 0, fmt.Errorf("GetColSums: %w", err)
			}
			sum += f
		}
		return sum, nil
	}
}

func MustGetColSums(cols ...int) func(BedEntry) (float64, error) {
	return func(b BedEntry) (float64, error) {
		sum := 0.0
		for _, col := range cols {
			if len(b.Fields) <= col {
				return math.NaN(), nil
			}
			f, err := strconv.ParseFloat(b.Fields[col], 64)
			if err != nil {
				return math.NaN(), nil
			}
			sum += f
		}
		return sum, nil
	}
}

func FilterByStdevs(col int, stdevs float64, entries ...BedEntry) ([]BedEntry, error) {
	h := Handle("GetBedVals: %w")

	getval := GetColFloat(col)

	vals, err := GetBedVals(getval, entries...)
	if err != nil { return nil, h(err) }

	stdev, err := stats.StandardDeviation(vals)
	if err != nil { return nil, h(err) }

	mean, err := stats.Mean(vals)
	if err != nil { return nil, h(err) }

	lowthresh := mean - (stdev * stdevs)
	hithresh := mean + (stdev * stdevs)

	filt := func(b BedEntry) bool {
		val, err := getval(b)
		if err != nil { return false }
		return (val >= lowthresh) && (val <= hithresh)
	}

	return FilterBedEntries(filt, entries...), nil
	// func Percentile(input Float64Data, percent float64) (percentile float64, err error) {
}

func FilterByColsumStdevs(cols []int, stdevs float64, entries ...BedEntry) ([]BedEntry, error) {
	h := Handle("GetBedVals: %w")

	getval := GetColSums(cols...)

	vals, err := GetBedVals(getval, entries...)
	if err != nil { return nil, h(err) }

	stdev, err := stats.StandardDeviation(vals)
	if err != nil { return nil, h(err) }

	mean, err := stats.Mean(vals)
	if err != nil { return nil, h(err) }

	lowthresh := mean - (stdev * stdevs)
	hithresh := mean + (stdev * stdevs)

	filt := func(b BedEntry) bool {
		val, err := getval(b)
		if err != nil { return false }
		return (val >= lowthresh) && (val <= hithresh)
	}

	return FilterBedEntries(filt, entries...), nil
	// func Percentile(input Float64Data, percent float64) (percentile float64, err error) {
}

func WriteBedEntries(w io.Writer, entries ...BedEntry) error {
	for _, entry := range entries {
		fmt.Fprintf(w, "%v\t%v\t%v", entry.Chr, entry.Start, entry.End)
		for _, field := range entry.Fields {
			fmt.Fprintf(w, "\t%v", field)
		}
		fmt.Fprintf(w, "\n")
	}
	return nil
}

func ReadCovFiltFlags() ([]int, float64) {
	colsp := flag.String("c", "", "comma-separated columns to filter by")
	stdevp := flag.String("s", "", "Standard deviations to keep")
	flag.Parse()
	if *colsp == "" {
		panic(fmt.Errorf("missing column specifier -c"))
	}
	stdevs, err := strconv.ParseFloat(*stdevp, 64)
	if err != nil {
		panic(fmt.Errorf("Could not parse stdev %v", *stdevp))
	}

	colstrs := strings.Split(*colsp, ",")
	var cols []int
	for _, colstr := range colstrs {
		col, err := strconv.ParseInt(colstr, 0, 64)
		if err != nil {
			panic(fmt.Errorf("Could not parse col %v", colstr))
		}
		cols = append(cols, int(col))
	}
	return cols, stdevs
}

func FullFilterCov() error {
	cols, stdevs := ReadCovFiltFlags()

	entries, err := ReadBed(os.Stdin)
	if err != nil { return err }

	fentries, err := FilterByColsumStdevs(cols, stdevs, entries...)
	if err != nil { return err }

	err = WriteBedEntries(os.Stdout, fentries...)
	if err != nil { return err }

	return nil
}
