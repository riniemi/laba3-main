#include "ds/singly_list.hpp"

using namespace std;

namespace ds {


SNode* SinglyList::findNode(const string& v) const {
    for (SNode* c = head; c; c = c->next) if (c->val == v) return c;
    return nullptr;
}

SNode* SinglyList::mustFind(const string& v) const {
    SNode* t = findNode(v);
    if (!t) fail("value not found: " + v);
    return t;
}

SinglyList::~SinglyList() { while (head) popHead(); }

void SinglyList::pushHead(const string& v) {
    SNode* x = new SNode{v, head};
    head = x;
    if (!tail) tail = x;
    ++n;
}

void SinglyList::pushTail(const string& v) {
    SNode* x = new SNode{v, nullptr};
    if (tail) tail->next = x; else head = x;
    tail = x;
    ++n;
}

void SinglyList::pushAfter(const string& target, const string& v) {
    SNode* t = mustFind(target);
    t->next = new SNode{v, t->next};
    if (t == tail) tail = t->next;
    ++n;
}

void SinglyList::pushBefore(const string& target, const string& v) {
    mustFind(target);
    if (head->val == target) return pushHead(v);
    SNode* p = head;
    while (p->next->val != target) p = p->next;  // p - узел перед целевым
    p->next = new SNode{v, p->next};
    ++n;
}

string SinglyList::popHead() {
    if (!head) fail("list is empty");
    SNode* t = head;
    string r = t->val;
    head = t->next;
    if (!head) tail = nullptr;
    delete t;
    --n;
    return r;
}

string SinglyList::popTail() {  // нужно найти предпоследний узел
    if (!head) fail("list is empty");
    if (head == tail) return popHead();
    SNode* p = head;
    while (p->next != tail) p = p->next;
    string r = tail->val;
    delete tail;
    tail = p;
    p->next = nullptr;
    --n;
    return r;
}

string SinglyList::delAfter(const string& target) {
    SNode* t = mustFind(target);
    if (!t->next) fail("no element after " + target);
    SNode* x = t->next;
    string r = x->val;
    t->next = x->next;
    if (x == tail) tail = t;
    delete x;
    --n;
    return r;
}

string SinglyList::delBefore(const string& target) {
    SNode* t = mustFind(target);
    if (t == head) fail("no element before " + target);
    if (head->next == t) return popHead();
    SNode* p = head;
    while (p->next->next != t) p = p->next;  // p - узел за два до целевого
    SNode* x = p->next;
    string r = x->val;
    p->next = t;
    delete x;
    --n;
    return r;
}

bool SinglyList::delValue(const string& v) {
    if (!findNode(v)) return false;
    if (head->val == v) { popHead(); return true; }
    SNode* p = head;
    while (p->next->val != v) p = p->next;
    SNode* x = p->next;
    p->next = x->next;
    if (x == tail) tail = p;
    delete x;
    --n;
    return true;
}

int SinglyList::find(const string& v) const {
    int i = 0;
    for (SNode* c = head; c; c = c->next, ++i) if (c->val == v) return i;
    return -1;
}

const string& SinglyList::at(size_t i) const {
    if (i >= n) fail("index out of range: " + to_string(i));
    SNode* c = head;
    while (i--) c = c->next;
    return c->val;
}

Strs SinglyList::items() const {
    Strs r;
    for (SNode* c = head; c; c = c->next) r.push_back(c->val);
    return r;
}

void SinglyList::load(const Strs& t) { for (const auto& s : t) pushTail(s); }

string SinglyList::show() const {
    if (!head) return "NULL";
    string r;
    for (SNode* c = head; c; c = c->next) {
        if (c != head) r += " -> ";
        r += c->val;
    }
    return r + " -> NULL";
}

}  // namespace ds
