# Moved

Source has moved to [git.sr.ht/~danieljamespost/nyne](https://git.sr.ht/~danieljamespost/nyne)

# Nyne
Nyne is an autoformatter for [acme](http://acme.cat-v.org/) that executes commands on "Put" and "New" events. Configuration is done through a json formatted config file located at $HOME/lib/nyne by default. You can copy example.config to this location to give it a try. The path to the config file can also be set by exporting the $NYNERULES variable in your environment like `NYNERULES=/home/daniel/.nyne`. As you can see in the example, Nyne supports running multiple commands for a group of file extensions and requires the $NAME macro to be placed in a sensible location within the args passed to each command.

# Dependencies
The "expand" option depends on [nynetab](https://git.sr.ht/~danieljamespost/nynetab) that implements tab expansion with the command `nynetab <width>`. Nynetab can be installed with `go get git.sr.ht/~danieljamespost/nynetab`.

# Install
`go get git.sr.ht/~danieljamespost/nyne` and then execute `nyne` in acme.
