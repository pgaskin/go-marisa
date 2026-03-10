### specs

- Go 1.26
- T14s Gen 6, Intel 258V, performance mode
- `CGO_ENABLED=0 go1.26.0 test -v -run='^$' -bench=BenchmarkTrie -count=6 github.com/pgaskin/go-marisa`

### wasm2go

075b64a04783d6c188ac9f66ea52688bc71edb9e

```
goos: linux
goarch: amd64
pkg: github.com/pgaskin/go-marisa
cpu: Intel(R) Core(TM) Ultra 7 258V
BenchmarkTrie
BenchmarkTrie/Words
BenchmarkTrie/Words/Build
BenchmarkTrie/Words/Build-8            3         438924498 ns/op          10.02 MB/s        466550 keys/op         1062940 keys/s              940.8 ns/key     166610776 B/op        66 allocs/op
BenchmarkTrie/Words/Build-8            3         385076195 ns/op          11.42 MB/s        466550 keys/op         1211579 keys/s              825.4 ns/key     166610776 B/op        66 allocs/op
BenchmarkTrie/Words/Build-8            3         402151379 ns/op          10.93 MB/s        466550 keys/op         1160136 keys/s              862.0 ns/key     166610776 B/op        66 allocs/op
BenchmarkTrie/Words/Build-8            3         388799674 ns/op          11.31 MB/s        466550 keys/op         1199976 keys/s              833.4 ns/key     166610813 B/op        66 allocs/op
BenchmarkTrie/Words/Build-8            3         411699621 ns/op          10.68 MB/s        466550 keys/op         1133230 keys/s              882.4 ns/key     166610776 B/op        66 allocs/op
BenchmarkTrie/Words/Build-8            3         423515889 ns/op          10.38 MB/s        466550 keys/op         1101612 keys/s              907.8 ns/key     166610781 B/op        66 allocs/op
BenchmarkTrie/Words/ReadFrom
BenchmarkTrie/Words/ReadFrom-8              1299            996146 ns/op        1418.82 MB/s           138.0 reads/op    4271420 B/op         53 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1147            910192 ns/op        1552.81 MB/s           138.0 reads/op    4271293 B/op         53 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1106           1031819 ns/op        1369.77 MB/s           138.0 reads/op    4271357 B/op         53 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1118            935135 ns/op        1511.39 MB/s           138.0 reads/op    4271499 B/op         53 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1154            995848 ns/op        1419.24 MB/s           138.0 reads/op    4271330 B/op         53 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1026            996012 ns/op        1419.01 MB/s           138.0 reads/op    4271395 B/op         53 allocs/op
BenchmarkTrie/Words/WriteTo
BenchmarkTrie/Words/WriteTo-8             462288              2558 ns/op        552574.10 MB/s         138.0 writes/op       277 B/op         19 allocs/op
BenchmarkTrie/Words/WriteTo-8             470254              2515 ns/op        561898.77 MB/s         138.0 writes/op       277 B/op         19 allocs/op
BenchmarkTrie/Words/WriteTo-8             426367              2562 ns/op        551689.62 MB/s         138.0 writes/op       277 B/op         19 allocs/op
BenchmarkTrie/Words/WriteTo-8             509163              2727 ns/op        518234.12 MB/s         138.0 writes/op       277 B/op         19 allocs/op
BenchmarkTrie/Words/WriteTo-8             589803              2518 ns/op        561300.05 MB/s         138.0 writes/op       277 B/op         19 allocs/op
BenchmarkTrie/Words/WriteTo-8             412988              2524 ns/op        560006.50 MB/s         138.0 writes/op       277 B/op         19 allocs/op
BenchmarkTrie/Words/UnmarshalBinary
BenchmarkTrie/Words/UnmarshalBinary-8       2018            616936 ns/op        2290.92 MB/s     3188392 B/op         38 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       1593            635795 ns/op        2222.97 MB/s     3188392 B/op         38 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       1659            645438 ns/op        2189.76 MB/s     3188392 B/op         38 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       1995            633766 ns/op        2230.08 MB/s     3188392 B/op         38 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       1969            630465 ns/op        2241.76 MB/s     3188392 B/op         38 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       1887            676657 ns/op        2088.73 MB/s     3188392 B/op         38 allocs/op
BenchmarkTrie/Words/MarshalBinary
BenchmarkTrie/Words/MarshalBinary-8         3290            443561 ns/op        3186.37 MB/s     1417240 B/op          2 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2502            426325 ns/op        3315.20 MB/s     1417241 B/op          2 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2450            413799 ns/op        3415.56 MB/s     1417240 B/op          2 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2606            427905 ns/op        3302.96 MB/s     1417240 B/op          2 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2792            437549 ns/op        3230.16 MB/s     1417241 B/op          2 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2730            429412 ns/op        3291.37 MB/s     1417240 B/op          2 allocs/op
BenchmarkTrie/Words/DumpSeq
BenchmarkTrie/Words/DumpSeq-8                 13          80093619 ns/op           5825074 keys/s              171.7 ns/key      5735346 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 13          85751448 ns/op           5440743 keys/s              183.8 ns/key      5735348 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 12          92787271 ns/op           5028182 keys/s              198.9 ns/key      5735346 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 13          84129345 ns/op           5545643 keys/s              180.3 ns/key      5735346 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 19          91355714 ns/op           5106970 keys/s              195.8 ns/key      5735346 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 12          93092534 ns/op           5011696 keys/s              199.5 ns/key      5735345 B/op     466528 allocs/op
BenchmarkTrie/Words/Lookup
BenchmarkTrie/Words/Lookup-8             3591576               303.9 ns/op         3290395 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3878244               306.8 ns/op         3259704 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3672697               306.8 ns/op         3259375 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3602612               314.4 ns/op         3180814 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3456769               308.4 ns/op         3242472 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3569264               291.9 ns/op         3426345 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01
BenchmarkTrie/Words/Lookup#01-8          5050815               259.7 ns/op         3850629 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          4201736               250.6 ns/op         3990896 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          4779867               256.6 ns/op         3896900 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          4090004               248.4 ns/op         4026495 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          4201806               251.9 ns/op         3969512 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          4174742               251.8 ns/op         3970836 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02
BenchmarkTrie/Words/Lookup#02-8          2940226               373.2 ns/op         2679589 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2954186               369.8 ns/op         2704362 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2873665               389.5 ns/op         2567267 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          3023365               373.9 ns/op         2674708 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2712358               385.1 ns/op         2596980 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2944796               374.7 ns/op         2668893 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup
BenchmarkTrie/Words/ReverseLookup-8     17242885                76.93 ns/op       12998403 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     14512213                73.73 ns/op       13563422 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     17269840                75.97 ns/op       13162605 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     16044164                72.85 ns/op       13726024 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     17326984                72.10 ns/op       13868735 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     16275063                73.53 ns/op       13599150 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup#01
BenchmarkTrie/Words/ReverseLookup#01-8   4911018               242.1 ns/op         4130659 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   5987607               208.0 ns/op         4806743 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   5288460               210.6 ns/op         4748177 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   6024813               192.9 ns/op         5182844 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   4939678               204.8 ns/op         4883386 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   7150965               195.9 ns/op         5105534 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3596            311403 ns/op              1887 keys/op         6059684 keys/s              165.0 ns/key        30832 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3747            338981 ns/op              1887 keys/op         5566681 keys/s              179.6 ns/key        30831 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3948            325432 ns/op              1887 keys/op         5798449 keys/s              172.5 ns/key        30832 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3763            338475 ns/op              1887 keys/op         5575003 keys/s              179.4 ns/key        30832 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3447            328944 ns/op              1887 keys/op         5736539 keys/s              174.3 ns/key        30832 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   4017            309685 ns/op              1887 keys/op         6093298 keys/s              164.1 ns/key        30832 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              404770              3059 ns/op                15.00 keys/op      4903359 keys/s              203.9 ns/key          408 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              388766              3135 ns/op                15.00 keys/op      4784330 keys/s              209.0 ns/key          408 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              372222              2883 ns/op                15.00 keys/op      5203091 keys/s              192.2 ns/key          408 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              416664              3081 ns/op                15.00 keys/op      4868902 keys/s              205.4 ns/key          408 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              315650              3495 ns/op                15.00 keys/op      4291590 keys/s              233.0 ns/key          408 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              399300              3218 ns/op                15.00 keys/op      4661288 keys/s              214.5 ns/key          408 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 349           3296566 ns/op             20325 keys/op         6165512 keys/s              162.2 ns/key       295496 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 309           3822962 ns/op             20325 keys/op         5316562 keys/s              188.1 ns/key       295496 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 330           3707643 ns/op             20325 keys/op         5481924 keys/s              182.4 ns/key       295496 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 304           3702192 ns/op             20325 keys/op         5489993 keys/s              182.1 ns/key       295496 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 326           3918272 ns/op             20325 keys/op         5187239 keys/s              192.8 ns/key       295496 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 320           3712399 ns/op             20325 keys/op         5474903 keys/s              182.7 ns/key       295496 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1118821              1076 ns/op                 3.000 keys/op     2787267 keys/s              358.8 ns/key          168 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8              963423              1061 ns/op                 3.000 keys/op     2828782 keys/s              353.5 ns/key          168 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1098853              1084 ns/op                 3.000 keys/op     2768614 keys/s              361.2 ns/key          168 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1000000              1047 ns/op                 3.000 keys/op     2865251 keys/s              349.0 ns/key          168 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8              995772              1054 ns/op                 3.000 keys/op     2846095 keys/s              351.4 ns/key          168 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1000000              1093 ns/op                 3.000 keys/op     2744299 keys/s              364.4 ns/key          168 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq
BenchmarkTrie/Words/CommonPrefixSearchSeq-8              1000000              1253 ns/op                 6.000 keys/op     4786926 keys/s              208.9 ns/key          224 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               917301              1225 ns/op                 6.000 keys/op     4897185 keys/s              204.2 ns/key          224 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               943780              1279 ns/op                 6.000 keys/op     4692276 keys/s              213.1 ns/key          224 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               996649              1280 ns/op                 6.000 keys/op     4688345 keys/s              213.3 ns/key          224 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               909326              1241 ns/op                 6.000 keys/op     4834952 keys/s              206.8 ns/key          224 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               857208              1183 ns/op                 6.000 keys/op     5070150 keys/s              197.2 ns/key          224 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1262419               928.7 ns/op               3.000 keys/op     3230441 keys/s              309.6 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1235666               918.4 ns/op               3.000 keys/op     3266595 keys/s              306.1 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1274683               955.8 ns/op               3.000 keys/op     3138736 keys/s              318.6 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1236302               947.3 ns/op               3.000 keys/op     3166899 keys/s              315.8 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1208889               963.7 ns/op               3.000 keys/op     3113166 keys/s              321.2 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1281438               889.9 ns/op               3.000 keys/op     3371259 keys/s              296.6 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/LookupAvg
BenchmarkTrie/Words/LookupAvg-8                                6         180811049 ns/op           2580325 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         172507004 ns/op           2704536 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         180639495 ns/op           2582777 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         171830419 ns/op           2715185 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         190132919 ns/op           2453817 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         195541247 ns/op           2385947 keys/s
BenchmarkTrie/Words/ReverseLookupAvg
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         177447354 ns/op           2629237 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         179790335 ns/op           2594973 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         193644190 ns/op           2409322 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         177309370 ns/op           2631282 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         7         177628487 ns/op           2626555 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         7         184673908 ns/op           2526350 keys/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         652793846 ns/op               515.4 ns/result      466550 queries/op      1940356 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         659625289 ns/op               520.8 ns/result      466550 queries/op      1920261 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         726194547 ns/op               573.3 ns/result      466550 queries/op      1744235 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         716498580 ns/op               565.7 ns/result      466550 queries/op      1767838 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         691871840 ns/op               546.2 ns/result      466550 queries/op      1830762 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         735115416 ns/op               580.4 ns/result      466550 queries/op      1723067 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 3         470499252 ns/op               371.5 ns/result      466550 queries/op      2692148 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 3         478943811 ns/op               378.1 ns/result      466550 queries/op      2644676 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 3         452405558 ns/op               357.2 ns/result      466550 queries/op      2799819 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 3         476457248 ns/op               376.2 ns/result      466550 queries/op      2658483 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 3         469134183 ns/op               370.4 ns/result      466550 queries/op      2699976 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 3         476506564 ns/op               376.2 ns/result      466550 queries/op      2658208 result/s
BenchmarkTrie/Go125
BenchmarkTrie/Go125/Build
BenchmarkTrie/Go125/Build-8                                   43          37646753 ns/op          17.23 MB/s         16092 keys/op          427448 keys/s             2339 ns/key       12502872 B/op         56 allocs/op
BenchmarkTrie/Go125/Build-8                                   31          37197105 ns/op          17.44 MB/s         16092 keys/op          432615 keys/s             2312 ns/key       12502882 B/op         56 allocs/op
BenchmarkTrie/Go125/Build-8                                   32          36635815 ns/op          17.71 MB/s         16092 keys/op          439243 keys/s             2277 ns/key       12502889 B/op         56 allocs/op
BenchmarkTrie/Go125/Build-8                                   32          37729820 ns/op          17.19 MB/s         16092 keys/op          426507 keys/s             2345 ns/key       12502879 B/op         56 allocs/op
BenchmarkTrie/Go125/Build-8                                   27          37973393 ns/op          17.08 MB/s         16092 keys/op          423771 keys/s             2360 ns/key       12502872 B/op         56 allocs/op
BenchmarkTrie/Go125/Build-8                                   33          37926918 ns/op          17.10 MB/s         16092 keys/op          424290 keys/s             2357 ns/key       12502875 B/op         56 allocs/op
BenchmarkTrie/Go125/ReadFrom
BenchmarkTrie/Go125/ReadFrom-8                             15763             82554 ns/op        1044.16 MB/s           137.0 reads/op     354285 B/op         48 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                             14658             88458 ns/op         974.47 MB/s           137.0 reads/op     354288 B/op         48 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                             14946             84089 ns/op        1025.10 MB/s           137.0 reads/op     354280 B/op         48 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                             14324             80139 ns/op        1075.63 MB/s           137.0 reads/op     354271 B/op         48 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                             14874             87386 ns/op         986.43 MB/s           137.0 reads/op     354303 B/op         48 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                             12841             88004 ns/op         979.50 MB/s           137.0 reads/op     354284 B/op         48 allocs/op
BenchmarkTrie/Go125/WriteTo
BenchmarkTrie/Go125/WriteTo-8                             412526              2497 ns/op        34517.92 MB/s          137.0 writes/op       248 B/op         17 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             484747              2194 ns/op        39280.25 MB/s          137.0 writes/op       248 B/op         17 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             511597              2389 ns/op        36080.35 MB/s          137.0 writes/op       248 B/op         17 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             527263              2588 ns/op        33307.89 MB/s          137.0 writes/op       248 B/op         17 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             489252              2499 ns/op        34491.95 MB/s          137.0 writes/op       248 B/op         17 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             403378              2628 ns/op        32795.70 MB/s          137.0 writes/op       248 B/op         17 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary
BenchmarkTrie/Go125/UnmarshalBinary-8                      14502             89933 ns/op         958.49 MB/s      452264 B/op         39 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      13028             93100 ns/op         925.88 MB/s      452264 B/op         39 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      12513             93344 ns/op         923.46 MB/s      452264 B/op         39 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      12400             89657 ns/op         961.44 MB/s      452264 B/op         39 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      12498             90767 ns/op         949.69 MB/s      452264 B/op         39 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      12140             95777 ns/op         900.00 MB/s      452264 B/op         39 allocs/op
BenchmarkTrie/Go125/MarshalBinary
BenchmarkTrie/Go125/MarshalBinary-8                        37618             30558 ns/op        2820.89 MB/s       90136 B/op          2 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        40099             28732 ns/op        3000.12 MB/s       90136 B/op          2 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        41032             29649 ns/op        2907.33 MB/s       90136 B/op          2 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        40633             29565 ns/op        2915.61 MB/s       90136 B/op          2 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        40734             30320 ns/op        2843.03 MB/s       90136 B/op          2 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        41013             31139 ns/op        2768.24 MB/s       90136 B/op          2 allocs/op
BenchmarkTrie/Go125/DumpSeq
BenchmarkTrie/Go125/DumpSeq-8                                180           6887813 ns/op           2336310 keys/s              428.0 ns/key       750272 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                182           7140290 ns/op           2253693 keys/s              443.7 ns/key       750272 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                152           7216598 ns/op           2229870 keys/s              448.5 ns/key       750272 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                160           7240833 ns/op           2222400 keys/s              450.0 ns/key       750272 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                169           7294653 ns/op           2206008 keys/s              453.3 ns/key       750272 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                163           6609583 ns/op           2434657 keys/s              410.7 ns/key       750272 B/op      16097 allocs/op
BenchmarkTrie/Go125/Lookup
BenchmarkTrie/Go125/Lookup-8                            16855837                73.15 ns/op       13671236 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            16687834                70.23 ns/op       14239209 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            17863987                67.58 ns/op       14796384 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            16065151                73.48 ns/op       13608476 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            16404421                66.86 ns/op       14955620 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            15225795                70.40 ns/op       14203702 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01
BenchmarkTrie/Go125/Lookup#01-8                          2000698               529.2 ns/op         1889830 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2378797               526.5 ns/op         1899185 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2111280               529.0 ns/op         1890285 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2106999               531.2 ns/op         1882485 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          1987053               529.5 ns/op         1888587 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2131227               534.2 ns/op         1871869 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02
BenchmarkTrie/Go125/Lookup#02-8                           673056              1520 ns/op            657851 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           759825              1493 ns/op            669933 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           817104              1508 ns/op            663234 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           721834              1521 ns/op            657507 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           708585              1541 ns/op            648754 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           748088              1574 ns/op            635377 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/ReverseLookup
BenchmarkTrie/Go125/ReverseLookup-8                      8963806               129.6 ns/op         7716348 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                     10120360               133.0 ns/op         7520819 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      9699117               147.7 ns/op         6771416 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      8578368               133.1 ns/op         7513962 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      9228799               133.9 ns/op         7466172 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      9258109               138.4 ns/op         7223583 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01
BenchmarkTrie/Go125/ReverseLookup#01-8                   1371218               820.9 ns/op         1218159 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1486902               780.2 ns/op         1281790 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1465855               844.6 ns/op         1183976 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1373847               886.4 ns/op         1128109 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1572044               778.8 ns/op         1284039 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1872037               803.6 ns/op         1244386 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3504            351330 ns/op               820.0 keys/op       2333991 keys/s              428.5 ns/key        54592 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3034            366107 ns/op               820.0 keys/op       2239782 keys/s              446.5 ns/key        54592 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3339            358049 ns/op               820.0 keys/op       2290192 keys/s              436.6 ns/key        54592 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3510            359968 ns/op               820.0 keys/op       2277985 keys/s              439.0 ns/key        54592 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3477            348642 ns/op               820.0 keys/op       2351982 keys/s              425.2 ns/key        54592 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3384            334540 ns/op               820.0 keys/op       2451132 keys/s              408.0 ns/key        54592 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              176230              5980 ns/op                17.00 keys/op      2843013 keys/s              351.7 ns/key          616 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              185502              6388 ns/op                17.00 keys/op      2661253 keys/s              375.8 ns/key          616 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              217963              6287 ns/op                17.00 keys/op      2703906 keys/s              369.8 ns/key          616 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              187311              6397 ns/op                17.00 keys/op      2657390 keys/s              376.3 ns/key          616 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              195015              6066 ns/op                17.00 keys/op      2802366 keys/s              356.8 ns/key          616 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              189099              6122 ns/op                17.00 keys/op      2776777 keys/s              360.1 ns/key          616 B/op         22 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               890347              1247 ns/op                 3.000 keys/op     2405066 keys/s              415.8 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               931302              1268 ns/op                 3.000 keys/op     2365265 keys/s              422.8 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               981854              1285 ns/op                 3.000 keys/op     2333760 keys/s              428.5 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               953187              1273 ns/op                 3.000 keys/op     2356165 keys/s              424.4 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               983785              1241 ns/op                 3.000 keys/op     2418167 keys/s              413.5 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               915574              1246 ns/op                 3.000 keys/op     2407394 keys/s              415.4 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            355838              3714 ns/op                11.00 keys/op      2962043 keys/s              337.6 ns/key          600 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            281996              3676 ns/op                11.00 keys/op      2992357 keys/s              334.2 ns/key          600 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            387500              3822 ns/op                11.00 keys/op      2877730 keys/s              347.5 ns/key          600 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            308749              3699 ns/op                11.00 keys/op      2973826 keys/s              336.3 ns/key          600 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            331814              3747 ns/op                11.00 keys/op      2935647 keys/s              340.6 ns/key          600 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            386952              3667 ns/op                11.00 keys/op      2999567 keys/s              333.4 ns/key          600 B/op         16 allocs/op
BenchmarkTrie/Go125/LookupAvg
BenchmarkTrie/Go125/LookupAvg-8                               74          15490486 ns/op           1038834 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               68          15547204 ns/op           1035045 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               75          15211767 ns/op           1057868 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               76          15136087 ns/op           1063157 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               74          15325484 ns/op           1050019 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               73          16206843 ns/op            992917 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg
BenchmarkTrie/Go125/ReverseLookupAvg-8                        67          19963397 ns/op            806079 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        52          20766560 ns/op            774902 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        46          22937932 ns/op            701548 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        72          21866252 ns/op            735930 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        58          22081039 ns/op            728771 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        60          19780670 ns/op            813524 keys/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  28          70549341 ns/op               761.0 ns/result       16092 queries/op      1314117 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  18          71370608 ns/op               769.8 ns/result       16092 queries/op      1298999 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  16          70361854 ns/op               758.9 ns/result       16092 queries/op      1317624 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  16          73773532 ns/op               795.7 ns/result       16092 queries/op      1256690 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  15          69863972 ns/op               753.6 ns/result       16092 queries/op      1327012 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  19          71381654 ns/op               769.9 ns/result       16092 queries/op      1298796 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                36          35123994 ns/op               378.9 ns/result       16092 queries/op      2639516 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                33          35571248 ns/op               383.7 ns/result       16092 queries/op      2606336 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                32          35752843 ns/op               385.6 ns/result       16092 queries/op      2593089 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                32          33161283 ns/op               357.7 ns/result       16092 queries/op      2795744 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                34          35350022 ns/op               381.3 ns/result       16092 queries/op      2622639 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                38          34938096 ns/op               376.9 ns/result       16092 queries/op      2653556 result/s
PASS
ok      github.com/pgaskin/go-marisa    441.215s
```

