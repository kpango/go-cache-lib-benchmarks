# go-cache-lib-benchmarks
Go's cache libraries benchmarking

## Benchmark Results
```
go test -count=1 -timeout=30m -run=NONE -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/kpango/go-cache-lib-benchmarks
cpu: AMD Ryzen Threadripper 3990X 64-Core Processor 
BenchmarkDefaultMapSetGetSmallDataNoTTL/P100-128 	  322064	      5835 ns/op	     512 B/op	      32 allocs/op
BenchmarkDefaultMapSetGetSmallDataNoTTL/P1000-128         	  436360	      9178 ns/op	     512 B/op	      32 allocs/op
BenchmarkDefaultMapSetGetSmallDataNoTTL/P10000-128        	  274617	      7272 ns/op	     520 B/op	      32 allocs/op
BenchmarkDefaultMapSetGetBigDataNoTTL/P100-128            	       8	 176415972 ns/op	 4197799 B/op	  262224 allocs/op
BenchmarkDefaultMapSetGetBigDataNoTTL/P1000-128           	       7	 175827836 ns/op	 4205730 B/op	  262438 allocs/op
BenchmarkDefaultMapSetGetBigDataNoTTL/P10000-128          	       6	 179424900 ns/op	 4328889 B/op	  265513 allocs/op
BenchmarkSyncMapSetGetSmallDataNoTTL/P100-128             	  156650	      7640 ns/op	    1280 B/op	      48 allocs/op
BenchmarkSyncMapSetGetSmallDataNoTTL/P1000-128            	  154138	      8061 ns/op	    1281 B/op	      48 allocs/op
BenchmarkSyncMapSetGetSmallDataNoTTL/P10000-128           	  168157	      6913 ns/op	    1295 B/op	      48 allocs/op
BenchmarkSyncMapSetGetBigDataNoTTL/P100-128               	     133	  16134218 ns/op	10486085 B/op	  393221 allocs/op
BenchmarkSyncMapSetGetBigDataNoTTL/P1000-128              	      84	  13254026 ns/op	10486721 B/op	  393240 allocs/op
BenchmarkSyncMapSetGetBigDataNoTTL/P10000-128             	      92	  12951673 ns/op	10494551 B/op	  393436 allocs/op
BenchmarkGacheV2SetGetSmallDataNoTTL/P100-128             	 3101184	       378.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2SetGetSmallDataNoTTL/P1000-128            	 3394678	       355.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2SetGetSmallDataNoTTL/P10000-128           	 3285120	       355.7 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2SetGetSmallDataWithTTL/P100-128           	 3039324	       379.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2SetGetSmallDataWithTTL/P1000-128          	 3369207	       359.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2SetGetSmallDataWithTTL/P10000-128         	 3216296	       363.5 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2SetGetBigDataNoTTL/P100-128               	     165	  12452627 ns/op	     360 B/op	       4 allocs/op
BenchmarkGacheV2SetGetBigDataNoTTL/P1000-128              	     201	  14133239 ns/op	     626 B/op	      11 allocs/op
BenchmarkGacheV2SetGetBigDataNoTTL/P10000-128             	     194	  13792372 ns/op	    4394 B/op	     105 allocs/op
BenchmarkGacheV2SetGetBigDataWithTTL/P100-128             	     223	   5737008 ns/op	  128577 B/op	     861 allocs/op
BenchmarkGacheV2SetGetBigDataWithTTL/P1000-128            	     190	   7967436 ns/op	   96816 B/op	     849 allocs/op
BenchmarkGacheV2SetGetBigDataWithTTL/P10000-128           	     367	  11602256 ns/op	   50582 B/op	     482 allocs/op
BenchmarkGacheSetGetSmallDataNoTTL/P100-128               	 1437166	       826.8 ns/op	     640 B/op	      32 allocs/op
BenchmarkGacheSetGetSmallDataNoTTL/P1000-128              	 1390311	       831.1 ns/op	     640 B/op	      32 allocs/op
BenchmarkGacheSetGetSmallDataNoTTL/P10000-128             	 1445632	       849.6 ns/op	     640 B/op	      32 allocs/op
BenchmarkGacheSetGetSmallDataWithTTL/P100-128             	 1494949	       804.6 ns/op	     640 B/op	      32 allocs/op
BenchmarkGacheSetGetSmallDataWithTTL/P1000-128            	 1424524	       814.0 ns/op	     640 B/op	      32 allocs/op
BenchmarkGacheSetGetSmallDataWithTTL/P10000-128           	 1430383	       828.6 ns/op	     640 B/op	      32 allocs/op
BenchmarkGacheSetGetBigDataNoTTL/P100-128                 	     104	  10361095 ns/op	 5243028 B/op	  262148 allocs/op
BenchmarkGacheSetGetBigDataNoTTL/P1000-128                	     123	  10061366 ns/op	 5243512 B/op	  262160 allocs/op
BenchmarkGacheSetGetBigDataNoTTL/P10000-128               	     105	  10351891 ns/op	 5250537 B/op	  262336 allocs/op
BenchmarkGacheSetGetBigDataWithTTL/P100-128               	      66	  15364809 ns/op	 5542737 B/op	  264687 allocs/op
BenchmarkGacheSetGetBigDataWithTTL/P1000-128              	     108	  10858465 ns/op	 5367322 B/op	  263673 allocs/op
BenchmarkGacheSetGetBigDataWithTTL/P10000-128             	     112	   9236041 ns/op	 5393340 B/op	  263780 allocs/op
BenchmarkTTLCacheSetGetSmallDataNoTTL/P100-128            	  156762	      8828 ns/op	       0 B/op	       0 allocs/op
BenchmarkTTLCacheSetGetSmallDataNoTTL/P1000-128           	  314836	     11616 ns/op	       0 B/op	       0 allocs/op
BenchmarkTTLCacheSetGetSmallDataNoTTL/P10000-128          	  232116	      7722 ns/op	       8 B/op	       0 allocs/op
BenchmarkTTLCacheSetGetSmallDataWithTTL/P100-128          	   98041	     17947 ns/op	       0 B/op	       0 allocs/op
BenchmarkTTLCacheSetGetSmallDataWithTTL/P1000-128         	   59180	     17482 ns/op	       3 B/op	       0 allocs/op
BenchmarkTTLCacheSetGetSmallDataWithTTL/P10000-128        	   94988	     17563 ns/op	      20 B/op	       0 allocs/op
BenchmarkTTLCacheSetGetBigDataNoTTL/P100-128              	       6	 195622018 ns/op	    4465 B/op	      98 allocs/op
BenchmarkTTLCacheSetGetBigDataNoTTL/P1000-128             	       5	 207363361 ns/op	   16369 B/op	     409 allocs/op
BenchmarkTTLCacheSetGetBigDataNoTTL/P10000-128            	       5	 213898513 ns/op	  162192 B/op	    4049 allocs/op
BenchmarkTTLCacheSetGetBigDataWithTTL/P100-128            	       4	 288744365 ns/op	    5330 B/op	     135 allocs/op
BenchmarkTTLCacheSetGetBigDataWithTTL/P1000-128           	       4	 291442235 ns/op	   21514 B/op	     526 allocs/op
BenchmarkTTLCacheSetGetBigDataWithTTL/P10000-128          	       5	 271916601 ns/op	  162896 B/op	    4065 allocs/op
BenchmarkGoCacheSetGetSmallDataNoTTL/P100-128             	  348920	      3233 ns/op	     256 B/op	      16 allocs/op
BenchmarkGoCacheSetGetSmallDataNoTTL/P1000-128            	  651602	      5707 ns/op	     256 B/op	      16 allocs/op
BenchmarkGoCacheSetGetSmallDataNoTTL/P10000-128           	  391192	      3244 ns/op	     260 B/op	      16 allocs/op
BenchmarkGoCacheSetGetSmallDataWithTTL/P100-128           	  165358	      7878 ns/op	     256 B/op	      16 allocs/op
BenchmarkGoCacheSetGetSmallDataWithTTL/P1000-128          	  160317	     11738 ns/op	     257 B/op	      16 allocs/op
BenchmarkGoCacheSetGetSmallDataWithTTL/P10000-128         	  180841	      7519 ns/op	     267 B/op	      16 allocs/op
BenchmarkGoCacheSetGetBigDataNoTTL/P100-128               	       6	 181843878 ns/op	 2100937 B/op	  131162 allocs/op
BenchmarkGoCacheSetGetBigDataNoTTL/P1000-128              	       6	 211241572 ns/op	 2110981 B/op	  131421 allocs/op
BenchmarkGoCacheSetGetBigDataNoTTL/P10000-128             	       6	 173761968 ns/op	 2231672 B/op	  134440 allocs/op
BenchmarkGoCacheSetGetBigDataWithTTL/P100-128             	       5	 203293607 ns/op	 2100564 B/op	  131168 allocs/op
BenchmarkGoCacheSetGetBigDataWithTTL/P1000-128            	       5	 215343070 ns/op	 2112820 B/op	  131476 allocs/op
BenchmarkGoCacheSetGetBigDataWithTTL/P10000-128           	       6	 219402965 ns/op	 2231612 B/op	  134447 allocs/op
BenchmarkBigCacheSetGetSmallDataNoTTL/P100-128            	  139681	      9072 ns/op	   10092 B/op	      32 allocs/op
BenchmarkBigCacheSetGetSmallDataNoTTL/P1000-128           	  216896	      6053 ns/op	    7574 B/op	      32 allocs/op
BenchmarkBigCacheSetGetSmallDataNoTTL/P10000-128          	  168037	      6709 ns/op	   16319 B/op	      32 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkBigCacheSetGetSmallDataWithTTL/P100-128          	  142272	      8160 ns/op	    5734 B/op	      32 allocs/op
BenchmarkBigCacheSetGetSmallDataWithTTL/P1000-128         	  197328	      5647 ns/op	    4705 B/op	      32 allocs/op
BenchmarkBigCacheSetGetSmallDataWithTTL/P10000-128        	  226819	      5507 ns/op	    2056 B/op	      32 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkBigCacheSetGetBigDataNoTTL/P100-128              	       1	1252380183 ns/op	1281736608 B/op	  264818 allocs/op
BenchmarkBigCacheSetGetBigDataNoTTL/P1000-128             	       1	1491902273 ns/op	1872409632 B/op	  265262 allocs/op
BenchmarkBigCacheSetGetBigDataNoTTL/P10000-128            	       1	1529820799 ns/op	2662698208 B/op	  283329 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkBigCacheSetGetBigDataWithTTL/P100-128            	       1	1146731984 ns/op	1205235112 B/op	  227443 allocs/op
BenchmarkBigCacheSetGetBigDataWithTTL/P1000-128           	       1	1285453353 ns/op	1427759176 B/op	  245789 allocs/op
BenchmarkBigCacheSetGetBigDataWithTTL/P10000-128          	       1	1145637657 ns/op	537675176 B/op	  282345 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkGCacheLRUSetGetSmallDataNoTTL/P100-128           	   78538	     14491 ns/op	    2770 B/op	      92 allocs/op
BenchmarkGCacheLRUSetGetSmallDataNoTTL/P1000-128          	   64812	     15904 ns/op	    2769 B/op	      92 allocs/op
BenchmarkGCacheLRUSetGetSmallDataNoTTL/P10000-128         	   80349	     14253 ns/op	    2780 B/op	      93 allocs/op
BenchmarkGCacheLRUSetGetSmallDataWithTTL/P100-128         	   85308	     14928 ns/op	    1152 B/op	      64 allocs/op
BenchmarkGCacheLRUSetGetSmallDataWithTTL/P1000-128        	   88101	     16564 ns/op	    1154 B/op	      64 allocs/op
BenchmarkGCacheLRUSetGetSmallDataWithTTL/P10000-128       	  126004	     12669 ns/op	    1167 B/op	      64 allocs/op
BenchmarkGCacheLRUSetGetBigDataNoTTL/P100-128             	       5	 228057323 ns/op	12377401 B/op	  576823 allocs/op
BenchmarkGCacheLRUSetGetBigDataNoTTL/P1000-128            	       5	 232424450 ns/op	12388700 B/op	  577122 allocs/op
BenchmarkGCacheLRUSetGetBigDataNoTTL/P10000-128           	       5	 233161058 ns/op	12534494 B/op	  580761 allocs/op
BenchmarkGCacheLRUSetGetBigDataWithTTL/P100-128           	       5	 238523025 ns/op	12376396 B/op	  576816 allocs/op
BenchmarkGCacheLRUSetGetBigDataWithTTL/P1000-128          	       5	 229796721 ns/op	12388780 B/op	  577123 allocs/op
BenchmarkGCacheLRUSetGetBigDataWithTTL/P10000-128         	       5	 220375684 ns/op	12535305 B/op	  580766 allocs/op
BenchmarkGCacheLFUSetGetSmallDataNoTTL/P100-128           	   80607	     14529 ns/op	    2077 B/op	      78 allocs/op
BenchmarkGCacheLFUSetGetSmallDataNoTTL/P1000-128          	  117544	     15680 ns/op	    2070 B/op	      78 allocs/op
BenchmarkGCacheLFUSetGetSmallDataNoTTL/P10000-128         	   94684	     14128 ns/op	    2095 B/op	      78 allocs/op
BenchmarkGCacheLFUSetGetSmallDataWithTTL/P100-128         	  111253	     13672 ns/op	    1232 B/op	      65 allocs/op
BenchmarkGCacheLFUSetGetSmallDataWithTTL/P1000-128        	   79286	     16552 ns/op	    1157 B/op	      64 allocs/op
BenchmarkGCacheLFUSetGetSmallDataWithTTL/P10000-128       	   92112	     13903 ns/op	    1175 B/op	      64 allocs/op
BenchmarkGCacheLFUSetGetBigDataNoTTL/P100-128             	       5	 219469875 ns/op	11118596 B/op	  550603 allocs/op
BenchmarkGCacheLFUSetGetBigDataNoTTL/P1000-128            	       6	 224812353 ns/op	10849384 B/op	  546487 allocs/op
BenchmarkGCacheLFUSetGetBigDataNoTTL/P10000-128           	       5	 225070403 ns/op	11275744 B/op	  554542 allocs/op
BenchmarkGCacheLFUSetGetBigDataWithTTL/P100-128           	       5	 225341043 ns/op	11118249 B/op	  550601 allocs/op
BenchmarkGCacheLFUSetGetBigDataWithTTL/P1000-128          	       5	 240683175 ns/op	11130520 B/op	  550910 allocs/op
BenchmarkGCacheLFUSetGetBigDataWithTTL/P10000-128         	       5	 219953071 ns/op	11275710 B/op	  554540 allocs/op
BenchmarkGCacheARCSetGetSmallDataNoTTL/P100-128           	   42540	     24902 ns/op	    3478 B/op	     107 allocs/op
BenchmarkGCacheARCSetGetSmallDataNoTTL/P1000-128          	   50721	     24631 ns/op	    3470 B/op	     107 allocs/op
BenchmarkGCacheARCSetGetSmallDataNoTTL/P10000-128         	   59578	     23067 ns/op	    3427 B/op	     106 allocs/op
BenchmarkGCacheARCSetGetSmallDataWithTTL/P100-128         	   59557	     17613 ns/op	    1153 B/op	      64 allocs/op
BenchmarkGCacheARCSetGetSmallDataWithTTL/P1000-128        	   74067	     20532 ns/op	    1156 B/op	      64 allocs/op
BenchmarkGCacheARCSetGetSmallDataWithTTL/P10000-128       	   69200	     16790 ns/op	    1179 B/op	      64 allocs/op
BenchmarkGCacheARCSetGetBigDataNoTTL/P100-128             	       2	 604006094 ns/op	33911000 B/op	  722193 allocs/op
BenchmarkGCacheARCSetGetBigDataNoTTL/P1000-128            	       3	 371503391 ns/op	16455026 B/op	  656047 allocs/op
BenchmarkGCacheARCSetGetBigDataNoTTL/P10000-128           	       3	 362508912 ns/op	16695981 B/op	  662092 allocs/op
BenchmarkGCacheARCSetGetBigDataWithTTL/P100-128           	       2	 599410824 ns/op	33910920 B/op	  722191 allocs/op
BenchmarkGCacheARCSetGetBigDataWithTTL/P1000-128          	       3	 342562004 ns/op	16453800 B/op	  656032 allocs/op
BenchmarkGCacheARCSetGetBigDataWithTTL/P10000-128         	       3	 358504649 ns/op	16695666 B/op	  662089 allocs/op
BenchmarkDefaultMapSetOnlySmallDataNoTTL/P100-128         	  218677	      5601 ns/op	     522 B/op	      32 allocs/op
BenchmarkDefaultMapSetOnlySmallDataNoTTL/P1000-128        	  270268	      5637 ns/op	     553 B/op	      32 allocs/op
BenchmarkDefaultMapSetOnlySmallDataNoTTL/P10000-128       	   89755	     12216 ns/op	    1662 B/op	      60 allocs/op
BenchmarkDefaultMapSetOnlyBigDataNoTTL/P100-128           	      10	 114630891 ns/op	 4298258 B/op	  264736 allocs/op
BenchmarkDefaultMapSetOnlyBigDataNoTTL/P1000-128          	      12	 108794966 ns/op	 5048175 B/op	  283499 allocs/op
BenchmarkDefaultMapSetOnlyBigDataNoTTL/P10000-128         	      10	 111709204 ns/op	14435452 B/op	  518172 allocs/op
BenchmarkSyncMapSetOnlySmallDataNoTTL/P100-128            	  238866	      5708 ns/op	    1290 B/op	      48 allocs/op
BenchmarkSyncMapSetOnlySmallDataNoTTL/P1000-128           	  190341	      6110 ns/op	    1338 B/op	      49 allocs/op
BenchmarkSyncMapSetOnlySmallDataNoTTL/P10000-128          	  133785	      8285 ns/op	    2052 B/op	      67 allocs/op
BenchmarkSyncMapSetOnlyBigDataNoTTL/P100-128              	     295	   7969627 ns/op	10489273 B/op	  393304 allocs/op
BenchmarkSyncMapSetOnlyBigDataNoTTL/P1000-128             	     283	   7594854 ns/op	10522007 B/op	  394122 allocs/op
BenchmarkSyncMapSetOnlyBigDataNoTTL/P10000-128            	     116	   8857558 ns/op	11368604 B/op	  415287 allocs/op
BenchmarkGacheV2SetOnlySmallDataNoTTL/P100-128            	  748718	      1467 ns/op	       4 B/op	       0 allocs/op
BenchmarkGacheV2SetOnlySmallDataNoTTL/P1000-128           	  809707	      1449 ns/op	      14 B/op	       0 allocs/op
BenchmarkGacheV2SetOnlySmallDataNoTTL/P10000-128          	  361364	      3123 ns/op	     287 B/op	       7 allocs/op
BenchmarkGacheV2SetOnlyBigDataNoTTL/P100-128              	     501	   5481720 ns/op	    2213 B/op	      53 allocs/op
BenchmarkGacheV2SetOnlyBigDataNoTTL/P1000-128             	     260	   5485101 ns/op	   39632 B/op	     987 allocs/op
BenchmarkGacheV2SetOnlyBigDataNoTTL/P10000-128            	     171	   7849288 ns/op	  599174 B/op	   14974 allocs/op
BenchmarkGacheSetOnlySmallDataNoTTL/P100-128              	  509191	      2237 ns/op	     643 B/op	      32 allocs/op
BenchmarkGacheSetOnlySmallDataNoTTL/P1000-128             	  535208	      2195 ns/op	     659 B/op	      32 allocs/op
BenchmarkGacheSetOnlySmallDataNoTTL/P10000-128            	  225684	      4711 ns/op	    1094 B/op	      43 allocs/op
BenchmarkGacheSetOnlyBigDataNoTTL/P100-128                	     213	   5927275 ns/op	 5247715 B/op	  262265 allocs/op
BenchmarkGacheSetOnlyBigDataNoTTL/P1000-128               	     194	   5848007 ns/op	 5295691 B/op	  263464 allocs/op
BenchmarkGacheSetOnlyBigDataNoTTL/P10000-128              	     136	   9526672 ns/op	 5995866 B/op	  280969 allocs/op
BenchmarkTTLCacheSetOnlySmallDataNoTTL/P100-128           	  435186	      3504 ns/op	       5 B/op	       0 allocs/op
BenchmarkTTLCacheSetOnlySmallDataNoTTL/P1000-128          	  483974	      3542 ns/op	      23 B/op	       0 allocs/op
BenchmarkTTLCacheSetOnlySmallDataNoTTL/P10000-128         	  383566	      3650 ns/op	     269 B/op	       6 allocs/op
BenchmarkTTLCacheSetOnlyBigDataNoTTL/P100-128             	      10	 114304801 ns/op	  102968 B/op	    2583 allocs/op
BenchmarkTTLCacheSetOnlyBigDataNoTTL/P1000-128            	       9	 114270845 ns/op	 1139754 B/op	   28480 allocs/op
BenchmarkTTLCacheSetOnlyBigDataNoTTL/P10000-128           	       9	 115359895 ns/op	11379386 B/op	  284477 allocs/op
BenchmarkGoCacheSetOnlySmallDataNoTTL/P100-128            	  195320	      5273 ns/op	     267 B/op	      16 allocs/op
BenchmarkGoCacheSetOnlySmallDataNoTTL/P1000-128           	  135201	      8800 ns/op	     335 B/op	      17 allocs/op
BenchmarkGoCacheSetOnlySmallDataNoTTL/P10000-128          	   58136	     19084 ns/op	    2044 B/op	      60 allocs/op
BenchmarkGoCacheSetOnlyBigDataNoTTL/P100-128              	      10	 106829564 ns/op	 2200071 B/op	  133653 allocs/op
BenchmarkGoCacheSetOnlyBigDataNoTTL/P1000-128             	      10	 101989007 ns/op	 3121649 B/op	  156692 allocs/op
BenchmarkGoCacheSetOnlyBigDataNoTTL/P10000-128            	      10	 106876751 ns/op	12337732 B/op	  387094 allocs/op
BenchmarkBigCacheSetOnlySmallDataNoTTL/P100-128           	  415227	      2883 ns/op	    8661 B/op	       0 allocs/op
BenchmarkBigCacheSetOnlySmallDataNoTTL/P1000-128          	  586422	      2735 ns/op	    8191 B/op	       0 allocs/op
BenchmarkBigCacheSetOnlySmallDataNoTTL/P10000-128         	  186198	      7049 ns/op	   52028 B/op	      13 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkBigCacheSetOnlyBigDataNoTTL/P100-128             	       4	 267472189 ns/op	2164387096 B/op	    6685 allocs/op
BenchmarkBigCacheSetOnlyBigDataNoTTL/P1000-128            	       4	 354652404 ns/op	3785397948 B/op	   64274 allocs/op
BenchmarkBigCacheSetOnlyBigDataNoTTL/P10000-128           	       5	 285097053 ns/op	3158425368 B/op	  512141 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkGCacheLRUSetOnlySmallDataNoTTL/P100-128          	  227121	      5805 ns/op	     906 B/op	      48 allocs/op
BenchmarkGCacheLRUSetOnlySmallDataNoTTL/P1000-128         	  233058	      6736 ns/op	     944 B/op	      49 allocs/op
BenchmarkGCacheLRUSetOnlySmallDataNoTTL/P10000-128        	  192927	      6954 ns/op	    1431 B/op	      61 allocs/op
BenchmarkGCacheLRUSetOnlyBigDataNoTTL/P100-128            	       9	 123638203 ns/op	 7455281 B/op	  396093 allocs/op
BenchmarkGCacheLRUSetOnlyBigDataNoTTL/P1000-128           	       9	 122065045 ns/op	 8479276 B/op	  421693 allocs/op
BenchmarkGCacheLRUSetOnlyBigDataNoTTL/P10000-128          	       9	 125589346 ns/op	18719096 B/op	  677692 allocs/op
BenchmarkGCacheLFUSetOnlySmallDataNoTTL/P100-128          	  175653	      7021 ns/op	     909 B/op	      48 allocs/op
BenchmarkGCacheLFUSetOnlySmallDataNoTTL/P1000-128         	  252823	      8278 ns/op	     941 B/op	      49 allocs/op
BenchmarkGCacheLFUSetOnlySmallDataNoTTL/P10000-128        	  108327	      9557 ns/op	    1848 B/op	      71 allocs/op
BenchmarkGCacheLFUSetOnlyBigDataNoTTL/P100-128            	       9	 124857252 ns/op	 7454565 B/op	  396087 allocs/op
BenchmarkGCacheLFUSetOnlyBigDataNoTTL/P1000-128           	       8	 126519025 ns/op	 8620665 B/op	  425242 allocs/op
BenchmarkGCacheLFUSetOnlyBigDataNoTTL/P10000-128          	       8	 128891442 ns/op	20141153 B/op	  713246 allocs/op
BenchmarkGCacheARCSetOnlySmallDataNoTTL/P100-128          	  130440	      7811 ns/op	     914 B/op	      48 allocs/op
BenchmarkGCacheARCSetOnlySmallDataNoTTL/P1000-128         	  206230	      8111 ns/op	     950 B/op	      49 allocs/op
BenchmarkGCacheARCSetOnlySmallDataNoTTL/P10000-128        	   12645	     79939 ns/op	    9091 B/op	     251 allocs/op
BenchmarkGCacheARCSetOnlyBigDataNoTTL/P100-128            	       6	 169627989 ns/op	 7512233 B/op	  397523 allocs/op
BenchmarkGCacheARCSetOnlyBigDataNoTTL/P1000-128           	       7	 156545207 ns/op	 8804595 B/op	  429827 allocs/op
BenchmarkGCacheARCSetOnlyBigDataNoTTL/P10000-128          	       6	 174078570 ns/op	24410124 B/op	  819941 allocs/op
BenchmarkDefaultMapGetSmallDataNoTTL/P100-128             	 1757305	       640.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkDefaultMapGetSmallDataNoTTL/P1000-128            	 1828549	       652.3 ns/op	       5 B/op	       0 allocs/op
BenchmarkDefaultMapGetSmallDataNoTTL/P10000-128           	  365136	      2952 ns/op	     280 B/op	       7 allocs/op
BenchmarkDefaultMapGetBigDataNoTTL/P100-128               	     193	   5956091 ns/op	    5311 B/op	     132 allocs/op
BenchmarkDefaultMapGetBigDataNoTTL/P1000-128              	     184	   8161142 ns/op	   55662 B/op	    1391 allocs/op
BenchmarkDefaultMapGetBigDataNoTTL/P10000-128             	      46	  21762217 ns/op	 2226193 B/op	   55656 allocs/op
BenchmarkSyncMapGetSmallDataNoTTL/P100-128                	197122238	         5.871 ns/op	       0 B/op	       0 allocs/op
BenchmarkSyncMapGetSmallDataNoTTL/P1000-128               	177865125	         6.183 ns/op	       0 B/op	       0 allocs/op
BenchmarkSyncMapGetSmallDataNoTTL/P10000-128              	52177215	        21.55 ns/op	       1 B/op	       0 allocs/op
BenchmarkSyncMapGetBigDataNoTTL/P100-128                  	     206	   6720944 ns/op	    4984 B/op	     124 allocs/op
BenchmarkSyncMapGetBigDataNoTTL/P1000-128                 	     196	   6506530 ns/op	   52254 B/op	    1306 allocs/op
BenchmarkSyncMapGetBigDataNoTTL/P10000-128                	     114	  11370485 ns/op	  898300 B/op	   22458 allocs/op
BenchmarkGacheV2GetSmallDataNoTTL/P100-128                	21000858	        57.18 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2GetSmallDataNoTTL/P1000-128               	20676634	        57.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheV2GetSmallDataNoTTL/P10000-128              	35721402	        29.37 ns/op	       2 B/op	       0 allocs/op
BenchmarkGacheV2GetBigDataNoTTL/P100-128                  	     434	   5695617 ns/op	    2363 B/op	      59 allocs/op
BenchmarkGacheV2GetBigDataNoTTL/P1000-128                 	     374	   6287994 ns/op	   27385 B/op	     684 allocs/op
BenchmarkGacheV2GetBigDataNoTTL/P10000-128                	     110	  10915280 ns/op	  930960 B/op	   23274 allocs/op
BenchmarkGacheGetSmallDataNoTTL/P100-128                  	169520691	         7.096 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheGetSmallDataNoTTL/P1000-128                 	144362780	         7.300 ns/op	       0 B/op	       0 allocs/op
BenchmarkGacheGetSmallDataNoTTL/P10000-128                	19560367	        51.18 ns/op	       5 B/op	       0 allocs/op
BenchmarkGacheGetBigDataNoTTL/P100-128                    	     391	   5192528 ns/op	    2622 B/op	      65 allocs/op
BenchmarkGacheGetBigDataNoTTL/P1000-128                   	     228	   5670977 ns/op	   44924 B/op	    1123 allocs/op
BenchmarkGacheGetBigDataNoTTL/P10000-128                  	     151	   9416057 ns/op	  678176 B/op	   16954 allocs/op
BenchmarkTTLCacheGetSmallDataNoTTL/P100-128               	  419108	      7723 ns/op	       6 B/op	       0 allocs/op
BenchmarkTTLCacheGetSmallDataNoTTL/P1000-128              	  352797	      3598 ns/op	      32 B/op	       0 allocs/op
BenchmarkTTLCacheGetSmallDataNoTTL/P10000-128             	   71990	     15975 ns/op	    1443 B/op	      35 allocs/op
BenchmarkTTLCacheGetBigDataNoTTL/P100-128                 	      12	 101410740 ns/op	   85828 B/op	    2154 allocs/op
BenchmarkTTLCacheGetBigDataNoTTL/P1000-128                	      12	  95905221 ns/op	  853800 B/op	   21352 allocs/op
BenchmarkTTLCacheGetBigDataNoTTL/P10000-128               	      10	 100553472 ns/op	10240714 B/op	  256022 allocs/op
BenchmarkGoCacheGetSmallDataNoTTL/P100-128                	 1907704	       617.0 ns/op	       0 B/op	       0 allocs/op
BenchmarkGoCacheGetSmallDataNoTTL/P1000-128               	 1628834	       623.4 ns/op	       6 B/op	       0 allocs/op
BenchmarkGoCacheGetSmallDataNoTTL/P10000-128              	  396585	      2732 ns/op	     258 B/op	       6 allocs/op
BenchmarkGoCacheGetBigDataNoTTL/P100-128                  	     170	   6227226 ns/op	    6029 B/op	     150 allocs/op
BenchmarkGoCacheGetBigDataNoTTL/P1000-128                 	     157	   6972640 ns/op	   65229 B/op	    1630 allocs/op
BenchmarkGoCacheGetBigDataNoTTL/P10000-128                	      62	  16375446 ns/op	 1651693 B/op	   41293 allocs/op
BenchmarkBigCacheGetSmallDataNoTTL/P100-128               	 3247472	       459.6 ns/op	    2048 B/op	      32 allocs/op
BenchmarkBigCacheGetSmallDataNoTTL/P1000-128              	 2481582	       423.5 ns/op	    2052 B/op	      32 allocs/op
BenchmarkBigCacheGetSmallDataNoTTL/P10000-128             	 1081971	       990.5 ns/op	    2142 B/op	      34 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkBigCacheGetBigDataNoTTL/P100-128                 	      10	 100412425 ns/op	536973802 B/op	  264724 allocs/op
BenchmarkBigCacheGetBigDataNoTTL/P1000-128                	      14	  79895912 ns/op	537602778 B/op	  280446 allocs/op
BenchmarkBigCacheGetBigDataNoTTL/P10000-128               	       7	 144524691 ns/op	551500187 B/op	  627887 allocs/op
ctx done, shutting down bigcache cleanup routine
BenchmarkGCacheLRUGetSmallDataNoTTL/P100-128              	  243140	      4367 ns/op	     265 B/op	      16 allocs/op
BenchmarkGCacheLRUGetSmallDataNoTTL/P1000-128             	  266984	      4564 ns/op	     296 B/op	      16 allocs/op
BenchmarkGCacheLRUGetSmallDataNoTTL/P10000-128            	  115518	      9800 ns/op	    1146 B/op	      38 allocs/op
BenchmarkGCacheLRUGetBigDataNoTTL/P100-128                	      43	  37402095 ns/op	 2121147 B/op	  131674 allocs/op
BenchmarkGCacheLRUGetBigDataNoTTL/P1000-128               	      33	  37313418 ns/op	 2407637 B/op	  138837 allocs/op
BenchmarkGCacheLRUGetBigDataNoTTL/P10000-128              	      26	  40407780 ns/op	 6035810 B/op	  229541 allocs/op
BenchmarkGCacheLFUGetSmallDataNoTTL/P100-128              	  261147	      4517 ns/op	     261 B/op	      16 allocs/op
BenchmarkGCacheLFUGetSmallDataNoTTL/P1000-128             	  232633	      4630 ns/op	     302 B/op	      17 allocs/op
BenchmarkGCacheLFUGetSmallDataNoTTL/P10000-128            	   82003	     12745 ns/op	    1510 B/op	      47 allocs/op
BenchmarkGCacheLFUGetBigDataNoTTL/P100-128                	      45	  36822274 ns/op	 2120086 B/op	  131648 allocs/op
BenchmarkGCacheLFUGetBigDataNoTTL/P1000-128               	      38	  36942957 ns/op	 2366804 B/op	  137816 allocs/op
BenchmarkGCacheLFUGetBigDataNoTTL/P10000-128              	      25	  40389236 ns/op	 6193349 B/op	  233480 allocs/op
BenchmarkGCacheARCGetSmallDataNoTTL/P100-128              	  184375	      5686 ns/op	     264 B/op	      16 allocs/op
BenchmarkGCacheARCGetSmallDataNoTTL/P1000-128             	  180766	      5600 ns/op	     314 B/op	      17 allocs/op
BenchmarkGCacheARCGetSmallDataNoTTL/P10000-128            	   82004	     12994 ns/op	    1510 B/op	      47 allocs/op
BenchmarkGCacheARCGetBigDataNoTTL/P100-128                	      55	  41313466 ns/op	 2115969 B/op	  131545 allocs/op
BenchmarkGCacheARCGetBigDataNoTTL/P1000-128               	      28	  40821381 ns/op	 2463065 B/op	  140223 allocs/op
BenchmarkGCacheARCGetBigDataNoTTL/P10000-128              	      24	  46870338 ns/op	 6364047 B/op	  237748 allocs/op
PASS
ok  	github.com/kpango/go-cache-lib-benchmarks	1449.621s
```

