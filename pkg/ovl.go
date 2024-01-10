package covplots

import (
	"fmt"
	"io"
)

func HicOvlColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	return GetMultipleCols(rs, []int{0,1,2,18}), nil
}

func HicOvlColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	return GetMultipleColsSome(rs, []int{0,1,2,18}, somereaders), nil
}

func HicNonovlColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicNonovlColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,19}), nil
}

func HicNonovlColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicNonovlColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,19}, somereaders), nil
}

func HicNonovlPropFpkmColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicNonovlPropFpkmColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,25}), nil
}

func HicNonovlPropFpkmColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicNonovlPropFpkmColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,25}, somereaders), nil
}

func HicNonovlPropColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicNonovlPropColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,21}), nil
}

func HicNonovlPropColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicNonovlPropColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,21}, somereaders), nil
}

func HicOvlPropFpkmColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicOvlPropFpkmColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,24}), nil
}

func HicOvlPropFpkmColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicOvlPropFpkmColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,24}, somereaders), nil
}

func HicOvlPropColumns(rs []io.Reader, args any) ([]io.Reader, error) {
	fmt.Printf("HicOvlPropColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleCols(rs, []int{0,1,2,20}), nil
}

func HicOvlPropColumnsSome(rs []io.Reader, args any) ([]io.Reader, error) {
	somereaders := ToIntSlice(args)
	fmt.Printf("HicOvlPropColumns: putting rs %v into GetMultiple Cols\n", rs)
	return GetMultipleColsSome(rs, []int{0,1,2,20}, somereaders), nil
}
