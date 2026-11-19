// Фабрика: по букве типа создает нужную структуру и возвращает указатель на интерфейс IStorable.
// Благодаря виртуальным методам с результатом работают одинаково, не зная его реального типа (полиморфизм).
#pragma once
#include <memory>

#include "ds/istorable.hpp"

namespace ds {

class Factory {
public:
    Factory() = delete;
    // M - массив, F - односвязный список, L - двусвязный, S - стек, Q - очередь, T - дерево, H - хеш-таблица
    static std::unique_ptr<IStorable> create(char type);
};

}  // namespace ds
