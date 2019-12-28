# Nyne
Nyne is an autoformatter for [acme](http://acme.cat-v.org/) that executes commands on "Put" and "New" events. Configuration is done through a json formatted config file located at $HOME/lib/nyne by default. You can copy example.config to this location to give it a try. The path to the config file can also be set by exporting the $NYNERULES variable in your environment like `NYNERULES=/home/daniel/.nyne`. As you can see in the example, Nyne supports running multiple commands for a group of file extensions and requires the $NAME macro to be placed in a sensible location within the args passed to each command. 

# Dependencies
The "expand" option in the configuration currently depends on [my fork](https://github.com/danieljamespost/edwood/tree/expandtab) of Edwood that implements tab expansion with the command "Tabexpand". If you aren't using my fork, set the "expand" option to false.

# Install
`go get git.sr.ht/~danieljamespost/nyne` and then execute `nyne` in acme.