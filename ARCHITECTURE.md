# MapMaps

MapMaps are really useful for iterating over roughly ordered and time-based items.

## MapMap performance

Without locking, loading can be done very fast:
```
Testing map bucket iterate
Loading the map bucket with 10000000 entries
Loaded in 2.126984778s
Iterating 1000000 values
Iterated in 38.335443ms
Deleting chunks 1000000 values
Deleted chunks in 27.604647ms
Deleting individual 1000000 values (need to reload tree)
Deleted individuals in 50.044488ms
```

However to allow it to be thread safe, we need to add locking on sub-map creation. This slows down adding a bit:
```
Adding 10000000 items...
MapMap filled 10000000 items in 3.173124865s
```

### Why not Btrees

Slow:

```
Testing tree iterate
Loading the tree with 10000000 entries
Loaded in 8.469774443s
Ascending 1000000 values
Ascended in 58.20793ms
Deleting 1000000 values
Deleted in 509.624107ms
```

This speed difference gets dramatically obvious at scale:

```
Testing tree iterate
Loading the tree with 100000000 entries
Loaded in 1m38.781428795s
Ascending 10000000 values
Ascended in 634.068737ms
Deleting 10000000 values
Deleted in 5.750347824s



Testing map bucket iterate
Loading the map bucket with 100000000 entries
Loaded in 21.148275508s
Iterating 10000000 values
Iterated in 435.18215ms
Deleting chunks 10000000 values
Deleted chunks in 297.872649ms
Deleting individual 10000000 values (need to reload tree)
Deleted individuals in 570.332297ms
```

That's the difference between O(log(n)) and O(n) operations.
