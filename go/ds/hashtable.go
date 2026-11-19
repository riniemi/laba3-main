package ds

import (
	"fmt"
	"strconv"
	"strings"
)

// Хеш-таблица с открытой адресацией (строка -> строка).
// Два метода разрешения коллизий (вариант 1):
//   Linear - линейное апробирование: i-я проба = (h1 + i) mod N, шаг всегда 1;
//   Double - двойное хеширование:    i-я проба = (h1 + i * h2) mod N, шаг зависит от ключа.
// Размер таблицы N - простое число (для двойного хеширования это гарантирует обход всех ячеек).
// h1 = djb2(ключ) mod N, h2 = 1 + fnv1a(ключ) mod (N - 1) (метод деления).
// Ячейки хранятся в массиве Array: "" - свободна, "\x01" - удалена (надгробие),
// иначе "ключ\x1Fзначение". При заполнении (с надгробиями) более 70% таблица пересобирается:
// удваивается, если живых элементов больше половины, иначе размер сохраняется и просто
// убираются надгробия.

// Method - способ разрешения коллизий.
type Method int

// Методы разрешения коллизий.
const (
	Linear Method = iota
	Double
)

const (
	sep     = '\x1F'  // разделяет ключ и значение внутри ячейки
	sepStr  = "\x1F"  // то же в виде строки
	deleted = "\x01"  // метка удаленной ячейки
	none    = -1      // "нет ячейки"
	maxLoad = 1 << 22 // верхняя граница размера таблицы при загрузке из файла
)

// HashTable - хеш-таблица (H).
type HashTable struct {
	method   Method
	slots    *Array
	cells    int // число ячеек (простое)
	count    int // живых элементов
	tombs    int // надгробий
	autoGrow bool

	probes int // сколько ячеек просмотрено (для эмпирического анализа)
}

var _ Storable = (*HashTable)(nil)

// NewHashTable создает таблицу с емкостью 11 и автоувеличением.
func NewHashTable(m Method) *HashTable {
	return NewHashTableWith(m, 11, true)
}

// NewHashTableWith создает таблицу с заданной начальной емкостью (будет округлена до простого)
// и флагом автоувеличения.
func NewHashTableWith(m Method, capacity int, grow bool) *HashTable {
	h := &HashTable{autoGrow: grow}
	h.init(m, capacity)
	return h
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}
	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}
	return true
}

func nextPrime(n int) int {
	for !isPrime(n) {
		n++
	}
	return n
}

// hasKey: ячейка s хранит ключ key? Сравнение без создания временных строк.
func hasKey(s, key string) bool {
	return len(s) > len(key) && s[len(key)] == sep && s[:len(key)] == key
}

func (h *HashTable) init(m Method, capacity int) {
	h.method = m
	if capacity < 3 {
		capacity = 3
	}
	h.cells = nextPrime(capacity)
	h.slots = NewArray()
	for i := 0; i < h.cells; i++ {
		h.slots.Push("")
	}
	h.count = 0
	h.tombs = 0
}

// slot и setSlot: индекс всегда в пределах [0, cells), поэтому ошибки быть не может.
func (h *HashTable) slot(i int) string {
	s, _ := h.slots.Get(i)
	return s
}

func (h *HashTable) setSlot(i int, s string) {
	_ = h.slots.Set(i, s)
}

func (h *HashTable) home(key string) int {
	return int(hashDjb2(key) % uint(h.cells))
}

func (h *HashTable) stepOf(key string) int {
	if h.method == Linear {
		return 1
	}
	return 1 + int(hashFnv1a(key)%uint(h.cells-1)) // от 1 до N-1: ячейка не повторится
}

