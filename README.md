<div style="text-align:center">
  <img src="https://raw.githubusercontent.com/dnjp/nyne/master/resources/glenda.jpg" alt="drawing" width="200"/>
  <p style="font-style: italic;">By Renée French</p>
</div>

# Nyne

Nyne automates what are typically manual tasks when using Acme. Think of nyne as
the `.vimrc` of Acme. Instead of needing to set custom indentation settings and
manually run external commands like clang-format against your file, nyne does
all of that for you. It can also optionally expand hard tabs to soft tabs, which
is a feature not included in Acme by default.

Included in a full install of nyne are bundled utilities for acme:

- `nyne`: The core autoformatting engine that is run from within acme
- `nynetab`: Implements tab expansion for the given buffer if not configured for
  nyne
- `a+`: Indent selected source code
- `a-`: Unindent selected source code
- `com`: Comment/uncomment selected source code

## Configuration

Nyne and the bundled utilities use a [configuration file](./config.go) during
the build to generate static code used for all formatting rules. Think of
`config.go` as roughly the equavalent of a `config.h` file used to configure
many C programs. Because this file is used to generate static configuration for
Nyne that is baked into the binary, any changes made to this file after build
will not be noticed. In order for the changes to be picked up, you must rebuild
nyne and restart the `nyne` executable if already running. This has the added
benefit that the bundled utilities can be executed without nyne running and
without having to re-read a config file while maintaining all of your
configuration options. The available configuration options are documented in
[./config.go](./config.go).

## Install

The build for nyne assumes that you have followed the configuration instructions
above and have [plan9port](https://github.com/9fans/plan9port) utilities
installed and available in your `$PATH`. If you cannot execute `mk` in
particular, head over to the main
[plan9port](https://9fans.github.io/plan9port/) page and follow the installation
guide for your system.

To install nyne, first clone this repository:

```
%: git clone https://github.com/dnjp/nyne --branch 0.1.1 --single-branch
```

Then use [mk](https://9fans.github.io/plan9port/man/man1/mk.html) to build the
nyne binaries:

```
%: mk
```

This will build `nyne`, `nynetab`, `a+`, `a-`, and `com` and place them in
`./bin`. Once they are built, you are ready to install them to a directory in
your `$PATH`. On my system, I keep commands used for Acme in `$home/bin` and
have that directory added to my path. Assuming your system is setup like mine,
you can install nyne and the bundled utilities with:

```
%: installdir=$home/bin mk install
```

If the commands above completed successfully, you should now be able to execute
any of the nyne utilities.

To cleanup the build files simply run `mk nuke`. To uninstall nyne and all
utilities, run `installdir=$home/bin mk uninstall`.

## Usage

### nyne

Once you have built and installed nyne, simply execute `nyne` in acme by middle
clicking on the text "nyne" typed in the upper most window tag. Nyne will watch
for windows to be opened that match any of the extensions you have configured.
If it finds a match, it will write the menu options you've configured to the
scratch area and begin listening for file save events received when you middle
click `Put`. When this event is received, it will format the buffer using your
configured external formatting programs. If the program does not print to
stdout, a new file will be written to `/tmp`, formatted using youc configured
commands, and the output applied to your active buffer in acme.

If `tabexpand` is enabled for a given file extension, `nynetab` will be used to
convert tabs to spaces when you enter `tab` with your keyboard.

### nynetab

Nynetab is what is used under the hood for tab expansion in nyne. If you are
editing a buffer that has an extension not configured for nyne, simply execute
`nynetab <tab size>` in an Acme buffer to begin tab expansion. Otherwise, simply
executing `nynetab` will start tab expansion using your configured settings.

### a+/a-

`a+` and `a-` use your indentation settings to indent or unindent your selection
in acme using either tabs or spaces depending on what is configured. To use
these commands, write `|a+` or `|a-` to the scratch area in your acme window,
select the text you want to indent, and then middle click on `|a+` to indent or
`|a-` to unindent your selection.

### com

`com` uses the `commentstyle` you've configured a given file extension to
comment or uncomment a given selection in acme. Just as with `a+` or `a-`, you
can use `com` by writing `|com` to your scratch area, selecting the text you
want to un/comment, and then middle click on `|com` to execute the command.

## Bugs or Feature Requests

If you find a bug or if there is a feature you would like to see in Nyne, please
post it in the official
[Issue Tracker](https://todo.sr.ht/~danieljamespost/nyne).

## Contributing

Please feel free to file an issue if you run into any bugs or problems during
normal usage. Should you have any questions about how to use or setup nyne, you
can start a new [discussion](https://github.com/dnjp/nyne/discussions) thread.
Of course, if you have a fix for a bug or a new feature you'd like added to
nyne, please fork this repository, commit your changes to a new branch, and
submit a PR with your changes.
