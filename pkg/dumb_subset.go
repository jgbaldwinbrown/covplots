package covplots

import (
	"flag"
	"os"
	"bufio"
	"io"
	"fmt"
)

func GetSpan(text string) (Span, error) {
	var s Span

	_, err := fmt.Sscanf(text, "%s	%d	%d", &s.Chr, &s.Start, &s.End)
	if err != nil { return s, fmt.Errorf("GetSpan: line %v: %w", text, err) }

	return s, nil
}

func SubsetDumbOne(r io.Reader, spanmap map[Span]struct{}) (io.Reader, error) {
	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	out := PipeWrite(func(w io.Writer) {
		for s.Scan() {
			span, err := GetSpan(s.Text())
			if err != nil {
				panic(err)
			}
			if _, ok := spanmap[span]; ok {
				fmt.Fprintln(w, s.Text())
			}
		}
	})
	return out, nil
}

func GetPathSpanMap(path string) (map[Span]struct{}, error) {
	h := Handle("GetPathSpanMap: %w")

	r, e := OpenMaybeGz(path)
	if e != nil { return nil, h(e) }
	defer r.Close()

	m := map[Span]struct{}{}

	s := bufio.NewScanner(r)
	s.Buffer([]byte{}, 1e12)
	for s.Scan() {
		span, err := GetSpan(s.Text())
		if err != nil { return nil, h(err) }
		m[span] = struct{}{}
	}

	return m, nil
}

func SubsetDumb(rs []io.Reader, args any) ([]io.Reader, error) {
	if len(rs) < 1 {
		return []io.Reader{}, nil
	}

	spanpath, ok := args.(string)
	if !ok {
		return nil, fmt.Errorf("SubsetDumb: args %v is not a string", args)
	}

	spanmap, err := GetPathSpanMap(spanpath)
	if err != nil {
		return nil, fmt.Errorf("SubsetDumb: %w", err)
	}

	var out []io.Reader
	for _, r := range rs {
		outr, err := SubsetDumbOne(r, spanmap)
		if err != nil {
			return nil, fmt.Errorf("SubsetDumb: %w", err)
		}
		out = append(out, outr)
	}
	return out, nil
}

func ToPathAndInts(a any) (string, []int) {
	as := a.([]any)
	if len(as) != 2 {
		panic(fmt.Errorf("ToPathAndInts: len(as) %v != 2", len(as)))
	}

	path := as[0].(string)
	colsa := as[1].([]any)
	var cols []int
	for _, cola := range colsa {
		cols = append(cols, int(cola.(float64)))
	}

	return path, cols
}

func SubsetDumbSome(rs []io.Reader, args any) ([]io.Reader, error) {
	if len(rs) < 1 {
		return []io.Reader{}, nil
	}
	fmt.Println("SubsetDumb len(rs):", len(rs))

	spanpath, cols := ToPathAndInts(args)

	spanmap, err := GetPathSpanMap(spanpath)
	if err != nil {
		return nil, fmt.Errorf("SubsetDumb: %w", err)
	}

	out := make([]io.Reader, len(rs))
	copy(out, rs)

	for _, col := range cols {
		outr, err := SubsetDumbOne(rs[col], spanmap)
		if err != nil {
			return nil, fmt.Errorf("SubsetDumb: %w", err)
		}
		out[col] = outr
	}
	return out, nil
}

func RunDumbSubset() {
	spanpathp := flag.String("s", "", "Subset spans")
	flag.Parse()
	if *spanpathp == "" { panic("missing -s") }

	outarr, err := SubsetDumb([]io.Reader{os.Stdin}, *spanpathp)
	if err != nil { panic(err) }
	if len(outarr) != 1 {
		panic(fmt.Errorf("len(outarr) %v != 1", len(outarr)))
	}

	_, err = io.Copy(os.Stdout, outarr[0])
	if err != nil { panic(err) }
}

