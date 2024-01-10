package covplots

import (
	"compress/gzip"
	"io"
	"fmt"
)

func Gunzip(rs []io.Reader, args any) ([]io.Reader, error) {
	var out []io.Reader
	fmt.Printf("Gunzip: %d readers going in\n", len(rs))
	for i, r := range rs {
		fmt.Printf("Gunzip: starting reader %d\n", i)
		outr, err := GunzipOne(r)
		if err != nil {
			return nil, err
		}
		out = append(out, outr)
	}
	return out, nil
}

func GunzipOne(r io.Reader) (io.Reader, error) {
	gr, err := gzip.NewReader(r)
	if err != nil {
		fmt.Println("GunzipOne error:", err)
		return nil, err
	}

	return PipeWrite(func(w io.Writer) {

		defer gr.Close()
		fmt.Printf("running gunzip func\n")

		n, err := io.Copy(w, gr)

		fmt.Printf("gzip wrote %v characters\n", n)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}), nil
}
