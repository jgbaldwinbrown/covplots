package covplots

import (
	"strconv"
	"bufio"
	"github.com/jgbaldwinbrown/lscan/pkg"
	"io"
	"fmt"
	"math"
)

func Log10(rs []io.Reader, args any) ([]io.Reader, error) {
	return OneArgArithMulti(rs, math.Log10), nil
}

func Abs(rs []io.Reader, args any) ([]io.Reader, error) {
	return OneArgArithMulti(rs, math.Abs), nil
}

func Add(rs []io.Reader, args any) ([]io.Reader, error) {
	return TwoArgArithMulti(rs, func(x, y float64) float64 { return x + y }), nil
}

func OneArgArithMulti(rs []io.Reader, f func(float64) float64) []io.Reader {
	var out []io.Reader
	for _, r := range rs {
		outr := OneArgArith(r, f)
		out = append(out, outr)
	}
	return out
}

func OneArgArith(r io.Reader, f func(float64) float64) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Printf("running Log10 internal func\n")
		var line []string
		split := lscan.ByByte('\t')

		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		i := 0
		for s.Scan() {
			line = lscan.SplitByFunc(line, s.Text(), split)
			var val float64
			if len(line) < 3 {
				continue
			}
			if len(line) < 4 {
				val = math.NaN()
			} else {
				var err error
				val, err = strconv.ParseFloat(line[3], 64)
				if err != nil {
					val = math.NaN()
				} else {
					val = f(val)
				}
			}

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", line[0], line[1], line[2], val)
			i++
		}
		fmt.Printf("OneArgArith lines: %v\n", i)
	})
}

func TwoArgArithMulti(rs []io.Reader, f func(float64, float64) float64) []io.Reader {
	var out []io.Reader
	for _, r := range rs {
		outr := TwoArgArith(r, f)
		out = append(out, outr)
	}
	return out
}

func AlwaysParseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return math.NaN()
	}
	return val
}

func TwoArgArith(r io.Reader, f func(float64, float64) float64) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Printf("running Log10 internal func\n")
		var line []string
		split := lscan.ByByte('\t')

		s := bufio.NewScanner(r)
		s.Buffer([]byte{}, 1e12)
		i := 0
		for s.Scan() {
			line = lscan.SplitByFunc(line, s.Text(), split)
			var val0, val1, outval float64
			if len(line) < 3 {
				continue
			}
			if len(line) < 5 {
				val0 = math.NaN()
				val1 = math.NaN()
			} else {
				val0 = AlwaysParseFloat(line[3])
				val1 = AlwaysParseFloat(line[4])
			}
			outval = f(val0, val1)

			fmt.Fprintf(w, "%v\t%v\t%v\t%v\n", line[0], line[1], line[2], outval)
			i++
		}
		fmt.Printf("TwoArgArith lines: %v\n", i)
	})
}
