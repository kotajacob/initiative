# initiative

A D&D assistant comprised of a tool to display on a player facing monitor and a
control tool for running combat.

## client
When opened, the client reads a config from:
`XDG_CONFIG_HOME/initiative/client.toml`. An example is included in this repo.

You can then add new entries by typing `n`, delete with `x`, navigate up/down
with `j/k`, enter initiative by typing a number, and when you're ready enter
battle mode with enter.

Entering battle mode sends the sorted list of combatants to the server and
Highlights the first combatants. While in battle mode, pressing `enter` or `j`
will move to the next combatant and pressing `backspace` or `k` will move to
the previous one.

## server
When opened, the client reads a config from:
`XDG_CONFIG_HOME/initiative/server.toml`. An example is included in this repo.

The server listens on `:6666` and uses a simple plain text TCP protocol. Each
command follows the form: `<command>,<option>,<option>...\n`.

### `start`
Sets the display variable, runs optional commands (such as turning on a monitor) and opens the window.

### `battle,<combatant>...`
Displays a list of combatants.

### `highlight,<index>`
Highlights a combatant by the 0 indexed list of combatants.

### `end`
Closes the window and runs optional commands (such as turning off a monitor),
but leaves the server listening.
