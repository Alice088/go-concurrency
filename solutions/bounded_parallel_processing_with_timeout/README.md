---
# Bounded Parallel Processing with Timeout

## Overview

This example implements a concurrent data processing pipeline with **bounded parallelism** and a **strict execution timeout**.

An input stream of jobs (`in`) is processed by a fixed number of workers. Each job is handled by a potentially slow function (`processData`). Work is distributed across workers, and results are aggregated into a single output channel.

**Constraint:** the system must run no longer than **5 seconds**.
Any jobs not completed within this window are discarded.

---

## Problem It Solves

* **Bounded parallelism**
  Limits the number of concurrent workers to prevent resource exhaustion.

* **Time-bounded execution**
  Guarantees completion within a fixed SLA using `context.WithTimeout`.

* **Partial results**
  Trades completeness for latency — unfinished jobs are ignored.

* **Isolation of slow operations**
  Prevents slow or blocking tasks from stalling the entire pipeline.

---

## Architecture

### Producer

* Sends jobs into the `in` channel
* Closes the channel when done

### Worker Pool (Fan-out)

* Spawns a fixed number of workers
* All workers consume from the same input channel
* Jobs are distributed dynamically

### Worker

* Reads from `in`
* Processes data via `processData`
* Sends results to `out`
* Stops on:

  * input channel close
  * context cancellation

### Aggregation

* All workers write to a shared `out` channel
* `out` is closed after all workers finish (`WaitGroup`)

### Time Control

```
ctx, cancel := context.WithTimeout(..., 5*time.Second)
```

* Enforces global timeout
* Propagates cancellation to all workers

---

## Implementation Notes

### Cancellation by Abandonment

* `processData` runs in its own goroutine
* It does not support cancellation
* Results are ignored if timeout occurs
* Underlying work may still continue

This is typical when dealing with:

* legacy code
* third-party APIs without context support

---

### Bounded Concurrency

```
numWorkers = 5
```

At most 5 jobs are processed concurrently.

---

### Safe Shutdown

* `WaitGroup` ensures all workers complete
* `out` is closed exactly once
* Consumers can safely `range` over `out`

---

## Behavior

* Jobs are processed concurrently and out of order
* Some jobs may not complete due to timeout
* Total runtime is bounded (~5 seconds)

---

## Real-World Use Cases

Used in systems where **latency is more important than completeness**:

* API aggregation (search, recommendations)
* External service calls / integrations
* Queue consumers and batch processing
* Streaming and real-time systems

---

## Core Idea

Process a **limited number of tasks**, within a **strict time budget**, using **controlled concurrency**.

---

## Summary

This pattern combines:

* Fan-out (task distribution)
* Worker pool (bounded concurrency)
* Result aggregation
* Context-based timeout

Suitable for systems that must remain responsive under unpredictable workloads.
