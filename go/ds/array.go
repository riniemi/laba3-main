package ds

import (
	"strconv"
	"strings"
)

// Array (M) - динамический массив строк с удвоением емкости.
type Array struct {
	data   []string // длина среза равна емкости
	length int
}

// Проверка на этапе компиляции, что *Array реализует интерфейс.
var _ Storable = (*Array)(nil)

// NewArray создает пустой массив с емкостью 4.
func NewArray() *Array {
	return &Array{data: make([]string, 4)}
}

func (a *Array) check(i int) error {
	if i < 0 || i >= a.length {
		return wrap(ErrIndex, strconv.Itoa(i))
	}
	return nil
}

// Insert вставляет значение на позицию idx, O(n).
func (a *Array) Insert(idx int, v string) error {
	if idx < 0 || idx > a.length {
		return wrap(ErrIndex, strconv.Itoa(idx))
	}
	if a.length == len(a.data) { // удвоение емкости: добавление в конец амортизированно O(1)
		nd := make([]string, len(a.data)*2)
		copy(nd, a.data[:a.length])
		a.data = nd
	}
	for i := a.length; i > idx; i-- {
		a.data[i] = a.data[i-1]
	}
	a.data[idx] = v
	a.length++
	return nil
}

// Push добавляет значение в конец, O(1) амортизированно.
func (a *Array) Push(v string) {
	_ = a.Insert(a.length, v) // индекс length всегда допустим
}

// Get возвращает элемент по индексу, O(1).
func (a *Array) Get(i int) (string, error) {
	if err := a.check(i); err != nil {
		return "", err
	}
	return a.data[i], nil
}

// Set заменяет элемент по индексу, O(1).
func (a *Array) Set(i int, v string) error {
	if err := a.check(i); err != nil {
		return err
	}
	a.data[i] = v
	return nil
}

// Erase удаляет элемент и возвращает его, O(n).
func (a *Array) Erase(i int) (string, error) {
	if err := a.check(i); err != nil {
		return "", err
	}
	r := a.data[i]
	for j := i; j+1 < a.length; j++ {
		a.data[j] = a.data[j+1]
	}
	a.length--
	a.data[a.length] = ""
	return r, nil
}

// Type возвращает букву типа.
func (a *Array) Type() byte { return 'M' }

// Size возвращает число элементов.
func (a *Array) Size() int { return a.length }

// Dump возвращает копию элементов.
func (a *Array) Dump() []string {
	r := make([]string, a.length)
	copy(r, a.data[:a.length])
	return r
}

// Load добавляет элементы в конец.
func (a *Array) Load(t []string) error {
	for _, s := range t {
		a.Push(s)
	}
	return nil
}

// Show рисует массив: [a, b, c] (len=3).
func (a *Array) Show() string {
	return "[" + strings.Join(a.Dump(), ", ") + "] (len=" + strconv.Itoa(a.length) + ")"
}
