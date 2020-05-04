#!/bin/zsh -f
for sync in cg mg gvstm; do
    for threads in 1 2 4 8; do
	for RO in 100 90 60; do
	    for i in 1 2 3 4 5; do
		stmbench7 -outDir="bench_results/${sync}_${threads}_${RO}" -sync=$sync -threads=$threads -duration=1s -roRatio=$RO;
	    done
	done
    done
done