// findIndex возвращает индекс ячейки с ключом или none.
func (h *HashTable) findIndex(key string) int {
	idx, step := h.home(key), h.stepOf(key)
	for n := 0; n < h.cells; n++ {
		h.probes++
		s := h.slot(idx)
		if s == "" {
			return none // свободная ячейка - цепочка кончилась, ключа нет
		}
		if s != deleted && hasKey(s, key) {
			return idx // надгробия пропускаем, идем дальше
		}
		idx = (idx + step) % h.cells
	}
	return none
}

// insertFresh вставляет ключ, которого точно нет; роста нет.
func (h *HashTable) insertFresh(key, value string) error {
	idx, step := h.home(key), h.stepOf(key)
	for n := 0; n < h.cells; n++ {
		h.probes++
		s := h.slot(idx)
		if s == "" || s == deleted {
			if s == deleted {
				h.tombs--
			}
			h.setSlot(idx, key+sepStr+value)
			h.count++
			return nil
		}
		idx = (idx + step) % h.cells
	}
	return ErrFull
}

func (h *HashTable) rehash(newCap int) error {
	type pair struct{ k, v string }
	live := make([]pair, 0, h.count)
	for i := 0; i < h.cells; i++ {
		s := h.slot(i)
		if s == "" || s == deleted {
			continue
		}
		p := strings.IndexByte(s, sep)
		live = append(live, pair{s[:p], s[p+1:]})
	}
	h.init(h.method, newCap) // новая пустая таблица
	for _, kv := range live {
		if err := h.insertFresh(kv.k, kv.v); err != nil { // все ключи хешируются заново
			return err
		}
	}
	return nil
}

// Put добавляет или обновляет пару. true - новый ключ, false - значение обновлено.
func (h *HashTable) Put(key, value string) (bool, error) {
	if strings.IndexByte(key, sep) >= 0 {
		return false, ErrBadKey
	}
	if h.autoGrow && (h.count+h.tombs+1)*10 > h.cells*7 { // слишком много занятых ячеек
		newCap := h.cells // живых много - удвоить, иначе только убрать надгробия
		if (h.count+1)*2 > h.cells {
			newCap = h.cells * 2
		}
		if err := h.rehash(newCap); err != nil {
			return false, err
		}
	}
	idx, step, firstFree := h.home(key), h.stepOf(key), none
	for n := 0; n < h.cells; n++ {
		h.probes++
		s := h.slot(idx)
		if s == "" { // ключа нет
			if firstFree == none {
				firstFree = idx
			}
			break
		}
		if s == deleted {
			if firstFree == none {
				firstFree = idx // запомним первое надгробие
			}
		} else if hasKey(s, key) {
			h.setSlot(idx, key+sepStr+value) // ключ уже есть: обновляем
			return false, nil
		}
		idx = (idx + step) % h.cells
	}
	if firstFree == none {
		return false, ErrFull
	}
	if h.slot(firstFree) == deleted {
		h.tombs--
	}
	h.setSlot(firstFree, key+sepStr+value)
	h.count++
	return true, nil
}

// Get ищет ключ. Второе значение - найден ли он.
func (h *HashTable) Get(key string) (string, bool) {
	idx := h.findIndex(key)
	if idx == none {
		return "", false
	}
	s := h.slot(idx)
	return s[strings.IndexByte(s, sep)+1:], true
}

// Contains проверяет наличие ключа.
func (h *HashTable) Contains(key string) bool {
	_, ok := h.Get(key)
	return ok
}

// Erase удаляет ключ; false, если его не было.
func (h *HashTable) Erase(key string) bool {
	idx := h.findIndex(key)
	if idx == none {
		return false
	}
	h.setSlot(idx, deleted) // нельзя просто очистить ячейку: оборвется цепочка проб других ключей
	h.count--
	h.tombs++
	return true
}

// Type возвращает букву типа.
func (h *HashTable) Type() byte { return 'H' }

// Size возвращает число живых элементов.
func (h *HashTable) Size() int { return h.count }

// Capacity возвращает число ячеек.
func (h *HashTable) Capacity() int { return h.cells }

