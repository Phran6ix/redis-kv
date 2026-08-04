# Redis-Compatible Key-Value Store Tasks

## 1. Network Setup
- [ ] Set up a TCP listener on port 6379 using the `net` package.
- [ ] Create an infinite loop to accept incoming client connections.
- [ ] Spawn a new goroutine for each accepted connection to handle concurrency.

## 2. RESP Parsing (Input)
- [ ] Implement a reader that reads the incoming byte stream from the connection.
- [ ] Write a parser to decode the RESP format (specifically arrays `*` and bulk strings `$`).
- [ ] Extract the command (e.g., `SET` or `GET`) and its arguments from the parsed RESP data.

## 3. Data Store
- [ ] Create a global memory store using a Go map (`map[string]string`).
- [ ] Wrap the map with a `sync.RWMutex` to prevent race conditions during concurrent access.

## 4. Command Implementation
- [ ] **SET Command:**
  - Acquire a Write Lock (`Lock()`).
  - Insert or update the key-value pair in the map.
  - Release the Write Lock (`Unlock()`).
- [ ] **GET Command:**
  - Acquire a Read Lock (`RLock()`).
  - Retrieve the value from the map using the key.
  - Release the Read Lock (`RUnlock()`).

## 5. RESP Serialization (Output)
- [ ] Implement a writer to encode responses back to RESP format.
- [ ] For `SET`: Return a simple string (`+OK\r\n`).
- [ ] For `GET` (Found): Return a bulk string (e.g., `$5\r\nvalue\r\n`).
- [ ] For `GET` (Not Found): Return a null bulk string (`$-1\r\n`).
