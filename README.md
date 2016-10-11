# IBF

`ibf` is a CLI for working with Invertible Bloom Filters. It allows you to
create, insert, delete, and subtract IBFs. More information on bloom filters in
general may be found on [wikipedia][bloom filters]. This implementation mostly
follows the algorithms presented in [What’s the Difference?  Efficient Set
Reconciliation without Prior Context][whats the difference].

An IBF is particularly useful for determining difference between arbitrarily
large sets where the expected difference between the sets is relatively small.
The complexity for listing the difference between `IBF/1` and `IBF/2` is O(Δ).

See the [walkthrough][walkthrough] for a brief overview of the topics.

## Install

Following the standard Go installation pattern:

```bash
$ go get github.com/calebcase/ibf
```

## Usage

### Creating a Set

This will create an IBF with a size of 80 cells. It takes approximately `1.5 *
delta` cells to have high probability of decoding all differences.

```bash
$ ibf create a.ibf 80
```

### Adding Elements

This will generate all numbers from 0 to 10000000, one per line, and add them
to the IBF.

```bash
$ seq 0 10000000 | ibf insert a.ibf
```

### Removing Elements

This will create a copy of the above set and remove numbers 0 through 9.

```bash
$ cp a.ibf b.ibf
$ seq 0 9 | ibf remove a.ibf
```

### Subtracting IBFs

This will compute the symmetric difference between two IBFs.

```bash
$ ibf subtract a.ibf b.ibf a-b.ibf
```

### Listing

This will list values from the computed symmetric difference.

```bash
$ ibf list a-b.ibf
7
0
4
8
6
1
5
9
3
2
```

If the listing is incomplete (which is possible in all cases, but more likely
as the difference exceeds the `cells / 1.5`), then the listing will return what
it can and exit non-zero (meaning it is always clear when a listing is
incomplete).

### Comparing

Similar to listing, but skipping the creation of the extra difference file and
following the style of the familiar `comm` tool:

```bash
$ ibf comm a.ibf b.ibf
```

### Chaining

By default, insert and delete will attempt to echo their stdin to stdout if
stdout appears to be a pipe. This can be used to chain inserting (or removing)
into multiple IBFs a single command.

It can be useful to build several IBFs at the same time with various sizes to
handle a range of possible differences efficiently. IBFs are well suited to
incremental updates.

This will create IBFs of a range of sizes:

```bash
$ for s in 64 128 256; do ibf create a.$s.ibf $s; done
```

This will insert the same set of data into all of them:

```bash
$ seq 0 10000000 | ibf insert a.64.ibf | ibf insert a.128.ibf | ibf insert a.256.ibf
```

This will copy those and remove the same elements from all of them:

```bash
$ for s in 64 128 256; do cp a.$s.ibf b.$s.ibf; done
$ seq 100 200 | ibf remove b.64.ibf | ibf remove b.128.ibf | ibf remove b.256.ibf
```

We can then compute the differences of each size:

```bash
$ for s in 64 128 256; do ibf subtract a.$s.ibf b.$s.ibf a-b.$s.ibf; done
```

Then observe that we get an incomplete listing (non-zero exit) from the IBF of
size 64, but from sizes 128 and 256 we get the complete difference:

```bash
$ ibf list a-b.64.ibf; echo Incomplete $?
110
116
Unable to list all elements (right).
Incomplete 1
```

```bash
$ ibf list a-b.128.ibf; echo Incomplete $?
139
107
195
148
...
177
117
101
Incomplete 0
```

```bash
$ ibf list a-b.256.ibf; echo Incomplete $?
139
102
117
169
...
149
158
157
Incomplete 0
```

### Seeding

The tool currently uses a fixed set of 3 hash functions. The parameters to the
hash functions are determined through a pseudo-random number generator. The
seed for the generator defaults to `0`, but can be provided when creating the
IBF.

It is necessary for the hash function parameters to match in order to subtract
two different IBFs. Therefor, if you intend to generate IBFs on different
systems and you do not use the default seed of 0, you must arrange that the
same seed is used in both sets.

### Arbitrary Data

The tool is designed such that it can easily insert any newline separate data.
The largest limitation on the size of each data element inserted into the set
is the memory of the system itself.

For example, assuming you don't have newlines in your file names (an assumption
you should be careful about), you can determine the difference between two file
listings very efficiently:

```bash
$ rm /home/$USER/foobar
$ ibf create home.1.ibf 10
$ find /home/$USER | ibf insert home.1.ibf
$ touch /home/$USER/foobar
$ ibf create home.2.ibf 10
$ find /home/$USER | ibf insert home.2.ibf
$ ibf subtract home.1.ibf home.2.ibf home.1-2.ibf
$ ibf subtract home.2.ibf home.1.ibf home.2-1.ibf
$ ibf list home.1-2.ibf
$ ibf list home.2-1.ibf
/home/ccase/foobar
```

## Perspective

### Runtime

As far as utility is concerned, this tool can answer questions similar to the
tool `comm` - That is, what is only in set 1 and 2 (but not in both)?

`comm` requires sorted inputs so we do some extra work here to get the inputs
into a usable format. We do similar prep work to insert the elements into the
IBFs. The incremental cost of updating the IBFs is lower than re-sorting would
be.

```bash
$ seq 0 10000000 > raw.seq.1 
$ seq 100 10000000 > raw.seq.2
```

```bash
$ sort raw.seq.1 > sorted.seq.1
$ sort raw.seq.2 > sorted.seq.2
$ time comm -3 sorted.seq.1 sorted.seq.2 | wc -l
100

real  0m2.984s
user  0m2.944s
sys 0m0.036s
```

```bash
$ ibf create seq.1.ibf 150
$ ibf create seq.2.ibf 150
$ cat raw.seq.1 | ibf insert seq.1.ibf
$ cat raw.seq.2 | ibf insert seq.2.ibf
$ time ibf comm seq.1.ibf seq.2.ibf | wc -l
100

real  0m0.010s
user  0m0.008s
sys   0m0.000s
```

### Size

Considering the example from the runtime, it is necessary to provide `comm` the
complete sets. With an IBF, however, the size is related only to the difference
expected.

Size of the sets:

```bash
$ du -sh sorted.seq.1
76M sorted.seq.1
$ du -sh sorted.seq.2
76M sorted.seq.2
```

Size of an IBF that can detect up to 100 changes with high probability:

```bash
$ du -sh seq.1.ibf
12K seq.1
$ du -sh seq.2.ibf
12K seq.2
```

---

[bloom filters]: https://en.wikipedia.org/wiki/Bloom_filter
[whats the difference]: https://www.ics.uci.edu/~eppstein/pubs/EppGooUye-SIGCOMM-11.pdf
[walkthrough]: walkthrough.md
