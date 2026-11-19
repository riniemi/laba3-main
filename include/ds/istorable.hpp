// Интерфейс (абстрактный класс) для всех структур данных проекта.
// Каждая структура умеет: назвать свой тип, сообщить размер, отдать элементы списком строк (dump),
// восстановиться из такого списка (load) и нарисовать себя для вывода (show).
// На этом построена сериализация: Serializer работает с любой структурой через этот интерфейс.
#pragma once
#include <cstddef>
#include <string>

#include "ds/common.hpp"

namespace ds {

class IStorable {
public:
    IStorable() = default;
    IStorable(const IStorable&) = default;
    IStorable& operator=(const IStorable&) = default;
    virtual ~IStorable() = default;   // LCOV_EXCL_LINE  виртуальный деструктор: удаление через указатель на интерфейс безопасно

    virtual char type() const = 0;                    // буква типа: M F L S Q T H
    virtual std::size_t size() const = 0;             // число элементов
    virtual Strs dump() const = 0;                    // элементы в порядке хранения
    virtual void load(const Strs& tokens) = 0;        // восстановление из элементов
    virtual std::string show() const = 0;             // текст для вывода на экран
};

}  // namespace ds
