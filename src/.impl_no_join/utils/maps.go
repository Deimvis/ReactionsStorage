package utils

// FilterMapIn filters map in-place.
// Keeps only those elements for which given filFn returns true.
func FilterMapIn[K comparable, V any](m map[K]V, filFn func(k K, v V) bool) {
	for k, v := range m {
		if !filFn(k, v) {
			delete(m, k)
		}
	}
}
