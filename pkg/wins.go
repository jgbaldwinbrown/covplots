package covplots

import (
	"io"
	"fmt"
	"github.com/jgbaldwinbrown/slide/pkg"
)

// func SlidingMeans(inconn io.Reader, outconn io.Writer, size float64, step float64)

type SlidingMeanArgs struct {
	WinSize float64
	WinStep float64
}

func SlidingMean(rs []io.Reader, args any) ([]io.Reader, error) {
	var wargs SlidingMeanArgs
	err := UnmarshalJsonOut(args, &wargs)
	if err != nil {
		return nil, fmt.Errorf("SlidingMeans: %w", err)
	}

	return SlidingMeansInternal(rs, wargs), nil
}

func SlidingMeansInternal(rs []io.Reader, wargs SlidingMeanArgs) []io.Reader {
	var out []io.Reader
	for _, r := range rs {
		outr := OneSlidingMean(r, wargs)
		out = append(out, outr)
	}
	return out
}

func OneSlidingMean(r io.Reader, wargs SlidingMeanArgs) io.Reader {
	return PipeWrite(func(w io.Writer) {
		fmt.Printf("running Log10 internal func\n")
		slide.SlidingMeans(r, w, wargs.WinSize, wargs.WinStep)
	})
}
