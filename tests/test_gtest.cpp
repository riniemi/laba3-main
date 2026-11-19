// Модульные тесты на Google Test: стек, очередь, хеш-таблица.
#include <gtest/gtest.h>

#include <stdexcept>
#include <string>

#include "ds/hashtable.hpp"
#include "ds/queue.hpp"
#include "ds/stack.hpp"
#include "test_utils.hpp"

using namespace ds;

// ---------- Стек ----------

TEST(StackTest, LifoOrder) {
    Stack s;
    s.push("1");
    s.push("2");
    s.push("3");
    EXPECT_EQ(s.size(), 3u);
    EXPECT_EQ(s.pop(), "3");
    EXPECT_EQ(s.pop(), "2");
    EXPECT_EQ(s.pop(), "1");
    EXPECT_EQ(s.size(), 0u);
}

TEST(StackTest, ItemsAndDumpAreMirrored) {
    Stack s;
    s.push("1");
    s.push("2");
    EXPECT_EQ(s.items(), S({"2", "1"}));   // от вершины к низу
    EXPECT_EQ(s.dump(), S({"1", "2"}));    // от низа к вершине (так записывается в файл)
    EXPECT_EQ(s.type(), 'S');
}

TEST(StackTest, ShowDrawsPileWithTopMark) {
    Stack s;
    EXPECT_EQ(s.show(), "(empty stack)");
    s.push("1");
    s.push("2");
    EXPECT_EQ(s.show(), "| 2 |  <- top\n| 1 |\n+---+");
}

TEST(StackTest, LoadRestoresTopElement) {
    Stack s;
    s.load(S({"1", "2", "3"}));
    EXPECT_EQ(s.pop(), "3");
}

TEST(StackTest, PopFromEmptyThrows) {
    Stack s;
    EXPECT_THROW(s.pop(), std::runtime_error);
}

// ---------- Очередь ----------

TEST(QueueTest, FifoOrder) {
    Queue q;
    q.push("a");
    q.push("b");
    q.push("c");
    EXPECT_EQ(q.size(), 3u);
    EXPECT_EQ(q.pop(), "a");
    EXPECT_EQ(q.pop(), "b");
    EXPECT_EQ(q.items(), S({"c"}));
    EXPECT_EQ(q.type(), 'Q');
}

TEST(QueueTest, ShowAndDumpAndLoad) {
    Queue q;
    EXPECT_EQ(q.show(), "(empty queue)");
    q.load(S({"1", "2"}));
    EXPECT_EQ(q.show(), "front -> 1, 2 <- back");
    EXPECT_EQ(q.dump(), S({"1", "2"}));
}

TEST(QueueTest, PopFromEmptyThrows) {
    Queue q;
    EXPECT_THROW(q.pop(), std::runtime_error);
}

// ---------- Хеш-таблица ----------

TEST(HashTableTest, PutGetEraseForBothMethods) {
    for (Method m : {LINEAR, DOUBLE}) {
        HashTable t(m);
        EXPECT_TRUE(t.put("apple", "10"));
        EXPECT_TRUE(t.put("banana", "20"));
        EXPECT_FALSE(t.put("apple", "15"));      // обновление, а не новый ключ
        std::string v;
        EXPECT_TRUE(t.get("apple", v));
        EXPECT_EQ(v, "15");
        EXPECT_TRUE(t.contains("banana"));
        EXPECT_FALSE(t.contains("cherry"));
        EXPECT_EQ(t.size(), 2u);
        EXPECT_TRUE(t.erase("banana"));
        EXPECT_FALSE(t.erase("banana"));
        EXPECT_FALSE(t.get("banana", v));
        EXPECT_EQ(t.size(), 1u);
        EXPECT_EQ(t.methodOf(), m);
        EXPECT_EQ(t.type(), 'H');
    }
}

TEST(HashTableTest, GrowsAndKeepsAllKeys) {
    for (Method m : {LINEAR, DOUBLE}) {
        HashTable t(m);
        for (int i = 0; i < 300; ++i) EXPECT_TRUE(t.put("key" + std::to_string(i), std::to_string(i * 2)));
        EXPECT_EQ(t.size(), 300u);
        EXPECT_TRUE(t.capacity() > 300u);
        for (int i = 0; i < 300; ++i) {
            std::string v;
            ASSERT_TRUE(t.get("key" + std::to_string(i), v));
            ASSERT_EQ(v, std::to_string(i * 2));
        }
    }
}

