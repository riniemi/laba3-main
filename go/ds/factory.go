package ds

import "fmt"

// Create создает пустую структуру по букве типа и возвращает ее как интерфейс Storable.
// M - массив, F - односвязный список, L - двусвязный, S - стек, Q - очередь,
// T - дерево, H - хеш-таблица. С результатом работают одинаково, не зная его реального типа (полиморфизм).
func Create(t byte) (Storable, error) {
	switch t {
	case 'M':
		return NewArray(), nil
	case 'F':
		return NewSinglyList(), nil
	case 'L':
		return NewDoublyList(), nil
	case 'S':
		return NewStack(), nil
	case 'Q':
		return NewQueue(), nil
	case 'T':
		return NewCompleteTree(), nil
	case 'H':
		return NewHashTable(Double), nil
	}
	return nil, fmt.Errorf("%w: structure type %q", ErrUnknown, t)
}
