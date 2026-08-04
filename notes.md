### The `cmd/` Directory: The Entry Point

This directory holds the main applications for the project.

* **Purpose:** It should contain very little logic. Its only job is to read configuration (like the port number to bind to), wire the internal components together, and start the system.
* **Structure:** You would typically have a `cmd/server/main.go` file. If you later decided to build a custom CLI client to talk to your server, you could easily add a `cmd/client/main.go` without muddying the server logic.

### The `internal/` Directory: Private Application Code

Go treats directories named `internal` specially: code inside cannot be imported by other projects. This is where the actual mechanics of your Redis clone will live, broken down by domain.

* **`internal/server/` (The Transport Layer):** This domain is strictly responsible for network I/O. It handles listening on port 6379, accepting connections, and spawning the goroutines. It doesn't know what a `SET` command is; it just knows how to push and pull raw bytes over a TCP socket.
* **`internal/resp/` (The Protocol Layer):** This is your translation engine. It takes the raw byte stream provided by the server layer and parses it into structured commands based on the RESP specification. It also handles taking your system's responses and serializing them back into RESP format before handing them back to the network layer.
* **`internal/store/` (The State Layer):** This domain manages the hash map and the `sync.RWMutex`. It is entirely agnostic to the network. You could theoretically swap out the TCP server for an HTTP server, or the RESP parser for a JSON parser, and the `store` logic would remain completely unchanged.

### The Interface Boundary

The beauty of this structure is how the components interact.

The `server` accepts a connection and passes the data to the `resp` parser. The parser identifies a `SET` command and its arguments. The system then calls a method on the `store` to safely mutate the hash map. Once complete, the system passes a success message back through the `resp` serializer, which hands the formatted bytes back to the `server` to send over the wire.

Structuring it this way ensures that if a bug occurs in your mutex locking, you know exactly which directory to look in, entirely bypassing the networking and parsing code.
