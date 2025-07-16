package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	utls "github.com/refraction-networking/utls"
	"golang.org/x/net/proxy"
 )

var proxyList []string

func main() {
	rand.Seed(time.Now().UnixNano())

	targetURL := flag.String("url", "https://httpbin.org/get", "Target URL to send the request to" )
	proxyFile := flag.String("proxies", "proxies.txt", "File containing a list of proxies (format: user:pass@ip:port)")
	numRequests := flag.Int("n", 1, "Number of requests to send")
	tlsProfile := flag.String("profile", "random", "TLS profile to use (chrome, firefox, safari, random)")
	flag.Parse()

	err := loadProxies(*proxyFile)
	if err != nil {
		log.Fatalf("FATAL: Could not load proxies from %s. Error: %v", *proxyFile, err)
	}
	if len(proxyList) == 0 {
		log.Fatal("FATAL: Proxy file is empty or could not be read. At least one proxy is required.")
	}
	log.Printf("INFO: Successfully loaded %d proxies.\n", len(proxyList))

	for i := 0; i < *numRequests; i++ {
		log.Printf("----------- Starting Request %d of %d -----------", i+1, *numRequests)
		sendRequest(*targetURL, *tlsProfile)
		log.Printf("--------------------------------------------------\n\n")
		time.Sleep(1 * time.Second)
	}
}

func loadProxies(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open proxy file '%s': %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			proxyList = append(proxyList, line)
		}
	}
	return scanner.Err()
}

func sendRequest(targetURL, profile string) {
	randomProxyStr := proxyList[rand.Intn(len(proxyList))]
	log.Printf("INFO: Selected Proxy: %s", strings.Split(randomProxyStr, "@")[1])

	proxyURL, err := url.Parse("http://" + randomProxyStr )
	if err != nil {
		log.Printf("ERROR: Failed to parse proxy URL: %v\n", err)
		return
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		log.Printf("ERROR: Failed to create proxy dialer: %v\n", err)
		return
	}

	transport := &http.Transport{
		DialContext: dialer.(proxy.ContextDialer ).DialContext,
		DialTLS: func(network, addr string) (net.Conn, error) {
			return dialTLSWithUTLS(network, addr, profile)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   20 * time.Second,
	}

	req, err := http.NewRequest("GET", targetURL, nil )
	if err != nil {
		log.Printf("ERROR: Failed to create request: %v\n", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	log.Println("INFO: Sending request...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("RESULT: Request Failed. Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("RESULT: Success! Status Code: %d\n", resp.StatusCode)
	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	bodyString := string(bodyBytes)

	if strings.Contains(strings.ToLower(bodyString), "captcha") || strings.Contains(strings.ToLower(bodyString), "challenge") {
		log.Println("STATUS: Potential challenge detected in response body.")
	} else {
		log.Println("STATUS: No obvious challenge detected.")
	}
}

func dialTLSWithUTLS(network, addr, profile string) (net.Conn, error) {
	var clientHello utls.ClientHelloID
	switch strings.ToLower(profile) {
	case "firefox":
		clientHello = utls.HelloFirefox_102
	case "safari":
		clientHello = utls.HelloSafari_16_0
	case "chrome":
		clientHello = utls.HelloChrome_106
	default: // random
		profiles := []utls.ClientHelloID{utls.HelloChrome_106, utls.HelloFirefox_102, utls.HelloSafari_16_0}
		clientHello = profiles[rand.Intn(len(profiles))]
	}

	log.Printf("INFO: Using TLS Profile: %s", clientHello.Str())

	config := &utls.Config{
		InsecureSkipVerify: true,
	}

	dialConn, err := net.DialTimeout(network, addr, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("net.DialTimeout error: %w", err)
	}

	uconn := utls.UClient(dialConn, config, clientHello)
	if err := uconn.Handshake(); err != nil {
		return nil, fmt.Errorf("uTLS Handshake error: %w", err)
	}

	return uconn, nil
}
