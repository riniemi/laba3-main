// Бенчмарк структур данных: среднее время одной операции (наносекунды) при разном числе элементов N.
// Показывает на практике асимптотику из теории: O(1) почти не растет с N, O(n) растет пропорционально N.
//
// Что где находится в файле:
//   kSizes                          - наборы размеров N (по одному столбцу в итоговой таблице);
//   nsPerOp, makeKeys, padTo, row   - вспомогательные функции (замер времени, ключи, выравнивание, вывод строки);
//   функции array*, singly*, doubly*, stack*, queue*, hash*, tree* - по одной на измеряемую операцию структуры;
//   main                            - печатает шапку таблицы и вызывает row(...) для каждой операции.
#include <chrono>
#include <cstdio>
#include <string>
#include <vector>

#include "ds/array.hpp"
#include "ds/complete_tree.hpp"
#include "ds/doubly_list.hpp"
#include "ds/hashtable.hpp"
#include "ds/queue.hpp"
#include "ds/singly_list.hpp"
#include "ds/stack.hpp"

using namespace ds;
using Clock = std::chrono::steady_clock;

namespace {

// Размеры структур, на которых замеряется время (по одному столбцу на каждое N).
const std::vector<int> kSizes = {1000, 4000, 16000};

// Выполняет действие и возвращает среднее время на одну из ops операций, в наносекундах.
template <class F>
double nsPerOp(F&& action, int ops) {
    const auto start = Clock::now();
    action();
    const auto end = Clock::now();
    return static_cast<double>(std::chrono::duration_cast<std::chrono::nanoseconds>(end - start).count()) / ops;
}

// Готовит n ключей вида "key0", "key1", ... для хеш-таблицы.
std::vector<std::string> makeKeys(int n) {
    std::vector<std::string> keys;
    keys.reserve(static_cast<std::size_t>(n));
    for (int i = 0; i < n; ++i) keys.push_back("key" + std::to_string(i));
    return keys;
}

// Дополняет строку пробелами до width символов. printf считает байты, а русская буква занимает два байта,
// поэтому длину считаем сами: байты-продолжения UTF-8 (10xxxxxx) в счет не идут.
std::string padTo(const std::string& s, std::size_t width) {
    std::size_t chars = 0;
    for (unsigned char c : s) if ((c & 0xC0U) != 0x80U) ++chars;
    return chars >= width ? s : s + std::string(width - chars, ' ');
}

// Печатает одну строку таблицы: название операции, теорию (O(...)) и время для каждого N.
void row(const char* name, const char* theory, double (*run)(int)) {
    std::printf("%s %-6s", padTo(name, 34).c_str(), theory);
    for (int n : kSizes) std::printf(" %12.1f", run(n));
    std::printf("\n");
}

// ---- массив ----
// Замер: добавление в конец (амортизированно O(1)).
double arrayPush(int n) {
    Array a;
    return nsPerOp([&] { for (int i = 0; i < n; ++i) a.push("x"); }, n);
}
// Замер: вставка в начало со сдвигом всех элементов (O(n)).
double arrayInsertFront(int n) {
    Array a;
    return nsPerOp([&] { for (int i = 0; i < n; ++i) a.insert(0, "x"); }, n);
}
// Замер: чтение элемента по индексу (O(1)).
double arrayGet(int n) {
    Array a;
    for (int i = 0; i < n; ++i) a.push("x");
    std::size_t sink = 0;
    const double t = nsPerOp([&] { for (int i = 0; i < n; ++i) sink += a.get(static_cast<std::size_t>(i)).size(); }, n);
    return sink == 0 ? 0 : t;
}

// ---- списки ----
// Замер: добавление в хвост односвязного списка по указателю tail (O(1)).
double singlyPushTail(int n) {
    SinglyList l;
    return nsPerOp([&] { for (int i = 0; i < n; ++i) l.pushTail("x"); }, n);
}
// Замер: удаление из хвоста односвязного списка - сначала нужно дойти до предпоследнего узла (O(n)).
double singlyPopTail(int n) {
    SinglyList l;
    for (int i = 0; i < n; ++i) l.pushTail("x");
    return nsPerOp([&] { for (int i = 0; i < n; ++i) l.popTail(); }, n);
}
// Замер: удаление из хвоста двусвязного списка по указателю tail (O(1)) - для сравнения с односвязным.
double doublyPopTail(int n) {
    DoublyList l;
    for (int i = 0; i < n; ++i) l.pushTail("x");
    return nsPerOp([&] { for (int i = 0; i < n; ++i) l.popTail(); }, n);
}
// Замер: поиск последнего элемента в односвязном списке (O(n)); повторяется 100 раз.
double singlyFind(int n) {
    SinglyList l;
    for (int i = 0; i < n; ++i) l.pushTail(std::to_string(i));
    const std::string last = std::to_string(n - 1);
    int sink = 0;
    const double t = nsPerOp([&] { for (int i = 0; i < 100; ++i) sink += l.find(last); }, 100);
    return sink < 0 ? 0 : t;
}

// ---- стек и очередь ----
// Замер: n добавлений и n извлечений стека (обе операции O(1)).
double stackPushPop(int n) {
    Stack s;
    return nsPerOp([&] { for (int i = 0; i < n; ++i) s.push("x"); for (int i = 0; i < n; ++i) s.pop(); }, 2 * n);
}
// Замер: n добавлений и n извлечений очереди (обе операции O(1)).
double queuePushPop(int n) {
    Queue q;
    return nsPerOp([&] { for (int i = 0; i < n; ++i) q.push("x"); for (int i = 0; i < n; ++i) q.pop(); }, 2 * n);
}

// ---- хеш-таблица ----
// Общая функция замера вставки: создает таблицу заданным методом и кладет n ключей (O(1) в среднем).
double hashPut(Method m, int n) {
    const std::vector<std::string> keys = makeKeys(n);
    HashTable t(m);
    return nsPerOp([&] { for (const auto& k : keys) t.put(k, "v"); }, n);
}
// Общая функция замера поиска: present - ищем существующие ключи, иначе отсутствующие (промахи).
double hashGet(Method m, int n, bool present) {
    const std::vector<std::string> keys = makeKeys(n);
    HashTable t(m, 11, true);
    for (const auto& k : keys) t.put(k, "v");
    const std::vector<std::string> probe = present ? keys : makeKeys(2 * n);
    const std::size_t from = present ? 0 : static_cast<std::size_t>(n);
    std::string v;
    int found = 0;
    const double r = nsPerOp([&] { for (std::size_t i = from; i < probe.size(); ++i) found += t.get(probe[i], v) ? 1 : 0; }, n);
    return found < 0 ? 0 : r;
}
// Обертки для row(...): заранее заданные комбинации метода поиска.
double hashPutLinear(int n) { return hashPut(LINEAR, n); }
double hashPutDouble(int n) { return hashPut(DOUBLE, n); }
double hashGetLinear(int n) { return hashGet(LINEAR, n, true); }
double hashGetDouble(int n) { return hashGet(DOUBLE, n, true); }
double hashMiss(int n) { return hashGet(DOUBLE, n, false); }

// ---- дерево ----
// Замер: вставка в полное дерево - поиск первого свободного места обходом в ширину (O(n)).
double treeInsert(int n) {
    CompleteTree t;
    return nsPerOp([&] { for (int i = 0; i < n; ++i) t.insert(std::to_string(i)); }, n);
}
// Замер: обход INORDER (время делится на число узлов n).
double treeTraverse(int n) {
    CompleteTree t;
    for (int i = 0; i < n; ++i) t.insert(std::to_string(i));
    std::size_t sink = 0;
    const double r = nsPerOp([&] { sink += t.traverse("INORDER").size(); }, n);   // время на один узел
    return sink == 0 ? 0 : r;
}

}  // namespace

