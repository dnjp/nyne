# Nyne
Nyne is a simple autoformatter for [acme](http://acme.cat-v.org/) built on top of [acmego](https://github.com/9fans/go/tree/master/acme/acmego) that executes on file save. Configuration is done through a json formatted config file located at $HOME/lib/nyne by default. You can copy example.config to this location to give it a try. The path to the config file can also be set by exporting the $NYNERULES variable in your environment like `NYNERULES=/home/daniel/.nyne`. As you can see in the example, Nyne supports running multiple commands for a group of file extensions and requires the $NAME macro to be placed in a sensible location within the args passed to each command. 

# Install
`go get git.sr.ht/~danieljamespost/nyne` and then execute `nyne` in acme.