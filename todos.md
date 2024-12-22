
- [ ] add an endpoint we can use to check if a server is actually another DWMBT instance 
  - [ ] eventually: check if it's a compatible version
- [ ] add authentication
- [x] use a Context with a timeout and `exec.CommandContext()` to avoid blocking forever waiting for external commands
- [ ] parse and consistently format MAC addresses
  - `bluetoothctl` formats MACs like `CC:98:8B:20:7D:DB` while `blueutil` formats them like `cc-98-8b-20-7d-db`, so output is inconsistent
  - we should probably parse these all into a generic MAC address type, format them consistently when displaying them to the user, and format them appropriately for each `BluetoothManager` backend as well
    - https://pkg.go.dev/net#HardwareAddr consistently formats its output with colons but otherwise doesn't seem ideal
    - possibly useful: https://github.com/thatmattlove/go-macaddr/tree/main
