package ds

import (
	"strconv"
	"testing"
)

// Бенчмарки: время одной операции при разном числе элементов N.
// Запуск: go test ./ds -run '^$' -bench . -benchmem
// O(1) почти не растет с N, O(n) растет пропорционально N.

var sizes = []int{1000, 4000, 16000}

func forSizes(b *testing.B, f func(b *testing.B, n int)) {
	for _, n := range sizes {
		n := n
		b.Run("N="+strconv.Itoa(n), func(b *testing.B) { f(b, n) })
	}
}

func makeKeys(n int) []string {
	keys := make([]string, n)
	for i := range keys {
		keys[i] = "key" + strconv.Itoa(i)
	}
	return keys
}

// ---- массив ----

func BenchmarkArrayPush(b *testing.B) { // O(1)
	a := NewArray()
	for i := 0; i < b.N; i++ {
		a.Push("x")
	}
}

func BenchmarkArrayInsertFront(b *testing.B) { // O(n)
	forSizes(b, func(b *testing.B, n int) {
		a := NewArray()
		for i := 0; i < n; i++ {
			a.Push("x")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = a.Insert(0, "x")
			_, _ = a.Erase(0)
		}
	})
}

func BenchmarkArrayGet(b *testing.B) { // O(1)
	forSizes(b, func(b *testing.B, n int) {
		a := NewArray()
		for i := 0; i < n; i++ {
			a.Push("x")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = a.Get(i % n)
		}
	})
}

// ---- списки ----

func BenchmarkSinglyListPushTail(b *testing.B) { // O(1)
	l := NewSinglyList()
	for i := 0; i < b.N; i++ {
		l.PushTail("x")
	}
}

func BenchmarkSinglyListPopTail(b *testing.B) { // O(n): ищем предпоследний узел
	forSizes(b, func(b *testing.B, n int) {
		l := NewSinglyList()
		for i := 0; i < n; i++ {
			l.PushTail("x")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = l.PopTail()
			l.PushTail("x")
		}
	})
}

func BenchmarkDoublyListPopTail(b *testing.B) { // O(1): есть указатель prev
	forSizes(b, func(b *testing.B, n int) {
		l := NewDoublyList()
		for i := 0; i < n; i++ {
			l.PushTail("x")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = l.PopTail()
			l.PushTail("x")
		}
	})
}

func BenchmarkSinglyListFindLast(b *testing.B) { // O(n)
	forSizes(b, func(b *testing.B, n int) {
		l := NewSinglyList()
		for i := 0; i < n; i++ {
			l.PushTail(strconv.Itoa(i))
		}
		last := strconv.Itoa(n - 1)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			l.Find(last)
		}
	})
}

// ---- стек и очередь ----

func BenchmarkStackPushPop(b *testing.B) { // O(1)
	s := NewStack()
	for i := 0; i < b.N; i++ {
		s.Push("x")
		_, _ = s.Pop()
	}
}

func BenchmarkQueuePushPop(b *testing.B) { // O(1)
	q := NewQueue()
	for i := 0; i < b.N; i++ {
		q.Push("x")
		_, _ = q.Pop()
	}
}

// ---- хеш-таблица ----

func benchHashPut(b *testing.B, m Method) {
	forSizes(b, func(b *testing.B, n int) {
		keys := makeKeys(n)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h := NewHashTable(m)
			for _, k := range keys {
				_, _ = h.Put(k, "v")
			}
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/put")
	})
}

func BenchmarkHashPutLinear(b *testing.B) { benchHashPut(b, Linear) } // O(1)
func BenchmarkHashPutDouble(b *testing.B) { benchHashPut(b, Double) } // O(1)

func benchHashGet(b *testing.B, m Method, present bool) {
	forSizes(b, func(b *testing.B, n int) {
		keys := makeKeys(n)
		h := NewHashTableWith(m, 11, true)
		for _, k := range keys {
			_, _ = h.Put(k, "v")
		}
		probe := keys
		if !present {
			probe = makeKeys(2 * n)[n:] // ключи, которых в таблице нет
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			h.Get(probe[i%n])
		}
	})
}

func BenchmarkHashGetLinear(b *testing.B) { benchHashGet(b, Linear, true) }
func BenchmarkHashGetDouble(b *testing.B) { benchHashGet(b, Double, true) }
func BenchmarkHashMiss(b *testing.B)      { benchHashGet(b, Double, false) }

// ---- дерево ----

func BenchmarkTreeInsert(b *testing.B) { // O(n) на вставку; метрика - среднее на одну вставку
	forSizes(b, func(b *testing.B, n int) {
		for i := 0; i < b.N; i++ {
			t := NewCompleteTree()
			for j := 0; j < n; j++ {
				_, _ = t.Insert(strconv.Itoa(j))
			}
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/insert")
	})
}

func BenchmarkTreeTraverse(b *testing.B) { // O(1) на узел
	forSizes(b, func(b *testing.B, n int) {
		t := NewCompleteTree()
		for j := 0; j < n; j++ {
			_, _ = t.Insert(strconv.Itoa(j))
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = t.Traverse("INORDER")
		}
		b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N)/float64(n), "ns/node")
	})
}
