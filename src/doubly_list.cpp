#include "ds/doubly_list.hpp"

using namespace std;

namespace ds {


DNode* DoublyList::findNode(const string& v) const {
    for (DNode* c = head; c; c = c->next) if (c->val == v) return c;
    return nullptr;
}

DNode* DoublyList::mustFind(const string& v) const {
    DNode* t = findNode(v);
    if (!t) fail("value not found: " + v);
    return t;
}

string DoublyList::unlink(DNode* x) {  // главное преимущество двусвязного списка
    string r = x->val;
    if (x->prev) x->prev->next = x->next; else head = x->next;
    if (x->next) x->next->prev = x->prev; else tail = x->prev;
    delete x;
    --n;
    return r;
}

void DoublyList::insertAfter(DNode* t, const string& v) {
    DNode* x = new DNode{v, t, t->next};
    if (t->next) t->next->prev = x; else tail = x;
    t->next = x;
    ++n;
}

DoublyList::~DoublyList() { while (head) popHead(); }

void DoublyList::pushHead(const string& v) {
    DNode* x = new DNode{v, nullptr, head};
    if (head) head->prev = x; else tail = x;
    head = x;
    ++n;
}

void DoublyList::pushTail(const string& v) {
    DNode* x = new DNode{v, tail, nullptr};
    if (tail) tail->next = x; else head = x;
    tail = x;
    ++n;
}

void DoublyList::pushAfter(const string& target, const string& v) { insertAfter(mustFind(target), v); }

void DoublyList::pushBefore(const string& target, const string& v) {
    DNode* t = mustFind(target);
    if (t == head) pushHead(v); else insertAfter(t->prev, v);
}

string DoublyList::popHead() { if (!head) fail("list is empty"); return unlink(head); }

string DoublyList::popTail() { if (!head) fail("list is empty"); return unlink(tail); }

string DoublyList::delAfter(const string& target) {
    DNode* t = mustFind(target);
    if (!t->next) fail("no element after " + target);
    return unlink(t->next);
}

string DoublyList::delBefore(const string& target) {
    DNode* t = mustFind(target);
    if (!t->prev) fail("no element before " + target);
    return unlink(t->prev);
}

bool DoublyList::delValue(const string& v) {
    DNode* x = findNode(v);
    if (x) unlink(x);
    return x != nullptr;
}

int DoublyList::find(const string& v) const {
    int i = 0;
    for (DNode* c = head; c; c = c->next, ++i) if (c->val == v) return i;
    return -1;
}

const string& DoublyList::at(size_t i) const {
    if (i >= n) fail("index out of range: " + to_string(i));
    DNode* c = head;
    while (i--) c = c->next;
    return c->val;
}

Strs DoublyList::items() const {
    Strs r;
    for (DNode* c = head; c; c = c->next) r.push_back(c->val);
    return r;
}

Strs DoublyList::itemsReversed() const {
    Strs r;
    for (DNode* c = tail; c; c = c->prev) r.push_back(c->val);
    return r;
}

void DoublyList::load(const Strs& t) { for (const auto& s : t) pushTail(s); }

string DoublyList::show() const {
    if (!head) return "NULL";
    string r = "NULL <-> ";
    for (DNode* c = head; c; c = c->next) {
        if (c != head) r += " <-> ";
        r += c->val;
    }
    return r + " <-> NULL";
}

}  // namespace ds
