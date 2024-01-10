package covplots

import (
	"os/exec"
	"io"
	"os"
	"fmt"
)

func ToStringSlice(arg any) ([]string, error) {
	h := Handle("ToStringSlice: %w")
	asl, ok := arg.([]any)
	if !ok { return nil, h(fmt.Errorf("arg %v not []any", arg)) }

	ssl := []string{}
	for _, a := range asl {
		s, ok := a.(string)
		if !ok { return nil, h(fmt.Errorf("a %v not slice", a)) }
		ssl = append(ssl, s)
	}
	return ssl, nil
}

func MustStringSlice(arg any) []string {
	s, e := ToStringSlice(arg)
	if e != nil { panic(e) }
	return s
}

func ShellOne(r io.Reader, args []string) (io.Reader, error) {
	h := Handle("ShellOne: %w")
	if len(args) < 1 { return nil, h(fmt.Errorf("len(args) %v < 1", len(args))) }

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = r
	cmd.Stderr = os.Stderr

	out, e := cmd.StdoutPipe()
	if e != nil { return nil, h(e) }

	e = cmd.Start()
	if e != nil { return nil, h(e) }

	return out, nil
}

func Shell(rs []io.Reader, args any) ([]io.Reader, error) {
	h := Handle("Shell: %w")

	argstrs, e := ToStringSlice(args)
	if e != nil { return nil, h(e) }

	var out []io.Reader
	for _, r := range rs {
		outr, e := ShellOne(r, argstrs)
		if e != nil { return nil, h(e) }
		out = append(out, outr)
	}

	return out, nil
}

func ToStrsAndInts(args any) ([]string, []int) {
	argsa := args.([]any)
	if len(argsa) != 2 { panic(fmt.Errorf("ToStrsAndInts: len(argsa) %v != 2", len(argsa))) }

	strs := MustStringSlice(argsa[0])
	ints := ToIntSlice(argsa[1])
	return strs, ints
}

func ShellSome(rs []io.Reader, args any) ([]io.Reader, error) {
	h := Handle("ShellSome: %w")

	argstrs, cols := ToStrsAndInts(args)

	out := make([]io.Reader, len(rs))
	copy(out, rs)

	for _, col := range cols {
		outr, e := ShellOne(rs[col], argstrs)
		if e != nil { return nil, h(e) }
		out[col] = outr
	}
	return out, nil
}
