# SOCKS5 Proxy Docker Container

A secure, IPv4-only SOCKS5 proxy for Docker, designed to run behind a VPN gateway or as a standalone service.

---

## 📦 Links
- **Docker Hub**: [boingbasti/nordvpn-socks5](https://hub.docker.com/r/boingbasti/nordvpn-socks5)  
- **GitHub Repository**: [boingbasti/docker-nordvpn-socks5](https://github.com/boingbasti/docker-nordvpn-socks5)

---

## ✨ Features
- **Minimal & Secure**: Built using a multi-stage Docker build → tiny (~10 MB) Alpine-based image with minimal attack surface.  
- **Access Control**:
  - Restrict access by IP address or subnet using `ALLOWED_IPS`.
  - Optional username + password authentication (`PROXY_USER`, `PROXY_PASSWORD`).  
- **Strictly IPv4**:
  - Blocks all incoming requests from IPv6 clients.
  - Prevents outgoing connections to IPv6 destinations → no IPv6 leaks.  
- **Easy Configuration**: Fully managed via environment variables.  
- **Informative Logging**: Clear logs for allowed connections, blocked attempts, and errors.  

---

## 🛠 Requirements
- Docker installed on your host system.

---

## 📦 Environment Variables

| Variable         | Default | Description                                                                 |
|------------------|---------|-----------------------------------------------------------------------------|
| `PROXY_PORT`     | `1080`  | The port on which the SOCKS5 proxy listens inside the container.            |
| `ALLOWED_IPS`    | *(unset)* | Comma-separated list of IPs or subnets allowed. If unset, all IPs are allowed. |
| `PROXY_USER`     | *(unset)* | Optional username for authentication (requires `PROXY_PASSWORD`).          |
| `PROXY_PASSWORD` | *(unset)* | Optional password for authentication (requires `PROXY_USER`).              |

---

## 🚀 Usage

### 1. Standalone Mode
Run the SOCKS5 proxy directly and expose it to your host machine.

```yaml
version: "3.9"

services:
  socks5-proxy:
    image: boingbasti/nordvpn-socks5:latest
    container_name: socks5-proxy
    ports:
      - "1080:1080"
    environment:
      # Allows access from the whole 192.168.1.0/24 subnet
      - ALLOWED_IPS=192.168.1.0/24
      # Optional authentication
      # - PROXY_USER=myuser
      # - PROXY_PASSWORD=mypassword
    restart: unless-stopped
```

➡️ Clients can now connect to the proxy at `HOST_IP:1080`.

---

### 2. VPN Gateway Mode (Recommended)
Attach the proxy to a VPN gateway container (e.g., [docker-nordvpn-gateway](https://github.com/boingbasti/docker-nordvpn-gateway)) so all traffic is routed through the VPN tunnel.

```yaml
# In your existing docker-compose.yml with the vpn-gateway service...
services:
  vpn-gateway:
    # ... your vpn-gateway configuration ...

  socks5-proxy:
    image: boingbasti/nordvpn-socks5:latest
    container_name: nordvpn-socks5
    # Shares the VPN gateway’s network stack
    network_mode: "service:vpn-gateway"
    depends_on:
      - vpn-gateway
    environment:
      - PROXY_PORT=1080
      - ALLOWED_IPS=192.168.1.0/24
    restart: unless-stopped
```

➡️ Clients should connect to the VPN gateway’s LAN IP (e.g., `192.168.1.240:1080`).

---

## 🔍 Troubleshooting

- **Connection refused / not working**  
  - Ensure the client’s IP is included in `ALLOWED_IPS`.  
  - In standalone mode, confirm port mapping (`1080:1080`).  

- **Authentication error**  
  - Both `PROXY_USER` and `PROXY_PASSWORD` must be set.  

- **Blocked IPv6 log messages**  
  - This is intentional. The proxy is strictly **IPv4-only**.  
