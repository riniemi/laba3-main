#include "ds/queue.hpp"

using namespace std;

namespace ds {


void Queue::push(const string& v) { l.pushTail(v); }

string Queue::pop() {
    if (!l.size()) fail("queue is empty");
    return l.popHead();
}

Strs Queue::items() const { return l.items(); }

void Queue::load(const Strs& t) { for (const auto& s : t) push(s); }

string Queue::show() const {
    const SNode* first = l.headNode();
    if (!first) return "(empty queue)";
    string r = "front -> ";
    for (const SNode* c = first; c; c = c->next) {
        if (c != first) r += ", ";
        r += c->val;
    }
    return r + " <- back";
}

}  // namespace ds
