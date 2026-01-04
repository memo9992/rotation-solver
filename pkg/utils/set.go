package utils

type Set[T comparable] struct {
	s map[T]struct{}
}

func NewSet[T comparable](t ...T) *Set[T] {
	s := &Set[T]{
		s: map[T]struct{}{},
	}
	for _, item := range t {
		s.s[item] = struct{}{}
	}

	return s
}

func (s *Set[T]) Add(t T) bool {
	out := true
	_, ok := s.s[t]
	if ok {
		out = false
	}
	s.s[t] = struct{}{}
	return out
}

func (s *Set[T]) Remove(t T) bool {
	_, ok := s.s[t]
	if ok {
		delete(s.s, t)
	}
	return ok
}

func (s *Set[T]) Contains(t T) bool {
	_, ok := s.s[t]
	return ok
}
