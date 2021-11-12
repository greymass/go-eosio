go-eosio
========

Fast EOSIO primitives for Go.


Benchmarks
----------

```
go version: 1.17.1
goos: linux
goarch: amd64
pkg: github.com/greymass/go-eosio/internal/benchmarks
cpu: Intel(R) Xeon(R) CPU E5-2673 v3 @ 2.40GHz
Benchmark_Decode_AbiDef-2               271942	    4143 ns/op	  1760 B/op	    52 allocs/op
Benchmark_Decode_AbiDef_EosCanada-2     221100	    5306 ns/op	  1192 B/op	    44 allocs/op
Benchmark_Encode_AbiDef-2               362860	    3143 ns/op	  1240 B/op	    39 allocs/op
Benchmark_Encode_AbiDef_EosCanada-2     149077	    7965 ns/op	  2016 B/op	    66 allocs/op
Benchmark_Decode-2                      409782	    2716 ns/op	  1016 B/op	    52 allocs/op
Benchmark_Decode_NoOptimize-2           117184	   10158 ns/op	  1352 B/op	    92 allocs/op
Benchmark_Decode_EosCanada-2             58075	   20342 ns/op	  3432 B/op	   164 allocs/op
Benchmark_Encode-2                     1000000	    1160 ns/op	   392 B/op	    39 allocs/op
Benchmark_Encode_NoOptimize-2           197502	    5988 ns/op	  1056 B/op	    88 allocs/op
Benchmark_Encode_EosCanada-2             95030	   11927 ns/op	  1696 B/op	   134 allocs/op
```

[All benchmark runs](https://github.com/greymass/go-eosio/actions/workflows/benchmark.yml)


License
-------

Copyright (C) 2021  Greymass Inc.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
