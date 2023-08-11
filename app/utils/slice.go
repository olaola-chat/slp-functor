package utils

import "strconv"

func DistinctUint64Slice(slice []uint64) []uint64 {
	m := make(map[uint64]struct{})
	for _, s := range slice {
		m[s] = struct{}{}
	}
	res := make([]uint64, 0)
	for k := range m {
		res = append(res, k)
	}
	return res
}

func DistinctUint32Slice(slice []uint32) []uint32 {
	m := make(map[uint32]struct{})
	for _, s := range slice {
		m[s] = struct{}{}
	}
	res := make([]uint32, 0)
	for k := range m {
		res = append(res, k)
	}
	return res
}

func StrSliceToUint64Slice(strs []string) []uint64 {
	res := make([]uint64, 0)
	for _, s := range strs {
		id, _ := strconv.Atoi(s)
		res = append(res, uint64(id))
	}
	return res
}

func Uint64SliceToStrSlice(ids []uint64) []string {
	res := make([]string, 0)
	for _, id := range ids {
		res = append(res, strconv.FormatUint(id, 10))
	}
	return res
}
