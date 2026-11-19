// Модульные тесты на Boost.Test: вспомогательные функции, массив, односвязный и двусвязный списки.
#define BOOST_TEST_MODULE DsBoostTests
#include <boost/test/included/unit_test.hpp>

#include <stdexcept>
#include <string>

#include "ds/array.hpp"
#include "ds/common.hpp"
#include "ds/doubly_list.hpp"
#include "ds/singly_list.hpp"
#include "test_utils.hpp"

using namespace ds;

BOOST_AUTO_TEST_SUITE(CommonSuite)

BOOST_AUTO_TEST_CASE(quote_wraps_only_when_needed) {
    BOOST_CHECK_EQUAL(quote("abc"), "abc");
    BOOST_CHECK_EQUAL(quote("a b"), "\"a b\"");
    BOOST_CHECK_EQUAL(quote(""), "\"\"");
    BOOST_CHECK_EQUAL(quote("a\"b"), "\"a\\\"b\"");
    BOOST_CHECK_EQUAL(quote("a\\b"), "\"a\\\\b\"");
}

BOOST_AUTO_TEST_CASE(tokenize_splits_and_honours_quotes) {
    BOOST_CHECK(tokenize("MPUSH arr \"hello world\"") == S({"MPUSH", "arr", "hello world"}));
    BOOST_CHECK(tokenize("  a   b ") == S({"a", "b"}));
    BOOST_CHECK(tokenize("\"a\\\"b\" c") == S({"a\"b", "c"}));
    BOOST_CHECK(tokenize("\"\" x") == S({"", "x"}));
    BOOST_CHECK(tokenize("").empty());
}

BOOST_AUTO_TEST_CASE(tokenize_rejects_unterminated_quote) {
    BOOST_CHECK_THROW(tokenize("\"abc"), std::runtime_error);
}

BOOST_AUTO_TEST_CASE(join_spaced_upper_endsWith) {
    BOOST_CHECK_EQUAL(join(S({"a", "b", "c"}), "-"), "a-b-c");
    BOOST_CHECK_EQUAL(join(S({}), "-"), "");
    BOOST_CHECK_EQUAL(spaced(S({})), "(empty)");
    BOOST_CHECK_EQUAL(spaced(S({"a b", "c"})), "\"a b\" c");
    BOOST_CHECK_EQUAL(upper("abC1"), "ABC1");
    BOOST_CHECK(endsWith("hello", "llo"));
    BOOST_CHECK(!endsWith("lo", "hello"));
    BOOST_CHECK(!endsWith("hello", "xlo"));
}

BOOST_AUTO_TEST_CASE(arg_and_toIndex) {
    BOOST_CHECK_EQUAL(arg(S({"A", "B"}), 1), "B");
    BOOST_CHECK_THROW(arg(S({"A"}), 3), std::runtime_error);
    BOOST_CHECK_EQUAL(toIndex("12"), 12u);
    BOOST_CHECK_EQUAL(toIndex("0"), 0u);
    BOOST_CHECK_THROW(toIndex("-1"), std::runtime_error);
    BOOST_CHECK_THROW(toIndex("abc"), std::runtime_error);
    BOOST_CHECK_THROW(toIndex(""), std::runtime_error);
    BOOST_CHECK_THROW(toIndex("1234567890"), std::runtime_error);
}