// Точка входа: печатает шапку таблицы, затем по одной строке на каждую операцию (см. row(...)).
int main() {
    std::printf("Время одной операции, нс (steady_clock). Столбцы: N = 1000, 4000, 16000\n\n");
    std::printf("%s %-6s %12d %12d %12d\n", padTo("операция", 34).c_str(), "теория", kSizes[0], kSizes[1], kSizes[2]);
    std::printf("%s\n", std::string(80, '-').c_str());
    // массив
    row("Array::push (в конец)", "O(1)", arrayPush);
    row("Array::insert (в начало)", "O(n)", arrayInsertFront);
    row("Array::get", "O(1)", arrayGet);
    // списки
    row("SinglyList::pushTail", "O(1)", singlyPushTail);
    row("SinglyList::popTail", "O(n)", singlyPopTail);
    row("DoublyList::popTail", "O(1)", doublyPopTail);
    row("SinglyList::find (последний)", "O(n)", singlyFind);
    // стек и очередь
    row("Stack push+pop", "O(1)", stackPushPop);
    row("Queue push+pop", "O(1)", queuePushPop);
    // хеш-таблица
    row("HashTable::put, линейное", "O(1)", hashPutLinear);
    row("HashTable::put, двойное", "O(1)", hashPutDouble);
    row("HashTable::get, линейное", "O(1)", hashGetLinear);
    row("HashTable::get, двойное", "O(1)", hashGetDouble);
    row("HashTable: нет ключа", "O(1)", hashMiss);
    // дерево
    row("CompleteTree::insert", "O(n)", treeInsert);
    row("CompleteTree::traverse (на узел)", "O(1)", treeTraverse);
    return 0;
}
