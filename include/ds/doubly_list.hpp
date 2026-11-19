// Двусвязный список (L) строк.
#pragma once
#include "ds/istorable.hpp"

namespace ds {


struct DNode { std::string val; DNode *prev, *next; };

class DoublyList : public IStorable {
private:
    DNode *head = nullptr, *tail = nullptr;
    size_t n = 0;
    DNode* findNode(const std::string& v) const;
    DNode* mustFind(const std::string& v) const;
    std::string unlink(DNode* x);               // O(1)
    void insertAfter(DNode* t, const std::string& v);
public:
    DoublyList() {}
    DoublyList(const DoublyList&) = delete;
    DoublyList& operator=(const DoublyList&) = delete;
    ~DoublyList();

    void pushHead(const std::string& v);                                  // O(1)
    void pushTail(const std::string& v);                                  // O(1)
    void pushAfter(const std::string& target, const std::string& v);      // O(n) - поиск
    void pushBefore(const std::string& target, const std::string& v);     // O(n) - поиск
    std::string popHead();                                                // O(1)
    std::string popTail();                                                // O(1)
    std::string delAfter(const std::string& target);
    std::string delBefore(const std::string& target);
    bool delValue(const std::string& v);
    int find(const std::string& v) const;
    const std::string& at(size_t i) const;
    Strs items() const;
    Strs itemsReversed() const;

    size_t size() const { return n; }
    char type() const { return 'L'; }
    Strs dump() const { return items(); }
    void load(const Strs& t);
    std::string show() const;
};

}  // namespace ds
