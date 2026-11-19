// Хеш-таблица с открытой адресацией (строка -> строка).
// Два метода разрешения коллизий (вариант 1):
//   LINEAR - линейное апробирование: i-я проба = (h1 + i) mod N, шаг всегда 1;
//   DOUBLE - двойное хеширование:    i-я проба = (h1 + i * h2) mod N, шаг зависит от ключа.
// Размер таблицы N - простое число (для двойного хеширования это гарантирует обход всех ячеек).
// h1 = djb2(ключ) mod N, h2 = 1 + fnv1a(ключ) mod (N - 1) (метод деления).
// Ячейки хранятся в массиве Array из ЛР1: "" - свободна, "\x01" - удалена (надгробие),
// иначе "ключ\x1Fзначение". При заполнении (с надгробиями) более 70% таблица пересобирается: удваивается,
// если живых элементов больше половины, иначе размер сохраняется и просто убираются надгробия.
#pragma once
#include "ds/istorable.hpp"
#include <string>

#include "ds/array.hpp"

namespace ds {


enum Method : unsigned char { LINEAR, DOUBLE };

class HashTable : public IStorable {
private:
    Method method;
    Array slots;
    size_t cap = 0;        // число ячеек (простое)
    size_t count = 0;      // живых элементов
    size_t tombs = 0;      // надгробий
    bool autoGrow;

    size_t home(const std::string& key) const;
    size_t stepOf(const std::string& key) const;
    size_t findIndex(const std::string& key);                   // индекс ячейки или (size_t)-1
    void insertFresh(const std::string& key, const std::string& value);
    void rehash(size_t newCap);
    void init(Method m, size_t capacity);
public:
    int probes = 0;              // сколько ячеек просмотрено (для эмпирического анализа)

    explicit HashTable(Method m, size_t capacity = 11, bool grow = true);
    HashTable(const HashTable&) = delete;
    HashTable& operator=(const HashTable&) = delete;

    bool put(const std::string& key, const std::string& value);  // true - новый ключ, false - значение обновлено
    bool get(const std::string& key, std::string& value);        // найден ли ключ
    bool contains(const std::string& key);
    bool erase(const std::string& key);
    size_t size() const { return count; }
    char type() const { return 'H'; }
    size_t capacity() const { return cap; }
    Method methodOf() const { return method; }
    void resetCounters() { probes = 0; }

    Strs dump() const;                // для записи в файл: метод, размер, ячейки с индексами
    void load(const Strs& t);         // восстановление в точности той же раскладки
    std::string show() const;         // все занятые ячейки
    std::string stats() const;          // параметры таблицы
};

}  // namespace ds
