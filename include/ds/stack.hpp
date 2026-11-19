// Стек (S) строк на основе односвязного списка (вершина - голова списка).
#pragma once
#include "ds/istorable.hpp"
#include "ds/singly_list.hpp"

namespace ds {


class Stack : public IStorable {
private:
    SinglyList l;
public:
    void push(const std::string& v);   // O(1)
    std::string pop();                 // O(1)
    Strs items() const;                // от вершины к низу

    size_t size() const { return l.size(); }
    char type() const { return 'S'; }
    Strs dump() const;        // от низа к вершине
    void load(const Strs& t);
    std::string show() const;
};

}  // namespace ds
