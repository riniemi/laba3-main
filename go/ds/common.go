// Package ds содержит абстрактные структуры данных ЛР3 на языке Go:
// массив, списки, стек, очередь, полное двоичное дерево, хеш-таблицу с фильтром Блума,
// а также сериализатор и фабрику, работающие с любой структурой через интерфейс Storable.
package ds

import (
	"errors"
	"fmt"
	"strings"
)

// Ошибки пакета. В Go вместо исключений принято возвращать error;
// конкретную причину можно проверить через errors.Is.
var (
	ErrEmpty      = errors.New("structure is empty")
	ErrIndex      = errors.New("index out of range")
	ErrNotFound   = errors.New("value not found")
	ErrNoNeighbor = errors.New("no neighbour element")
	ErrNotInteger = errors.New("tree stores integers only")
	ErrUnknown    = errors.New("unknown value")
	ErrBadKey     = errors.New("key must not contain the character 0x1F")
	ErrFull       = errors.New("hash table is full")
	ErrCorrupted  = errors.New("corrupted data")
	ErrNotEmpty   = errors.New("object must be empty before loading")
	ErrTypeMatch  = errors.New("type mismatch")
)

// wrap добавляет к базовой ошибке пояснение: "index out of range: 7".
func wrap(base error, detail string) error {
	return fmt.Errorf("%w: %s", base, detail)
}

// isSpace повторяет isspace из C для ASCII-символов.
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r'
}

// Quote превращает строку в токен: в двойных кавычках, если есть пробелы, кавычки или обратная косая.
func Quote(s string) string {
	need := s == ""
	for i := 0; i < len(s); i++ {
		if isSpace(s[i]) || s[i] == '"' || s[i] == '\\' {
			need = true
		}
	}
	if !need {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		if s[i] == '"' || s[i] == '\\' {
			b.WriteByte('\\')
		}
		b.WriteByte(s[i])
	}
	b.WriteByte('"')
	return b.String()
}

// Tokenize разбирает строку на токены. Токен в двойных кавычках может содержать пробелы;
// внутри кавычек \" и \\ означают символы " и \.
func Tokenize(s string) ([]string, error) {
	res := make([]string, 0)
	i, n := 0, len(s)
	for i < n {
		for i < n && isSpace(s[i]) {
			i++
		}
		if i >= n {
			break
		}
		var tok strings.Builder
		if s[i] == '"' {
			i++
			closed := false
			for i < n {
				c := s[i]
				i++
				if c == '\\' && i < n {
					tok.WriteByte(s[i])
					i++
				} else if c == '"' {
					closed = true
					break
				} else {
					tok.WriteByte(c)
				}
			}
			if !closed {
				return nil, wrap(ErrCorrupted, "unterminated quote in: "+s)
			}
		} else {
			for i < n && !isSpace(s[i]) {
				tok.WriteByte(s[i])
				i++
			}
		}
		res = append(res, tok.String())
	}
	return res, nil
}

// IsInteger: строка - целое число (возможно, отрицательное).
func IsInteger(s string) bool {
	i := 0
	if s != "" && s[0] == '-' {
		i = 1
	}
	if i >= len(s) {
		return false
	}
	for ; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// NormInt приводит число к каноническому виду: "-007" -> "-7". Не число - ошибка.
func NormInt(s string) (string, error) {
	if !IsInteger(s) {
		return "", wrap(ErrNotInteger, "got: "+s)
	}
	minus := false
	i := 0
	if s[0] == '-' {
		minus = true
		i = 1
	}
	for i+1 < len(s) && s[i] == '0' {
		i++
	}
	res := s[i:]
	if res == "0" {
		return "0", nil
	}
	if minus {
		return "-" + res, nil
	}
	return res, nil
}
