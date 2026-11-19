package ds

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strings"
)

// Сериализация и десериализация любой структуры проекта через интерфейс Storable.
//
// Текстовый формат (читаем человеком): одна строка "<тип> <элемент> <элемент> ...",
// элементы с пробелами берутся в двойные кавычки. Пример:  S 1 2 "hello world"
// Бинарный формат (компактный): 4 байта сигнатуры "DSB1", 1 байт типа, uint32 число элементов,
// затем для каждого элемента uint32 длина и сами байты. Числа записываются в порядке
// little-endian, поэтому файл одинаково читается на любой машине.

var magic = [4]byte{'D', 'S', 'B', '1'}

const maxItemSize = 64 * 1024 * 1024 // защита от повреждённых файлов

func requireEmpty(obj Storable) error {
	if obj.Size() != 0 {
		return ErrNotEmpty
	}
	return nil
}

// SaveText записывает структуру одной строкой текста.
func SaveText(obj Storable, w io.Writer) error {
	var b strings.Builder
	b.WriteByte(obj.Type())
	for _, s := range obj.Dump() {
		b.WriteByte(' ')
		b.WriteString(Quote(s))
	}
	b.WriteByte('\n')
	_, err := io.WriteString(w, b.String())
	return err
}

// LoadText читает структуру из текстовой строки. obj должен быть пустым.
func LoadText(obj Storable, r io.Reader) error {
	if err := requireEmpty(obj); err != nil {
		return err
	}
	line, err := bufio.NewReader(r).ReadString('\n')
	if err != nil && (!errors.Is(err, io.EOF) || line == "") {
		return wrap(ErrCorrupted, "text data is empty")
	}
	tok, err := Tokenize(line)
	if err != nil {
		return err
	}
	if len(tok) == 0 || tok[0] != string([]byte{obj.Type()}) {
		return wrap(ErrTypeMatch, "in text data")
	}
	return obj.Load(tok[1:])
}

func putU32(w io.Writer, v uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}

func getU32(r io.Reader) (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, wrap(ErrCorrupted, "binary data is truncated")
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

// SaveBinary записывает структуру в бинарном формате.
func SaveBinary(obj Storable, w io.Writer) error {
	items := obj.Dump()
	if _, err := w.Write(magic[:]); err != nil {
		return err
	}
	if _, err := w.Write([]byte{obj.Type()}); err != nil {
		return err
	}
	if err := putU32(w, uint32(len(items))); err != nil {
		return err
	}
	for _, s := range items {
		if err := putU32(w, uint32(len(s))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, s); err != nil {
			return err
		}
	}
	return nil
}

// LoadBinary читает структуру из бинарного формата. obj должен быть пустым.
func LoadBinary(obj Storable, r io.Reader) error {
	if err := requireEmpty(obj); err != nil {
		return err
	}
	var head [5]byte
	if _, err := io.ReadFull(r, head[:]); err != nil || [4]byte(head[:4]) != magic {
		return wrap(ErrCorrupted, "not a binary data file")
	}
	if head[4] != obj.Type() {
		return wrap(ErrTypeMatch, "in binary data")
	}
	count, err := getU32(r)
	if err != nil {
		return err
	}
	items := make([]string, 0)
	for i := uint32(0); i < count; i++ {
		n, err := getU32(r)
		if err != nil {
			return err
		}
		if n > maxItemSize {
			return wrap(ErrCorrupted, "binary data is corrupted")
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(r, buf); err != nil {
			return wrap(ErrCorrupted, "binary data is truncated")
		}
		items = append(items, string(buf))
	}
	return obj.Load(items)
}

// SaveTextFile записывает структуру в текстовый файл.
func SaveTextFile(obj Storable, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	return SaveText(obj, f)
}

// LoadTextFile читает структуру из текстового файла.
func LoadTextFile(obj Storable, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return LoadText(obj, f)
}

// SaveBinaryFile записывает структуру в бинарный файл.
func SaveBinaryFile(obj Storable, path string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { err = errors.Join(err, f.Close()) }()
	return SaveBinary(obj, f)
}

// LoadBinaryFile читает структуру из бинарного файла.
func LoadBinaryFile(obj Storable, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return LoadBinary(obj, f)
}
