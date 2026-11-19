package ds

import (
	"strconv"
	"strings"
)

type dnode struct {
	val        string
	prev, next *dnode
}

// DoublyList (L) - двусвязный список строк.
type DoublyList struct {
	head, tail *dnode
	n          int
}

var _ Storable = (*DoublyList)(nil)

// NewDoublyList создает пустой список.
func NewDoublyList() *DoublyList { return &DoublyList{} }

func (l *DoublyList) findNode(v string) *dnode {
	for c := l.head; c != nil; c = c.next {
		if c.val == v {
			return c
		}
	}
	return nil
}

func (l *DoublyList) mustFind(v string) (*dnode, error) {
	t := l.findNode(v)
	if t == nil {
		return nil, wrap(ErrNotFound, v)
	}
	return t, nil
}

// unlink вынимает узел за O(1) - главное преимущество двусвязного списка.
func (l *DoublyList) unlink(x *dnode) string {
	if x.prev != nil {
		x.prev.next = x.next
	} else {
		l.head = x.next
	}
	if x.next != nil {
		x.next.prev = x.prev
	} else {
		l.tail = x.prev
	}
	l.n--
	return x.val
}

func (l *DoublyList) insertAfter(t *dnode, v string) {
	x := &dnode{val: v, prev: t, next: t.next}
	if t.next != nil {
		t.next.prev = x
	} else {
		l.tail = x
	}
	t.next = x
	l.n++
}

// PushHead добавляет в начало, O(1).
func (l *DoublyList) PushHead(v string) {
	x := &dnode{val: v, next: l.head}
	if l.head != nil {
		l.head.prev = x
	} else {
		l.tail = x
	}
	l.head = x
	l.n++
}

// PushTail добавляет в конец, O(1).
func (l *DoublyList) PushTail(v string) {
	x := &dnode{val: v, prev: l.tail}
	if l.tail != nil {
		l.tail.next = x
	} else {
		l.head = x
	}
	l.tail = x
	l.n++
}

// PushAfter вставляет после target; O(n) из-за поиска.
func (l *DoublyList) PushAfter(target, v string) error {
	t, err := l.mustFind(target)
	if err != nil {
		return err
	}
	l.insertAfter(t, v)
	return nil
}

// PushBefore вставляет перед target; O(n) из-за поиска.
func (l *DoublyList) PushBefore(target, v string) error {
	t, err := l.mustFind(target)
	if err != nil {
		return err
	}
	if t == l.head {
		l.PushHead(v)
	} else {
		l.insertAfter(t.prev, v)
	}
	return nil
}

// PopHead удаляет и возвращает первый элемент, O(1).
func (l *DoublyList) PopHead() (string, error) {
	if l.head == nil {
		return "", ErrEmpty
	}
	return l.unlink(l.head), nil
}

// PopTail удаляет и возвращает последний элемент, O(1).
func (l *DoublyList) PopTail() (string, error) {
	if l.head == nil {
		return "", ErrEmpty
	}
	return l.unlink(l.tail), nil
}

// DelAfter удаляет элемент после target.
func (l *DoublyList) DelAfter(target string) (string, error) {
	t, err := l.mustFind(target)
	if err != nil {
		return "", err
	}
	if t.next == nil {
		return "", wrap(ErrNoNeighbor, "after "+target)
	}
	return l.unlink(t.next), nil
}

// DelBefore удаляет элемент перед target.
func (l *DoublyList) DelBefore(target string) (string, error) {
	t, err := l.mustFind(target)
	if err != nil {
		return "", err
	}
	if t.prev == nil {
		return "", wrap(ErrNoNeighbor, "before "+target)
	}
	return l.unlink(t.prev), nil
}

// DelValue удаляет первый узел со значением v; false, если такого нет.
func (l *DoublyList) DelValue(v string) bool {
	x := l.findNode(v)
	if x == nil {
		return false
	}
	l.unlink(x)
	return true
}

// Find возвращает индекс первого узла со значением v или -1.
func (l *DoublyList) Find(v string) int {
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
func (l *DoublyList) At(i int) (string, error) {
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
func (l *DoublyList) Items() []string {
	r := make([]string, 0, l.n)
	for c := l.head; c != nil; c = c.next {
		r = append(r, c.val)
	}
	return r
}

// ItemsReversed возвращает элементы от хвоста к голове (по указателям prev).
func (l *DoublyList) ItemsReversed() []string {
	r := make([]string, 0, l.n)
	for c := l.tail; c != nil; c = c.prev {
		r = append(r, c.val)
	}
	return r
}

// Type возвращает букву типа.
func (l *DoublyList) Type() byte { return 'L' }

// Size возвращает число элементов.
func (l *DoublyList) Size() int { return l.n }

// Dump возвращает элементы в порядке хранения.
func (l *DoublyList) Dump() []string { return l.Items() }

// Load добавляет элементы в хвост.
func (l *DoublyList) Load(t []string) error {
	for _, s := range t {
		l.PushTail(s)
	}
	return nil
}

// Show рисует список с двойными стрелками.
func (l *DoublyList) Show() string {
	if l.head == nil {
		return "NULL"
	}
	return "NULL <-> " + strings.Join(l.Items(), " <-> ") + " <-> NULL"
}
