package utils

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
