package ds

import "strings"

// Stack (S) - стек строк на основе односвязного списка (вершина - голова списка).
// Список входит в стек по значению (композиция): стек владеет им целиком.
type Stack struct {
	l SinglyList
}

var _ Storable = (*Stack)(nil)

// NewStack создает пустой стек.
func NewStack() *Stack { return &Stack{} }

// Push кладет значение на вершину, O(1).
func (s *Stack) Push(v string) { s.l.PushHead(v) }

// Pop снимает значение с вершины, O(1).
func (s *Stack) Pop() (string, error) {
	if s.l.Size() == 0 {
		return "", wrap(ErrEmpty, "stack")
	}
	return s.l.PopHead()
}

// Items возвращает элементы от вершины к низу.
func (s *Stack) Items() []string { return s.l.Items() }

// Type возвращает букву типа.
func (s *Stack) Type() byte { return 'S' }

// Size возвращает число элементов.
func (s *Stack) Size() int { return s.l.Size() }

// Dump возвращает элементы от низа к вершине (в таком порядке они записываются в файл).
func (s *Stack) Dump() []string {
	v := s.l.Items()
	for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
		v[i], v[j] = v[j], v[i]
	}
	return v
}

// Load кладет элементы по очереди; последний окажется на вершине.
func (s *Stack) Load(t []string) error {
	for _, x := range t {
		s.Push(x)
	}
	return nil
}

// Show рисует стопку с отметкой вершины.
func (s *Stack) Show() string {
	v := s.l.Items()
	if len(v) == 0 {
		return "(empty stack)"
	}
	var b strings.Builder
	for i, x := range v {
		b.WriteString("| " + x + " |")
		if i == 0 {
			b.WriteString("  <- top")
		}
		b.WriteString("\n")
	}
	b.WriteString("+---+")
	return b.String()
}