TEST(HashTableTest, TombstonesDoNotBreakSearchChains) {
    HashTable t(LINEAR, 11, false);
    for (int i = 0; i < 6; ++i) t.put("k" + std::to_string(i), "v");
    for (int i = 0; i < 6; i += 2) EXPECT_TRUE(t.erase("k" + std::to_string(i)));
    for (int i = 1; i < 6; i += 2) EXPECT_TRUE(t.contains("k" + std::to_string(i)));   // цепочки проб целы
    EXPECT_TRUE(t.put("k0", "again"));                                                  // надгробие переиспользуется
    EXPECT_TRUE(t.contains("k0"));
}

TEST(HashTableTest, ChurnRebuildsInsteadOfGrowingForever) {
    HashTable t(DOUBLE);
    for (int i = 0; i < 2000; ++i) {                  // постоянные вставки и удаления копят надгробия
        t.put("k" + std::to_string(i), "v");
        if (i >= 5) t.erase("k" + std::to_string(i - 5));
    }
    EXPECT_EQ(t.size(), 5u);
    EXPECT_TRUE(t.capacity() < 200u);                 // таблица очищалась от надгробий, а не росла
}

TEST(HashTableTest, FullTableThrows) {
    HashTable t(LINEAR, 3, false);
    EXPECT_EQ(t.capacity(), 3u);
    EXPECT_TRUE(t.put("a", "1"));
    EXPECT_TRUE(t.put("b", "2"));
    EXPECT_TRUE(t.put("c", "3"));
    EXPECT_THROW(t.put("d", "4"), std::runtime_error);
    EXPECT_FALSE(t.put("a", "updated"));              // обновить существующий ключ можно и в полной таблице
}

TEST(HashTableTest, KeyWithSeparatorIsRejected) {
    HashTable t(DOUBLE);
    std::string bad = std::string("a") + '\x1f' + "b";
    EXPECT_THROW(t.put(bad, "v"), std::runtime_error);
}

TEST(HashTableTest, DumpAndLoadRestoreSameLayout) {
    HashTable a(DOUBLE);
    for (int i = 0; i < 30; ++i) a.put("k" + std::to_string(i), "v" + std::to_string(i));
    a.erase("k3");
    a.erase("k7");
    HashTable b(LINEAR);
    b.load(a.dump());
    EXPECT_EQ(b.size(), a.size());
    EXPECT_EQ(b.capacity(), a.capacity());
    EXPECT_EQ(b.methodOf(), DOUBLE);
    EXPECT_EQ(b.dump(), a.dump());                    // ячейки и надгробия на тех же местах
    std::string v;
    EXPECT_TRUE(b.get("k5", v));
    EXPECT_EQ(v, "v5");
    EXPECT_FALSE(b.get("k3", v));
}

TEST(HashTableTest, LoadRejectsCorruptedData) {
    HashTable t(DOUBLE);
    EXPECT_THROW(t.load(S({"bad"})), std::runtime_error);
    EXPECT_THROW(t.load(S({"weird", "11", "0", "0"})), std::runtime_error);
    EXPECT_THROW(t.load(S({"linear", "10", "0", "0"})), std::runtime_error);                // размер не простое число
    EXPECT_THROW(t.load(S({"linear", "11", "1", "0", "99", "k", "v"})), std::runtime_error);  // индекс вне таблицы
    EXPECT_THROW(t.load(S({"linear", "11", "0", "1", "99"})), std::runtime_error);            // надгробие вне таблицы
}

TEST(HashTableTest, ShowAndStatsDescribeTable) {
    HashTable t(LINEAR);
    t.put("k", "v");
    t.erase("k");
    t.put("a b", "c d");
    std::string shown = t.show();
    EXPECT_TRUE(shown.find("method: linear probing") != std::string::npos);
    EXPECT_TRUE(shown.find("\"a b\" => \"c d\"") != std::string::npos);
    EXPECT_TRUE(shown.find("(deleted)") != std::string::npos);
    EXPECT_TRUE(HashTable(DOUBLE).stats().find("double hashing") != std::string::npos);
}
