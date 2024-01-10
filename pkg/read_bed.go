package covplots

import (
	"os"
	"strconv"
	"io"
	"encoding/csv"
	"fmt"
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
	Fields []string
}

func ParseBedEntry(line []string) (BedEntry, error) {
	h := Handle("ParseBedEntry: %w")

	if len(line) < 3 {
		return BedEntry{}, fmt.Errorf("ParseBedEntry: line %v too short", line)
	}

	var b BedEntry
	var e error

	b.Chr = line[0]

	b.Start, e = strconv.ParseInt(line[1], 0, 64)
	if E(e) { return b, h(e) }

	b.End, e = strconv.ParseInt(line[2], 0, 64)
	if E(e) { return b, h(e) }

	b.Fields = line[3:]

	return b, nil
}

func ReadBed(r io.Reader) ([]BedEntry, error) {
	h := Handle("ReadBed: %w")

	cr := csv.NewReader(r)
	cr.LazyQuotes = true
	cr.Comma = rune('\t')

	var out []BedEntry
	for line, e := cr.Read(); e != io.EOF; line, e = cr.Read() {
		if E(e) { return nil, h(e) }

		entry, e := ParseBedEntry(line)
		if E(e) { return nil, h(e) }

		out = append(out, entry)
	}
	return out, nil
}

func ReadBedPath(path string) ([]BedEntry, error) {
	h := Handle("ReadBedPath: %w")

	r, e := os.Open(path)
	if E(e) { return nil, h(e) }
	defer r.Close()

	bed, e := ReadBed(r)
	if E(e) { return nil, h(e) }

	return bed, nil
}
