# SuperQueue <!-- omit in toc -->

Super Simple, Super Scalable, Super Speedy, Super Queue.

## Table of Contents <!-- omit in toc -->
- [Motivation](#motivation)
- [Inspiration](#inspiration)
- [Using as a Primitive](#using-as-a-primitive)
- [Why MapMaps](#why-mapmaps)
- [Get Started](#get-started)
  - [Running ScyllaDB](#running-scylladb)
  - [Running SuperQueue (single partition)](#running-superqueue-single-partition)
- [API Docs](#api-docs)
  - [POST /record - Create a new record](#post-record---create-a-new-record)
  - [GET /record - Get a record](#get-record---get-a-record)
  - [POST /ack/:recordID](#post-ackrecordid)
  - [POST /nack/:recordID](#post-nackrecordid)
  - [Get /metrics](#get-metrics)

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

A WIP Request Router is available here: https://github.com/clusterspace-io/SuperQueueRequestRouter

## Inspiration

In order to build such a system, we need to address some elephants in the room.

First, cloud provider specific solutions such as SQS and GCP Pub/Sub can achieve many of these scalability features... at an outrageous cost at scale.

In order to make these queues, I spent a long time thinking about how to design the architecture in such a way that is not limited by throughput, cardinality, etc. - I wanted to make something that scales like SQS.

SQS happens to have no ordering guarantees. (I assume) By this they mean they partition your queue (just like DynamoDB), and items within a partition should be within the same order. Based on how full each partition is and how they get load balanced, things in other partitions may be out of order.

While searching for how to design the running application itself, I came across [Segment's Centrifuge](https://segment.com/blog/introducing-centrifuge/), thank you to them for discussing their architecture because it confirmed what I wanted to do: DB backed queues with in-memory processing. This also helped shape the design for the non-blocking aspect of the system.

However I've implemented some stark changes from their architecture:

I've designed these job queues to be pull-based, meaning any worker can integrate with it (this dramatically reduces the amount of code required). This required designing more efficient data structures for parsing timestamp based data (see MapMaps for more).

In their design they use MySQL, I chose ScyllaDB (I've also done tests with CRDB). I think using something that is write-optimized is obvious here, as for all usage by disaster recovery we are write only. Furthermore being able to scale out linearly is awesome. Consistency is not needed, because the concept of a single processor per DB (namespace/keyspace in the context of Scylla) guarantees that only we will be writing to it and every write has a unique primary key (no conflict issues). Furthermore during disaster recovery we can read with `ALL` consistency if needed, but by the time recovery would start records should have fully propagated, but we can run with `ALL` anyway just to be sure. A `TWO` write consistency gives us sufficient durability.

I tried using CockroachDB with both range and hash partitioning, but running on my laptop any real load would induce 100ms> inserts, Scylla could go up to 10ms. Both should be within single digit ms (Scylla even going under 1ms) on real DB clusters on real hardware, however ScyllaDB should scale better with the same hardware so that is the choice for now, plus consistency is not needed. Changing the DB is very easy since only a few queries need to be changed (2 writes, and 1 read). In theory something in-memory like Redis could work well too, but you'll have to concern yourself with data set size (adding TTLs and archiving data could be reasonable since there is a max lifetime to records).

## Using as a Primitive

SuperQueue can be used as is, but doesn't offer much protection. It exposes metrics, but does not enforce high availability by itself. It is designed to be a primitive that can be wrapped by other service to enable a wide array of uses.

For example, managed SuperQueue will use the metrics, custom request router, and service discovery to offer serverless SuperQueues. It will scale up based on metrics, and route randomly to partitions. It will create more partitions as needed, as well as drain existing ones into others during scale down. The request router will handle auth and rate limiting as well.

Another example is handling when certain partitions are empty or full. If a partition is empty, it can ask around other partitions before returning the client request. Same with whether partitions are full (although we should avoid getting there at all costs).

SuperQueues being virtual also allows it to scale extremely quickly. If the resources are ready, a new partition can come up in less than 200ms (including registration with service discovery, 200ms from nothing to accepting requests).

SuperQueues could also be used in a self hosted manner to allow flexibility how load balancing and scaling is done. It also allows establishing what ever limits are desired, and what happens during extreme back pressure.

It could even be used internally as a golang package.

## Why MapMaps

In order to maintain high performance at scale, a data structure was needed that could efficiently scan over timestamp-based data that was not append only, while also allowing for O(1) access (for ack/nack of in-flight timeouts). In order to accommodate this, I created the MapMap (I assume I created this, I haven't seen it used anywhere else like this).

While the idea of nested maps is not novel, I believe my implementation is. The reason it is called a MapMap is because it consists of nested maps. In Go, it looks like this:

```go
type MapMap map[int64]map[string]interface{}
```

The outside map serves as a configurable time bucket system. The keys are bucket epoch timestamps in milliseconds. This allows for rough ordering (not maintained within the bucket) of bucketed data. So for example, you might want to bucket your data every 5ms, meaning that in a single second you would have buckets of `...0`, `...5`, `...10`, and so on.

Within these buckets exist maps in which the keys are unique document identifiers.

In Go, we can iterate over a map in O(N) list like a list, so the basic concept is this:

1. We have a MapMapConsumer that on some interval (ideally matching the bucket interval) consume all bucket from the last iterate time up to now. This will consume everything up to the current timestamp no matter any delay or latency in processing.
2. For eac bucket it runs some `ConsumerFunc`, which in this case will queue up the items by iterating over the map in O(N) time.
3. Consumer will delete the bucket

In-flight items also get their timeouts placed in this MapMap. When an item is acked or nacked, we can remove it from the MapMap in O(1) by calculating it's bucket, then removing the item from the map in that bucket. We also delete it from an in-flight map we also keep track of.

So this means we can iterate over data at a configurable interval, data is ordered by interval (not ordered within), can be placed arbitrarily in the future (to really any time in the future), and all parsed in O(N) where N is the number of items that exist up to now in the map. By doing up to current time we also account for any pauses or increased latency in processing (we never miss anything).

Oh yeah and we do this way faster than any b-tree or LSM tree could.

Besides the downside of rough ordering (there is no reason we need exact ordering, besides its configurable by changing the bucket interval), we get O(N) where we want O(N), and O(1) where we want O(1), a pretty great tradeoff. The term MapMap also conveniently lines up with the similar HyperLogLog name, which also does rough calculations rather than exact.

## Get Started

To run as standalone, there are a few environment variables that need to be setup, as well as running ScyllaDB somewhere.

### Running ScyllaDB

The easiest way to do this is with docker:

```
docker run -p 9042:9042 --name some-scylla -d scylladb/scylla
```

_Pro tip: Change `9042:9042` to `localhost:9042:9042` if you only want Scylla exposed on your localhost interface_

### Running SuperQueue (single partition)

As it currently stands, the whole binary supports a single partition of a single queue. This is due to how the `main.go` file is setup.

```
HTTP_PORT=8080 PARTITION=part1 go run .
```

_Pro tip: At high load, writing to stdout becomes a bottleneck, so do a `> out.txt` if you are scale testing!_

## API Docs

There are only a few endpoints, which make the system super simple and robust. For examples I will be using the `httpie` cli.

### POST /record - Create a new record

This will create a new record for processing.

Headers:
- `content-type: application/json`

Body:
```
{
  payload: string, // The string payload, typically stringified JSON
  delay_ms?: int // Delay queueing, this should be a reasonable time in the future (at least 100ms) since there is no validation here currently. A value in the past will get (nearly) immediately queued.
}
```

Expected Response code: `204`

Example:

```
http post http://localhost:8080/record payload=hey!
```

### GET /record - Get a record

This will fetch the next record that is available in the queue (partition). Currently there is a hard-coded 30s in-flight timeout, meaning that after 30 seconds if you do not ack or nack the record it will requeue (and unable to ack or nack).

Expected Response body (code `200`):

```
{
  id: string,
  payload: string,
  attempts: int
}
```

You will get a `204` response if the queue is empty.

Example:
```
http get http://localhost:8080/record
```

### POST /ack/:recordID

Acknowledge a record to prevent further processing.

Expected Response code: `200`

Example:
```
http post http://localhost:8080/ack/partition1_21yFbkxyFx6AjihUA2CN0WkrfJD
```

### POST /nack/:recordID

Negatively acknowledge a record to immediately requeue the record. An optional delay can be added to override any immediate re-enqueue or back off.

Body:
```
{
  delay_ms?: int // Manual override of exponential back off
}
```

Expected Response code: `200`


Example:
```
http post http://localhost:8080/nack/partition1_21yFbkxyFx6AjihUA2CN0WkrfJD
```

### Get /metrics

Get metrics about the partition in prometheus format

Expected Response code: `200`
