# Rotoxy

A rotating tor proxy service that starts a configurable number of tor socks proxies and expose them under one reverse proxy (http).
The tor socks proxies are randomly selected by the reverse proxy

### Prerequisites
In order to use the tool you need have Tor installed on the machine

### Usage
Download the latest release from github

```bash
./rotoxy --help # Show usage
./rotoxy # Run with default parameters
./rotoxy -tors=1 -port=8080 -circuitInterval=30 # Run with custom parameters
```

### TODOS
* Tests
* Better documentation
* Docker image
