#include "ds/stack.hpp"

#include <algorithm>

using namespace std;

namespace ds {


void Stack::push(const string& v) { l.pushHead(v); }

string Stack::pop() {
    if (!l.size()) fail("stack is empty");
    return l.popHead();
}

Strs Stack::items() const { return l.items(); }

Strs Stack::dump() const {
    Strs v = l.items();
    reverse(v.begin(), v.end());
    return v;
}

void Stack::load(const Strs& t) { for (const auto& s : t) push(s); }

string Stack::show() const {
    const SNode* first = l.headNode();
    if (!first) return "(empty stack)";
    string r;
    for (const SNode* c = first; c; c = c->next) {
        r += "| " + c->val + " |";
        if (c == first) r += "  <- top";
        r += "\n";
    }
    return r + "+---+";
}

}  // namespace ds
