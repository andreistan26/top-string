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

So this app should be resumed to interfaces as the main abstraction.
Essentially, we will need 2 interfaces, a `sender` and a `reciver`.

### Receiver

- StartServer
- Accept

### Sender

## Local data and local counter

For this we need a basic hashmap. 


## Usage

For local usage:
```console
top-string start local [PATH] [OPTIONS]
```

For remote server usage:
```console
top-string start server [OPTIONS]
```

For remote sender usage:
```console
top-string start send [IP:PORT] [PATH] [OPTIONS]
```


