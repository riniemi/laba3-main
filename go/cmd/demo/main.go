// Демонстрация: фабрика, полиморфизм через интерфейс Storable, текстовая и бинарная сериализация.
package main

import (
	"fmt"
	"os"
	"slices"
	"unsafe"

	"lab3go/ds"
)

func must(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}

func fileSize(path string) int {
	st, err := os.Stat(path)
	must(err)
	return int(st.Size())
}

func main() {
	// Один срез интерфейсов хранит структуры разных типов (полиморфизм).
	all := make([]ds.Storable, 0)
	for _, t := range []byte("MFLSQTH") {
		s, err := ds.Create(t)
		must(err)
		all = append(all, s)
	}

	must(all[0].Load([]string{"10", "hello world", "-5"})) // массив
	must(all[1].Load([]string{"1", "2", "3"}))             // односвязный список
	must(all[2].Load([]string{"a", "b", "c"}))             // двусвязный список
	must(all[3].Load([]string{"1", "2", "3"}))             // стек
	must(all[4].Load([]string{"first", "second"}))         // очередь
	tree, ok := all[5].(*ds.CompleteTree)                  // приведение интерфейса к реальному типу
	if !ok {
		must(ds.ErrUnknown)
	}
	for _, v := range []string{"10", "20", "30", "15", "25", "5", "1"} {
		_, err := tree.Insert(v)
		must(err)
	}
	table, ok := all[6].(*ds.HashTable)
	if !ok {
		must(ds.ErrUnknown)
	}
	_, err := table.Put("apple", "10")
	must(err)
	_, err = table.Put("banana", "20")
	must(err)

	for _, p := range all {
		fmt.Printf("=== %c (элементов: %d)\n%s\n\n", p.Type(), p.Size(), p.Show())
	}

	fmt.Println("=== Сериализация ===")
	for _, p := range all {
		txt := fmt.Sprintf("demo_%c.txt", p.Type())
		bin := fmt.Sprintf("demo_%c.bin", p.Type())
		must(ds.SaveTextFile(p, txt))
		must(ds.SaveBinaryFile(p, bin))
		fromText, err := ds.Create(p.Type())
		must(err)
		fromBin, err := ds.Create(p.Type())
		must(err)
		must(ds.LoadTextFile(fromText, txt))
		must(ds.LoadBinaryFile(fromBin, bin))
		same := slices.Equal(fromText.Dump(), p.Dump()) && slices.Equal(fromBin.Dump(), p.Dump())
		verdict := "НЕТ"
		if same {
			verdict = "да"
		}
		fmt.Printf("%c: текст %d байт, бинарный %d байт, восстановлено верно: %s\n",
			p.Type(), fileSize(txt), fileSize(bin), verdict)
		must(os.Remove(txt))
		must(os.Remove(bin))
	}

	// В Go нет виртуальных таблиц в привычном виде: значение интерфейса - пара (указатель на тип, указатель на данные).
	var st ds.Storable
	fmt.Println("\n=== unsafe.Sizeof (на 64-битной платформе) ===")
	fmt.Printf("интерфейс Storable=%d байт (2 указателя), Array=%d, SinglyList=%d, DoublyList=%d, Stack=%d, Queue=%d, CompleteTree=%d, HashTable=%d\n",
		unsafe.Sizeof(st), unsafe.Sizeof(ds.Array{}), unsafe.Sizeof(ds.SinglyList{}), unsafe.Sizeof(ds.DoublyList{}),
		unsafe.Sizeof(ds.Stack{}), unsafe.Sizeof(ds.Queue{}), unsafe.Sizeof(ds.CompleteTree{}), unsafe.Sizeof(ds.HashTable{}))
}
