// Модульные тесты на Catch2 (v3): полное двоичное дерево, сериализация, фабрика и полиморфизм.
#include <catch2/catch_test_macros.hpp>

#include <cstdio>
#include <memory>
#include <sstream>
#include <stdexcept>
#include <string>
#include <vector>

#include "ds/array.hpp"
#include "ds/complete_tree.hpp"
#include "ds/factory.hpp"
#include "ds/hashtable.hpp"
#include "ds/istorable.hpp"
#include "ds/queue.hpp"
#include "ds/serializer.hpp"
#include "ds/stack.hpp"
#include "test_utils.hpp"

using namespace ds;

namespace {

void fill(CompleteTree& t) {
    for (const char* v : {"10", "20", "30", "15", "25", "5", "1"}) t.insert(v);
}

}  // namespace

// ---------- Полное двоичное дерево ----------

TEST_CASE("tree: insert fills levels left to right") {
    CompleteTree t;
    fill(t);
    REQUIRE(t.size() == 7u);
    REQUIRE(t.traverse("BFS") == S({"10", "20", "30", "15", "25", "5", "1"}));
    REQUIRE(t.type() == 'T');
}

TEST_CASE("tree: four traversal orders") {
    CompleteTree t;
    fill(t);
    REQUIRE(t.traverse("PREORDER") == S({"10", "20", "15", "25", "30", "5", "1"}));
    REQUIRE(t.traverse("INORDER") == S({"15", "20", "25", "10", "5", "30", "1"}));
    REQUIRE(t.traverse("POSTORDER") == S({"15", "25", "20", "5", "1", "30", "10"}));
    REQUIRE(t.traverse("SORTED") == S({"1", "5", "10", "15", "20", "25", "30"}));
    REQUIRE_THROWS_AS(t.traverse("ZIGZAG"), std::runtime_error);
}

TEST_CASE("tree: numbers are normalized and negative numbers work") {
    CompleteTree t;
    REQUIRE(t.insert("007") == "7");
    REQUIRE(t.insert("-5") == "-5");
    REQUIRE(t.insert("-10") == "-10");
    REQUIRE(t.contains("-10"));
    REQUIRE(t.contains("7"));
    REQUIRE_FALSE(t.contains("99"));
    REQUIRE(t.traverse("SORTED") == S({"-10", "-5", "7"}));
}

TEST_CASE("tree: only integers are accepted") {
    CompleteTree t;
    REQUIRE_THROWS_AS(t.insert("abc"), std::runtime_error);
    REQUIRE_THROWS_AS(t.insert("1.5"), std::runtime_error);
    REQUIRE_THROWS_AS(t.contains("x"), std::runtime_error);
    REQUIRE(t.size() == 0u);
}

TEST_CASE("tree: isComplete") {
    CompleteTree t;
    REQUIRE(t.isComplete());   // пустое дерево считается полным
    fill(t);
    REQUIRE(t.isComplete());
    t.insert("7");
    REQUIRE(t.isComplete());   // вставка всегда занимает первое свободное место

    CompleteTree broken;       // у 3 есть ребенок 4, а у 2 детей нет: дырка на нижнем уровне
    broken.load(S({"1", "2", "#", "#", "3", "4", "#", "#", "#"}));
    REQUIRE_FALSE(broken.isComplete());
}

TEST_CASE("tree: dump and load round trip") {
    CompleteTree a;
    fill(a);
    REQUIRE(a.dump() == S({"10", "20", "15", "#", "#", "25", "#", "#", "30", "5", "#", "#", "1", "#", "#"}));
    CompleteTree b;
    b.insert("999");                // load заменяет прежнее содержимое
    b.load(a.dump());
    REQUIRE(b.dump() == a.dump());
    REQUIRE(b.size() == 7u);
    CompleteTree empty;
    b.load(S({}));
    REQUIRE(b.size() == 0u);
    REQUIRE_THROWS_AS(empty.load(S({"abc"})), std::runtime_error);
}

