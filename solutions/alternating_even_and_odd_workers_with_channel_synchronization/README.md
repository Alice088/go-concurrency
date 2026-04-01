## Alternating Even/Odd Workers with Channel Synchronization

This task demonstrates coordination between multiple goroutines using channels, ensuring ordered output while enforcing a time limit on execution.

### Problem Description

You need to implement a concurrent system that writes numbers from 1 to 10 into a shared channel. The numbers must be produced in strict order:

1, 2, 3, 4, ..., 10

Two worker goroutines are responsible for generating the numbers:

* one worker handles even numbers
* another worker handles odd numbers

The key constraint is that workers must alternate execution. At any given moment, only one worker is allowed to write to the shared channel, and control must switch after each successful write.

### Synchronization Mechanism

A dedicated channel (`switchCh`) is used as a control signal:

* `true` indicates that the even worker is allowed to proceed
* `false` indicates that the odd worker is allowed to proceed

Each worker:

* waits for a signal from `switchCh`
* checks if it is allowed to act
* if yes, writes the next number and flips the signal
* if not, returns the signal unchanged

This creates a cooperative scheduling mechanism between goroutines without using mutexes.

### Execution Flow

The main routine:

* initializes channels:

  * `nums` for output
  * `switchCh` for coordination
* starts a context with a timeout of 6 seconds
* launches:

  * a goroutine to close `nums` when the timeout expires
  * a goroutine to seed the initial state of `switchCh`
  * two worker goroutines (`wEven`, `wOdd`)
* collects results from `nums` into a slice
* prints the final slice and execution time

### Workers Behavior

* `wEven` writes values where index is even (producing 2, 4, 6, ...)
* `wOdd` writes values where index is odd (producing 1, 3, 5, ...)

Each worker loops over a fixed range and participates in coordination via `switchCh`.

### Termination

The system stops when:

* either all numbers are processed
* or the context timeout is reached

Additionally, one of the workers closes the `nums` channel, while a separate goroutine may also attempt closure based on timeout, which introduces a potential race condition and is part of the exercise.

### Key Concepts Practiced

* goroutine coordination without mutexes
* channel-based signaling
* controlled alternation between workers
* context-based cancellation
* safe and unsafe channel closing patterns
* potential race conditions in concurrent systems

### Notes

This implementation intentionally exposes edge cases:

* multiple goroutines attempting to close the same channel
* reliance on cooperative scheduling via shared channel state
* absence of explicit worker pool control despite configuration hints

These aspects make it suitable for analysis, debugging, and improvement exercises in concurrent Go systems.

### Result
````
wEven: write: 1 
wOdd: write: 2
wEven: write: 3
wOdd: write: 4
wEven: write: 5
wOdd: write: 6
wEven: write: 7
wOdd: write: 8
wEven: write: 9
wOdd: write: 10
[1 2 3 4 5 6 7 8 9 10]
164.175µs
````