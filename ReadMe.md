# **TCP Connection Management: Faulty vs. Fixed Clients & Server Behavior**

This repository demonstrates **TCP connection handling issues** and fixes when dealing with **HTTP connections**. It includes:

- **Faulty Client (`client_faulty.go`)**: Current HTTP client causing excessive TCP socket accumulation.
- **Fixed Client (`client_fix.go`)**: Optimized client that prevents excessive open connections.
- **Server (`server.go`)**: Simulates slow responses and tracks active TCP connections.

---

## **1. Faulty Client (`client_faulty.go`)**

### **Description**
The `client_faulty.go` script represents **current prometurbo HTTP client flow** that leads to **excessive open TCP sockets** over time. This happens due to:

- Creating a **new `http.Client` for every request** instead of reusing an existing one.
- **Improper keep-alive connection handling**, leading to accumulated TCP sockets.
- **Idle connections are not closed efficiently**, causing **resource exhaustion**.

### **Behavior**
- Sends **periodic HTTP requests** to `server.go`, with **randomized additional requests**.
- Fails to reuse connections, **continuously increasing** the number of open TCP sockets.
- Does not efficiently handle **slow server responses**.

### **Expected Behavior**
**High number of open TCP sockets.**  
**Possible exhaustion of system resources.**  
**Increasing delays or failures due to accumulated connections.**  

---

## **2. Fixed Client (`client_fix.go`)**

### **Description**
The `client_fix.go` script is an **optimized version** of the faulty client that **properly manages HTTP connections**, preventing **socket accumulation**.

### **Key Fixes**
Uses a **Shared HTTP client** instead of creating a new instance per request.  
**Enables Keep-Alive connections** to reuse existing TCP sockets.  
**Implements `IdleConnTimeout`** to close idle connections efficiently. 
Optimizes connection pooling with:
MaxIdleConns: 5 → Allows up to 5 idle connections globally, ensuring connection reuse when handling multiple slow-responding servers, reducing the need to create new connections.
MaxIdleConnsPerHost: 2 → Limits each host to 2 idle connections, preventing a single slow-responding host from occupying too many resources, ensuring fair connection distribution across multiple hosts.
Calls **`CloseIdleConnections()`** when a request fails to prevent stale connections.  

### **Behavior**
- Sends **periodic HTTP requests** while ensuring that **excessive TCP sockets do not accumulate**.
- Handles **server delays gracefully** without keeping unnecessary connections open.
- Ensures **only required connections remain active**, improving **resource utilization**.

### **Expected Improvements**
**Reduced number of open TCP sockets** by reusing existing connections.  
**Better performance** due to efficient connection handling.  
**Lower resource consumption**, preventing system exhaustion.  

---

## **3. Server (`server.go`)**

### **Overview**
The server simulates **delayed HTTP responses** while tracking **active TCP connections**.

### **Key Features**
- **Simulated Delay** – Introduces a **20-second delay** before responding to requests.
- **Connection Tracking** – Monitors **active TCP connections** using the `netstat` command.
- **Client Disconnection Handling** – Checks if the client **disconnects before** the response is sent.
- **Keep-Alive Support** – Uses `Connection: keep-alive` to allow **persistent connections**.

### **How It Works**
1. **`trackServerSockets()`** runs in a **separate goroutine** to log active connections **every 5 seconds**.
2. **`slowResponseHandler()`**:
   - Increments **active connection count**.
   - Waits **20 seconds** before responding.
   - Sends **HTTP 200 OK** with a response message.
3. **`main()`** starts the **HTTP server** and calls `trackServerSockets()`.

### **Use Case**
This setup is **useful for testing client behavior** under **slow server responses** and ensuring **proper connection persistence**.

---

## **How to Run**
1. **Start the server**:
   ```sh
   go run server.go
   ```

2. **Run the faulty client** (observe increasing TCP connections):
   ```sh
   go run client_faulty.go
   ```

3. **Run the fixed client** (notice stable and efficient connection handling):
   ```sh
   go run client_fix.go
   ```

---

## **Comparison Table: Faulty vs. Fixed Client**

| Feature                 | Faulty Client | Fixed Client |
|-------------------------|-----------------|---------------|
| **Persistent Connections** | No (Creates a new TCP handshake every time) | Yes (Reuses connections) |
| **Handles Idle Connections** | No (Sockets accumulate) | Yes (Closes unused sockets) |
| **Closes Connection on Error** | No (Can leave open sockets) | Yes (Calls `CloseIdleConnections()`) |
| **Netstat TCP Count** | Keeps increasing | Stays stable |

---

## **Key Takeaways**
**Proper connection management** prevents excessive open connections.  
**Keep-Alive** improves **performance** by reusing TCP sockets.  
**Efficient handling of idle connections** reduces **resource exhaustion**.  
**Poorly managed clients can overload a server, even at low request rates**.  

This project provides a practical test scenario for **TCP socket behavior**, **slow server responses**, and **efficient connection handling in Go**.

