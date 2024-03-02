package services

func GetKeys[K comparable, V any](m map[K]V) []K {
	res := make([]K, 0, len(m))
	for key := range m {
		res = append(res, key)
	}
	return res
}
