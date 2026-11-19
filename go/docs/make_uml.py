#!/usr/bin/env python3
"""UML-диаграмма классов Go-версии: docs/uml_classes_go.dot -> .png (Graphviz) и текст PlantUML .puml.
Запуск из корня проекта: python3 docs/make_uml.py"""
import subprocess

CLASSES = {
    'Storable': ('<<interface>>', [], [
        '+ Type() byte', '+ Size() int', '+ Dump() []string', '+ Load(tokens []string) error', '+ Show() string']),
    'Array': ('', ['- data : []string', '- length : int'],
              ['+ Push(v)', '+ Insert(i, v) error', '+ Get(i) (string, error)', '+ Set(i, v) error', '+ Erase(i) (string, error)']),
    'SinglyList': ('', ['- head, tail : *snode', '- n : int'],
                   ['+ PushHead(v) / PushTail(v)', '+ PushAfter(t, v) / PushBefore(t, v) error', '+ PopHead() / PopTail() (string, error)',
                    '+ DelAfter(t) / DelBefore(t) (string, error)', '+ DelValue(v) bool', '+ Find(v) int', '+ At(i) (string, error)']),
    'DoublyList': ('', ['- head, tail : *dnode', '- n : int'],
                   ['+ PushHead / PushTail / PushAfter / PushBefore', '+ PopHead() / PopTail() (string, error)',
                    '+ DelAfter / DelBefore / DelValue', '+ Find(v) int', '+ At(i) (string, error)', '+ ItemsReversed() []string']),
    'Stack': ('', ['- l : SinglyList'], ['+ Push(v)', '+ Pop() (string, error)', '+ Items() []string']),
    'Queue': ('', ['- l : SinglyList'], ['+ Push(v)', '+ Pop() (string, error)', '+ Items() []string']),
    'CompleteTree': ('', ['- root : *tnode'],
                     ['+ Insert(k) (string, error)', '+ Contains(k) (bool, error)', '+ IsComplete() bool', '+ Traverse(order) ([]string, error)']),
    'HashTable': ('', ['- method : Method', '- slots : *Array', '- cells, count, tombs : int', '- autoGrow, useBloom : bool',
                       '- bloom : *BloomFilter'],
                  ['+ Put(k, v) (bool, error)', '+ Get(k) (string, bool)', '+ Contains(k) bool', '+ Erase(k) bool',
                   '+ Capacity() int', '+ Stats() string']),
    'BloomFilter': ('', ['- words : []uint64', '- m : int', '- k = 3 (const)'],
                    ['+ Add(key)', '+ MightContain(key) bool', '+ Reset(bits)', '+ Bits() int', '+ SetBits() int']),
    'serializer.go': ('<<functions>>', [], ['+ SaveText / LoadText(obj, io)', '+ SaveBinary / LoadBinary(obj, io)',
                                             '+ Save/Load...File(obj, path)']),
    'factory.go': ('<<functions>>', [], ['+ Create(t byte) (Storable, error)']),
    'snode': ('<<struct>>', ['+ val : string', '+ next : *snode'], []),
    'dnode': ('<<struct>>', ['+ val : string', '+ prev, next : *dnode'], []),
    'tnode': ('<<struct>>', ['+ key : string', '+ left, right : *tnode'], []),
}
IMPLEMENT = ['Array', 'SinglyList', 'DoublyList', 'Stack', 'Queue', 'CompleteTree', 'HashTable']
COMPOSE = [('Stack', 'SinglyList'), ('Queue', 'SinglyList'), ('HashTable', 'Array'), ('HashTable', 'BloomFilter')]
ASSOC = [('SinglyList', 'snode'), ('DoublyList', 'dnode'), ('CompleteTree', 'tnode')]
DEPEND = [('serializer.go', 'Storable'), ('factory.go', 'Storable')]


def esc(s):
    for a, b in (('\\', '\\\\'), ('{', '\\{'), ('}', '\\}'), ('|', '\\|'), ('<', '\\<'), ('>', '\\>')):
        s = s.replace(a, b)
    return s


dot = ['digraph UML {', '  rankdir=BT; nodesep=0.35; ranksep=0.7; splines=true;',
       '  node [shape=record, fontname="DejaVu Sans", fontsize=10, style=filled, fillcolor="#FFFFEE"];',
       '  edge [fontname="DejaVu Sans", fontsize=9];']
for name, (stereo, attrs, methods) in CLASSES.items():
    node = name.replace('.', '_')
    title = (esc(stereo) + '\\n' if stereo else '') + name
    parts = [title, ''.join(esc(a) + '\\l' for a in attrs), ''.join(esc(m) + '\\l' for m in methods)]
    fill = '#E8F0FF' if name == 'Storable' else '#FFFFEE'
    dot.append('  %s [label="{%s}", fillcolor="%s"];' % (node, '|'.join(parts), fill))
for c in IMPLEMENT:
    dot.append('  %s -> Storable [arrowhead=empty, style=dashed];' % c)   # реализация интерфейса (неявная в Go)
for whole, part in COMPOSE:
    dot.append('  %s -> %s [arrowhead=diamond, label="1"];' % (part, whole))
for a, b in ASSOC:
    dot.append('  %s -> %s [dir=back, arrowtail=vee, label="0..*"];' % (b, a))
for a, b in DEPEND:
    dot.append('  %s -> %s [style=dashed, arrowhead=vee];' % (a.replace('.', '_'), b))
dot.append('}')
open('docs/uml_classes_go.dot', 'w', encoding='utf-8').write('\n'.join(dot) + '\n')
subprocess.run(['dot', '-Tpng', '-Gdpi=130', 'docs/uml_classes_go.dot', '-o', 'docs/uml_classes_go.png'], check=True)

pu = ['@startuml', 'skinparam classAttributeIconSize 0', 'hide empty members']
for name, (stereo, attrs, methods) in CLASSES.items():
    kind = 'interface' if stereo == '<<interface>>' else 'class'
    extra = ' <<struct>>' if stereo == '<<struct>>' else (' <<functions>>' if stereo == '<<functions>>' else '')
    ident = '"%s"' % name if '.' in name else name
    pu.append('%s %s%s {' % (kind, ident, extra))
    pu += ['  ' + a for a in attrs] + ['  ' + m for m in methods]
    pu.append('}')
for c in IMPLEMENT:
    pu.append('Storable <|.. %s' % c)
for w, p in COMPOSE:
    pu.append('%s *-- "1" %s' % (w, p))
for a, b in ASSOC:
    pu.append('%s o-- "0..*" %s' % (a, b))
for a, b in DEPEND:
    pu.append('"%s" ..> %s' % (a, b))
pu.append('@enduml')
open('docs/uml_classes_go.puml', 'w', encoding='utf-8').write('\n'.join(pu) + '\n')
print('ok')
