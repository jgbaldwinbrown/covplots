package covplots

import (
	"os"
	"io"
	"bufio"
	"github.com/jgbaldwinbrown/lscan/pkg"
	"fmt"
	"strconv"
)

func MultiplePerBpNormalize(rs []io.Reader, args any) ([]io.Reader, error) {
	return MultiplePerBp(rs), nil
}

func MultiplePerBp(rs []io.Reader) []io.Reader {
	out := make([]io.Reader, len(rs))
	for i, r := range rs {
		out[i] = PerBp(r)
	}
	return out
}

func ScanInts(dest *[]float64, src []string) error {
	*dest = (*dest)[:0]
	for _, s := range src {
		i, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*dest = append(*dest, i)
	}
	return nil
}

func PerBp(r io.Reader) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Println("running PerBp internal func")
		var line []string
		var floats []float64
		split := lscan.ByByte('\t')

		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		for s.Scan() {
			line = lscan.SplitByFunc(line, s.Text(), split)
			if len(line) < 4 {
				fmt.Fprintf(os.Stderr, "PerBp: line %v too short\n", line)
				continue
			}

			err := ScanInts(&floats, line[1:4])
			if err != nil || len(floats) < 3 {
				fmt.Fprintf(os.Stderr, "PerBp: err %v or len(floats) < 3 %v\n", err, floats)
				continue
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", line[0], line[1], line[2], floats[2] / (floats[1] - floats[0]))
		}
	})
}
