# XOR

In the beginning...

## Truth Table

```
  x y
x 0 1
y 1 0
```

## Missing Element "Trick"

XOR is communicative, associative, and not idempotent:

```
A ^ B == B ^ A
A ^ (B ^ C) == (A ^ B) ^ C
A ^ A == 0 != A 
```

This peculiar combination of properties can be used to find the missing element
of two otherwise identical sets:

```python
>>> import operator
>>> a = [5, 7, 9]
>>> b = [3, 5, 7, 9]
>>> A = reduce(operator.xor, a, 0)
>>> A
11
>>> B = reduce(operator.xor, b, 0)
>>> B
8
>>> A ^ B
3
```

This quality does not depend on the size of the sets (we could do the last step
with very large sets in O(1) time).

# Bloom Filter Example

```python
>>> from pybloom import BloomFilter
>>> f = BloomFilter(capacity=10, error_rate=0.5)
>>> f.add(1)
False
>>> f.bitarray
bitarray('00000000000100000000000000000')
>>> f.add(2)
False
>>> f.bitarray
bitarray('00000000000100000000000010000')
+>>> 1 in f
True
+>>> 2 in f
True
+>>> 3 in f # False Positive
True
```

# Counting Bloom Filter

Just like a bloom filter, but expand each bit into a cell that contains a counter.

```json
[{"count": 0},{"count": 0},...]
```

The counts permit deleting elements from the filter (as long as you don't
exceed max count size).

# Invertible Bloom Filter

An invertible bloom filter combines the curious properties of XOR with a
counting bloom filter. In addition to the `count` field we now also have an
`id`. The `id` will contain the actual keys inserted into the filter.

```json
[{"count": 0, "id": 0},{"count": 0, "id": 0},...]
```

When we add elements to the set:

1. hash(value) -> N positions
2. for each position, increment the `count` and xor `id` with value.

Finally there is a `hash` field that is used to reduce the error rate of miss
identifying an `id` when the cell count is 1.

So the final cells look a bit like:

```json
[{"count": 0, "id": 0, "hash": 0},{"count": 0, "id": 0, "hash": 0},...]
```

OK, let's do a example with the `ibf` tool. It produces the IBFs in JSON format
so we can easily inspect what is happening.

```
$ ibf create demo.ibf 10
$ jq '.' demo.ibf
```

This first block is the serialized parameters to our hash functions that pick
positions for entries.

```json
{
  "positioners": [
    {
      "key": [
        8717895732742166000,
        2259404117704393200
      ]
    },
    {
      "key": [
        6050128673802996000,
        501233450539197800
      ]
    },
    {
      "key": [
        3390393562759376400,
        2669985732393126000
      ]
    }
  ],
```

The hasher section is the serialized parameters to our hash function that fills
in the `hash` field of a cell.

```json
  "hasher": {
    "key": [
      1774932891286980000,
      6044372234677422000
    ]
  },
```

The actual IBF cells, composed as discussed of an `id`, `hash`, and `count`.

```json
  "size": 10,
  "cells": [
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    },
    {
      "id": 0,
      "hash": 0,
      "count": 0
    }
  ],
```

The cardinality is the total number of elements in the IBF.

```json
  "cardinality": 0
}
```

Looking more closely at the cells when we insert an element:

```
$ ibf insert demo.ibf 'A Value'
$ jq '.cells' demo.ibf
```

```json
[
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 18331428859966820,
    "hash": 11333042293537241000,
    "count": 1
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 18331428859966820,
    "hash": 11333042293537241000,
    "count": 1
  },
  {
    "id": 18331428859966820,
    "hash": 11333042293537241000,
    "count": 1
  }
]
```

When we attempt to list entries we are looking for "pure" cells: Cells that
have a count of 1 and the `id` hashes to the same value as `hash`.

```
$ ibf list demo.ibf
A Value
```

As long as we can pull out "pure" cells we can unravel the IBF. But if we run
out of "pure" cells and there are non-empty ones left we know we weren't able
to get them all.

```
$ ibf insert demo.ibf 'B Value'
$ ibf insert demo.ibf 'C Value'
$ ibf list demo.ibf
A Value
B Value
C Value
$ jq '.cells' demo.ibf 
```

```json
[
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 18331428859966820,
    "hash": 11333042293537241000,
    "count": 1
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 18612903836677476,
    "hash": 10437079027557620000,
    "count": 1
  },
  {
    "id": 0,
    "hash": 0,
    "count": 0
  },
  {
    "id": 18894378813388132,
    "hash": 9477556251531188000,
    "count": 1
  },
  {
    "id": 281474976710656,
    "hash": 1391892804557128400,
    "count": 2
  },
  {
    "id": 844424930131968,
    "hash": 977555963097886000,
    "count": 2
  },
  {
    "id": 562949953421312,
    "hash": 2215778548118941700,
    "count": 2
  }
]
```
