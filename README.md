# Nyne 
Nyne is an autoformatter for [acme](http://acme.cat-v.org/) that executes
commands on "Put" and "New" events. Configuration is done through a
[TOML](https://github.com/toml-lang/toml) formatted config file located
at $HOME/.config/nyne/nyne.toml by default. You can copy example.toml
to this location to give it a try. The path to the config file can also
be set by exporting the $NYNERULES variable in your environment like
`NYNERULES=/home/daniel/.nyne`. As you can see in the example, Nyne
supports running multiple commands for a group of file extensions and
requires the $NAME macro to be placed in a sensible location within the
args passed to each command.

# Install 
`go get git.sr.ht/~danieljamespost/nyne/cmd/nyne` and then execute
`nyne` in acme.

# Nynetab
Nynetab is what is used under the hood for tab expansion in Nyne. To only
install Nynetab, run `go get git.sr.ht/~danieljamespost/nyne/cmd/nynetab`
and then execute `nynetab <tab size>` in Acme to begin tab expansion.

 


