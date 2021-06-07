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
./rotoxy --tors 1 --port 8080 --circuitInterval 30 # Run with custom parameters
```
### Docker
```bash
docker run -p 8080:8080 gtuk/rotoxy:0.2.0 # Run with default parameters
docker run -p 8088:8088 gtuk/rotoxy:0.2.0 --tors 1 --port 8080 --circuitInterval 30 # Run with custom parameters
```

### TODOS
* Tests
* Better documentation
