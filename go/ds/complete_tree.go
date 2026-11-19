package ds

import (
	"sort"
	"strconv"
	"strings"
)

type tnode struct {
	key         string
	left, right *tnode
}

// CompleteTree (T) - полное двоичное дерево. Хранит только целые числа (в том числе отрицательные).
//
// Все уровни заполнены полностью, кроме, возможно, последнего; последний заполняется
// слева направо без пропусков. Вставка кладет узел на первое свободное место при обходе
// в ширину, поэтому дерево из одних вставок всегда complete. В файле дерево хранится как
// прямой обход с '#' на месте пустых мест, так что там можно записать и неполное дерево.
type CompleteTree struct {
	root *tnode
}

var _ Storable = (*CompleteTree)(nil)

// NewCompleteTree создает пустое дерево.
func NewCompleteTree() *CompleteTree { return &CompleteTree{} }

// nodes возвращает все узлы в порядке обхода в ширину (сам срез служит очередью).
func (t *CompleteTree) nodes() []*tnode {
	var o []*tnode
	if t.root != nil {
		o = append(o, t.root)
	}
	for i := 0; i < len(o); i++ {
		if o[i].left != nil {
			o = append(o, o[i].left)
		}
		if o[i].right != nil {
			o = append(o, o[i].right)
		}
	}
	return o
}

// Три обхода в глубину отличаются только местом, где записывается сам узел.

func preorder(x *tnode, o *[]string) { // прямой: узел, левое, правое
	if x == nil {
		return
	}
	*o = append(*o, x.key)
	preorder(x.left, o)
	preorder(x.right, o)
}

func inorder(x *tnode, o *[]string) { // симметричный: левое, узел, правое
	if x == nil {
		return
	}
	inorder(x.left, o)
	*o = append(*o, x.key)
	inorder(x.right, o)
}

func postorder(x *tnode, o *[]string) { // обратный: левое, правое, узел
	if x == nil {
		return
	}
	postorder(x.left, o)
	postorder(x.right, o)
	*o = append(*o, x.key)
}

func dumpRec(x *tnode, o *[]string) {
	if x == nil {
		*o = append(*o, "#")
		return
	}
	*o = append(*o, x.key)
	dumpRec(x.left, o)
	dumpRec(x.right, o)
}

// build восстанавливает поддерево из прямого обхода; i - текущая позиция в t.
func build(t []string, i *int) (*tnode, error) {
	if *i >= len(t) {
		return nil, nil
	}
	key := t[*i]
	*i++
	if key == "#" {
		return nil, nil
	}
	norm, err := NormInt(key) // в дереве только числа
	if err != nil {
		return nil, err
	}
	x := &tnode{key: norm}
	if x.left, err = build(t, i); err != nil {
		return nil, err
	}
	if x.right, err = build(t, i); err != nil {
		return nil, err
	}
	return x, nil
}

// Insert кладет число на первое свободное место при обходе в ширину, O(n).
// Возвращает сохраненное (нормализованное) число; не число - ошибка, дерево не меняется.
func (t *CompleteTree) Insert(k string) (string, error) {
	key, err := NormInt(k)
	if err != nil {
		return "", err
	}
	z := &tnode{key: key}
	if t.root == nil {
		t.root = z
		return key, nil
	}
	for _, x := range t.nodes() {
		if x.left == nil {
			x.left = z
			return key, nil
		}
		if x.right == nil {
			x.right = z
			return key, nil
		}
	}
	return key, nil
}

// Contains проверяет наличие числа, O(n): дерево не упорядочено, смотрим все узлы.
func (t *CompleteTree) Contains(k string) (bool, error) {
	key, err := NormInt(k)
	if err != nil {
		return false, err
	}
	for _, x := range t.nodes() {
		if x.key == key {
			return true, nil
		}
	}
	return false, nil
}

// IsComplete проверяет свойство complete: обход в ширину вместе с пустыми местами (nil).
// Если после первого пустого места встретился настоящий узел - дерево не complete.
func (t *CompleteTree) IsComplete() bool {
	o := []*tnode{t.root}
	gap := false
	for i := 0; i < len(o); i++ {
		if o[i] == nil {
			gap = true
			continue
		}
		if gap {
			return false
		}
		o = append(o, o[i].left, o[i].right)
	}
	return true
}

// Traverse возвращает ключи в порядке PREORDER, INORDER, POSTORDER, BFS или SORTED (по возрастанию).
func (t *CompleteTree) Traverse(order string) ([]string, error) {
	o := make([]string, 0)
	switch order {
	case "PREORDER":
		preorder(t.root, &o)
	case "INORDER":
		inorder(t.root, &o)
	case "POSTORDER":
		postorder(t.root, &o)
	case "BFS":
		for _, x := range t.nodes() {
			o = append(o, x.key)
		}
	case "SORTED": // значения по возрастанию; само дерево не меняется
		for _, x := range t.nodes() {
			o = append(o, x.key)
		}
		sort.Slice(o, func(a, b int) bool {
			va, _ := strconv.ParseInt(o[a], 10, 64)
			vb, _ := strconv.ParseInt(o[b], 10, 64)
			return va < vb
		})
	default:
		return nil, wrap(ErrUnknown, "order "+order+" (inorder, preorder, postorder, bfs, sorted)")
	}
	return o, nil
}