TEST_CASE("tree: show draws the tree") {
    CompleteTree t;
    REQUIRE(t.show() == "(empty tree)");
    t.insert("10");
    REQUIRE(t.show() == "10");
    t.insert("20");
    t.insert("30");
    std::string two = t.show();
    REQUIRE(two.find("10") != std::string::npos);
    REQUIRE(two.find("20") != std::string::npos);
    REQUIRE(two.find("30") != std::string::npos);

    CompleteTree odd;               // деревья с одним ребенком рисуются уголками
    odd.load(S({"1", "2", "3", "#", "#", "#", "#"}));
    REQUIRE(odd.show().find("3") != std::string::npos);
    CompleteTree right;
    right.load(S({"1", "#", "2", "#", "#"}));
    REQUIRE(right.show().find("2") != std::string::npos);
    CompleteTree wide;
    wide.load(S({"-100", "5", "#", "#", "123456", "#", "#"}));
    REQUIRE(wide.show().find("123456") != std::string::npos);
}

// ---------- Фабрика и полиморфизм ----------

TEST_CASE("factory: creates every structure by letter") {
    for (char c : std::string("MFLSQTH")) {
        std::unique_ptr<IStorable> p = Factory::create(c);
        REQUIRE(p != nullptr);
        REQUIRE(p->type() == c);
        REQUIRE(p->size() == 0u);
        REQUIRE_FALSE(p->show().empty());
    }
    REQUIRE_THROWS_AS(Factory::create('Z'), std::runtime_error);
}

TEST_CASE("polymorphism: virtual calls reach the real class") {
    std::unique_ptr<IStorable> s = Factory::create('S');
    s->load(S({"1", "2"}));
    REQUIRE(dynamic_cast<Stack*>(s.get()) != nullptr);      // за указателем на интерфейс - настоящий стек
    REQUIRE(s->show() == "| 2 |  <- top\n| 1 |\n+---+");    // вызвалась версия show() именно стека
    REQUIRE(dynamic_cast<CompleteTree*>(s.get()) == nullptr);
    REQUIRE(s->dump() == S({"1", "2"}));
}

// ---------- Сериализация ----------

namespace {

Strs sampleFor(char type) {
    switch (type) {
        case 'M': return S({"a", "b c", "-5", ""});
        case 'F': return S({"1", "2", "3"});
        case 'L': return S({"x", "y"});
        case 'S': return S({"1", "2", "3"});
        case 'Q': return S({"p", "q"});
        case 'T': return S({"10", "20", "15", "#", "#", "25", "#", "#", "30", "5", "#", "#", "1", "#", "#"});
        default: break;
    }
    HashTable h(LINEAR);
    for (int i = 0; i < 20; ++i) h.put("k" + std::to_string(i), "v " + std::to_string(i));
    h.erase("k4");
    return h.dump();
}

std::unique_ptr<IStorable> filled(char type) {
    std::unique_ptr<IStorable> p = Factory::create(type);
    p->load(sampleFor(type));
    return p;
}

}  // namespace

TEST_CASE("serializer: text round trip for every structure") {
    for (char c : std::string("MFLSQTH")) {
        std::unique_ptr<IStorable> a = filled(c);
        std::stringstream ss;
        Serializer::saveText(*a, ss);
        std::unique_ptr<IStorable> b = Factory::create(c);
        Serializer::loadText(*b, ss);
        REQUIRE(b->dump() == a->dump());
        REQUIRE(b->size() == a->size());
    }
}

TEST_CASE("serializer: binary round trip for every structure") {
    for (char c : std::string("MFLSQTH")) {
        std::unique_ptr<IStorable> a = filled(c);
        std::stringstream ss(std::ios::in | std::ios::out | std::ios::binary);
        Serializer::saveBinary(*a, ss);
        std::unique_ptr<IStorable> b = Factory::create(c);
        Serializer::loadBinary(*b, ss);
        REQUIRE(b->dump() == a->dump());
        REQUIRE(b->size() == a->size());
    }
}

TEST_CASE("serializer: text format is readable") {
    std::unique_ptr<IStorable> a = filled('M');
    std::stringstream ss;
    Serializer::saveText(*a, ss);
    REQUIRE(ss.str() == "M a \"b c\" -5 \"\"\n");
}