BOOST_AUTO_TEST_CASE(integers) {
    BOOST_CHECK(isInteger("5"));
    BOOST_CHECK(isInteger("-5"));
    BOOST_CHECK(isInteger("123456789012345678"));
    BOOST_CHECK(isInteger("1234567890123456789"));
    BOOST_CHECK(!isInteger("-"));
    BOOST_CHECK(!isInteger(""));
    BOOST_CHECK(!isInteger("1.5"));
    BOOST_CHECK(!isInteger("+5"));
    BOOST_CHECK_EQUAL(normInt("-007"), "-7");
    BOOST_CHECK_EQUAL(normInt("0"), "0");
    BOOST_CHECK_EQUAL(normInt("0001234567890123456789"), "1234567890123456789");
    BOOST_CHECK_THROW(normInt("abc"), std::runtime_error);
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(ArraySuite)

BOOST_AUTO_TEST_CASE(push_insert_get_set_erase) {
    Array a;
    a.push("10");
    a.push("20");
    a.insert(1, "15");
    BOOST_CHECK_EQUAL(a.size(), 3u);
    BOOST_CHECK_EQUAL(a.get(0), "10");
    BOOST_CHECK_EQUAL(a.get(1), "15");
    BOOST_CHECK_EQUAL(a.get(2), "20");
    a.set(0, "99");
    BOOST_CHECK_EQUAL(a.get(0), "99");
    BOOST_CHECK_EQUAL(a.erase(1), "15");
    BOOST_CHECK(a.dump() == S({"99", "20"}));
    BOOST_CHECK_EQUAL(a.show(), "[99, 20] (len=2)");
    BOOST_CHECK_EQUAL(a.type(), 'M');
}

BOOST_AUTO_TEST_CASE(grows_beyond_initial_capacity) {
    Array a;
    for (int i = 0; i < 100; ++i) a.push(std::to_string(i));
    BOOST_CHECK_EQUAL(a.size(), 100u);
    BOOST_CHECK_EQUAL(a.get(0), "0");
    BOOST_CHECK_EQUAL(a.get(99), "99");
    a.insert(0, "first");
    BOOST_CHECK_EQUAL(a.get(0), "first");
    BOOST_CHECK_EQUAL(a.get(1), "0");
}

BOOST_AUTO_TEST_CASE(errors) {
    Array a;
    a.push("x");
    BOOST_CHECK_THROW(a.get(5), std::runtime_error);
    BOOST_CHECK_THROW(a.set(5, "y"), std::runtime_error);
    BOOST_CHECK_THROW(a.erase(5), std::runtime_error);
    BOOST_CHECK_THROW(a.insert(9, "y"), std::runtime_error);
}

BOOST_AUTO_TEST_CASE(load_and_show_empty) {
    Array a;
    BOOST_CHECK_EQUAL(a.show(), "[] (len=0)");
    a.load(S({"x", "y z"}));
    BOOST_CHECK_EQUAL(a.size(), 2u);
    BOOST_CHECK_EQUAL(a.show(), "[x, y z] (len=2)");
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(SinglyListSuite)

BOOST_AUTO_TEST_CASE(four_ways_to_insert) {
    SinglyList l;
    l.pushTail("2");
    l.pushHead("1");
    l.pushAfter("2", "4");
    l.pushBefore("4", "3");
    BOOST_CHECK(l.items() == S({"1", "2", "3", "4"}));
    BOOST_CHECK_EQUAL(l.size(), 4u);
    BOOST_CHECK_EQUAL(l.show(), "1 -> 2 -> 3 -> 4 -> NULL");
    BOOST_CHECK_EQUAL(l.type(), 'F');
}

BOOST_AUTO_TEST_CASE(insert_keeps_head_and_tail_correct) {
    SinglyList l;
    l.pushTail("a");
    l.pushBefore("a", "first");   // перед головой
    l.pushAfter("a", "mid");      // после хвоста: хвост должен сместиться
    l.pushTail("last");
    BOOST_CHECK(l.items() == S({"first", "a", "mid", "last"}));
}

BOOST_AUTO_TEST_CASE(find_and_at) {
    SinglyList l;
    l.load(S({"a", "b", "c"}));
    BOOST_CHECK_EQUAL(l.find("c"), 2);
    BOOST_CHECK_EQUAL(l.find("zz"), -1);
    BOOST_CHECK_EQUAL(l.at(1), "b");
    BOOST_CHECK_THROW(l.at(9), std::runtime_error);
}

BOOST_AUTO_TEST_CASE(deletions) {
    SinglyList l;
    l.load(S({"1", "2", "3", "4", "5"}));
    BOOST_CHECK_EQUAL(l.delBefore("3"), "2");      // 1 3 4 5
    BOOST_CHECK_EQUAL(l.delAfter("3"), "4");       // 1 3 5
    BOOST_CHECK(l.delValue("3"));                  // 1 5
    BOOST_CHECK(!l.delValue("zz"));
    BOOST_CHECK(l.delValue("5"));                  // хвост: 1
    BOOST_CHECK_EQUAL(l.popTail(), "1");           // единственный элемент
    BOOST_CHECK_EQUAL(l.size(), 0u);
    BOOST_CHECK_EQUAL(l.show(), "NULL");
}

BOOST_AUTO_TEST_CASE(deleting_head_and_neighbours) {
    SinglyList l;
    l.load(S({"a", "b", "c"}));
    BOOST_CHECK_EQUAL(l.delBefore("b"), "a");      // перед вторым - это голова
    BOOST_CHECK(l.delValue("b"));                  // теперь b - голова
    BOOST_CHECK_EQUAL(l.items().size(), 1u);
    l.pushTail("d");
    BOOST_CHECK_EQUAL(l.popTail(), "d");
    BOOST_CHECK_EQUAL(l.popHead(), "c");
}

BOOST_AUTO_TEST_CASE(errors) {
    SinglyList l;
    BOOST_CHECK_THROW(l.popHead(), std::runtime_error);
    BOOST_CHECK_THROW(l.popTail(), std::runtime_error);
    BOOST_CHECK_THROW(l.pushAfter("none", "x"), std::runtime_error);
    BOOST_CHECK_THROW(l.pushBefore("none", "x"), std::runtime_error);
    l.load(S({"a", "b"}));
    BOOST_CHECK_THROW(l.delAfter("b"), std::runtime_error);    // после последнего нечего удалять
    BOOST_CHECK_THROW(l.delBefore("a"), std::runtime_error);   // перед первым нечего удалять
    BOOST_CHECK_THROW(l.delAfter("none"), std::runtime_error);
}

BOOST_AUTO_TEST_CASE(pop_head_and_tail_of_longer_list) {
    SinglyList l;
    l.load(S({"1", "2", "3"}));
    BOOST_CHECK_EQUAL(l.popTail(), "3");
    BOOST_CHECK_EQUAL(l.popHead(), "1");
    BOOST_CHECK(l.dump() == S({"2"}));
}

BOOST_AUTO_TEST_SUITE_END()

BOOST_AUTO_TEST_SUITE(DoublyListSuite)

BOOST_AUTO_TEST_CASE(four_ways_to_insert) {
    DoublyList l;
    l.pushTail("b");
    l.pushHead("a");
    l.pushAfter("b", "d");
    l.pushBefore("d", "c");
    BOOST_CHECK(l.items() == S({"a", "b", "c", "d"}));
    BOOST_CHECK(l.itemsReversed() == S({"d", "c", "b", "a"}));
    BOOST_CHECK_EQUAL(l.size(), 4u);
    BOOST_CHECK_EQUAL(l.show(), "NULL <-> a <-> b <-> c <-> d <-> NULL");
    BOOST_CHECK_EQUAL(l.type(), 'L');
}

BOOST_AUTO_TEST_CASE(insert_before_head_and_after_tail) {
    DoublyList l;
    l.pushTail("m");
    l.pushBefore("m", "first");
    l.pushAfter("m", "last");
    BOOST_CHECK(l.items() == S({"first", "m", "last"}));
    BOOST_CHECK(l.itemsReversed() == S({"last", "m", "first"}));
}

BOOST_AUTO_TEST_CASE(find_at_empty_show) {
    DoublyList l;
    BOOST_CHECK_EQUAL(l.show(), "NULL");
    l.load(S({"x", "y", "z"}));
    BOOST_CHECK_EQUAL(l.find("z"), 2);
    BOOST_CHECK_EQUAL(l.find("q"), -1);
    BOOST_CHECK_EQUAL(l.at(0), "x");
    BOOST_CHECK_THROW(l.at(5), std::runtime_error);
}

BOOST_AUTO_TEST_CASE(deletions_are_o1_at_both_ends) {
    DoublyList l;
    l.load(S({"1", "2", "3", "4", "5", "6"}));
    BOOST_CHECK_EQUAL(l.popHead(), "1");
    BOOST_CHECK_EQUAL(l.popTail(), "6");
    BOOST_CHECK_EQUAL(l.delBefore("3"), "2");
    BOOST_CHECK_EQUAL(l.delAfter("3"), "4");
    BOOST_CHECK(l.delValue("3"));
    BOOST_CHECK(!l.delValue("3"));
    BOOST_CHECK(l.dump() == S({"5"}));
}

BOOST_AUTO_TEST_CASE(errors) {
    DoublyList l;
    BOOST_CHECK_THROW(l.popHead(), std::runtime_error);
    BOOST_CHECK_THROW(l.popTail(), std::runtime_error);
    BOOST_CHECK_THROW(l.pushAfter("none", "x"), std::runtime_error);
    l.load(S({"a", "b"}));
    BOOST_CHECK_THROW(l.delAfter("b"), std::runtime_error);
    BOOST_CHECK_THROW(l.delBefore("a"), std::runtime_error);
}

BOOST_AUTO_TEST_SUITE_END()
