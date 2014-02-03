#!/bin/bash

echo "Diff:"
diff -u <(./level0 < ultralong.txt) <(./original < ultralong.txt)

echo "Original:"
time ./original < long.txt >/dev/null

echo "First:"
time ./first < long.txt >/dev/null

echo "Second:"
time ./second < long.txt >/dev/null

echo "Hash:"
time ./hash < long.txt >/dev/null

echo "BTree:"
time ./btree < long.txt >/dev/null

echo "Level0:"
time ./level0 < long.txt >/dev/null