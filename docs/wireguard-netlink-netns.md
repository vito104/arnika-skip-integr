# wireguard-netlink-netns

**Key writer module - installs the WireGuard PSK into a local WireGuard interface inside a network namespace through the kernel's netlink API.**

This is the single document for the `wireguard-netlink-netns` module. For the generic architecture all key reader and key writer modules follow, see [`KEYCONTROL.md`](../KEYCONTROL.md).

---

## At a Glance

| | |
|---|---|
| **Module name** | `wireguard-netlink-netns` |
| **Kind** | Key writer (sink) |
| **Build tag** | `wireguard_netlink_netns` |
| **Adapter** | [`repositories/wireguard-netlink-netns.go`](../repositories/wireguard-netlink-netns.go) |
| **Wiring** | [`wireguardnetlinknetns.go`](../wireguardnetlinknetns.go) |
| **Target** | A **local** WireGuard interface in a network namespace |
| **Transport** | `wgctrl` over netlink inside the namespace |
| **Dependencies** | `golang.zx2c4.com/wireguard/wgctrl`, `github.com/containernetworking/plugins/pkg/ns` |
| **Privileges** | `CAP_NET_ADMIN`, `CAP_SYS_ADMIN` — it reconfigures a network device inside a namespace |

---

## How the Module Works

Same as [`wireguard-netlink`](wireguard-netlink.md), but executes inside a network namespace. Uses `containernetworking/plugins/pkg/ns` to enter the namespace and then delegates to `WireguardNetlinkRepository`.

---

## Part 1 — Prepare the Host

Same as [`wireguard-netlink`](wireguard-netlink.md), but the WireGuard interface must exist inside the specified network namespace.

---

## Part 2 — Configuration Reference

| Env var | Required | Description |
|---|:---:|---|
| `WIREGUARD_INTERFACE` | yes | Name of the WireGuard interface inside the namespace |
| `WIREGUARD_PEER_PUBLIC_KEY` | yes | Public key of the peer whose PSK is rotated |
| `WIREGUARD_NETNS_PATH` | yes | Path to the network namespace (e.g., `/var/run/netns/myns`) |

---

## Part 3 — Compile

```bash
GOEXPERIMENT=runtimesecret go build -tags wireguard_netlink_netns .
```

Via the Makefile:
```bash
make build BUILD_TAGS=wireguard_netlink_netns
```

---

## Part 4 — Run

Arnika must start after the interface and namespace exist. Requires `CAP_NET_ADMIN` and `CAP_SYS_ADMIN` (for namespace operations).

---

## References

- Module architecture: [`KEYCONTROL.md`](../KEYCONTROL.md)
- Base module: [`wireguard-netlink.md`](wireguard-netlink.md)
- `wgctrl` package: <https://pkg.go.dev/golang.zx2c4.com/wireguard/wgctrl>
- `containernetworking/plugins/pkg/ns`: <https://pkg.go.dev/github.com/containernetworking/plugins/pkg/ns>
- WireGuard network namespace documentation: <https://www.wireguard.com/netns/#the-new-namespace-solution>