// Type возвращает букву типа.
func (t *CompleteTree) Type() byte { return 'T' }

// Size возвращает число узлов.
func (t *CompleteTree) Size() int { return len(t.nodes()) }

// Dump возвращает прямой обход с '#' на месте пустых мест.
func (t *CompleteTree) Dump() []string {
	o := make([]string, 0)
	dumpRec(t.root, &o)
	return o
}

// Load восстанавливает дерево из прямого обхода. При ошибке дерево остается пустым.
func (t *CompleteTree) Load(tokens []string) error {
	t.root = nil
	i := 0
	root, err := build(tokens, &i)
	if err != nil {
		return err
	}
	t.root = root
	return nil
}

// ---------- Красивый вертикальный вывод ----------
// Каждое поддерево рисуется как прямоугольный блок строк. Блоки левого и правого ребенка
// ставятся рядом, а родитель центрируется над ними и соединяется перекладиной:
//        10
//    ┌────┴────┐
//    20        30

type block struct {
	rows  [][]rune
	width int // ширина блока
	root  int // колонка, в которой стоит корень блока
}

func spaces(n int) []rune {
	if n < 0 {
		n = 0
	}
	r := make([]rune, n)
	for i := range r {
		r[i] = ' '
	}
	return r
}

// fit дополняет строку пробелами до ширины w.
func fit(s []rune, w int) []rune {
	r := make([]rune, 0, w)
	r = append(r, s...)
	for len(r) < w {
		r = append(r, ' ')
	}
	return r
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func buildBlock(x *tnode) block {
	label := []rune(x.key) // в дереве только цифры и минус
	n := len(label)
	if x.left == nil && x.right == nil {
		return block{rows: [][]rune{label}, width: n, root: (n - 1) / 2}
	}

	const gap = 3 // пустое место между блоками детей
	var l, r block
	hl, hr := x.left != nil, x.right != nil
	if hl {
		l = buildBlock(x.left)
	}
	if hr {
		r = buildBlock(x.right)
	}

	rOff := 0 // где начинается блок правого ребенка
	if hl {
		rOff = l.width + gap
	}
	lr, rr := l.root, rOff+r.root
	var c int // колонка родителя
	switch {
	case hl && hr:
		c = (lr + rr) / 2
	case hl:
		c = lr + 2 // один левый ребенок: родитель правее и выше
	default:
		c = rr - 2 // один правый ребенок: родитель левее и выше
	}
	start := c - (n-1)/2 // первая буква подписи родителя
	shift := 0           // если вылезли влево - сдвигаем все вправо
	if start < 0 {
		shift = -start
	}
	lr += shift
	rr += shift
	c += shift
	start += shift

	var childW int
	switch {
	case hl && hr:
		childW = l.width + gap + r.width
	case hl:
		childW = l.width
	default:
		childW = r.width
	}
	width := maxInt(shift+childW, start+n)

	res := block{width: width, root: c}
	top := spaces(width)
	copy(top[start:], label)
	res.rows = append(res.rows, top)

	bar := spaces(width) // перекладина между родителем и детьми
	from, to := c, c
	if hl {
		from = lr
	}
	if hr {
		to = rr
	}
	for i := from; i <= to; i++ {
		bar[i] = '─'
	}
	if hl {
		bar[lr] = '┌'
	}
	if hr {
		bar[rr] = '┐'
	}
	switch {
	case hl && hr:
		bar[c] = '┴'
	case hl:
		bar[c] = '┘'
	default:
		bar[c] = '└'
	}
	res.rows = append(res.rows, bar)

	h := maxInt(len(l.rows), len(r.rows))
	for i := 0; i < h; i++ {
		line := spaces(shift)
		if hl {
			var s []rune
			if i < len(l.rows) {
				s = l.rows[i]
			}
			line = append(line, fit(s, l.width)...)
		}
		if hl && hr {
			line = append(line, spaces(gap)...)
		}
		if hr {
			var s []rune
			if i < len(r.rows) {
				s = r.rows[i]
			}
			line = append(line, fit(s, r.width)...)
		}
		res.rows = append(res.rows, fit(line, width))
	}
	return res
}

// Show рисует дерево сверху вниз.
func (t *CompleteTree) Show() string {
	if t.root == nil {
		return "(empty tree)"
	}
	b := buildBlock(t.root)
	lines := make([]string, 0, len(b.rows))
	for _, r := range b.rows {
		lines = append(lines, strings.TrimRight(string(r), " "))
	}
	return strings.Join(lines, "\n")
}
