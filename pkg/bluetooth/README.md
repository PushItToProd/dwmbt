
`BluetoothManager` implementation standards:

- Above all else, remember that we're working with untrusted inputs.
- Prevent shell injection by using `exec.Command` to invoke external commands
  directly, not via the shell.
  - TODO: use `exec.CommandContext` with a timeout context passed from the
    daemon so we won't hang forever waiting for commands to finish.
- Always validate and sanitize input before passing it to external commands.
  Even though using `exec.Command` prevents shell injection, an attacker could
  craft an input that exploits a vulnerability in the underlying command.
  - This is done automatically by `saferBluetoothManager`, so in principle this
    is handled for you already, but it's worth keeping in mind.
