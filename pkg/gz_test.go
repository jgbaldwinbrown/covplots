package covplots

import (
	"io"
	"testing"
	"strings"
)

// all_singlebp_multiline.go:func OpenMaybeGz(path string) (io.ReadCloser, error) {

// const input := `Hello world!
// hi
// hi
// `
// 
// func gzipped() string {
// 	var b strings.Builder
// 	r := strings.NewReader(input)
// 	w := gzip.NewWriter(&b)
// 
// 	e := io.Copy(w, r)
// 	if e != nil { panic(e) }
// 
// 	w.Flush()
// 	return b.String()
// }


func TestOpenMaybeGz(t *testing.T) {
	expect := `Hello world!
hi
hi
`
	r, e := OpenMaybeGz("a.txt")
	if e != nil { panic(e) }
	defer r.Close()

	var b strings.Builder
	_, e = io.Copy(&b, r)
	if e != nil { panic(e) }

	out := b.String()

	if out != expect {
		t.Errorf("out %v != expect %v", out, expect)
	}
}

func TestOpenMaybeGz2(t *testing.T) {
	expect := `Hello world!
hi
hi
`
	r, e := OpenMaybeGz("a.txt.gz")
	if e != nil { panic(e) }
	defer r.Close()

	var b strings.Builder
	_, e = io.Copy(&b, r)
	if e != nil { panic(e) }

	out := b.String()

	if out != expect {
		t.Errorf("out %v != expect %v", out, expect)
	}
}

func TestHeader(t *testing.T) {
	expect := `hi
hi
`
	r, e := OpenMaybeGz("a.txt.gz")
	if e != nil { panic(e) }
	defer r.Close()

	rs2, e := StripHeader([]io.Reader{r}, nil)
	if e != nil { panic(e) }
	defer CloseAny(rs2...)

	var b strings.Builder
	_, e = io.Copy(&b, rs2[0])
	if e != nil { panic(e) }

	out := b.String()

	if out != expect {
		t.Errorf("out %v != expect %v", out, expect)
	}
}

func TestOneHeader(t *testing.T) {
	expect := `hi
hi
`
	r, e := OpenMaybeGz("a.txt.gz")
	if e != nil { panic(e) }
	defer r.Close()

	r2 := StripOneHeader(r)
	defer r2.(io.ReadCloser).Close()

	var b strings.Builder
	_, e = io.Copy(&b, r2)
	if e != nil { panic(e) }

	out := b.String()

	if out != expect {
		t.Errorf("out %v != expect %v", out, expect)
	}
}
