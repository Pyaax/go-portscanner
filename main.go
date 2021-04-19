package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
)

func worker(ports, results chan int, address string) {
	for p := range ports {
		address := fmt.Sprintf(address+":%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

var (
	flagStartPort = flag.Int("start", 0, "Set this flag to provide the start port.")
	flagEndPort   = flag.Int("end", 65535, "Set this flag to provide the end port.")
	flagAddress   = flag.String("address", "scanme.nmap.org", "Set this flag to provide an address.")
)

func main() {
	flag.Parse()

	if *flagStartPort < 0 || *flagStartPort > 65535 {
		log.Fatalln("The allowed port range is from 0 to 65535")
	}
	if *flagEndPort <= 0 || *flagEndPort > 65535 {
		log.Fatalln("The allowed port range is from 0 to 65535")
	}

	if *flagEndPort < *flagStartPort {
		log.Fatalln("The end port range must be greater than the start port")
	}

	ports := make(chan int, 100)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results, *flagAddress)
	}

	go func() {
		for i := *flagStartPort; i <= *flagEndPort; i++ {
			ports <- i
		}
	}()
	for i := *flagStartPort; i <= *flagEndPort; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)

	// Sorting the open ports
	sort.Ints(openports)

	for _, port := range openports {
		fmt.Printf("[+] Port %d is open\n", port)
	}
}