<!-- benchmark-chart-section-start -->

## Benchmark Charts

### BigData NoTTL

![BigData NoTTL](images/BigData_NoTTL_3d_chart.svg)

[📊 View Interactive 3D Chart](https://htmlpreview.github.io/?https://github.com/kpango/go-cache-lib-benchmarks/blob/main/images/BigData_NoTTL_3d_chart.html)

### BigData WithTTL

![BigData WithTTL](images/BigData_WithTTL_3d_chart.svg)

[📊 View Interactive 3D Chart](https://htmlpreview.github.io/?https://github.com/kpango/go-cache-lib-benchmarks/blob/main/images/BigData_WithTTL_3d_chart.html)

### SmallData NoTTL

![SmallData NoTTL](images/SmallData_NoTTL_3d_chart.svg)

[📊 View Interactive 3D Chart](https://htmlpreview.github.io/?https://github.com/kpango/go-cache-lib-benchmarks/blob/main/images/SmallData_NoTTL_3d_chart.html)

### SmallData WithTTL

![SmallData WithTTL](images/SmallData_WithTTL_3d_chart.svg)

[📊 View Interactive 3D Chart](https://htmlpreview.github.io/?https://github.com/kpango/go-cache-lib-benchmarks/blob/main/images/SmallData_WithTTL_3d_chart.html)

<!-- benchmark-chart-section-end -->