// MethodOf возвращает метод разрешения коллизий.
func (h *HashTable) MethodOf() Method { return h.method }

// Probes возвращает число просмотренных ячеек с последнего сброса счетчиков.
func (h *HashTable) Probes() int { return h.probes }

// ResetCounters обнуляет счетчики.
func (h *HashTable) ResetCounters() { h.probes = 0 }

// Dump: метод, размер, число живых, число надгробий, затем тройки (индекс ключ значение),
// затем индексы надгробий. Нужен для записи в файл.
func (h *HashTable) Dump() []string {
	var live, tomb []string
	for i := 0; i < h.cells; i++ {
		s := h.slot(i)
		if s == "" {
			continue
		}
		if s == deleted {
			tomb = append(tomb, strconv.Itoa(i))
			continue
		}
		p := strings.IndexByte(s, sep)
		live = append(live, strconv.Itoa(i), s[:p], s[p+1:])
	}
	name := "double"
	if h.method == Linear {
		name = "linear"
	}
	r := []string{name, strconv.Itoa(h.cells), strconv.Itoa(h.count), strconv.Itoa(h.tombs)}
	r = append(r, live...)
	return append(r, tomb...)
}

// Load восстанавливает в точности ту же раскладку. Таблица должна быть пустой.
func (h *HashTable) Load(t []string) error {
	bad := wrap(ErrCorrupted, "hash table file")
	if len(t) < 4 || (t[0] != "linear" && t[0] != "double") {
		return bad
	}
	c, err1 := strconv.Atoi(t[1])
	nLive, err2 := strconv.Atoi(t[2])
	nTomb, err3 := strconv.Atoi(t[3])
	if err1 != nil || err2 != nil || err3 != nil || c < 3 || c > maxLoad || nLive < 0 || nTomb < 0 ||
		nLive > len(t) || nTomb > len(t) {
		return bad
	}
	if len(t) != 4+3*nLive+nTomb || !isPrime(c) {
		return bad
	}
	m := Double
	if t[0] == "linear" {
		m = Linear
	}
	h.init(m, c)
	for k := 0; k < nLive; k++ {
		idx, err := strconv.Atoi(t[4+3*k])
		if err != nil || idx < 0 || idx >= h.cells {
			return bad
		}
		h.setSlot(idx, t[5+3*k]+sepStr+t[6+3*k])
		h.count++
	}
	for k := 0; k < nTomb; k++ {
		idx, err := strconv.Atoi(t[4+3*nLive+k])
		if err != nil || idx < 0 || idx >= h.cells {
			return bad
		}
		h.setSlot(idx, deleted)
		h.tombs++
	}
	return nil
}

// Show выводит параметры таблицы и все занятые ячейки.
func (h *HashTable) Show() string {
	var b strings.Builder
	b.WriteString(h.Stats())
	for i := 0; i < h.cells; i++ {
		s := h.slot(i)
		if s == "" {
			continue
		}
		if s == deleted {
			b.WriteString("\n[" + strconv.Itoa(i) + "] (deleted)")
			continue
		}
		p := strings.IndexByte(s, sep)
		b.WriteString("\n[" + strconv.Itoa(i) + "] " + Quote(s[:p]) + " => " + Quote(s[p+1:]))
	}
	return b.String()
}

// pct: процент целый, а для малых значений с долями.
func pct(x float64) string {
	if x < 0.1 {
		return fmt.Sprintf("%.2f%%", x*100)
	}
	return fmt.Sprintf("%.0f%%", x*100)
}

// Stats возвращает параметры таблицы.
func (h *HashTable) Stats() string {
	name := "double hashing"
	if h.method == Linear {
		name = "linear probing"
	}
	return fmt.Sprintf("method: %s, capacity: %d, size: %d, deleted: %d, load: %s",
		name, h.cells, h.count, h.tombs, pct(float64(h.count)/float64(h.cells)))
}