### wazero

37cde21a937175ff0928e5e6035b00a08be91daa

```
goos: linux
goarch: amd64
pkg: github.com/pgaskin/go-marisa
cpu: Intel(R) Core(TM) Ultra 7 258V
BenchmarkTrie
BenchmarkTrie/Words
BenchmarkTrie/Words/Build
BenchmarkTrie/Words/Build-8            3         347165782 ns/op          12.66 MB/s        466550 keys/op         1343884 keys/s              744.1 ns/key     166750509 B/op       117 allocs/op
BenchmarkTrie/Words/Build-8            3         339028773 ns/op          12.97 MB/s        466550 keys/op         1376138 keys/s              726.7 ns/key     166750370 B/op       115 allocs/op
BenchmarkTrie/Words/Build-8            3         364908097 ns/op          12.05 MB/s        466550 keys/op         1278542 keys/s              782.1 ns/key     166750461 B/op       116 allocs/op
BenchmarkTrie/Words/Build-8            4         323412373 ns/op          13.59 MB/s        466550 keys/op         1442587 keys/s              693.2 ns/key     166750452 B/op       116 allocs/op
BenchmarkTrie/Words/Build-8            3         351637035 ns/op          12.50 MB/s        466550 keys/op         1326795 keys/s              753.7 ns/key     166750461 B/op       116 allocs/op
BenchmarkTrie/Words/Build-8            3         359995845 ns/op          12.21 MB/s        466550 keys/op         1295988 keys/s              771.6 ns/key     166750461 B/op       116 allocs/op
BenchmarkTrie/Words/ReadFrom
BenchmarkTrie/Words/ReadFrom-8              1255           1080353 ns/op        1308.23 MB/s           138.0 reads/op    4295796 B/op         91 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1010           1057952 ns/op        1335.93 MB/s           138.0 reads/op    4295829 B/op         91 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1063           1078561 ns/op        1310.40 MB/s           138.0 reads/op    4295906 B/op         91 allocs/op
BenchmarkTrie/Words/ReadFrom-8               940           1097303 ns/op        1288.02 MB/s           138.0 reads/op    4295671 B/op         91 allocs/op
BenchmarkTrie/Words/ReadFrom-8               957           1148261 ns/op        1230.86 MB/s           138.0 reads/op    4295839 B/op         91 allocs/op
BenchmarkTrie/Words/ReadFrom-8              1128           1061916 ns/op        1330.95 MB/s           138.0 reads/op    4295753 B/op         91 allocs/op
BenchmarkTrie/Words/WriteTo
BenchmarkTrie/Words/WriteTo-8              96801             11263 ns/op        125485.49 MB/s         138.0 writes/op       325 B/op         20 allocs/op
BenchmarkTrie/Words/WriteTo-8              94635             12881 ns/op        109722.59 MB/s         138.0 writes/op       325 B/op         20 allocs/op
BenchmarkTrie/Words/WriteTo-8             106137             11382 ns/op        124173.54 MB/s         138.0 writes/op       325 B/op         20 allocs/op
BenchmarkTrie/Words/WriteTo-8             124113             11266 ns/op        125457.07 MB/s         138.0 writes/op       325 B/op         20 allocs/op
BenchmarkTrie/Words/WriteTo-8              98773             11384 ns/op        124148.13 MB/s         138.0 writes/op       325 B/op         20 allocs/op
BenchmarkTrie/Words/WriteTo-8             109686             11885 ns/op        118916.54 MB/s         138.0 writes/op       325 B/op         20 allocs/op
BenchmarkTrie/Words/UnmarshalBinary
BenchmarkTrie/Words/UnmarshalBinary-8       3109            338096 ns/op        4180.33 MB/s     1684130 B/op         77 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       3481            315437 ns/op        4480.61 MB/s     1684127 B/op         77 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       3343            314815 ns/op        4489.47 MB/s     1684128 B/op         77 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       3640            303326 ns/op        4659.51 MB/s     1684129 B/op         77 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       3454            331769 ns/op        4260.04 MB/s     1684128 B/op         77 allocs/op
BenchmarkTrie/Words/UnmarshalBinary-8       3301            350604 ns/op        4031.19 MB/s     1684131 B/op         77 allocs/op
BenchmarkTrie/Words/MarshalBinary
BenchmarkTrie/Words/MarshalBinary-8         2852            418678 ns/op        3375.75 MB/s     1417292 B/op          3 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2484            448975 ns/op        3147.95 MB/s     1417293 B/op          3 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2781            437084 ns/op        3233.59 MB/s     1417293 B/op          3 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2600            437176 ns/op        3232.91 MB/s     1417293 B/op          3 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2548            421597 ns/op        3352.37 MB/s     1417293 B/op          3 allocs/op
BenchmarkTrie/Words/MarshalBinary-8         2630            438415 ns/op        3223.78 MB/s     1417293 B/op          3 allocs/op
BenchmarkTrie/Words/DumpSeq
BenchmarkTrie/Words/DumpSeq-8                 13          85500766 ns/op           5456698 keys/s              183.3 ns/key      5735363 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 16          90228634 ns/op           5170766 keys/s              193.4 ns/key      5735364 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 13          91489801 ns/op           5099491 keys/s              196.1 ns/key      5735366 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 15          99458345 ns/op           4690919 keys/s              213.2 ns/key      5735376 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 13          94614513 ns/op           4931075 keys/s              202.8 ns/key      5735363 B/op     466528 allocs/op
BenchmarkTrie/Words/DumpSeq-8                 14          84677997 ns/op           5509724 keys/s              181.5 ns/key      5735366 B/op     466528 allocs/op
BenchmarkTrie/Words/Lookup
BenchmarkTrie/Words/Lookup-8             3417460               342.6 ns/op         2918549 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3669522               379.2 ns/op         2637111 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             2964512               344.9 ns/op         2899212 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3194532               343.5 ns/op         2910947 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3273039               341.8 ns/op         2925830 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup-8             3109461               367.0 ns/op         2725000 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01
BenchmarkTrie/Words/Lookup#01-8          4471070               281.7 ns/op         3549832 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          3806240               283.0 ns/op         3533914 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          3597780               278.8 ns/op         3586596 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          3549194               282.7 ns/op         3536811 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          3799068               282.0 ns/op         3545998 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#01-8          4098265               278.6 ns/op         3589373 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02
BenchmarkTrie/Words/Lookup#02-8          2870958               414.1 ns/op         2415144 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2838796               390.6 ns/op         2560332 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2857773               392.2 ns/op         2550035 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2840725               423.9 ns/op         2358988 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          2809348               399.3 ns/op         2504343 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/Lookup#02-8          3092194               387.7 ns/op         2579180 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup
BenchmarkTrie/Words/ReverseLookup-8      8808794               118.0 ns/op         8475988 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8      9604789               123.7 ns/op         8081880 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8      8281345               129.0 ns/op         7754710 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8      7353034               137.0 ns/op         7299838 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     10127605               116.2 ns/op         8604043 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup-8     10634683               123.6 ns/op         8092974 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Words/ReverseLookup#01
BenchmarkTrie/Words/ReverseLookup#01-8   4781504               274.2 ns/op         3646543 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   4835442               249.2 ns/op         4012618 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   4443799               275.4 ns/op         3631724 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   5070644               228.0 ns/op         4386820 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   4460179               252.5 ns/op         3959635 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/ReverseLookup#01-8   4153220               269.1 ns/op         3716488 keys/s              8 B/op          1 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq
BenchmarkTrie/Words/PredictiveSearchSeq-8                   2912            385811 ns/op              1887 keys/op         4890999 keys/s              204.5 ns/key        30848 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3091            346153 ns/op              1887 keys/op         5451359 keys/s              183.4 ns/key        30848 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3022            371087 ns/op              1887 keys/op         5085064 keys/s              196.7 ns/key        30848 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3147            436337 ns/op              1887 keys/op         4324645 keys/s              231.2 ns/key        30848 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   2956            385769 ns/op              1887 keys/op         4891526 keys/s              204.4 ns/key        30848 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq-8                   3248            400803 ns/op              1887 keys/op         4708055 keys/s              212.4 ns/key        30848 B/op       1892 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              317250              3773 ns/op                15.00 keys/op      3975123 keys/s              251.6 ns/key          424 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              334020              3816 ns/op                15.00 keys/op      3930882 keys/s              254.4 ns/key          424 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              306553              3820 ns/op                15.00 keys/op      3926867 keys/s              254.7 ns/key          424 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              341082              3755 ns/op                15.00 keys/op      3994556 keys/s              250.3 ns/key          424 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              428509              3916 ns/op                15.00 keys/op      3830191 keys/s              261.1 ns/key          424 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#01-8              298834              3726 ns/op                15.00 keys/op      4025774 keys/s              248.4 ns/key          424 B/op         20 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 262           4069874 ns/op             20325 keys/op         4994017 keys/s              200.2 ns/key       295512 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 325           3692686 ns/op             20325 keys/op         5504131 keys/s              181.7 ns/key       295512 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 270           4014345 ns/op             20325 keys/op         5063100 keys/s              197.5 ns/key       295512 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 286           3835237 ns/op             20325 keys/op         5299549 keys/s              188.7 ns/key       295512 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 288           3816629 ns/op             20325 keys/op         5325387 keys/s              187.8 ns/key       295512 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#02-8                 391           3992396 ns/op             20325 keys/op         5090930 keys/s              196.4 ns/key       295512 B/op      20330 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03
BenchmarkTrie/Words/PredictiveSearchSeq#03-8              975985              1207 ns/op                 3.000 keys/op     2484627 keys/s              402.5 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8              972820              1157 ns/op                 3.000 keys/op     2593287 keys/s              385.6 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1000000              1175 ns/op                 3.000 keys/op     2552388 keys/s              391.8 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1019797              1176 ns/op                 3.000 keys/op     2550308 keys/s              392.1 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8              989252              1173 ns/op                 3.000 keys/op     2558525 keys/s              390.9 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/PredictiveSearchSeq#03-8             1000000              1237 ns/op                 3.000 keys/op     2425939 keys/s              412.2 ns/key          184 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               668370              1515 ns/op                 6.000 keys/op     3959140 keys/s              252.6 ns/key          240 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               770920              1516 ns/op                 6.000 keys/op     3958132 keys/s              252.6 ns/key          240 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               694478              1507 ns/op                 6.000 keys/op     3980908 keys/s              251.2 ns/key          240 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               774152              1602 ns/op                 6.000 keys/op     3746133 keys/s              266.9 ns/key          240 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               789375              1606 ns/op                 6.000 keys/op     3736191 keys/s              267.7 ns/key          240 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq-8               884073              1430 ns/op                 6.000 keys/op     4195895 keys/s              238.3 ns/key          240 B/op         11 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1135462              1148 ns/op                 3.000 keys/op     2613977 keys/s              382.6 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1090623              1169 ns/op                 3.000 keys/op     2567205 keys/s              389.5 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1059226              1094 ns/op                 3.000 keys/op     2741520 keys/s              364.8 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8            922218              1161 ns/op                 3.000 keys/op     2582886 keys/s              387.2 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1000000              1100 ns/op                 3.000 keys/op     2727827 keys/s              366.6 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Words/CommonPrefixSearchSeq#01-8           1000000              1122 ns/op                 3.000 keys/op     2673726 keys/s              374.0 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Words/LookupAvg
BenchmarkTrie/Words/LookupAvg-8                                6         187162312 ns/op           2492764 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         187442736 ns/op           2489035 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         195998919 ns/op           2380376 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         201886940 ns/op           2310954 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         197859084 ns/op           2357998 keys/s
BenchmarkTrie/Words/LookupAvg-8                                6         208578872 ns/op           2236810 keys/s
BenchmarkTrie/Words/ReverseLookupAvg
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         210897214 ns/op           2212221 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         205713620 ns/op           2267964 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         206685270 ns/op           2257302 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         228692416 ns/op           2040081 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         213005181 ns/op           2190327 keys/s
BenchmarkTrie/Words/ReverseLookupAvg-8                         6         214671266 ns/op           2173328 keys/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         750665428 ns/op               592.6 ns/result      466550 queries/op      1687375 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         755504404 ns/op               596.5 ns/result      466550 queries/op      1676567 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         731095042 ns/op               577.2 ns/result      466550 queries/op      1732545 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         745081004 ns/op               588.2 ns/result      466550 queries/op      1700024 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         738629062 ns/op               583.1 ns/result      466550 queries/op      1714868 result/s
BenchmarkTrie/Words/PredictiveSearchSeqAvg-8                   2         728606562 ns/op               575.2 ns/result      466550 queries/op      1738457 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 2         578049694 ns/op               456.4 ns/result      466550 queries/op      2191258 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 2         551381234 ns/op               435.3 ns/result      466550 queries/op      2297265 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 2         573241556 ns/op               452.6 ns/result      466550 queries/op      2209638 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 2         546797460 ns/op               431.7 ns/result      466550 queries/op      2316498 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 2         567420736 ns/op               448.0 ns/result      466550 queries/op      2232305 result/s
BenchmarkTrie/Words/CommonPrefixSearchSeqAvg-8                 2         583669761 ns/op               460.8 ns/result      466550 queries/op      2170156 result/s
BenchmarkTrie/Go125
BenchmarkTrie/Go125/Build
BenchmarkTrie/Go125/Build-8                                   45          28367329 ns/op          22.87 MB/s         16092 keys/op          567273 keys/s             1763 ns/key       12592739 B/op        102 allocs/op
BenchmarkTrie/Go125/Build-8                                   44          26728483 ns/op          24.27 MB/s         16092 keys/op          602055 keys/s             1661 ns/key       12592722 B/op        102 allocs/op
BenchmarkTrie/Go125/Build-8                                   39          27905113 ns/op          23.25 MB/s         16092 keys/op          576669 keys/s             1734 ns/key       12592752 B/op        102 allocs/op
BenchmarkTrie/Go125/Build-8                                   44          28497723 ns/op          22.76 MB/s         16092 keys/op          564677 keys/s             1771 ns/key       12592751 B/op        102 allocs/op
BenchmarkTrie/Go125/Build-8                                   39          28065462 ns/op          23.12 MB/s         16092 keys/op          573375 keys/s             1744 ns/key       12592696 B/op        102 allocs/op
BenchmarkTrie/Go125/Build-8                                   38          26633911 ns/op          24.36 MB/s         16092 keys/op          604193 keys/s             1655 ns/key       12592743 B/op        102 allocs/op
BenchmarkTrie/Go125/ReadFrom
BenchmarkTrie/Go125/ReadFrom-8                              9865            126217 ns/op         682.95 MB/s           137.0 reads/op     378529 B/op         85 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                              9194            122698 ns/op         702.54 MB/s           137.0 reads/op     378516 B/op         85 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                              8086            130801 ns/op         659.01 MB/s           137.0 reads/op     378527 B/op         85 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                              8530            125438 ns/op         687.19 MB/s           137.0 reads/op     378518 B/op         85 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                              6871            146454 ns/op         588.58 MB/s           137.0 reads/op     378527 B/op         85 allocs/op
BenchmarkTrie/Go125/ReadFrom-8                              9618            124809 ns/op         690.66 MB/s           137.0 reads/op     378525 B/op         85 allocs/op
BenchmarkTrie/Go125/WriteTo
BenchmarkTrie/Go125/WriteTo-8                             101606             13865 ns/op        6217.09 MB/s           137.0 writes/op       296 B/op         18 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             100120             11114 ns/op        7756.09 MB/s           137.0 writes/op       296 B/op         18 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             102912             11737 ns/op        7343.99 MB/s           137.0 writes/op       296 B/op         18 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             111945             10735 ns/op        8029.99 MB/s           137.0 writes/op       296 B/op         18 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             108744             11473 ns/op        7513.33 MB/s           137.0 writes/op       296 B/op         18 allocs/op
BenchmarkTrie/Go125/WriteTo-8                             103572             11980 ns/op        7195.10 MB/s           137.0 writes/op       296 B/op         18 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary
BenchmarkTrie/Go125/UnmarshalBinary-8                      11949             86174 ns/op        1000.30 MB/s      389786 B/op         77 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      13069             98017 ns/op         879.43 MB/s      389787 B/op         77 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      12727            103417 ns/op         833.52 MB/s      389787 B/op         77 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      10000            105034 ns/op         820.69 MB/s      389789 B/op         77 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      13172             95781 ns/op         899.97 MB/s      389787 B/op         77 allocs/op
BenchmarkTrie/Go125/UnmarshalBinary-8                      12710             92122 ns/op         935.71 MB/s      389786 B/op         77 allocs/op
BenchmarkTrie/Go125/MarshalBinary
BenchmarkTrie/Go125/MarshalBinary-8                        24432             52504 ns/op        1641.78 MB/s       90184 B/op          3 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        24882             51861 ns/op        1662.15 MB/s       90184 B/op          3 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        23378             50080 ns/op        1721.24 MB/s       90184 B/op          3 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        22108             51077 ns/op        1687.66 MB/s       90184 B/op          3 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        22761             54276 ns/op        1588.18 MB/s       90184 B/op          3 allocs/op
BenchmarkTrie/Go125/MarshalBinary-8                        24374             51141 ns/op        1685.54 MB/s       90184 B/op          3 allocs/op
BenchmarkTrie/Go125/DumpSeq
BenchmarkTrie/Go125/DumpSeq-8                                154           7272338 ns/op           2212773 keys/s              451.9 ns/key       750288 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                157           7354434 ns/op           2188077 keys/s              457.0 ns/key       750288 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                153           7726308 ns/op           2082764 keys/s              480.1 ns/key       750288 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                163           7646978 ns/op           2104369 keys/s              475.2 ns/key       750287 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                158           7251026 ns/op           2219284 keys/s              450.6 ns/key       750288 B/op      16097 allocs/op
BenchmarkTrie/Go125/DumpSeq-8                                157           7483748 ns/op           2150268 keys/s              465.1 ns/key       750288 B/op      16097 allocs/op
BenchmarkTrie/Go125/Lookup
BenchmarkTrie/Go125/Lookup-8                            11883002                96.79 ns/op       10331429 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            11572773               101.0 ns/op         9900463 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            13040584                96.75 ns/op       10335615 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            10979232                97.09 ns/op       10299652 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            11791927                96.61 ns/op       10351408 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup-8                            11684952                96.41 ns/op       10372106 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01
BenchmarkTrie/Go125/Lookup#01-8                          2134261               535.3 ns/op         1868106 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2116249               522.5 ns/op         1913954 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2054878               538.9 ns/op         1855679 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2071251               537.0 ns/op         1862038 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2048132               548.3 ns/op         1823870 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#01-8                          2336852               561.5 ns/op         1781040 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02
BenchmarkTrie/Go125/Lookup#02-8                           793310              1477 ns/op            676825 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           783474              1434 ns/op            697434 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           793610              1442 ns/op            693548 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           718929              1438 ns/op            695392 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           771572              1444 ns/op            692676 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/Lookup#02-8                           803154              1551 ns/op            644851 keys/s              0 B/op          0 allocs/op
BenchmarkTrie/Go125/ReverseLookup
BenchmarkTrie/Go125/ReverseLookup-8                      7496821               172.3 ns/op         5803521 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      6452373               178.9 ns/op         5590395 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      6670972               172.9 ns/op         5782188 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      5665378               183.8 ns/op         5441096 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      6529042               174.8 ns/op         5720796 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup-8                      6525120               173.8 ns/op         5754792 keys/s              3 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01
BenchmarkTrie/Go125/ReverseLookup#01-8                   1134693               882.2 ns/op         1133567 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1465351               861.4 ns/op         1160838 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1482246               918.3 ns/op         1089017 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1336860               898.5 ns/op         1112921 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1309036               883.9 ns/op         1131302 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/ReverseLookup#01-8                   1502832               968.6 ns/op         1032437 keys/s             32 B/op          1 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   4424            390508 ns/op               820.0 keys/op       2099833 keys/s              476.2 ns/key        54608 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   2592            401573 ns/op               820.0 keys/op       2041972 keys/s              489.7 ns/key        54608 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   2618            387775 ns/op               820.0 keys/op       2114629 keys/s              472.9 ns/key        54608 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3138            407456 ns/op               820.0 keys/op       2012492 keys/s              496.9 ns/key        54608 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3200            377784 ns/op               820.0 keys/op       2170554 keys/s              460.7 ns/key        54608 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq-8                   3088            384917 ns/op               820.0 keys/op       2130334 keys/s              469.4 ns/key        54608 B/op        825 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              156410              7036 ns/op                17.00 keys/op      2416178 keys/s              413.9 ns/key          632 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              176371              6940 ns/op                17.00 keys/op      2449462 keys/s              408.3 ns/key          632 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              192555              6375 ns/op                17.00 keys/op      2666729 keys/s              375.0 ns/key          632 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              173628              6914 ns/op                17.00 keys/op      2458727 keys/s              406.7 ns/key          632 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              185493              7535 ns/op                17.00 keys/op      2256225 keys/s              443.2 ns/key          632 B/op         22 allocs/op
BenchmarkTrie/Go125/PredictiveSearchSeq#01-8              181202              7025 ns/op                17.00 keys/op      2420055 keys/s              413.2 ns/key          632 B/op         22 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               765229              1428 ns/op                 3.000 keys/op     2100260 keys/s              476.1 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               841008              1374 ns/op                 3.000 keys/op     2183579 keys/s              458.0 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               838351              1438 ns/op                 3.000 keys/op     2085805 keys/s              479.4 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8              1000000              1455 ns/op                 3.000 keys/op     2061955 keys/s              485.0 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               820123              1487 ns/op                 3.000 keys/op     2017802 keys/s              495.6 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq-8               736533              1409 ns/op                 3.000 keys/op     2129311 keys/s              469.6 ns/key          200 B/op          8 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            275232              4379 ns/op                11.00 keys/op      2512042 keys/s              398.1 ns/key          616 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            248614              4364 ns/op                11.00 keys/op      2520627 keys/s              396.7 ns/key          616 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            221556              4518 ns/op                11.00 keys/op      2434643 keys/s              410.7 ns/key          616 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            219063              4591 ns/op                11.00 keys/op      2395829 keys/s              417.4 ns/key          616 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            229464              4386 ns/op                11.00 keys/op      2508090 keys/s              398.7 ns/key          616 B/op         16 allocs/op
BenchmarkTrie/Go125/CommonPrefixSearchSeq#01-8            279040              4270 ns/op                11.00 keys/op      2575917 keys/s              388.2 ns/key          616 B/op         16 allocs/op
BenchmarkTrie/Go125/LookupAvg
BenchmarkTrie/Go125/LookupAvg-8                               81          14572292 ns/op           1104290 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               85          14559373 ns/op           1105270 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               76          14643194 ns/op           1098944 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               78          14589969 ns/op           1102951 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               78          14976815 ns/op           1074463 keys/s
BenchmarkTrie/Go125/LookupAvg-8                               76          14221821 ns/op           1131503 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg
BenchmarkTrie/Go125/ReverseLookupAvg-8                        68          22446567 ns/op            716904 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        46          21930739 ns/op            733767 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        51          21473996 ns/op            749374 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        69          20550296 ns/op            783056 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        68          21942890 ns/op            733360 keys/s
BenchmarkTrie/Go125/ReverseLookupAvg-8                        69          18985489 ns/op            847597 keys/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  28          68150305 ns/op               735.1 ns/result       16092 queries/op      1360378 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  19          68448149 ns/op               738.3 ns/result       16092 queries/op      1354460 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  27          71262687 ns/op               768.7 ns/result       16092 queries/op      1300965 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  14          71846096 ns/op               775.0 ns/result       16092 queries/op      1290400 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  18          73097336 ns/op               788.5 ns/result       16092 queries/op      1268314 result/s
BenchmarkTrie/Go125/PredictiveSearchSeqAvg-8                  15          68734911 ns/op               741.4 ns/result       16092 queries/op      1348812 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                26          40257507 ns/op               434.2 ns/result       16092 queries/op      2302939 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                30          40327540 ns/op               435.0 ns/result       16092 queries/op      2298937 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                25          41651443 ns/op               449.3 ns/result       16092 queries/op      2225865 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                26          39699305 ns/op               428.2 ns/result       16092 queries/op      2335313 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                28          36308154 ns/op               391.6 ns/result       16092 queries/op      2553430 result/s
BenchmarkTrie/Go125/CommonPrefixSearchSeqAvg-8                28          41281406 ns/op               445.3 ns/result       16092 queries/op      2245816 result/s
PASS
ok      github.com/pgaskin/go-marisa    428.632s
```

