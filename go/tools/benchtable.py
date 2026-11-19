#!/usr/bin/env python3
"""Превращает вывод `go test -bench` в читаемую таблицу (читает stdin, пишет stdout)."""
import re
import sys

rows = []
for line in sys.stdin:
    t = line.split()
    if len(t) < 4 or not t[0].startswith('Benchmark') or not t[1].isdigit():
        continue
    name = re.sub(r'-\d+$', '', t[0][len('Benchmark'):])       # убираем "Benchmark" и суффикс -12
    m = re.match(r'(.*)/N=(\d+)$', name)
    op, n = (m.group(1), m.group(2)) if m else (name, '-')
    vals = {}
    for i in range(2, len(t) - 1, 2):
        vals[t[i + 1]] = t[i]
    rows.append((op, n, vals))


def human(ns):
    ns = float(ns)
    if ns >= 1e9:
        return '%.2f s' % (ns / 1e9)
    if ns >= 1e6:
        return '%.2f ms' % (ns / 1e6)
    if ns >= 1e3:
        return '%.2f us' % (ns / 1e3)
    return '%.1f ns' % ns


def per_item(v):
    for unit in ('ns/put', 'ns/insert', 'ns/node'):
        if unit in v:
            return human(v[unit])
    return ''


head = ('Операция', 'N', 'время/операцию', 'на 1 элемент', 'B/op', 'allocs/op')
data = [(op, n, human(v.get('ns/op', 0)), per_item(v), v.get('B/op', ''), v.get('allocs/op', ''))
        for op, n, v in rows]
w = [max(len(str(x[i])) for x in data + [head]) for i in range(6)]


def fmt(r):
    return '  '.join(str(r[i]).ljust(w[i]) if i == 0 else str(r[i]).rjust(w[i]) for i in range(6))


print(fmt(head))
print('-' * len(fmt(head)))
prev = None
for r in data:
    if prev is not None and r[0] != prev:
        print()
    prev = r[0]
    print(fmt(r))
print()
print('"на 1 элемент": для хеш-таблицы - одна вставка, для дерева - одна вставка/один узел обхода.')
