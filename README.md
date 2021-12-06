# SuperQueue <!-- omit in toc -->

Super Simple, Super Scalable, Super Speedy, Super Queue.

## Table of Contents <!-- omit in toc -->
- [Motivation](#motivation)
- [Inspiration](#inspiration)

## Motivation

In microservice and monolithic architectures, queues have been a go-to method for decoupling services and increasing durability of individual work items. However for a long time, queues have never reached their full potential.

In most instances, we do not need many of the features provided by other queues, such as FIFO and strict ordering. For queues that build away from these kinds of features, they introduce new issues to achieve their low latency and throughput, such as in-memory only, or at-most-once processing guarantees.

All of these caveats mean that there are currently very few solutions existing to provide highly flexible and scalable queue systems.

The goal of SuperQueues was to appeal to as many use cases as possible by building a simple architecture. To allow for:

- Extreme scale (true horizontal linear scaling)
- Extremely low latency (both end to end, and per request)
- Super simple to use (there are only 4 endpoints to interact with)
- Super useful features (custom exponential backoff, delayed sending, all per-message)
- Non-blocking (messages don't prevent each other from being sent)
- At-least-once processing (messages will requeue until they are acknowledged or exceed their lifetime)
- Durability (when we say we've got your message, it is written durably)

## Inspiration

In order to build such a system, we need to address some elephants in the room.

First, cloud provider specific solutions such as SQS and GCP Pub/Sub can achieve many of these scalability features... at an outrageous cost at scale.

In order to make these queues, I spent a long time thinking about how to design the architecture in such a way that is not limited by throughput, cardinality, etc. - I wanted to make something that scales like SQS.

SQS happens to have no ordering guarantees. By this they basically mean they partition your queue (just like DynamoDB), and items within a partition should be within the same order. Based on how full each partition is and how they get load balanced, things in other partitions may be out of order.

While searching for a good way to design within the application itself, I cam across [Segment's](https://segment.com/blog/introducing-centrifuge/), big credit to them for discussing their architecture because it basically confirmed what I wanted to do: DB backed queues with in-memory processing. This also helped shape the design for the non-blocking aspect of the system.

However I've implemented some stark changes from their architecture:

I've designed these job queues to be pull-based, meaning any worker can integrate with it (this dramatically reduces the amount of code required). This required designing more efficient data structures for parsing timestamp based data (see MapMaps for more).

In their design they use MySQL, I chose ScyllaDB (I've also done tests with CRDB). I think using something that is write-optimized is obvious here, as for all usage by disaster recovery we are write only. Furthermore being able to scale out linearly is awesome. Consistency is not needed, because the concept of a single processor per DB (namespace/keyspace in the context of Scylla) guarantees that only we will be writing to it and every write has a unique primary key (no conflict issues). Furthermore during disaster recovery we can read with `ALL` consistency if needed, but by the time recovery would start records should have fully propagated, but we can run with `ALL` anyway just to be sure. `QUORUM` (with RF=3) or `TWO` write consistency gives us sufficient durability (we use `TWO` in practice).

I tried using CockroachDB with both range and hash partitioning, but running on my laptop any real load would induce 100ms> inserts, Scylla could go up to 10ms. Both should be within single digit ms (Scylla even going under 1ms) on real DB clusters on real hardware, however ScyllaDB should scale better with the same hardware so that is the choice for now, plus consistency is not needed. Changing the DB is very easy since only a few queries need to be changed (2 writes, and 1 read). In theory something in-memory like Redis could work well too, but you'll have to concern yourself with data set size (adding TTLs and archiving data could be reasonable since there is a max lifetime to records).
