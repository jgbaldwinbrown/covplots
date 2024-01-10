package main

import (
	"bufio"
	"sort"
	"math"
	"flag"
	"io"
	"encoding/csv"
	"fmt"
	"strconv"
	"os"
)

func Handle(format string) func(...any) error {
	return func(args ...any) error {
		return fmt.Errorf(format, args...)
	}
}

func E(e error) bool {
	return e != nil
}

type BedEntry struct {
	Chr string
	Start int64
	End int64
	Vals []float64
	Fields []string
}

func ParseBedEntry(line []string, col int) (BedEntry, error) {
	h := Handle("ParseBedEntry: %w")

	if len(line) <= col || len(line) < 3 {
		return BedEntry{}, fmt.Errorf("ParseBedEntry: line %v too short", line)
	}

	var b BedEntry
	var e error

	b.Chr = line[0]

	b.Start, e = strconv.ParseInt(line[1], 0, 64)
	if E(e) { return b, h(e) }

	b.End, e = strconv.ParseInt(line[2], 0, 64)
	if E(e) { return b, h(e) }

	val, e := strconv.ParseFloat(line[col], 64)
	if E(e) { return b, h(e) }
	b.Vals = append(b.Vals, val)

	b.Fields = line[3:]

	return b, nil
}

func ReadBed(r io.Reader, col int) ([]BedEntry, error) {
	h := Handle("ReadBed: %w")

	cr := csv.NewReader(r)
	cr.LazyQuotes = true
	cr.Comma = rune('\t')

	var out []BedEntry
	for line, e := cr.Read(); e != io.EOF; line, e = cr.Read() {
		if E(e) { return nil, h(e) }

		entry, e := ParseBedEntry(line, col)
		if E(e) { return nil, h(e) }

		out = append(out, entry)
	}
	return out, nil
}

func Quantile(dist []float64, percentile float64) float64 {
	l := float64(len(dist))
	pos := int(math.Floor(l * percentile))
	return dist[pos]
}

func SortedVals(bed []BedEntry) []float64 {
	vals := make([]float64, len(bed))
	for i, entry := range bed {
		vals[i] = entry.Vals[0]
	}
	sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })
	return vals
}

func FindThresh(bed []BedEntry, thresh float64) float64 {
	dist := SortedVals(bed)
	return Quantile(dist, 1 - thresh)
}

func PosLess(b1, b2 BedEntry) bool {
	if b1.Chr < b2.Chr {
		return true
	}
	if b1.Start < b2.Start {
		return true
	}
	return b1.End < b2.End
}

func FindPeaks(bed []BedEntry, rawThresh float64) []BedEntry {
	var peaks []BedEntry
	for _, entry := range bed {
		if entry.Vals[0] >= rawThresh {
			peaks = append(peaks, entry)
		}
	}
	return peaks
}

func Min(i, j int64) int64 {
	if i < j {
		return i
	}
	return j
}

func Max(i, j int64) int64 {
	if i > j {
		return i
	}
	return j
}

func Intersect(start1, end1, start2, end2 int64) bool {
	left := Max(start1, start2)
	right := Min(end1, end2)
	return left < right
}

func CloseEnough(b1, b2 BedEntry, mergeWidth int) bool {
	if b1.Chr != b2.Chr {
		return false
	}

	s1 := b1.Start - int64(mergeWidth)
	e1 := b1.End + int64(mergeWidth)
	return Intersect(s1, e1, b2.Start, b2.End)
}

func Mean(fs ...float64) float64 {
	sum := 0.0
	for _, f := range fs {
		sum += f
	}
	return sum / float64(len(fs))
}

func Merge(b1, b2 BedEntry) BedEntry {
	var out BedEntry
	out.Chr = b1.Chr
	out.Start = Min(b1.Start, b2.Start)
	out.End = Min(b1.End, b2.End)
	out.Vals = append(out.Vals, math.Max(b1.Vals[0], b2.Vals[0]))
	out.Vals = append(out.Vals, Mean(b1.Vals[0], b2.Vals[0]))
	return out
}

func MergePeaks(bed []BedEntry, mergeWidth int) []BedEntry {
	var out []BedEntry
	if len(bed) < 1 { return out }

	current := bed[0]
	for _, next := range bed[1:] {
		if CloseEnough(current, next, mergeWidth) {
			current = Merge(current, next)
		} else {
			out = append(out, current)
			current = next
		}
	}
	out = append(out, current)
	return out
}

func (b BedEntry) Copy() BedEntry {
	out := b
	out.Vals = make([]float64, len(b.Vals))
	copy(out.Vals, b.Vals)
	out.Fields = make([]string, len(b.Fields))
	copy(out.Fields, b.Fields)
	return out
}

func SetEntryWidth(b BedEntry, width int) BedEntry {
	out := b.Copy()
	midpoint := (out.Start + out.End) / 2
	out.Start = midpoint - int64(width)
	out.End = midpoint + int64(width)
	return out
}

func SetWidths(peaks []BedEntry, width int) []BedEntry {
	var out []BedEntry
	for _, peak := range peaks {
		out = append(out, SetEntryWidth(peak, width))
	}
	return out
}

func WriteBedEntry(w io.Writer, b BedEntry) error {
	h := Handle("WriteBedEntry: %w")

	_, e := fmt.Fprintf(w, "%v\t%v\t%v", b.Chr, b.Start, b.End)
	if E(e) { return h(e) }

	for _, val := range b.Vals {
		_, e = fmt.Fprintf(w, "\t%v", val)
		if E(e) { return h(e) }
	}

	_, e = fmt.Fprintf(w, "\n")
	if E(e) { return h(e) }

	return nil
}

func WriteBed(w io.Writer, bed []BedEntry) error {
	for _, b := range bed {
		e := WriteBedEntry(w, b)
		if E(e) {
			return fmt.Errorf("WriteBedEntry: %w", e)
		}
	}
	return nil
}

func FindBedPeaks(r io.Reader, w io.Writer, width, mergeWidth int, thresh float64, col int) error {
	h := Handle("FindPeaks: %w")

	bed, e := ReadBed(r, col)
	if E(e) { return h(e) }
	sort.Slice(bed, func(i, j int) bool { return PosLess(bed[i], bed[j]) })

	rawThresh := FindThresh(bed, thresh)
	peaks := FindPeaks(bed, rawThresh)
	mpeaks := MergePeaks(peaks, mergeWidth)
	bigPeaks := SetWidths(mpeaks, width)

	e = WriteBed(w, bigPeaks)
	if E(e) { return h(e) }

	return nil
}

func main() {
	widthp := flag.Int("w", 0, "Distance in bp to extend around peak")
	mergewidthp := flag.Int("m", 0, "Distance between peaks such that they are considered one peak")
	threshp := flag.String("t", "", "Peak finding threshold (fraction to keep)")
	colp := flag.Int("c", 3, "0-indexed column containing values to threshold")
	flag.Parse()

	if *threshp == "" {
		panic(fmt.Errorf("Must specify threshold (-t)"))
	}
	thresh, e := strconv.ParseFloat(*threshp, 64)
	if E(e) {
		panic(fmt.Errorf("Could not parse threshold %v", *threshp))
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	e = FindBedPeaks(os.Stdin, os.Stdout, *widthp, *mergewidthp, thresh, *colp)
	if E(e) {
		panic(e)
	}
}
