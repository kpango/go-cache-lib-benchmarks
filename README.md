# go-cache-lib-benchmarks
Go's cache libraries benchmarking

## Benchmark Results
```
go test -count=3 -run=NONE -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/kpango/go-cache-lib-benchmarks
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkDefaultMapSetSmallDataNoTTL-16     	 3012236	       692.0 ns/op	       4 B/op	       0 allocs/op
BenchmarkDefaultMapSetSmallDataNoTTL-16     	 2802657	       711.5 ns/op	       4 B/op	       0 allocs/op
BenchmarkDefaultMapSetSmallDataNoTTL-16     	 2824888	       753.2 ns/op	       4 B/op	       0 allocs/op
BenchmarkDefaultMapSetBigDataNoTTL-16       	   14371	    205540 ns/op	    1168 B/op	      23 allocs/op
BenchmarkDefaultMapSetBigDataNoTTL-16       	   12906	    183396 ns/op	    1489 B/op	      26 allocs/op
BenchmarkDefaultMapSetBigDataNoTTL-16       	   13338	    131959 ns/op	    1056 B/op	      25 allocs/op
BenchmarkSyncMapSetSmallDataNoTTL-16        	 5122106	       220.9 ns/op	     194 B/op	      12 allocs/op
BenchmarkSyncMapSetSmallDataNoTTL-16        	 5062830	       221.4 ns/op	     194 B/op	      12 allocs/op
BenchmarkSyncMapSetSmallDataNoTTL-16        	 4679046	       232.7 ns/op	     195 B/op	      12 allocs/op
BenchmarkSyncMapSetBigDataNoTTL-16          	   13371	     77154 ns/op	   25631 B/op	    1561 allocs/op
BenchmarkSyncMapSetBigDataNoTTL-16          	   14481	     76220 ns/op	   25559 B/op	    1559 allocs/op
BenchmarkSyncMapSetBigDataNoTTL-16          	   14397	     97863 ns/op	   25565 B/op	    1559 allocs/op
BenchmarkGacheSetSmallDataNoTTL-16          	 1713711	       614.3 ns/op	     232 B/op	      12 allocs/op
BenchmarkGacheSetSmallDataNoTTL-16          	 1694001	       605.9 ns/op	     232 B/op	      12 allocs/op
BenchmarkGacheSetSmallDataNoTTL-16          	 1561509	       649.4 ns/op	     232 B/op	      12 allocs/op
BenchmarkGacheSetSmallDataWithTTL-16        	 1698248	       682.1 ns/op	     232 B/op	      12 allocs/op
BenchmarkGacheSetSmallDataWithTTL-16        	 1527564	       664.0 ns/op	     233 B/op	      12 allocs/op
BenchmarkGacheSetSmallDataWithTTL-16        	 1621671	       641.9 ns/op	     234 B/op	      12 allocs/op
BenchmarkGacheSetBigDataNoTTL-16            	   26379	     40756 ns/op	   29201 B/op	    1548 allocs/op
BenchmarkGacheSetBigDataNoTTL-16            	   24628	     41999 ns/op	   29237 B/op	    1549 allocs/op
BenchmarkGacheSetBigDataNoTTL-16            	   23731	     44838 ns/op	   29300 B/op	    1550 allocs/op
BenchmarkGacheSetBigDataWithTTL-16          	   28090	     39239 ns/op	   29176 B/op	    1547 allocs/op
BenchmarkGacheSetBigDataWithTTL-16          	   29594	     39251 ns/op	   29154 B/op	    1547 allocs/op
BenchmarkGacheSetBigDataWithTTL-16          	   29846	     40237 ns/op	   29154 B/op	    1547 allocs/op
BenchmarkTTLCacheSetSmallDataNoTTL-16       	  871501	      1414 ns/op	     207 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataNoTTL-16       	  722188	      1527 ns/op	     211 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataNoTTL-16       	  893289	      1508 ns/op	     207 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataWithTTL-16     	  343142	      4339 ns/op	     232 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataWithTTL-16     	  332598	      5279 ns/op	     234 B/op	       5 allocs/op
BenchmarkTTLCacheSetSmallDataWithTTL-16     	  353947	      4695 ns/op	     231 B/op	       4 allocs/op
BenchmarkTTLCacheSetBigDataNoTTL-16         	    4568	    315428 ns/op	   27510 B/op	     583 allocs/op
BenchmarkTTLCacheSetBigDataNoTTL-16         	    4700	    370003 ns/op	   27427 B/op	     581 allocs/op
BenchmarkTTLCacheSetBigDataNoTTL-16         	    4794	    304255 ns/op	   27375 B/op	     580 allocs/op
BenchmarkTTLCacheSetBigDataWithTTL-16       	    1999	    815092 ns/op	   31158 B/op	     674 allocs/op
BenchmarkTTLCacheSetBigDataWithTTL-16       	    2115	    805168 ns/op	   30808 B/op	     665 allocs/op
BenchmarkTTLCacheSetBigDataWithTTL-16       	    1910	    720750 ns/op	   31459 B/op	     681 allocs/op
BenchmarkGoCacheSetSmallDataNoTTL-16        	 1805215	      1325 ns/op	      71 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataNoTTL-16        	 2226256	      1272 ns/op	      70 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataNoTTL-16        	 1936894	      1199 ns/op	      71 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataWithTTL-16      	 1362333	      2219 ns/op	      74 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataWithTTL-16      	  712969	      2130 ns/op	      83 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataWithTTL-16      	 1253125	      1994 ns/op	      75 B/op	       4 allocs/op
BenchmarkGoCacheSetBigDataNoTTL-16          	    5029	    355360 ns/op	   10862 B/op	     576 allocs/op
BenchmarkGoCacheSetBigDataNoTTL-16          	    4454	    370977 ns/op	   11188 B/op	     585 allocs/op
BenchmarkGoCacheSetBigDataNoTTL-16          	    4893	    320557 ns/op	   10931 B/op	     578 allocs/op
BenchmarkGoCacheSetBigDataWithTTL-16        	    3470	    375104 ns/op	   12010 B/op	     605 allocs/op
BenchmarkGoCacheSetBigDataWithTTL-16        	    3877	    436264 ns/op	   11622 B/op	     595 allocs/op
BenchmarkGoCacheSetBigDataWithTTL-16        	    3525	    425568 ns/op	   11945 B/op	     603 allocs/op
BenchmarkBigCacheSetSmallDataNoTTL-16       	  581522	      1945 ns/op	     227 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataNoTTL-16       	  588075	      2115 ns/op	     321 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataNoTTL-16       	  583336	      2018 ns/op	     227 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataWithTTL-16     	  557883	      2154 ns/op	     336 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataWithTTL-16     	  514293	      2039 ns/op	     251 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataWithTTL-16     	  553608	      2011 ns/op	     338 B/op	       8 allocs/op
BenchmarkBigCacheSetBigDataNoTTL-16         	    1778	    678977 ns/op	 1864206 B/op	    1717 allocs/op
BenchmarkBigCacheSetBigDataNoTTL-16         	    2577	    485414 ns/op	 1663685 B/op	    1661 allocs/op
BenchmarkBigCacheSetBigDataNoTTL-16         	    4208	    538862 ns/op	 1774536 B/op	    1613 allocs/op
BenchmarkBigCacheSetBigDataWithTTL-16       	    4063	    537079 ns/op	 1742579 B/op	    1586 allocs/op
BenchmarkBigCacheSetBigDataWithTTL-16       	    4334	    513776 ns/op	 1741159 B/op	    1610 allocs/op
BenchmarkBigCacheSetBigDataWithTTL-16       	    5673	    466756 ns/op	 1580498 B/op	    1593 allocs/op
BenchmarkFastCacheSetSmallDataNoTTL-16      	  948162	      1404 ns/op	      54 B/op	       4 allocs/op
BenchmarkFastCacheSetSmallDataNoTTL-16      	  850028	      1375 ns/op	      56 B/op	       4 allocs/op
BenchmarkFastCacheSetSmallDataNoTTL-16      	  904214	      1291 ns/op	      55 B/op	       4 allocs/op
BenchmarkFastCacheSetBigDataNoTTL-16        	    8720	    300558 ns/op	  796163 B/op	    2081 allocs/op
BenchmarkFastCacheSetBigDataNoTTL-16        	    7180	    287501 ns/op	  796503 B/op	    2092 allocs/op
BenchmarkFastCacheSetBigDataNoTTL-16        	    8733	    289308 ns/op	  796183 B/op	    2083 allocs/op
BenchmarkFreeCacheSetSmallDataNoTTL-16      	  908479	      1370 ns/op	     137 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataNoTTL-16      	  901101	      1337 ns/op	     137 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataNoTTL-16      	  905768	      1303 ns/op	     137 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataWithTTL-16    	  874141	      1343 ns/op	     138 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataWithTTL-16    	  895959	      1325 ns/op	     137 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataWithTTL-16    	  854461	      1396 ns/op	     138 B/op	       8 allocs/op
BenchmarkFreeCacheSetBigDataNoTTL-16        	    9096	    338709 ns/op	 1062326 B/op	    2595 allocs/op
BenchmarkFreeCacheSetBigDataNoTTL-16        	    8077	    372212 ns/op	 1062510 B/op	    2600 allocs/op
BenchmarkFreeCacheSetBigDataNoTTL-16        	    8865	    350124 ns/op	 1062371 B/op	    2596 allocs/op
BenchmarkFreeCacheSetBigDataWithTTL-16      	    8080	    359767 ns/op	 1060852 B/op	    2594 allocs/op
BenchmarkFreeCacheSetBigDataWithTTL-16      	    7945	    370194 ns/op	 1062261 B/op	    2600 allocs/op
BenchmarkFreeCacheSetBigDataWithTTL-16      	    7216	    322340 ns/op	 1062302 B/op	    2603 allocs/op
BenchmarkGCacheLRUSetSmallDataNoTTL-16      	  373542	      5408 ns/op	     731 B/op	      24 allocs/op
BenchmarkGCacheLRUSetSmallDataNoTTL-16      	  394716	      4710 ns/op	     731 B/op	      24 allocs/op
BenchmarkGCacheLRUSetSmallDataNoTTL-16      	  416959	      4960 ns/op	     730 B/op	      24 allocs/op
BenchmarkGCacheLRUSetSmallDataWithTTL-16    	  468116	      2937 ns/op	     317 B/op	      16 allocs/op
BenchmarkGCacheLRUSetSmallDataWithTTL-16    	  506085	      2765 ns/op	     315 B/op	      16 allocs/op
BenchmarkGCacheLRUSetSmallDataWithTTL-16    	  446923	      2863 ns/op	     319 B/op	      16 allocs/op
BenchmarkGCacheLRUSetBigDataNoTTL-16        	    2013	    693086 ns/op	  102769 B/op	    3233 allocs/op
BenchmarkGCacheLRUSetBigDataNoTTL-16        	    2065	    502612 ns/op	  102578 B/op	    3229 allocs/op
BenchmarkGCacheLRUSetBigDataNoTTL-16        	    2073	    626335 ns/op	  102613 B/op	    3228 allocs/op
BenchmarkGCacheLRUSetBigDataWithTTL-16      	    3151	    606643 ns/op	   99512 B/op	    3172 allocs/op
BenchmarkGCacheLRUSetBigDataWithTTL-16      	    2310	    640002 ns/op	  102492 B/op	    3213 allocs/op
BenchmarkGCacheLRUSetBigDataWithTTL-16      	    2881	    629191 ns/op	   99196 B/op	    3180 allocs/op
BenchmarkGCacheLFUSetSmallDataNoTTL-16      	  392299	      3996 ns/op	     559 B/op	      20 allocs/op
BenchmarkGCacheLFUSetSmallDataNoTTL-16      	  333810	      4152 ns/op	     568 B/op	      20 allocs/op
BenchmarkGCacheLFUSetSmallDataNoTTL-16      	  398437	      4092 ns/op	     559 B/op	      20 allocs/op
BenchmarkGCacheLFUSetSmallDataWithTTL-16    	  458202	      2677 ns/op	     318 B/op	      16 allocs/op
BenchmarkGCacheLFUSetSmallDataWithTTL-16    	  548630	      2357 ns/op	     313 B/op	      16 allocs/op
BenchmarkGCacheLFUSetSmallDataWithTTL-16    	  415496	      3588 ns/op	     322 B/op	      16 allocs/op
BenchmarkGCacheLFUSetBigDataNoTTL-16        	    1821	    719854 ns/op	   78837 B/op	    2742 allocs/op
BenchmarkGCacheLFUSetBigDataNoTTL-16        	    2148	    880979 ns/op	   77766 B/op	    2715 allocs/op
BenchmarkGCacheLFUSetBigDataNoTTL-16        	    1639	    737358 ns/op	   79591 B/op	    2761 allocs/op
BenchmarkGCacheLFUSetBigDataWithTTL-16      	    2304	    588965 ns/op	   74926 B/op	    2698 allocs/op
BenchmarkGCacheLFUSetBigDataWithTTL-16      	    2918	    699426 ns/op	   73735 B/op	    2669 allocs/op
BenchmarkGCacheLFUSetBigDataWithTTL-16      	    2840	    505000 ns/op	   73808 B/op	    2671 allocs/op
BenchmarkGCacheARCSetSmallDataNoTTL-16      	  285345	      6789 ns/op	     940 B/op	      28 allocs/op
BenchmarkGCacheARCSetSmallDataNoTTL-16      	  286674	      5909 ns/op	     942 B/op	      28 allocs/op
BenchmarkGCacheARCSetSmallDataNoTTL-16      	  287047	      7328 ns/op	     938 B/op	      28 allocs/op
BenchmarkGCacheARCSetSmallDataWithTTL-16    	  447154	      3143 ns/op	     319 B/op	      16 allocs/op
BenchmarkGCacheARCSetSmallDataWithTTL-16    	  412765	      3393 ns/op	     321 B/op	      16 allocs/op
BenchmarkGCacheARCSetSmallDataWithTTL-16    	  430100	      3362 ns/op	     320 B/op	      16 allocs/op
BenchmarkGCacheARCSetBigDataNoTTL-16        	    1188	    972261 ns/op	  132499 B/op	    3831 allocs/op
BenchmarkGCacheARCSetBigDataNoTTL-16        	     874	   1265679 ns/op	  136301 B/op	    3925 allocs/op
BenchmarkGCacheARCSetBigDataNoTTL-16        	    1563	   1237530 ns/op	  130100 B/op	    3766 allocs/op
BenchmarkGCacheARCSetBigDataWithTTL-16      	     919	   1199208 ns/op	  131411 B/op	    3883 allocs/op
BenchmarkGCacheARCSetBigDataWithTTL-16      	     728	   1444408 ns/op	  134604 B/op	    3969 allocs/op
BenchmarkGCacheARCSetBigDataWithTTL-16      	    1281	   1168172 ns/op	  127136 B/op	    3772 allocs/op
BenchmarkMCacheSetSmallDataNoTTL-16         	  153124	      8016 ns/op	    2476 B/op	      34 allocs/op
BenchmarkMCacheSetSmallDataNoTTL-16         	  178076	      7835 ns/op	    2487 B/op	      33 allocs/op
BenchmarkMCacheSetSmallDataNoTTL-16         	  157195	      8011 ns/op	    2489 B/op	      34 allocs/op
BenchmarkMCacheSetSmallDataWithTTL-16       	  130850	      8456 ns/op	    2171 B/op	      35 allocs/op
BenchmarkMCacheSetSmallDataWithTTL-16       	  128506	      8453 ns/op	    2168 B/op	      35 allocs/op
BenchmarkMCacheSetSmallDataWithTTL-16       	  129012	      8433 ns/op	    2167 B/op	      35 allocs/op
BenchmarkMCacheSetBigDataNoTTL-16           	    1369	   1041445 ns/op	  263449 B/op	    4331 allocs/op
BenchmarkMCacheSetBigDataNoTTL-16           	    1322	   1042766 ns/op	  263772 B/op	    4339 allocs/op
BenchmarkMCacheSetBigDataNoTTL-16           	    1636	    828306 ns/op	  278884 B/op	    4292 allocs/op
BenchmarkMCacheSetBigDataWithTTL-16         	    1262	   1130765 ns/op	  298609 B/op	    4475 allocs/op
BenchmarkMCacheSetBigDataWithTTL-16         	    1452	    870269 ns/op	  272295 B/op	    4415 allocs/op
BenchmarkMCacheSetBigDataWithTTL-16         	    1436	    951797 ns/op	  272399 B/op	    4417 allocs/op
BenchmarkBitcaskSetSmallDataNoTTL-16        	   48121	     25712 ns/op	     730 B/op	      30 allocs/op
BenchmarkBitcaskSetSmallDataNoTTL-16        	   51704	     20635 ns/op	     712 B/op	      30 allocs/op
BenchmarkBitcaskSetSmallDataNoTTL-16        	   64257	     25971 ns/op	     659 B/op	      29 allocs/op
BenchmarkBitcaskSetBigDataNoTTL-16          	    9225	    238228 ns/op	  787820 B/op	    1570 allocs/op
BenchmarkBitcaskSetBigDataNoTTL-16          	   11043	    252338 ns/op	  787591 B/op	    1564 allocs/op
BenchmarkBitcaskSetBigDataNoTTL-16          	    7423	    222093 ns/op	  788157 B/op	    1579 allocs/op
PASS
ok  	github.com/kpango/go-cache-lib-benchmarks	312.358s
go clean ./...
go clean -modcache
rm -rf ./*.log
rm -rf ./*.svg
rm -rf ./go.*
rm -rf pprof
rm -rf bench
rm -rf vendor
GO111MODULE=on go mod init
GO111MODULE=on go mod tidy
sleep 3
go test -count=3 -run=NONE -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/kpango/go-cache-lib-benchmarks
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkDefaultMapSetSmallDataNoTTL-16     	 2643079	      1285 ns/op	       6 B/op	       0 allocs/op
BenchmarkDefaultMapSetSmallDataNoTTL-16     	 2657940	      1114 ns/op	       5 B/op	       0 allocs/op
BenchmarkDefaultMapSetSmallDataNoTTL-16     	 2550610	       907.0 ns/op	       5 B/op	       0 allocs/op
BenchmarkDefaultMapSetBigDataNoTTL-16       	     846	   2015840 ns/op	   15223 B/op	     379 allocs/op
BenchmarkDefaultMapSetBigDataNoTTL-16       	     676	   1743838 ns/op	   19027 B/op	     474 allocs/op
BenchmarkDefaultMapSetBigDataNoTTL-16       	     640	   1725203 ns/op	   20111 B/op	     501 allocs/op
BenchmarkSyncMapSetSmallDataNoTTL-16        	 5187840	       236.3 ns/op	     194 B/op	      12 allocs/op
BenchmarkSyncMapSetSmallDataNoTTL-16        	 3360160	       300.4 ns/op	     198 B/op	      12 allocs/op
BenchmarkSyncMapSetSmallDataNoTTL-16        	 5226325	       338.4 ns/op	     194 B/op	      12 allocs/op
BenchmarkSyncMapSetBigDataNoTTL-16          	    1441	    771124 ns/op	  107483 B/op	    6370 allocs/op
BenchmarkSyncMapSetBigDataNoTTL-16          	    1374	    800952 ns/op	  107923 B/op	    6381 allocs/op
BenchmarkSyncMapSetBigDataNoTTL-16          	    1474	    796597 ns/op	  107280 B/op	    6365 allocs/op
BenchmarkGacheSetSmallDataNoTTL-16          	 4943024	       227.3 ns/op	      98 B/op	       4 allocs/op
BenchmarkGacheSetSmallDataNoTTL-16          	 4099231	       246.8 ns/op	      99 B/op	       4 allocs/op
BenchmarkGacheSetSmallDataNoTTL-16          	 5309982	       213.8 ns/op	      98 B/op	       4 allocs/op
BenchmarkGacheSetSmallDataWithTTL-16        	 6140960	       190.9 ns/op	      98 B/op	       4 allocs/op
BenchmarkGacheSetSmallDataWithTTL-16        	 6007318	       198.2 ns/op	      98 B/op	       4 allocs/op
BenchmarkGacheSetSmallDataWithTTL-16        	 6009673	       209.8 ns/op	      98 B/op	       4 allocs/op
BenchmarkGacheSetBigDataNoTTL-16            	    7651	    156950 ns/op	   50929 B/op	    2091 allocs/op
BenchmarkGacheSetBigDataNoTTL-16            	    7473	    161629 ns/op	   50948 B/op	    2091 allocs/op
BenchmarkGacheSetBigDataNoTTL-16            	    7255	    159363 ns/op	   50997 B/op	    2093 allocs/op
BenchmarkGacheSetBigDataWithTTL-16          	    7449	    150347 ns/op	   50905 B/op	    2091 allocs/op
BenchmarkGacheSetBigDataWithTTL-16          	    6404	    176525 ns/op	   51243 B/op	    2099 allocs/op
BenchmarkGacheSetBigDataWithTTL-16          	    7219	    155465 ns/op	   51027 B/op	    2093 allocs/op
BenchmarkTTLCacheSetSmallDataNoTTL-16       	  769513	      2008 ns/op	     210 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataNoTTL-16       	  745357	      2155 ns/op	     210 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataNoTTL-16       	  779490	      1924 ns/op	     209 B/op	       4 allocs/op
BenchmarkTTLCacheSetSmallDataWithTTL-16     	  219840	      6125 ns/op	     255 B/op	       5 allocs/op
BenchmarkTTLCacheSetSmallDataWithTTL-16     	  209316	      5533 ns/op	     258 B/op	       5 allocs/op
BenchmarkTTLCacheSetSmallDataWithTTL-16     	  216009	      5414 ns/op	     256 B/op	       5 allocs/op
BenchmarkTTLCacheSetBigDataNoTTL-16         	     723	   2226131 ns/op	  116978 B/op	    2498 allocs/op
BenchmarkTTLCacheSetBigDataNoTTL-16         	     579	   2095111 ns/op	  121593 B/op	    2610 allocs/op
BenchmarkTTLCacheSetBigDataNoTTL-16         	     595	   2115237 ns/op	  120978 B/op	    2595 allocs/op
BenchmarkTTLCacheSetBigDataWithTTL-16       	     466	   4296047 ns/op	  127253 B/op	    2747 allocs/op
BenchmarkTTLCacheSetBigDataWithTTL-16       	     444	   3744060 ns/op	  128648 B/op	    2781 allocs/op
BenchmarkTTLCacheSetBigDataWithTTL-16       	     403	   3708909 ns/op	  131710 B/op	    2855 allocs/op
BenchmarkGoCacheSetSmallDataNoTTL-16        	 1638081	      1422 ns/op	      72 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataNoTTL-16        	 1708333	      1476 ns/op	      72 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataNoTTL-16        	 1866837	      1661 ns/op	      71 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataWithTTL-16      	 1108484	      3057 ns/op	      76 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataWithTTL-16      	 1184559	      2593 ns/op	      75 B/op	       4 allocs/op
BenchmarkGoCacheSetSmallDataWithTTL-16      	  971601	      2505 ns/op	      78 B/op	       4 allocs/op
BenchmarkGoCacheSetBigDataNoTTL-16          	     541	   2117240 ns/op	   57284 B/op	    2641 allocs/op
BenchmarkGoCacheSetBigDataNoTTL-16          	     549	   2410998 ns/op	   56913 B/op	    2632 allocs/op
BenchmarkGoCacheSetBigDataNoTTL-16          	     604	   2167320 ns/op	   54711 B/op	    2579 allocs/op
BenchmarkGoCacheSetBigDataWithTTL-16        	     534	   2204202 ns/op	   57605 B/op	    2649 allocs/op
BenchmarkGoCacheSetBigDataWithTTL-16        	     489	   2212552 ns/op	   59884 B/op	    2704 allocs/op
BenchmarkGoCacheSetBigDataWithTTL-16        	     544	   2198998 ns/op	   57141 B/op	    2638 allocs/op
BenchmarkBigCacheSetSmallDataNoTTL-16       	  495056	      2231 ns/op	     259 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataNoTTL-16       	  531310	      2455 ns/op	     244 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataNoTTL-16       	  497934	      2230 ns/op	     333 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataWithTTL-16     	  550522	      2192 ns/op	     340 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataWithTTL-16     	  533551	      2278 ns/op	     349 B/op	       8 allocs/op
BenchmarkBigCacheSetSmallDataWithTTL-16     	  525224	      2262 ns/op	     247 B/op	       8 allocs/op
BenchmarkBigCacheSetBigDataNoTTL-16         	     160	   7337015 ns/op	27279060 B/op	    8165 allocs/op
BenchmarkBigCacheSetBigDataNoTTL-16         	     339	   6420407 ns/op	24246362 B/op	    7102 allocs/op
BenchmarkBigCacheSetBigDataNoTTL-16         	     358	   6458930 ns/op	24741548 B/op	    7018 allocs/op
BenchmarkBigCacheSetBigDataWithTTL-16       	     362	   6725875 ns/op	24106051 B/op	    6998 allocs/op
BenchmarkBigCacheSetBigDataWithTTL-16       	     350	   6706111 ns/op	25348638 B/op	    7069 allocs/op
BenchmarkBigCacheSetBigDataWithTTL-16       	     385	   6617622 ns/op	25009681 B/op	    6988 allocs/op
BenchmarkFastCacheSetSmallDataNoTTL-16      	  635128	      1582 ns/op	      62 B/op	       4 allocs/op
BenchmarkFastCacheSetSmallDataNoTTL-16      	  833314	      1634 ns/op	      56 B/op	       4 allocs/op
BenchmarkFastCacheSetSmallDataNoTTL-16      	  714148	      1538 ns/op	      59 B/op	       4 allocs/op
BenchmarkFastCacheSetBigDataNoTTL-16        	     379	   2857419 ns/op	12640979 B/op	    8477 allocs/op
BenchmarkFastCacheSetBigDataNoTTL-16        	     446	   2342194 ns/op	12634699 B/op	    8277 allocs/op
BenchmarkFastCacheSetBigDataNoTTL-16        	     478	   2674766 ns/op	12635025 B/op	    8372 allocs/op
BenchmarkFreeCacheSetSmallDataNoTTL-16      	  822547	      1390 ns/op	     140 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataNoTTL-16      	  809064	      1401 ns/op	     140 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataNoTTL-16      	  827200	      1386 ns/op	     139 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataWithTTL-16    	  830641	      1360 ns/op	     141 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataWithTTL-16    	  824226	      1409 ns/op	     139 B/op	       8 allocs/op
BenchmarkFreeCacheSetSmallDataWithTTL-16    	  823095	      1418 ns/op	     139 B/op	       8 allocs/op
BenchmarkFreeCacheSetBigDataNoTTL-16        	     399	   4260697 ns/op	16859147 B/op	   11044 allocs/op
BenchmarkFreeCacheSetBigDataNoTTL-16        	     385	   3957501 ns/op	16860350 B/op	   11073 allocs/op
BenchmarkFreeCacheSetBigDataNoTTL-16        	     368	   3881589 ns/op	16861921 B/op	   11112 allocs/op
BenchmarkFreeCacheSetBigDataWithTTL-16      	     364	   3527338 ns/op	16160911 B/op	   10449 allocs/op
BenchmarkFreeCacheSetBigDataWithTTL-16      	     415	   4047884 ns/op	16755916 B/op	   10919 allocs/op
BenchmarkFreeCacheSetBigDataWithTTL-16      	     368	   3837183 ns/op	16809060 B/op	   11065 allocs/op
BenchmarkGCacheLRUSetSmallDataNoTTL-16      	  264030	      6241 ns/op	     752 B/op	      24 allocs/op
BenchmarkGCacheLRUSetSmallDataNoTTL-16      	  253401	      7416 ns/op	     756 B/op	      24 allocs/op
BenchmarkGCacheLRUSetSmallDataNoTTL-16      	  239924	      5776 ns/op	     763 B/op	      24 allocs/op
BenchmarkGCacheLRUSetSmallDataWithTTL-16    	  550334	      2772 ns/op	     313 B/op	      16 allocs/op
BenchmarkGCacheLRUSetSmallDataWithTTL-16    	  291813	      3750 ns/op	     335 B/op	      17 allocs/op
BenchmarkGCacheLRUSetSmallDataWithTTL-16    	  491138	      3066 ns/op	     316 B/op	      16 allocs/op
BenchmarkGCacheLRUSetBigDataNoTTL-16        	     422	   4613691 ns/op	  417486 B/op	   13073 allocs/op
BenchmarkGCacheLRUSetBigDataNoTTL-16        	     417	   3928745 ns/op	  417473 B/op	   13080 allocs/op
BenchmarkGCacheLRUSetBigDataNoTTL-16        	     451	   4140610 ns/op	  416242 B/op	   13025 allocs/op
BenchmarkGCacheLRUSetBigDataWithTTL-16      	     416	   3963186 ns/op	  417799 B/op	   13082 allocs/op
BenchmarkGCacheLRUSetBigDataWithTTL-16      	     423	   3312618 ns/op	  416987 B/op	   13069 allocs/op
BenchmarkGCacheLRUSetBigDataWithTTL-16      	     421	   4188191 ns/op	  419304 B/op	   13080 allocs/op
BenchmarkGCacheLFUSetSmallDataNoTTL-16      	  245124	      5284 ns/op	     582 B/op	      21 allocs/op
BenchmarkGCacheLFUSetSmallDataNoTTL-16      	  293169	      5699 ns/op	     573 B/op	      20 allocs/op
BenchmarkGCacheLFUSetSmallDataNoTTL-16      	  242949	      5484 ns/op	     582 B/op	      21 allocs/op
BenchmarkGCacheLFUSetSmallDataWithTTL-16    	  513308	      3304 ns/op	     315 B/op	      16 allocs/op
BenchmarkGCacheLFUSetSmallDataWithTTL-16    	  490880	      2830 ns/op	     316 B/op	      16 allocs/op
BenchmarkGCacheLFUSetSmallDataWithTTL-16    	  274797	      4260 ns/op	     339 B/op	      17 allocs/op
BenchmarkGCacheLFUSetBigDataNoTTL-16        	     370	   3798405 ns/op	  323144 B/op	   11151 allocs/op
BenchmarkGCacheLFUSetBigDataNoTTL-16        	     380	   5331309 ns/op	  322630 B/op	   11130 allocs/op
BenchmarkGCacheLFUSetBigDataNoTTL-16        	     380	   4241376 ns/op	  322416 B/op	   11129 allocs/op
BenchmarkGCacheLFUSetBigDataWithTTL-16      	     400	   4825822 ns/op	  314999 B/op	   11059 allocs/op
BenchmarkGCacheLFUSetBigDataWithTTL-16      	     405	   4265683 ns/op	  314665 B/op	   11049 allocs/op
BenchmarkGCacheLFUSetBigDataWithTTL-16      	     387	   3207447 ns/op	  314595 B/op	   11076 allocs/op
BenchmarkGCacheARCSetSmallDataNoTTL-16      	  220232	      5641 ns/op	     958 B/op	      28 allocs/op
BenchmarkGCacheARCSetSmallDataNoTTL-16      	  210572	      8173 ns/op	     953 B/op	      28 allocs/op
BenchmarkGCacheARCSetSmallDataNoTTL-16      	  204333	      7552 ns/op	     958 B/op	      28 allocs/op
BenchmarkGCacheARCSetSmallDataWithTTL-16    	  245208	      4448 ns/op	     344 B/op	      17 allocs/op
BenchmarkGCacheARCSetSmallDataWithTTL-16    	  248721	      4467 ns/op	     344 B/op	      17 allocs/op
BenchmarkGCacheARCSetSmallDataWithTTL-16    	  231891	      5040 ns/op	     348 B/op	      17 allocs/op
BenchmarkGCacheARCSetBigDataNoTTL-16        	     204	   4958948 ns/op	  508843 B/op	   14879 allocs/op
BenchmarkGCacheARCSetBigDataNoTTL-16        	     252	   7074482 ns/op	  504757 B/op	   14717 allocs/op
BenchmarkGCacheARCSetBigDataNoTTL-16        	     220	   6063293 ns/op	  512894 B/op	   14920 allocs/op
BenchmarkGCacheARCSetBigDataWithTTL-16      	     243	   6324228 ns/op	  494903 B/op	   14743 allocs/op
BenchmarkGCacheARCSetBigDataWithTTL-16      	     228	   6782321 ns/op	  489347 B/op	   14584 allocs/op
BenchmarkGCacheARCSetBigDataWithTTL-16      	     253	   6377563 ns/op	  486799 B/op	   14475 allocs/op
BenchmarkMCacheSetSmallDataNoTTL-16         	  151471	      8412 ns/op	    2424 B/op	      34 allocs/op
BenchmarkMCacheSetSmallDataNoTTL-16         	  168376	      8946 ns/op	    2341 B/op	      33 allocs/op
BenchmarkMCacheSetSmallDataNoTTL-16         	  163966	      8748 ns/op	    2389 B/op	      34 allocs/op
BenchmarkMCacheSetSmallDataWithTTL-16       	  129789	      9139 ns/op	    2177 B/op	      35 allocs/op
BenchmarkMCacheSetSmallDataWithTTL-16       	  128959	      9070 ns/op	    2165 B/op	      35 allocs/op
BenchmarkMCacheSetSmallDataWithTTL-16       	  131551	      8933 ns/op	    2162 B/op	      35 allocs/op
BenchmarkMCacheSetBigDataNoTTL-16           	     344	   4396703 ns/op	 1199301 B/op	   17317 allocs/op
BenchmarkMCacheSetBigDataNoTTL-16           	     338	   3144073 ns/op	 1054836 B/op	   17334 allocs/op
BenchmarkMCacheSetBigDataNoTTL-16           	     339	   4248077 ns/op	 1054742 B/op	   17331 allocs/op
BenchmarkMCacheSetBigDataWithTTL-16         	     306	   5091143 ns/op	 1098553 B/op	   17844 allocs/op
BenchmarkMCacheSetBigDataWithTTL-16         	     327	   5154701 ns/op	 1161857 B/op	   17757 allocs/op
BenchmarkMCacheSetBigDataWithTTL-16         	     334	   3780142 ns/op	 1091573 B/op	   17723 allocs/op
BenchmarkBitcaskSetSmallDataNoTTL-16        	   62145	     22113 ns/op	     666 B/op	      29 allocs/op
BenchmarkBitcaskSetSmallDataNoTTL-16        	   56181	     31656 ns/op	     693 B/op	      29 allocs/op
BenchmarkBitcaskSetSmallDataNoTTL-16        	   41202	     28505 ns/op	     783 B/op	      32 allocs/op
BenchmarkBitcaskSetBigDataNoTTL-16          	     547	   1943178 ns/op	12606326 B/op	    6729 allocs/op
BenchmarkBitcaskSetBigDataNoTTL-16          	     609	   2189628 ns/op	12603941 B/op	    6669 allocs/op
BenchmarkBitcaskSetBigDataNoTTL-16          	     684	   2330591 ns/op	12601636 B/op	    6612 allocs/op
PASS
ok  	github.com/kpango/go-cache-lib-benchmarks	300.021s
```
