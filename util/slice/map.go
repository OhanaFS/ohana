package slice

func Map[A any, B any](inp []A, mapper func(in A) B) []B {
	out := make([]B, len(inp))
	for i, v := range inp {
		out[i] = mapper(v)
	}
	return out
}

func Filter[A any](inp []A, filter func(in A) bool) []A {
	var out []A
	for _, v := range inp {
		if filter(v) {
			out = append(out, v)
		}
	}
	return out
}
