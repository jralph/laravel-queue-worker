# Laravel Queue Worker

A simple queue worker spawner for written in `Go` for `Laravel`.

Basically, this command runs `php artisan queue:work` a given number of times asynchronously.

I build this with the main purpose being to run multiple queue workers asyncronously in a development environment, without having to setup sometuing like supervisor.

## Installation

```
go get github.com/jralph/laravel-queue-worker
```

## Usage

```
laravel-queue-worker --artisan="path/to/artisan" --processes=numberOfProcessesToRun
```

```
Usage: laravel-queue-worker [-a value] [-d value] [-m value] [-p value] [-q value] [-r value] [-s value] [-t value] [parameters ...]
 -a, --artisan=value    The path to artisan executable. [Default: "artisan"]
 -d, --delay=value      Amount of time to delay failed jobs. [Default: 0]
 -m, --memory=value     The memory limit in megabytes. [Default: 128]
 -p, --processes=value  The number of works to run. [Default: 5]
 -q, --queue=value      The queue to listen on. [Default: "default"]
 -r, --tries=value      The number of times to attempt a job. [Default: 0]
 -s, --sleep=value      Number of seconds to sleep when no jobs are available. [Default: 3]
 -t, --timeout=value    The number of seconds a child process can run for. [Default: 60]
```

### Example

```
// Spin up 20 laravel queue workers running in the current directory.
laravel-queue-worker --processes=20
```
