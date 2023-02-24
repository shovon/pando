package maputils

func Map[K1 comparable, V1 any, K2 comparable, V2 any](
	m map[K1]V1,
	t func(k K1, v V1) (K2, V2),
) map[K2]V2 {
	m2 := make(map[K2]V2)
	for k1, v1 := range m {
		k2, v2 := t(k1, v1)
		m2[k2] = v2
	}
	return m2
}
