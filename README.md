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
- `nynetab`: Implements tab expansion and indentation
- `save`: Utility to execute Put via keyboard bindings
- `a+`: Indent selected source code
- `a-`: Unindent selected source code
- `com`: Comment/uncomment selected source code
- `xcom`: Wrapper arround `com` intended to be invoked from a tool like skhd.
- `md`: Shortcuts for working with markdown

```
% md -h
Usage of md:
  -op string
    	the operation to perform: link, bold, italic, preview
```

## Configuration

Nyne and the bundled utilities use a [configuration
file](https://github.com/dnjp/nyne/blob/master/config.go) to configure
how it reacts to different file types, what to write to the menu, etc.
Alter this file to your liking before building and installing nyne.

`nynetab` is intended to be called from a tool like
[skhd](https://github.com/koekeishiya/skhd) which allows for
overriding the application handlers for particular key bindings. To
use nynetab with skhd, add something like the following to your
[skhdrc](https://github.com/koekeishiya/skhd/blob/master/examples/skhdrc):

```
tab [
  "acme" : nynetab
  "edwood" : nynetab
]

shift - tab [
  "acme" : nynetab -unindent=true
  "edwood" : nynetab -unindent=true
]
```

With these keybindings, entering `tab` without text selected will
insert either a hard or soft tab depending on the file type
configuration. If text is selected, entering `tab` will indent the
selected text and entering `shift+tab` will unindent the selected
text.

`save` is also intended to be called from something like skhd. `save`
simply executes `Put` on the focused window when invoked. This is how
I use it with skhd:

```
cmd - s [
  "acme" : save
  "edwood" : save
]
```

Similarly, you can map `com` to `cmd+/` with the following skhd settings:

```
cmd - 0x2C  [
  "acme" : xcom
  "edwood" : xcom
]
```

These are shortcuts for the `md` tool that are similar to that of
other text editors:

```
cmd - k [
  "acme" : md -op link
  "edwood" : md -op link
]

cmd - b [
  "acme" : md -op bold
  "edwood" : md -op bold
]

cmd - i [
  "acme" : md -op italic
  "edwood" : md -op italic
]

shift + cmd - r [
  "acme" : md -op preview
  "edwood" : md -op preview
]
```

## Install

To install nyne, first make sure that you have properly installed
[Go](https://go.dev/learn/) and then execute the following commands:

```
% git clone https://github.com/dnjp/nyne
% cd nyne
% go install ./...
```

This will build `nyne`, `nynetab`, `a+`, `a-`, and `com` and place them in
`$GOPATH/bin`.

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

Please feel free to file an [issue](https://github.com/dnjp/nyne/issues) if you
run into any bugs or problems during normal usage. Should you have any questions
about how to use or setup nyne, you can start a new
[discussion](https://github.com/dnjp/nyne/discussions) thread.

## Contributing

Please do! If you have a fix for a bug or a new feature you'd like added to
nyne, please fork this repository, commit your changes to a new branch, and
submit a [PR](https://github.com/dnjp/nyne/pulls) with your changes.
