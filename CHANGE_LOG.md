# Change Log

## [v1.5](https://github.com/thewizardplusplus/go-cache/tree/v1.5) (2020-11-07)

## [v1.4](https://github.com/thewizardplusplus/go-cache/tree/v1.4) (2020-10-20)

## [v1.3](https://github.com/thewizardplusplus/go-cache/tree/v1.3) (2019-08-15)

- implementation of garbage collection:
  - improve support of interruption via a context:
    - pass a context to the `gc.GC.Clean()` method;
    - additional interruption via a context:
      - in the `gc.TotalGC.Clean()` method;
      - in the `gc.PartialGC.Clean()` method;
- improve benchmarks:
  - add to benchmarks:
    - different storage sizes;
    - different expired percents;
  - stop at the end of each benchmark:
    - garbage collecting;
    - additional concurrent loading;
  - slow down additional concurrent loading;
- refactoring:
  - extract from the `gc.TotalGC.Clean()` method:
    - the `gc.TotalGC.handleIteration()` method;
  - extract from the `gc.PartialGC.Clean()` method:
    - the `gc.counter` structure;
    - the `gc.iterator` structure.

## [v1.2](https://github.com/thewizardplusplus/go-cache/tree/v1.2) (2019-07-09)

- implementation of garbage collection:
  - independent implementation of garbage collection running;
  - implementation of partial garbage collection (based on [expiration in Redis](https://redis.io/commands/expire#how-redis-expires-keys)).

### Features

- implementation of an in-memory cache:
  - operations:
    - getting a value by a key:
      - signaling a reason for the absence of a key - missed or expired;
    - getting a value by a key with deletion of expired values:
      - signaling a reason for the absence of a key - missed or expired;
    - setting a key-value pair with a specified time to live:
      - support of key-value pairs without a set time to live (persistent);
    - deletion;
- implementation of garbage collection:
  - independent implementation of garbage collection running:
    - support interruption via a context;
    - support specification of a running period;
  - implementation of total garbage collection (based on a full scan);
  - implementation of partial garbage collection (based on [expiration in Redis](https://redis.io/commands/expire#how-redis-expires-keys)).

## [v1.1](https://github.com/thewizardplusplus/go-cache/tree/v1.1) (2019-07-08)

- implementation of an in-memory cache:
  - make public the expired value model;
- implementation of garbage collection:
  - implementation of total garbage collection (based on a full scan):
    - support interruption via a context;
    - support specification of a running period.

### Features

- implementation of an in-memory cache:
  - operations:
    - getting a value by a key:
      - signaling a reason for the absence of a key - missed or expired;
    - getting a value by a key with deletion of expired values:
      - signaling a reason for the absence of a key - missed or expired;
    - setting a key-value pair with a specified time to live:
      - support of key-value pairs without a set time to live (persistent);
    - deletion;
- implementation of garbage collection:
  - implementation of total garbage collection (based on a full scan):
    - support interruption via a context;
    - support specification of a running period.

## [v1.0](https://github.com/thewizardplusplus/go-cache/tree/v1.0) (2019-07-05)

### Features

- implementation of an in-memory cache:
  - operations:
    - getting a value by a key:
      - signaling a reason for the absence of a key - missed or expired;
    - getting a value by a key with deletion of expired values:
      - signaling a reason for the absence of a key - missed or expired;
    - setting a key-value pair with a specified time to live:
      - support of key-value pairs without a set time to live (persistent);
    - deletion.
