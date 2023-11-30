# Simple Go Load Balancer

The Simple Go Load Balancer is a basic load balancing application written in Go. It uses a round-robin algorithm to distribute incoming requests across multiple backend servers.

## Features

- Round-robin load balancing
- Customizable backend servers
- Basic error handling

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) installed on your machine

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/newlinedeveloper/loadbalancer.git
    ```

2. Change into the project directory:

    ```bash
    cd loadbalancer
    ```

3. Build the project:

    ```bash
    go build
    ```

### Running the Load Balancer

1. Start the load balancer:

    ```bash
    ./loadbalancer
    ```

2. The load balancer will be running on `localhost:8000`.

## Configuration

You can customize the load balancer by modifying the `main.go` file. You can add or remove backend servers and adjust weights as needed.

```go
// Example of adding a new backend server
servers := []Server{
    newSimpleServer("https://www.example.com"),
    // Add more servers as needed
}

```
## Contributing

Contributions are welcome! If you find any issues or have improvements, feel free to open a [GitHub issue](https://github.com/your-username/simple-go-load-balancer/issues) or submit a pull request.


