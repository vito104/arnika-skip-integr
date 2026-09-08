# SKIP Protocol Integration

The `arnika` service supports integration with Key Management Systems (KMS) using the **SKIP** protocol alongside the standard ETSI QKD API (ETSI GS QKD 014).

## Configuration

To use the SKIP protocol, configure the following environment variables:

| Environment Variable | Description | Default | Required for SKIP |
| :--- | :--- | :--- | :--- |
| `KMS_PROTOCOL` | Protocol to use for KMS (`etsi014` or `skip`) | `etsi014` (implicit) | Yes (`skip`) |
| `SKIP_REMOTE_SYSTEM_ID` | System ID of the peer (the peer's KP id) | `""` | Yes |
| `KMS_URL` | URL of the SKIP-compatible KMS server | — | Yes |

## Technical Details

* **Key ID Handling**: All peer-supplied `keyID` values are strictly hex-encoded and sanitized before transmission.
* **Repository Implementation**: Uses `SKIPRepository` built via `NewSKIPRepository` to fetch and manage keys.
* **Build Constraint**: Requires the Go experiment runtime secret feature enabled (`GOEXPERIMENT=runtimesecret`).