package covplots

import (
	"bufio"
	"io"
	"fmt"
)

func StripOneHeader(r io.Reader) io.Reader {
	b := bufio.NewScanner(r)
	b.Buffer([]byte{}, 1e12)

	return PipeWrite(func(w io.Writer) {
		if !b.Scan() { return }

		for b.Scan() {
			fmt.Fprintln(w, b.Text())
		}
	})
}


func StripHeader(rs []io.Reader, args any) ([]io.Reader, error) {
	var out []io.Reader

	for _, r := range rs {
		out = append(out, StripOneHeader(r))
	}
	return out, nil
}

func StripHeaderSome(rs []io.Reader, args any) ([]io.Reader, error) {
	readers := ToIntSlice(args)
	out := make([]io.Reader, len(rs))
	copy(out, rs)

	for _, idx := range readers {
		out[idx] = StripOneHeader(rs[idx])
	}
	return out, nil
}
