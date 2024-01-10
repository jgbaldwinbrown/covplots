package covplots

import (
	"encoding/csv"
	"regexp"
	"io"
	"os"
	"fmt"
)

type ColSedArgs struct {
	Col int
	Pattern string
	Replace string
}

type ColSedSomeArgs struct {
	Files []int
	Col int
	Pattern string
	Replace string
}

func ColSed(rs []io.Reader, anyargs any) ([]io.Reader, error) {
	h := Handle("ColSed: %w")

	fmt.Fprintln(os.Stderr, "one")
	var args ColSedArgs
	err := UnmarshalJsonOut(anyargs, &args)

	if err != nil {
		return nil, h(err)
	}

	fmt.Fprintln(os.Stderr, "two")
	re, err := regexp.Compile(args.Pattern)
	if err != nil {
		return nil, fmt.Errorf("ChrSed: could not compile pattern %v with error %w", args.Pattern, err)
	}

	fmt.Fprintln(os.Stderr, "three")
	var outs []io.Reader
	for _, r := range rs {
		outs = append(outs, ColSedSingle(r, args.Col, re, args.Replace))
	}

	fmt.Fprintln(os.Stderr, "four")
	return outs, nil
}

func ColSedSome(rs []io.Reader, anyargs any) ([]io.Reader, error) {
	h := Handle("ColSedSome: %w")

	var args ColSedSomeArgs
	err := UnmarshalJsonOut(anyargs, &args)

	if err != nil {
		return nil, h(err)
	}

	re, err := regexp.Compile(args.Pattern)
	if err != nil {
		return nil, fmt.Errorf("ColSed: could not compile pattern %v with error %w", args.Pattern, err)
	}

	out := make([]io.Reader, len(rs))
	for i, r := range rs {
		out[i] = r
	}
	for _, ridx := range args.Files {
		out[ridx] = ColSedSingle(rs[ridx], args.Col, re, args.Replace)
	}
	return out, nil
}

func ColSedSingle(r io.Reader, col int, re *regexp.Regexp, replace string) (io.Reader) {
	h := Handle("ColSedSingle: %w")

	rout := PipeWrite(func(w io.Writer) {
		cr := csv.NewReader(r)
		cr.LazyQuotes = true
		cr.ReuseRecord = true
		cr.FieldsPerRecord = -1
		cr.Comma = rune('\t')

		cw := csv.NewWriter(w)
		cw.Comma = rune('\t')
		defer cw.Flush()

		i := 0
		j := 0

		for l, e := cr.Read() ; e != io.EOF; l, e = cr.Read() {
			if e != nil {
				panic(h(e))
			}
			if len(l) <= col {
				panic(h(fmt.Errorf("len(l) %v <= col %v", len(l), col)))
			}
			l[col] = re.ReplaceAllString(l[col], replace)
			cw.Write(l)
			j++
			i++
		}
		fmt.Fprintf(os.Stderr, "ColSed: printed %v of %v lines\n", j, i)
	})
	return rout
}

