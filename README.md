# Nyne
Nyne is a simple autoformatter for [acme](http://acme.cat-v.org/) built on top of [acmego](https://github.com/9fans/go/tree/master/acme/acmego) that executes on file save. Configuration is done through a json formatted config file located at $HOME/.nyne. You can copy example.config to this location to give it a try. As you can see in the example, Nyne supports running multiple commands for a group of file extensions and requires the $NAME macro to be placed in a sensible location within the args passed to each command. 

# Install
`go get git.sr.ht/~danieljamespost/nyne` and then execute `nyne` in acme.