TEST_CASE("serializer: binary layout is signature, type, count, length-prefixed items") {
    Stack s;
    s.load(S({"1", "2"}));
    std::stringstream ss(std::ios::in | std::ios::out | std::ios::binary);
    Serializer::saveBinary(s, ss);
    const std::string bytes = ss.str();
    REQUIRE(bytes.size() == 19u);                  // 4 + 1 + 4 + 2 * (4 + 1)
    REQUIRE(bytes.substr(0, 4) == "DSB1");
    REQUIRE(bytes[4] == 'S');
    REQUIRE(bytes[5] == 2);                        // число элементов, младший байт первым
    REQUIRE(bytes[9] == 1);                        // длина первого элемента
    REQUIRE(bytes[13] == '1');
}

TEST_CASE("serializer: text errors") {
    std::unique_ptr<IStorable> st = filled('S');
    std::stringstream ss;
    Serializer::saveText(*st, ss);
    Queue q;
    REQUIRE_THROWS_AS(Serializer::loadText(q, ss), std::runtime_error);       // записан стек, читаем в очередь
    std::stringstream empty;
    Queue q2;
    REQUIRE_THROWS_AS(Serializer::loadText(q2, empty), std::runtime_error);   // пустой ввод
    std::stringstream again;
    Serializer::saveText(*st, again);
    REQUIRE_THROWS_AS(Serializer::loadText(*st, again), std::runtime_error);  // читать можно только в пустой объект
}

TEST_CASE("serializer: binary errors") {
    std::unique_ptr<IStorable> st = filled('S');
    std::stringstream good(std::ios::in | std::ios::out | std::ios::binary);
    Serializer::saveBinary(*st, good);
    const std::string bytes = good.str();

    Queue q;
    std::stringstream wrongType(bytes);
    REQUIRE_THROWS_AS(Serializer::loadBinary(q, wrongType), std::runtime_error);

    Stack s1;
    std::stringstream badMagic(std::string("XXXX") + bytes.substr(4));
    REQUIRE_THROWS_AS(Serializer::loadBinary(s1, badMagic), std::runtime_error);

    Stack s2;
    std::stringstream noCount(bytes.substr(0, 6));                            // оборвано внутри числа элементов
    REQUIRE_THROWS_AS(Serializer::loadBinary(s2, noCount), std::runtime_error);

    Stack s3;
    std::stringstream cutItem(bytes.substr(0, bytes.size() - 1));             // оборвано внутри данных элемента
    REQUIRE_THROWS_AS(Serializer::loadBinary(s3, cutItem), std::runtime_error);

    Array a;
    std::string corrupted = "DSB1";
    corrupted += 'M';
    corrupted += std::string("\x01\x00\x00\x00", 4);                          // один элемент
    corrupted += std::string("\xff\xff\xff\xff", 4);                          // длина 4 ГБ: явно повреждение
    std::stringstream huge(corrupted);
    REQUIRE_THROWS_AS(Serializer::loadBinary(a, huge), std::runtime_error);

    Stack s4;
    std::stringstream again(bytes);
    s4.push("already here");
    REQUIRE_THROWS_AS(Serializer::loadBinary(s4, again), std::runtime_error);
}

TEST_CASE("serializer: files") {
    const std::string textPath = "ds_test_text.tmp", binPath = "ds_test_bin.tmp";
    std::unique_ptr<IStorable> a = filled('H');
    Serializer::saveTextFile(*a, textPath);
    Serializer::saveBinaryFile(*a, binPath);

    std::unique_ptr<IStorable> fromText = Factory::create('H');
    std::unique_ptr<IStorable> fromBin = Factory::create('H');
    Serializer::loadTextFile(*fromText, textPath);
    Serializer::loadBinaryFile(*fromBin, binPath);
    REQUIRE(fromText->dump() == a->dump());
    REQUIRE(fromBin->dump() == a->dump());
    std::remove(textPath.c_str());
    std::remove(binPath.c_str());

    Stack s;
    REQUIRE_THROWS_AS(Serializer::loadTextFile(s, "no_such_file.tmp"), std::runtime_error);
    REQUIRE_THROWS_AS(Serializer::loadBinaryFile(s, "no_such_file.tmp"), std::runtime_error);
    REQUIRE_THROWS_AS(Serializer::saveTextFile(*a, "/nonexistent_dir_xyz/f.tmp"), std::runtime_error);
    REQUIRE_THROWS_AS(Serializer::saveBinaryFile(*a, "/nonexistent_dir_xyz/f.tmp"), std::runtime_error);
}
