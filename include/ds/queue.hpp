// Очередь (Q) строк на основе односвязного списка (начало - голова, конец - хвост).
#pragma once
#include "ds/istorable.hpp"
#include "ds/singly_list.hpp"

namespace ds {


class Queue : public IStorable {
private:
    SinglyList l;
public:
    void push(const std::string& v);   // O(1)
    std::string pop();                 // O(1)
    Strs items() const;                // от начала к концу

    size_t size() const { return l.size(); }
    char type() const { return 'Q'; }
    Strs dump() const { return l.items(); }
    void load(const Strs& t);
    std::string show() const;
};

}  // namespace ds
