// === README.md ===
# GSocket Clone

## Components
- `gserver.go`: Relay server.
- `gs-client.go`: Client shell (reverse shell).
- `gs-listener.go`: Admin shell (interactive listener).
- `tun.go`: TUN/TAP interface handler (placeholder).
- `builder.go`: Generates key and injects into files.

## Usage
1. Build all components:
```sh
go build gserver.go
./gserver
```

```sh
go run builder.go
```

2. Deploy `gs-client` on target.
3. Run `gs-listener` to control session.

**Note:** Ensure port 443 is accessible on `gserver` side. You can use WebSocket tunneling or port 443 TCP directly.

## Security
- AES-CFB encryption
- Shared key auth

## Roadmap
- XOR toggle flag
- Web panel in Go
- WebSocket upgrade
- Auto reconnection
- VPN bridge via TUN/TAP
