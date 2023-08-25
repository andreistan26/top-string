# Top-String

## Objective

How to implement an efficient system for counting the most frequent strings out of many records.
We will implemment and then benchmark the following solutions:

- local data and local counter
- remote data and local counter
- remote data and distributed counter

We will test for different data volumes

- a lot of small strings (10 - 100 milion 256 char strings)
- not as many big strings (10k 10mb strings)

## The Base Model and Easy Benchmarking

Base model of string counting:

- `Sender`   : GetHashStream(options)        => hashes  : chan FileHash{md5, filename}
- `Receiver` : CountStrings(hashes, options) => results : pqueue

### Receiver

- StartServer
- Accept

### Sender

## Local data and local counter

For this we need a basic hashmap. 


## Usage

For local usage:
```console
top-string local [PATH] [OPTIONS]
```

For remote server usage:
```console
top-string server [OPTIONS]
```

For remote sender usage:
```console
top-string send [IP:PORT] [PATH] [OPTIONS]
```

## Remote

For the remote option we will have multiple ways of doing this, we can 
