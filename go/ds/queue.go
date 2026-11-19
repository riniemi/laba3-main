package ds

import "strings"

// Queue (Q) - очередь строк на основе односвязного списка (начало - голова, конец - хвост).
type Queue struct {
	l SinglyList
}

var _ Storable = (*Queue)(nil)

// NewQueue создает пустую очередь.
func NewQueue() *Queue { return &Queue{} }

// Push добавляет значение в конец, O(1).
func (q *Queue) Push(v string) { q.l.PushTail(v) }

// Pop извлекает значение из начала, O(1).
func (q *Queue) Pop() (string, error) {
	if q.l.Size() == 0 {
		return "", wrap(ErrEmpty, "queue")
	}
	return q.l.PopHead()
}

// Items возвращает элементы от начала к концу.
func (q *Queue) Items() []string { return q.l.Items() }

// Type возвращает букву типа.
func (q *Queue) Type() byte { return 'Q' }

// Size возвращает число элементов.
func (q *Queue) Size() int { return q.l.Size() }

// Dump возвращает элементы от начала к концу.
func (q *Queue) Dump() []string { return q.l.Items() }

// Load добавляет элементы в конец очереди.
func (q *Queue) Load(t []string) error {
	for _, x := range t {
		q.Push(x)
	}
	return nil
}

// Show рисует очередь.
func (q *Queue) Show() string {
	if q.l.Size() == 0 {
		return "(empty queue)"
	}
	return "front -> " + strings.Join(q.l.Items(), ", ") + " <- back"
}