### benchstat

```
goos: linux
goarch: amd64
pkg: github.com/pgaskin/go-marisa
cpu: Intel(R) Core(TM) Ultra 7 258V
                                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt        │
                                      │        sec/op         │    sec/op     vs base               │
Trie/Words/Build-8                               349.4m ±  7%   406.9m ±  8%  +16.46% (p=0.002 n=6)
Trie/Words/ReadFrom-8                           1079.5µ ±  6%   995.9µ ±  9%   -7.74% (p=0.002 n=6)
Trie/Words/WriteTo-8                            11.383µ ± 13%   2.541µ ±  7%  -77.68% (p=0.002 n=6)
Trie/Words/UnmarshalBinary-8                     323.6µ ±  8%   634.8µ ±  7%  +96.16% (p=0.002 n=6)
Trie/Words/MarshalBinary-8                       437.1µ ±  4%   428.7µ ±  3%        ~ (p=0.699 n=6)
Trie/Words/DumpSeq-8                             90.86m ±  9%   88.55m ± 10%        ~ (p=0.485 n=6)
Trie/Words/Lookup-8                              344.2n ± 10%   306.8n ±  5%  -10.87% (p=0.002 n=6)
Trie/Words/Lookup#01-8                           281.8n ±  1%   251.8n ±  3%  -10.64% (p=0.002 n=6)
Trie/Words/Lookup#02-8                           395.8n ±  7%   374.3n ±  4%   -5.42% (p=0.004 n=6)
Trie/Words/ReverseLookup-8                      123.65n ± 11%   73.63n ±  4%  -40.45% (p=0.002 n=6)
Trie/Words/ReverseLookup#01-8                    260.8n ± 13%   206.4n ± 17%  -20.86% (p=0.004 n=6)
Trie/Words/PredictiveSearchSeq-8                 385.8µ ± 13%   327.2µ ±  5%  -15.19% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#01-8              3.795µ ±  3%   3.108µ ± 12%  -18.09% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#02-8              3.914m ±  6%   3.710m ± 11%        ~ (p=0.132 n=6)
Trie/Words/PredictiveSearchSeq#03-8              1.175µ ±  5%   1.069µ ±  2%   -9.10% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq-8               1.516µ ±  6%   1.247µ ±  5%  -17.72% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq#01-8           1135.0n ±  4%   938.0n ±  5%  -17.36% (p=0.002 n=6)
Trie/Words/LookupAvg-8                           196.9m ±  6%   180.7m ±  8%   -8.23% (p=0.026 n=6)
Trie/Words/ReverseLookupAvg-8                    212.0m ±  8%   178.7m ±  8%  -15.68% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeqAvg-8              741.9m ±  2%   704.2m ±  7%   -5.08% (p=0.009 n=6)
Trie/Words/CommonPrefixSearchSeqAvg-8            570.3m ±  4%   473.5m ±  4%  -16.98% (p=0.002 n=6)
Trie/Go125/Build-8                               27.99m ±  5%   37.69m ±  3%  +34.67% (p=0.002 n=6)
Trie/Go125/ReadFrom-8                           125.83µ ± 16%   85.74µ ±  7%  -31.86% (p=0.002 n=6)
Trie/Go125/WriteTo-8                            11.605µ ± 19%   2.498µ ± 12%  -78.47% (p=0.002 n=6)
Trie/Go125/UnmarshalBinary-8                     96.90µ ± 11%   91.93µ ±  4%        ~ (p=0.180 n=6)
Trie/Go125/MarshalBinary-8                       51.50µ ±  5%   29.98µ ±  4%  -41.78% (p=0.002 n=6)
Trie/Go125/DumpSeq-8                             7.419m ±  4%   7.178m ±  8%   -3.24% (p=0.009 n=6)
Trie/Go125/Lookup-8                              96.77n ±  4%   70.31n ±  5%  -27.34% (p=0.002 n=6)
Trie/Go125/Lookup#01-8                           538.0n ±  4%   529.4n ±  1%        ~ (p=0.065 n=6)
Trie/Go125/Lookup#02-8                           1.443µ ±  7%   1.521µ ±  4%   +5.37% (p=0.041 n=6)
Trie/Go125/ReverseLookup-8                       174.3n ±  5%   133.5n ± 11%  -23.41% (p=0.002 n=6)
Trie/Go125/ReverseLookup#01-8                    891.2n ±  9%   812.2n ±  9%   -8.86% (p=0.015 n=6)
Trie/Go125/PredictiveSearchSeq-8                 389.1µ ±  5%   354.7µ ±  6%   -8.85% (p=0.002 n=6)
Trie/Go125/PredictiveSearchSeq#01-8              6.983µ ±  9%   6.205µ ±  4%  -11.14% (p=0.009 n=6)
Trie/Go125/CommonPrefixSearchSeq-8               1.433µ ±  4%   1.258µ ±  2%  -12.25% (p=0.002 n=6)
Trie/Go125/CommonPrefixSearchSeq#01-8            4.383µ ±  5%   3.707µ ±  3%  -15.42% (p=0.002 n=6)
Trie/Go125/LookupAvg-8                           14.58m ±  3%   15.41m ±  5%   +5.67% (p=0.002 n=6)
Trie/Go125/ReverseLookupAvg-8                    21.70m ± 13%   21.32m ±  8%        ~ (p=1.000 n=6)
Trie/Go125/PredictiveSearchSeqAvg-8              70.00m ±  4%   70.96m ±  4%        ~ (p=0.485 n=6)
Trie/Go125/CommonPrefixSearchSeqAvg-8            40.29m ± 10%   35.24m ±  6%  -12.55% (p=0.002 n=6)
geomean                                          109.5µ         92.58µ        -15.45%

                             │ /tmp/bench.wazero.txt │         /tmp/bench.wasm2go.txt         │
                             │          B/s          │      B/s        vs base                │
Trie/Words/Build-8                     12.00Mi ±  8%    10.30Mi ±  7%   -14.11% (p=0.002 n=6)
Trie/Words/ReadFrom-8                  1.219Gi ±  6%    1.322Gi ±  9%    +8.39% (p=0.002 n=6)
Trie/Words/WriteTo-8                   115.6Gi ± 12%    518.1Gi ±  7%  +348.04% (p=0.002 n=6)
Trie/Words/UnmarshalBinary-8           4.070Gi ±  8%    2.074Gi ±  6%   -49.05% (p=0.002 n=6)
Trie/Words/MarshalBinary-8             3.011Gi ±  4%    3.071Gi ±  4%         ~ (p=0.699 n=6)
Trie/Go125/Build-8                     22.11Mi ±  5%    16.41Mi ±  3%   -25.77% (p=0.002 n=6)
Trie/Go125/ReadFrom-8                  653.3Mi ± 14%    959.2Mi ±  7%   +46.81% (p=0.002 n=6)
Trie/Go125/WriteTo-8                   6.918Gi ± 16%   32.135Gi ± 14%  +364.48% (p=0.002 n=6)
Trie/Go125/UnmarshalBinary-8           848.5Mi ± 12%    894.3Mi ±  4%         ~ (p=0.180 n=6)
Trie/Go125/MarshalBinary-8             1.559Gi ±  5%    2.678Gi ±  4%   +71.77% (p=0.002 n=6)
geomean                                1.096Gi          1.477Gi         +34.83%

                                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt        │
                                      │        keys/op        │   keys/op    vs base                │
Trie/Words/Build-8                                466.6k ± 0%   466.6k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq-8                  1.887k ± 0%   1.887k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq#01-8                15.00 ± 0%    15.00 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq#02-8               20.32k ± 0%   20.32k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq#03-8                3.000 ± 0%    3.000 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/CommonPrefixSearchSeq-8                 6.000 ± 0%    6.000 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/CommonPrefixSearchSeq#01-8              3.000 ± 0%    3.000 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/Build-8                                16.09k ± 0%   16.09k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/PredictiveSearchSeq-8                   820.0 ± 0%    820.0 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/PredictiveSearchSeq#01-8                17.00 ± 0%    17.00 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/CommonPrefixSearchSeq-8                 3.000 ± 0%    3.000 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/CommonPrefixSearchSeq#01-8              11.00 ± 0%    11.00 ± 0%       ~ (p=1.000 n=6) ¹
geomean                                            147.6         147.6       +0.00%
¹ all samples are equal

                                      │ /tmp/bench.wazero.txt │        /tmp/bench.wasm2go.txt        │
                                      │        keys/s         │    keys/s      vs base               │
Trie/Words/Build-8                               1.335M ±  8%    1.147M ±  7%  -14.13% (p=0.002 n=6)
Trie/Words/DumpSeq-8                             5.135M ±  9%    5.274M ± 10%        ~ (p=0.485 n=6)
Trie/Words/Lookup-8                              2.905M ±  9%    3.260M ±  5%  +12.20% (p=0.002 n=6)
Trie/Words/Lookup#01-8                           3.548M ±  1%    3.970M ±  3%  +11.90% (p=0.002 n=6)
Trie/Words/Lookup#02-8                           2.527M ±  7%    2.672M ±  4%   +5.72% (p=0.004 n=6)
Trie/Words/ReverseLookup-8                       8.087M ± 10%   13.581M ±  4%  +67.93% (p=0.002 n=6)
Trie/Words/ReverseLookup#01-8                    3.838M ± 14%    4.845M ± 15%  +26.24% (p=0.004 n=6)
Trie/Words/PredictiveSearchSeq-8                 4.891M ± 12%    5.767M ±  6%  +17.91% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#01-8              3.953M ±  3%    4.827M ± 11%  +22.10% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#02-8              5.195M ±  6%    5.478M ± 13%        ~ (p=0.132 n=6)
Trie/Words/PredictiveSearchSeq#03-8              2.551M ±  5%    2.808M ±  2%  +10.06% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq-8               3.959M ±  6%    4.811M ±  5%  +21.53% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq#01-8            2.644M ±  4%    3.199M ±  5%  +20.99% (p=0.002 n=6)
Trie/Words/LookupAvg-8                           2.369M ±  6%    2.582M ±  8%   +8.96% (p=0.026 n=6)
Trie/Words/ReverseLookupAvg-8                    2.201M ±  7%    2.611M ±  8%  +18.60% (p=0.002 n=6)
Trie/Go125/Build-8                               575.0k ±  5%    427.0k ±  3%  -25.75% (p=0.002 n=6)
Trie/Go125/DumpSeq-8                             2.169M ±  4%    2.242M ±  9%   +3.35% (p=0.009 n=6)
Trie/Go125/Lookup-8                              10.33M ±  4%    14.22M ±  5%  +37.62% (p=0.002 n=6)
Trie/Go125/Lookup#01-8                           1.859M ±  4%    1.889M ±  1%        ~ (p=0.065 n=6)
Trie/Go125/Lookup#02-8                           693.1k ±  7%    657.7k ±  3%   -5.11% (p=0.041 n=6)
Trie/Go125/ReverseLookup-8                       5.738M ±  5%    7.490M ± 10%  +30.54% (p=0.002 n=6)
Trie/Go125/ReverseLookup#01-8                    1.122M ±  8%    1.231M ±  8%   +9.73% (p=0.015 n=6)
Trie/Go125/PredictiveSearchSeq-8                 2.107M ±  4%    2.312M ±  6%   +9.72% (p=0.002 n=6)
Trie/Go125/PredictiveSearchSeq#01-8              2.435M ± 10%    2.740M ±  4%  +12.55% (p=0.009 n=6)
Trie/Go125/CommonPrefixSearchSeq-8               2.093M ±  4%    2.385M ±  2%  +13.96% (p=0.002 n=6)
Trie/Go125/CommonPrefixSearchSeq#01-8            2.510M ±  5%    2.968M ±  3%  +18.24% (p=0.002 n=6)
Trie/Go125/LookupAvg-8                           1.104M ±  3%    1.044M ±  5%   -5.36% (p=0.002 n=6)
Trie/Go125/ReverseLookupAvg-8                    741.6k ± 14%    755.4k ±  8%        ~ (p=1.000 n=6)
geomean                                          2.525M          2.802M        +10.97%

                                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt        │
                                      │        sec/key        │   sec/key     vs base               │
Trie/Words/Build-8                               748.9n ±  7%   872.2n ±  8%  +16.46% (p=0.002 n=6)
Trie/Words/DumpSeq-8                             194.8n ±  9%   189.8n ± 10%        ~ (p=0.485 n=6)
Trie/Words/PredictiveSearchSeq-8                 204.4n ± 13%   173.4n ±  5%  -15.19% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#01-8              253.0n ±  3%   207.2n ± 12%  -18.10% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#02-8              192.6n ±  6%   182.5n ± 11%        ~ (p=0.132 n=6)
Trie/Words/PredictiveSearchSeq#03-8              391.9n ±  5%   356.1n ±  2%   -9.13% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq-8               252.6n ±  6%   207.9n ±  5%  -17.72% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq#01-8            378.3n ±  4%   312.7n ±  5%  -17.34% (p=0.002 n=6)
Trie/Go125/Build-8                               1.739µ ±  5%   2.342µ ±  3%  +34.68% (p=0.002 n=6)
Trie/Go125/DumpSeq-8                             461.1n ±  4%   446.1n ±  8%   -3.24% (p=0.009 n=6)
Trie/Go125/PredictiveSearchSeq-8                 474.6n ±  5%   432.5n ±  6%   -8.85% (p=0.002 n=6)
Trie/Go125/PredictiveSearchSeq#01-8              410.8n ±  9%   365.0n ±  4%  -11.15% (p=0.009 n=6)
Trie/Go125/CommonPrefixSearchSeq-8               477.7n ±  4%   419.3n ±  2%  -12.23% (p=0.002 n=6)
Trie/Go125/CommonPrefixSearchSeq#01-8            398.4n ±  5%   336.9n ±  3%  -15.42% (p=0.002 n=6)
geomean                                          385.8n         358.8n         -7.00%

                                      │ /tmp/bench.wazero.txt │        /tmp/bench.wasm2go.txt         │
                                      │         B/op          │     B/op      vs base                 │
Trie/Words/Build-8                             159.0Mi ± 0%     158.9Mi ± 0%   -0.08% (p=0.002 n=6)
Trie/Words/ReadFrom-8                          4.097Mi ± 0%     4.074Mi ± 0%   -0.57% (p=0.002 n=6)
Trie/Words/WriteTo-8                             325.0 ± 0%       277.0 ± 0%  -14.77% (p=0.002 n=6)
Trie/Words/UnmarshalBinary-8                   1.606Mi ± 0%     3.041Mi ± 0%  +89.32% (p=0.002 n=6)
Trie/Words/MarshalBinary-8                     1.352Mi ± 0%     1.352Mi ± 0%   -0.00% (p=0.002 n=6)
Trie/Words/DumpSeq-8                           5.470Mi ± 0%     5.470Mi ± 0%   -0.00% (p=0.002 n=6)
Trie/Words/Lookup-8                              0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/Lookup#01-8                           0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/Lookup#02-8                           0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/ReverseLookup-8                       0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/ReverseLookup#01-8                    8.000 ± 0%       8.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq-8               30.12Ki ± 0%     30.11Ki ± 0%   -0.05% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#01-8              424.0 ± 0%       408.0 ± 0%   -3.77% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#02-8            288.6Ki ± 0%     288.6Ki ± 0%   -0.01% (p=0.002 n=6)
Trie/Words/PredictiveSearchSeq#03-8              184.0 ± 0%       168.0 ± 0%   -8.70% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq-8               240.0 ± 0%       224.0 ± 0%   -6.67% (p=0.002 n=6)
Trie/Words/CommonPrefixSearchSeq#01-8            200.0 ± 0%       184.0 ± 0%   -8.00% (p=0.002 n=6)
Trie/Go125/Build-8                             12.01Mi ± 0%     11.92Mi ± 0%   -0.71% (p=0.002 n=6)
Trie/Go125/ReadFrom-8                          369.7Ki ± 0%     346.0Ki ± 0%   -6.40% (p=0.002 n=6)
Trie/Go125/WriteTo-8                             296.0 ± 0%       248.0 ± 0%  -16.22% (p=0.002 n=6)
Trie/Go125/UnmarshalBinary-8                   380.7Ki ± 0%     441.7Ki ± 0%  +16.03% (p=0.002 n=6)
Trie/Go125/MarshalBinary-8                     88.07Ki ± 0%     88.02Ki ± 0%   -0.05% (p=0.002 n=6)
Trie/Go125/DumpSeq-8                           732.7Ki ± 0%     732.7Ki ± 0%   -0.00% (p=0.002 n=6)
Trie/Go125/Lookup-8                              0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/Lookup#01-8                           0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/Lookup#02-8                           0.000 ± 0%       0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/ReverseLookup-8                       3.000 ± 0%       3.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/ReverseLookup#01-8                    32.00 ± 0%       32.00 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/PredictiveSearchSeq-8               53.33Ki ± 0%     53.31Ki ± 0%   -0.03% (p=0.002 n=6)
Trie/Go125/PredictiveSearchSeq#01-8              632.0 ± 0%       616.0 ± 0%   -2.53% (p=0.002 n=6)
Trie/Go125/CommonPrefixSearchSeq-8               200.0 ± 0%       184.0 ± 0%   -8.00% (p=0.002 n=6)
Trie/Go125/CommonPrefixSearchSeq#01-8            616.0 ± 0%       600.0 ± 0%   -2.60% (p=0.002 n=6)
geomean                                                     ²                  -0.15%               ²
¹ all samples are equal
² summaries must be >0 to compute geomean

                                      │ /tmp/bench.wazero.txt │        /tmp/bench.wasm2go.txt        │
                                      │       allocs/op       │  allocs/op   vs base                 │
Trie/Words/Build-8                              116.00 ± 1%      66.00 ± 0%  -43.10% (p=0.002 n=6)
Trie/Words/ReadFrom-8                            91.00 ± 0%      53.00 ± 0%  -41.76% (p=0.002 n=6)
Trie/Words/WriteTo-8                             20.00 ± 0%      19.00 ± 0%   -5.00% (p=0.002 n=6)
Trie/Words/UnmarshalBinary-8                     77.00 ± 0%      38.00 ± 0%  -50.65% (p=0.002 n=6)
Trie/Words/MarshalBinary-8                       3.000 ± 0%      2.000 ± 0%  -33.33% (p=0.002 n=6)
Trie/Words/DumpSeq-8                            466.5k ± 0%     466.5k ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/Lookup-8                              0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/Lookup#01-8                           0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/Lookup#02-8                           0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/ReverseLookup-8                       0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/ReverseLookup#01-8                    1.000 ± 0%      1.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq-8                1.892k ± 0%     1.892k ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq#01-8              20.00 ± 0%      20.00 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq#02-8             20.33k ± 0%     20.33k ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/PredictiveSearchSeq#03-8              8.000 ± 0%      8.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/CommonPrefixSearchSeq-8               11.00 ± 0%      11.00 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Words/CommonPrefixSearchSeq#01-8            8.000 ± 0%      8.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/Build-8                              102.00 ± 0%      56.00 ± 0%  -45.10% (p=0.002 n=6)
Trie/Go125/ReadFrom-8                            85.00 ± 0%      48.00 ± 0%  -43.53% (p=0.002 n=6)
Trie/Go125/WriteTo-8                             18.00 ± 0%      17.00 ± 0%   -5.56% (p=0.002 n=6)
Trie/Go125/UnmarshalBinary-8                     77.00 ± 0%      39.00 ± 0%  -49.35% (p=0.002 n=6)
Trie/Go125/MarshalBinary-8                       3.000 ± 0%      2.000 ± 0%  -33.33% (p=0.002 n=6)
Trie/Go125/DumpSeq-8                            16.10k ± 0%     16.10k ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/Lookup-8                              0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/Lookup#01-8                           0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/Lookup#02-8                           0.000 ± 0%      0.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/ReverseLookup-8                       1.000 ± 0%      1.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/ReverseLookup#01-8                    1.000 ± 0%      1.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/PredictiveSearchSeq-8                 825.0 ± 0%      825.0 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/PredictiveSearchSeq#01-8              22.00 ± 0%      22.00 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/CommonPrefixSearchSeq-8               8.000 ± 0%      8.000 ± 0%        ~ (p=1.000 n=6) ¹
Trie/Go125/CommonPrefixSearchSeq#01-8            16.00 ± 0%      16.00 ± 0%        ~ (p=1.000 n=6) ¹
geomean                                                     ²                -13.34%               ²
¹ all samples are equal
² summaries must be >0 to compute geomean

                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt       │
                      │       reads/op        │  reads/op   vs base                │
Trie/Words/ReadFrom-8              138.0 ± 0%   138.0 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/ReadFrom-8              137.0 ± 0%   137.0 ± 0%       ~ (p=1.000 n=6) ¹
geomean                            137.5        137.5       +0.00%
¹ all samples are equal

                     │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt       │
                     │       writes/op       │ writes/op   vs base                │
Trie/Words/WriteTo-8              138.0 ± 0%   138.0 ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/WriteTo-8              137.0 ± 0%   137.0 ± 0%       ~ (p=1.000 n=6) ¹
geomean                           137.5        137.5       +0.00%
¹ all samples are equal

                                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt       │
                                      │      sec/result       │ sec/result   vs base               │
Trie/Words/PredictiveSearchSeqAvg-8              585.6n ±  2%   556.0n ± 7%   -5.07% (p=0.009 n=6)
Trie/Words/CommonPrefixSearchSeqAvg-8            450.3n ±  4%   373.9n ± 4%  -16.98% (p=0.002 n=6)
Trie/Go125/PredictiveSearchSeqAvg-8              755.1n ±  4%   765.4n ± 4%        ~ (p=0.485 n=6)
Trie/Go125/CommonPrefixSearchSeqAvg-8            434.6n ± 10%   380.1n ± 6%  -12.54% (p=0.002 n=6)
geomean                                          542.4n         495.9n        -8.57%

                                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt        │
                                      │      queries/op       │ queries/op   vs base                │
Trie/Words/PredictiveSearchSeqAvg-8               466.6k ± 0%   466.6k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Words/CommonPrefixSearchSeqAvg-8             466.6k ± 0%   466.6k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/PredictiveSearchSeqAvg-8               16.09k ± 0%   16.09k ± 0%       ~ (p=1.000 n=6) ¹
Trie/Go125/CommonPrefixSearchSeqAvg-8             16.09k ± 0%   16.09k ± 0%       ~ (p=1.000 n=6) ¹
geomean                                           86.65k        86.65k       +0.00%
¹ all samples are equal

                                      │ /tmp/bench.wazero.txt │       /tmp/bench.wasm2go.txt       │
                                      │       result/s        │  result/s    vs base               │
Trie/Words/PredictiveSearchSeqAvg-8              1.707M ±  2%   1.799M ± 8%   +5.38% (p=0.009 n=6)
Trie/Words/CommonPrefixSearchSeqAvg-8            2.221M ±  4%   2.675M ± 5%  +20.46% (p=0.002 n=6)
Trie/Go125/PredictiveSearchSeqAvg-8              1.325M ±  4%   1.307M ± 4%        ~ (p=0.485 n=6)
Trie/Go125/CommonPrefixSearchSeqAvg-8            2.301M ± 11%   2.631M ± 6%  +14.35% (p=0.002 n=6)
geomean                                          1.844M         2.017M        +9.38%
```
