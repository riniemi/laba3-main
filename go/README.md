# ЛР3. Объектно-ориентированное программирование (Go)

Вторая реализация (первая - C++, `lab3_cpp`): те же структуры данных, тот же интерфейс, те же форматы сериализации.

## Структура

```
go.mod                     модуль lab3go (зависимость: github.com/stretchr/testify)
ds/                        пакет со структурами
  storable.go                интерфейс Storable (аналог IStorable)
  common.go                  ошибки, Quote, Tokenize, NormInt
  array.go singly_list.go doubly_list.go stack.go queue.go complete_tree.go
  hashfn.go hashtable.go          хеш-функции, хеш-таблица
  serializer.go factory.go                 сериализация (текст, бинарный), фабрика
  *_test.go                  тесты (testing + testify), bench_test.go - бенчмарки
cmd/demo/main.go           демонстрация
docs/                      UML (uml_classes_go.png/.puml/.dot), make_uml.py
Makefile  .golangci.yml
```

## Команды

```bash
go mod tidy        # один раз, нужен интернет (создает go.sum)
make fmt           # gofmt
make test          # тесты
make cover         # покрытие + coverage.html
make bench         # бенчмарки
make lint          # gofmt, go vet, staticcheck, golangci-lint
make demo
```

## Как Go выражает ООП

| Понятие | C++ | Go |
|---|---|---|
| Интерфейс | абстрактный класс `IStorable` | `interface Storable`, реализуется неявно |
| Инкапсуляция | `private` / `public` | имя с маленькой / с большой буквы |
| Наследование | `: public IStorable` | нет; композиция (`Stack` содержит `SinglyList`) и интерфейсы |
| Полиморфизм | виртуальные методы | вызов метода через значение интерфейса |
| Освобождение памяти | деструкторы, RAII | сборщик мусора |
| Ошибки | исключения | возврат `error`, проверка через `errors.Is` |
