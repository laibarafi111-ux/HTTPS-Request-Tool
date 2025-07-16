# HTTPS Request Tool

This is a command-line tool designed to perform HTTPS GET requests with advanced anti-bot evasion features, including TLS fingerprint (ClientHello) spoofing and proxy rotation.

## Core Features

-   **Customizable HTTPS Requests**: Send GET requests to any target URL.
-   **TLS Fingerprint Spoofing**: Mimics the TLS handshakes of popular browsers like Chrome, Firefox, and Safari.
-   **Proxy Support**: Routes all traffic through SOCKS5 proxies.
-   **Proxy Rotation**: Automatically rotates through a provided list of proxies for each request.

## Prerequisites

-   Go (version 1.20 or newer) installed on your system.
-   Git installed on your system.

## How to Run

1.  **Clone the repository:** Open your terminal or command prompt and run the following command. Replace `[URL of your repository]` with the actual URL from your browser's address bar.
    ```sh
    git clone [URL of your repository]
    ```

2.  **Navigate into the directory:**
    ```sh
    cd HTTPS-Request-Tool
    ```

3.  **Create a proxy file:** Create a new file in this directory named `proxies.txt`. Add your SOCKS5 proxies to this file, one per line, in the format `user:pass@ip:port`.

4.  **Run the tool:** Use the `go run` command from your terminal.

    **Basic Example (sends 1 request to the default test URL):**
    ```sh
    go run .
    ```

    **Advanced Example (sends 5 requests to a specific URL using the Firefox TLS profile):**
    ```sh
    go run . -url="https://www.example.com" -n=5 -profile="firefox"
    ```

### Command-Line Flags

-   `-url`: The target URL for the HTTPS request. (Default: `https://httpbin.org/get` )
-   `-proxies`: The path to your proxy list file. (Default: `proxies.txt`)
-   `-n`: The total number of requests to send. (Default: `1`)
-   `-profile`: The TLS profile to use. Options are `chrome`, `firefox`, `safari`, or `random`. (Default: `random`)
