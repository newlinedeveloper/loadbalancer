package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
)

// Server interface defines the methods required for a server in the load balancer
type Server interface {
	Address() string
	IsAlive() bool
	HealthCheck() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

// simpleServer represents a backend server
type simpleServer struct {
	addr        string
	healthCheck string
	proxy       *httputil.ReverseProxy
	mu          sync.Mutex
}

func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.HealthCheck()
}

func (s *simpleServer) HealthCheck() bool {
	// Perform a simple HTTP health check
	resp, err := http.Get(s.healthCheck)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func newSimpleServer(addr string, healthCheck string) *simpleServer {
	serverURL, err := url.Parse(addr)
	handleErr(err)

	return &simpleServer{
		addr:        addr,
		healthCheck: healthCheck,
		proxy:       httputil.NewSingleHostReverseProxy(serverURL),
	}
}

// LoadBalancer manages multiple backend servers and distributes requests
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func NewLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++

	return server
}

func (lb *LoadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address %q\n", targetServer.Address())

	targetServer.Serve(rw, req)
}

func main() {
	// Create three backend servers with health check endpoints
	servers := []Server{
		newSimpleServer("https://www.facebook.com", "https://www.facebook.com/health"),
		newSimpleServer("https://www.bing.com", "https://www.bing.com/health"),
		newSimpleServer("https://www.duckduckgo.com", "https://www.duckduckgo.com/health"),
	}

	// Create a load balancer
	lb := NewLoadBalancer("8000", servers)

	// Handle all requests by forwarding them to the load balancer
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req)
	}

	// Register the proxy handler for all requests
	http.HandleFunc("/", handleRedirect)

	// Start the load balancer
	fmt.Printf("Serving requests at 'localhost:%s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}

// handleErr prints the error and exits the program
func handleErr(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
