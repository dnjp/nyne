<div style="text-align:center">
  <img src="https://git.sr.ht/~danieljamespost/nyne/blob/master/resources/glenda.jpg" alt="drawing" width="200"/>
  <p style="font-style: italic;">By Renée French</p>
</div>

# Nyne

Nyne automates what are typically manual tasks when using Acme. Instead
of needing to set custom indentation settings and manually run external
commands like clang-format against your file, Nyne does all of that for
you. It can also optionally expand hard tabs to soft tabs, which is a
feature not included in Acme by default.

Included in a full install of nyne are bundled utilities for acme:

- `nyne`: The core autoformatting engine that is run from within acme
- `nynetab`: Implements tab expansion for the given buffer if not configured
  for nyne
- `a+`: Indent selected source code `a-`: Unindent selected source
  code
- `com`: Comment/uncomment selected source code

## Configuration

Nyne and the bundled utilities use a
[TOML](https://github.com/toml-lang/toml) configuration file during the
build to generate static code used for all formatting rules. This file is
expected to be located at $HOME/.config/nyne/nyne.toml by default, but can
be overriden by by exporting the $NYNERULES variable in your environment
like `NYNERULES=/home/daniel/.nyne`. Copy the [example](./example.toml)
to this location and modify it for your needs before building nyne.

Looking at the `format` block in the example configuration file, you
will see a block for each language that looks like the following. The
configuration options are documented below.

```
# "go" is an arbitrary name given to this configuration block. The
# name must be unique.
[format.go]

# An array of strings that include file extensions that nyne should
# apply the given formatting rules to.
extensions = [".go"]

# An integer representing the tab width used for indentation
indent = 8

# A boolean that determines whether to use hard tabs or spaces for
# indentation
tabexpand = false

# A string that contains the comment style for the given language.
commentstyle = "// "

    # The "commands" blocks is used to define the external program to
    # be run against against your buffer on file save. Any number of these
    # blocks may be defined.
    [[format.go.commands]]

    # A string representing the executable used to format the buffer
    exec = "gofmt"

    # An array of strings containing the arguments to the
    # executable. $NAME is a macro that will be replaced with the absolute
    # path to the file you are working on. This is a required argument.
    args = [ "$NAME" ]

    # A boolean representing whether the executable will print to
    # stdout. If the command writes the file in place, be sure to set this
    # to false.
    printsToStdout = true
```


Think of `nyne.toml` as the equavalent of a `config.def.h` file used to
configure many C programs. Because this file is used to generate static
configuration for Nyne that is baked into the binary, any changes made
to this file after build will not be noticed. In order for the changes
to be picked up, you must rebuild nyne and restart the `nyne` executable
if already running.

## Install

The build for nyne assumes that you have
followed the configuration instructions above and have
[plan9port](https://github.com/9fans/plan9port) utilities installed and
available in your `$PATH`. If you cannot execute `mk` in particular,
head over to the main [plan9port](https://9fans.github.io/plan9port/)
page and follow the installation guide for your system.

To install nyne, first clone this repository:

```
%: git clone https://git.sr.ht/~danieljamespost/nyne
```

Then use [mk](https://9fans.github.io/plan9port/man/man1/mk.html) to
build the nyne binaries:

```
%: mk
```

This will build `nyne`, `nynetab`, `a+`, `a-`, and `com` and place them
in `./bin`. Once they are built, you are ready to install them to a
directory in your `$PATH`. On my system, I keep commands used for Acme in
`$home/bin` and have that directory added to my path. Assuming your system
is setup like mine, you can install nyne and the bundled utilities with:

```
%: installdir=$home/bin mk install
```

If the commands above completed successfully, you should now be able to
execute any of the nyne utilities.

To cleanup the build files simply run `mk nuke`. To uninstall nyne and
all utilities, run `installdir=$home/bin mk uninstall`.

## Usage

### nyne

Once you have built and installed nyne, simply execute `nyne` in acme
by middle clicking on the text "nyne" typed in the upper most window
tag. Nyne will watch for windows to be opened that match any of the
extensions you have configured. If it finds a match, it will write the
menu options you've configured to the scratch area and begin listening
for file save events received when you middle click `Put`. When this
event is received, it will format the buffer using your configured
external formatting programs. If the program does not print to stdout,
a new file will be written to `/tmp`, formatted using youc configured
commands, and the output applied to your active buffer in acme.

If `tabexpand` is enabled for a given file extension, `nynetab` will be
used to convert tabs to spaces when you enter `tab` with your keyboard.

### nynetab

Nynetab is what is used under the hood for tab expansion in nyne. If
you are editing a buffer that is not managed by nyne, simply execute
`nynetab <tab size>` in an Acme buffer to begin tab expansion.

### a+/a-

`a+` and `a-` use your indentation settings to indent or unindent your
selection in acme using either tabs or spaces depending on what is
configured. To use these commands, write `|a+` or `|a-` to the scratch
area in your acme window, select the text you want to indent, and then
middle click on `|a+` to indent or `|a-` to unindent your selection.

### com

`com` is uses the `commentstyle` you've configured a given file extension
to comment or uncomment a given selection in acme. Just as with `a+`
or `a-`, you can use `com` by writing `|com` to your scratch area,
selecting the text you want to un/comment, and then middle click on
`|com` to execute the command.

## Bugs or Feature Requests

If you find a bug or if there is a feature you would like
to see in Nyne, please post it in the official [Issue
Tracker](https://todo.sr.ht/~danieljamespost/nyne).

## Contributing

Please do! Nyne uses an email based workflow for managing patches. If
you've never used the `git-send-email` command before, checkout
this [interactive guide](https://git-send-email.io/) for how
to set it up and get comfortable with the workflow. Patches
or questions about the submission process can be sent to
[~danieljamespost/nyne@lists.sr.ht](mailto:~danieljamespost/nyne@lists.sr.ht).
