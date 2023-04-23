# muni

>muni (無二), lit. "no two" in Japanese

A collection of algorithms for generating unique IDs.

## Supported algorithms

- [Snowflake ID](#snowflake-id)

# Snowflake ID

Snowflake IDs are unique identifiers used in distributed settings. The snowflake generates ~uint64 values, which means it can generate type aliases for easier marshalling.

There is a description of the default implementation from Wikipedia:
>Snowflakes are 64 bits in binary. (Only 63 are used to fit in a signed integer.) The first 41 bits are a timestamp, representing milliseconds since the chosen epoch. The next 10 bits represent a machine ID, preventing clashes. Twelve more bits represent a per-machine sequence number, to allow creation of multiple snowflakes in the same millisecond. The final number is generally serialized in decimal.

The muni implementation allows customization of the bit-lengths for all parts, as well as the tick duration. The implementation is lock-less.
The packages provides standard generators for Twitter snowflakes, Discord snowflakes and Instagram snowflakes (though it is untested whether the Instagram snowflakes are fully compatible, as their implementation might be using mod 1024 to generate the sequence part of the ID).

**Note that the generator panics if you try to generate an ID with an expired epoch.**

## Examples
```go
epoch := time.Date(2023, time.April, 1, 0, 0, 0, 0, time.UTC)

// 3 node bits for up to 8 nodes
// 19 sequence bits for up to 524288 IDs per millisecond
gen := snowflake.New[uint64](epoch, 1, 1*time.Millisecond, 41, 3, 19)

id := gen.ID()
```

## Benchmarks
```
goos: windows
goarch: amd64
pkg: github.com/sollniss/muni/snowflake
cpu: Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz
BenchmarkID-12    	44370328	        25.46 ns/op	       0 B/op	       0 allocs/op
```