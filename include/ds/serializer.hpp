// Сериализация и десериализация любой структуры проекта через интерфейс IStorable.
//
// Текстовый формат (читаем человеком): одна строка "<тип> <элемент> <элемент> ...",
//   элементы с пробелами берутся в двойные кавычки. Пример:  S 1 2 "hello world"
// Бинарный формат (компактный): 4 байта сигнатуры "DSB1", 1 байт типа, uint32 число элементов,
//   затем для каждого элемента uint32 длина и сами байты. Числа записываются в порядке little-endian,
//   поэтому файл одинаково читается на любой машине.
#pragma once
#include <iosfwd>
#include <string>

#include "ds/istorable.hpp"

namespace ds {

class Serializer {
public:
    Serializer() = delete;   // класс-утилита: только статические методы

    static void saveText(const IStorable& obj, std::ostream& out);
    static void loadText(IStorable& obj, std::istream& in);       // obj должен быть пустым

    static void saveBinary(const IStorable& obj, std::ostream& out);
    static void loadBinary(IStorable& obj, std::istream& in);     // obj должен быть пустым

    static void saveTextFile(const IStorable& obj, const std::string& path);
    static void loadTextFile(IStorable& obj, const std::string& path);
    static void saveBinaryFile(const IStorable& obj, const std::string& path);
    static void loadBinaryFile(IStorable& obj, const std::string& path);
};

}  // namespace ds
