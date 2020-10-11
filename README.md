<div style="text-align:center">
  <img src="https://git.sr.ht/~danieljamespost/nyne/blob/0.1/resources/glenda.jpg" alt="drawing" width="200"/>
</div>

*By Renée French*

# Nyne
Nyne automates what are typically manual tasks when using Acme. Instead
of needing to set a custom indentation settings and manually run external
commands like clang-format against your file, Nyne does all of that for
you. It can also optionally expand hard tabs to soft tabs, which is a
feature not included in Acme by default.

# Nynetab
Nynetab is what is used under the hood for tab expansion in Nyne. To only
install Nynetab, run `go get git.sr.ht/~danieljamespost/nyne/cmd/nynetab`
and then execute `nynetab <tab size>` in an Acme buffer to begin tab
expansion.

# Install 
Assuming you have Go installed, simply execute the following command:
```
go get git.sr.ht/~danieljamespost/nyne/cmd/nyne
```

# Usage
Configuration for Nyne is done through a
[TOML](https://github.com/toml-lang/toml) formatted config file located
at $HOME/.config/nyne/nyne.toml by default. You can copy example.toml
to this location and modify it for your needs to give it a try. The path
to the config file can also be set by exporting the $NYNERULES variable
in your environment like `NYNERULES=/home/daniel/.nyne`. As you can see
in the example, Nyne supports running multiple commands for a group of
file extensions and requires the $NAME macro to be placed where the file
name would be in any given formatting command.

Once your configuration file is in place, simply execute `nyne` in Acme
by middle clicking on the text "nyne" typed in the upper most window tag.




 


