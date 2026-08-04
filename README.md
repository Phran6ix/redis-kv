# redis-kv

A lightweight, concurrent, in-memory key-value store built from scratch in Go. 

This project implements a custom TCP server and a parser for the **RESP** (REdis Serialization Protocol). It is designed to act as a drop-in replacement for a subset of Redis commands, allowing standard tools like `redis-cli` to interact with it natively.

## 🎯 Purpose

I built this project to deepen my understanding of core cloud-native systems engineering concepts:
*   **Network Programming:** Managing raw TCP sockets and byte streams using Go's `net` package.
*   **Protocol Parsing:** Reading RFCs and writing a custom byte-level parser for the RESP binary protocol.
*   **Concurrency & Thread Safety:** Managing shared state across thousands of concurrent connections using goroutines and `sync.RWMutex`.

## ✨ Features

*   **Custom TCP Server:** Listens on port `6379` and handles multiple clients concurrently.
*   **RESP Parser:** Decodes incoming RESP arrays and bulk strings into Go structs.
*   **Thread-Safe In-Memory Store:** Uses a guarded Go map to process `GET` and `SET` commands without race conditions.
*   **Redis-CLI Compatible:** You can connect to it using standard Redis tools.

## 🚀 Getting Started

### Prerequisites
* Go 1.22+
* `redis-cli` (optional, for testing)

### Installation & Running

1. Clone the repository:
   ```bash
   git clone https://github.com/Phran6ix/redis-kv.git
   cd redis-kv
   ```

2. Start the server:
   ```bash
   go run cmd/server/main.go
   ```
   *The server will start listening for TCP connections on `127.0.0.1:6379`.*

### Testing with `redis-cli`

Open a new terminal window and connect to the server using the standard Redis CLI tool:

```bash
$ redis-cli
127.0.0.1:6379> SET name phran6ix
OK
127.0.0.1:6379> GET name
"phran6ix"
```

## 🏗️ Architecture

- `cmd/server/`: The entry point for the TCP server. Handles accepting connections and spawning goroutines.
- `internal/resp/`: Contains the custom RESP parser that decodes byte streams (`*`, `$`, `+`, `-`, `:`) into Go interfaces and structs.
- *(More to come as the project grows)*

## 🔮 Roadmap

- [ ] Implement `SETEX` (Keys with TTL / Expiration).
- [ ] Implement Append-Only File (AOF) persistence to save the in-memory map to disk.
- [ ] Implement `DEL` and `EXISTS` commands.
