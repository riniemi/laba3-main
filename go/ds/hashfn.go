package ds

// Три простые хеш-функции для строк.

func hashDjb2(s string) uint {
	h := uint(5381)
	for i := 0; i < len(s); i++ {
		h = h*33 + uint(s[i])
	}
	return h
}

func hashFnv1a(s string) uint {
	h := uint(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint(s[i])
		h *= 16777619
	}
	return h
}

func hashSdbm(s string) uint {
	h := uint(0)
	for i := 0; i < len(s); i++ {
		h = uint(s[i]) + (h << 6) + (h << 16) - h
	}
	return h
}
