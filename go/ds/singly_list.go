package ds

import (
	"strconv"
	"strings"
)

type snode struct {
	val  string
	next *snode
}

// SinglyList (F) - односвязный список строк.
type SinglyList struct {
	head, tail *snode
	n          int
}

var _ Storable = (*SinglyList)(nil)

// NewSinglyList создает пустой список.
func NewSinglyList() *SinglyList { return &SinglyList{} }

func (l *SinglyList) findNode(v string) *snode {
	for c := l.head; c != nil; c = c.next {
		if c.val == v {
			return c
		}
	}
	return nil
}

func (l *SinglyList) mustFind(v string) (*snode, error) {
	t := l.findNode(v)
	if t == nil {
		return nil, wrap(ErrNotFound, v)
	}
	return t, nil
}

// PushHead добавляет в начало, O(1).
func (l *SinglyList) PushHead(v string) {
	l.head = &snode{val: v, next: l.head}
	if l.tail == nil {
		l.tail = l.head
	}
	l.n++
}

// PushTail добавляет в конец, O(1).
func (l *SinglyList) PushTail(v string) {
	x := &snode{val: v}
	if l.tail != nil {
		l.tail.next = x
	} else {
		l.head = x
	}
	l.tail = x
	l.n++
}

// PushAfter вставляет после первого узла со значением target, O(n).
func (l *SinglyList) PushAfter(target, v string) error {
	t, err := l.mustFind(target)
	if err != nil {
		return err
	}
	t.next = &snode{val: v, next: t.next}
	if t == l.tail {
		l.tail = t.next
	}
	l.n++
	return nil
}

// PushBefore вставляет перед первым узлом со значением target, O(n).
func (l *SinglyList) PushBefore(target, v string) error {
	if _, err := l.mustFind(target); err != nil {
		return err
	}
	if l.head.val == target {
		l.PushHead(v)
		return nil
	}
	p := l.head
	for p.next.val != target { // p - узел перед целевым
		p = p.next
	}
	p.next = &snode{val: v, next: p.next}
	l.n++
	return nil
}

// PopHead удаляет и возвращает первый элемент, O(1).
func (l *SinglyList) PopHead() (string, error) {
	if l.head == nil {
		return "", ErrEmpty
	}
	t := l.head
	l.head = t.next
	if l.head == nil {
		l.tail = nil
	}
	l.n--
	return t.val, nil
}

// PopTail удаляет и возвращает последний элемент; нужно найти предпоследний узел, O(n).
func (l *SinglyList) PopTail() (string, error) {
	if l.head == nil {
		return "", ErrEmpty
	}
	if l.head == l.tail {
		return l.PopHead()
	}
	p := l.head
	for p.next != l.tail {
		p = p.next
	}
	r := l.tail.val
	l.tail = p
	p.next = nil
	l.n--
	return r, nil
}

// DelAfter удаляет элемент после target, O(n).
func (l *SinglyList) DelAfter(target string) (string, error) {
	t, err := l.mustFind(target)
	if err != nil {
		return "", err
	}
	if t.next == nil {
		return "", wrap(ErrNoNeighbor, "after "+target)
	}
	x := t.next
	t.next = x.next
	if x == l.tail {
		l.tail = t
	}
	l.n--
	return x.val, nil
}

// DelBefore удаляет элемент перед target, O(n).
func (l *SinglyList) DelBefore(target string) (string, error) {
	t, err := l.mustFind(target)
	if err != nil {
		return "", err
	}
	if t == l.head {
		return "", wrap(ErrNoNeighbor, "before "+target)
	}
	if l.head.next == t {
		return l.PopHead()
	}
	p := l.head
	for p.next.next != t { // p - узел за два до целевого
		p = p.next
	}
	x := p.next
	p.next = t
	l.n--
	return x.val, nil
}

// DelValue удаляет первый узел со значением v; false, если такого нет, O(n).
func (l *SinglyList) DelValue(v string) bool {
	if l.findNode(v) == nil {
		return false
	}
	if l.head.val == v {
		_, _ = l.PopHead()
		return true
	}
	p := l.head
	for p.next.val != v {
		p = p.next
	}
	x := p.next
	p.next = x.next
	if x == l.tail {
		l.tail = p
	}
	l.n--
	return true
}

// Find возвращает индекс первого узла со значением v или -1, O(n).
func (l *SinglyList) Find(v string) int {
	i := 0
	for c := l.head; c != nil; c = c.next {
		if c.val == v {
			return i
		}
		i++
	}
	return -1
}

// At возвращает элемент по индексу, O(n).
func (l *SinglyList) At(i int) (string, error) {
	if i < 0 || i >= l.n {
		return "", wrap(ErrIndex, strconv.Itoa(i))
	}
	c := l.head
	for ; i > 0; i-- {
		c = c.next
	}
	return c.val, nil
}

// Items возвращает элементы от головы к хвосту.
func (l *SinglyList) Items() []string {
	r := make([]string, 0, l.n)
	for c := l.head; c != nil; c = c.next {
		r = append(r, c.val)
	}
	return r
}

// Type возвращает букву типа.
func (l *SinglyList) Type() byte { return 'F' }

// Size возвращает число элементов.
func (l *SinglyList) Size() int { return l.n }

// Dump возвращает элементы в порядке хранения.
func (l *SinglyList) Dump() []string { return l.Items() }

// Load добавляет элементы в хвост.
func (l *SinglyList) Load(t []string) error {
	for _, s := range t {
		l.PushTail(s)
	}
	return nil
}

// Show рисует список со стрелками.
func (l *SinglyList) Show() string {
	if l.head == nil {
		return "NULL"
	}
	return strings.Join(l.Items(), " -> ") + " -> NULL"
}
