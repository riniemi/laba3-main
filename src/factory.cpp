#include "ds/factory.hpp"

#include <memory>
#include <string>

#include "ds/array.hpp"
#include "ds/complete_tree.hpp"
#include "ds/doubly_list.hpp"
#include "ds/hashtable.hpp"
#include "ds/queue.hpp"
#include "ds/singly_list.hpp"
#include "ds/stack.hpp"

namespace ds {

std::unique_ptr<IStorable> Factory::create(char type) {
    switch (type) {
        case 'M': return std::make_unique<Array>();
        case 'F': return std::make_unique<SinglyList>();
        case 'L': return std::make_unique<DoublyList>();
        case 'S': return std::make_unique<Stack>();
        case 'Q': return std::make_unique<Queue>();
        case 'T': return std::make_unique<CompleteTree>();
        case 'H': return std::make_unique<HashTable>(DOUBLE);
        default: break;
    }
    fail(std::string("unknown structure type: ") + type);
}

}  // namespace ds
