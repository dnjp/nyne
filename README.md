<div style="text-align:center">
  <img src="https://raw.githubusercontent.com/dnjp/nyne/master/resources/glenda.jpg" alt="drawing" width="200"/>
  <p style="font-style: italic;">By Renée French</p>
</div>

# nyne

nyne is a library and a collection of tools that enables interacting
with acme in ways that are more intuitive when coming from more
traditional text editors. As a library, nyne provides many abstractions
on top of the [9fans acme library](https://pkg.go.dev/9fans.net/go/acme)
to make event handling, finding focused windows, and other actions
significantly easier (see
[Event](https://pkg.go.dev/github.com/dnjp/nyne#Event) and the
provided [variables](https://pkg.go.dev/github.com/dnjp/nyne#pkg-variables)
for example). The included commands make working in acme faster and more user friendly, especially when combined with a keyboard mapping tool like [skhd](https://github.com/koekeishiya/skhd).

These are the commands that are included:

* [cmd/a+](./cmd/a+): Indent selected source code
* [cmd/a-](./cmd/a-): Unindent selected source code
* [cmd/aspell](./cmd/aspell): A spell checker for acme
* [cmd/com](./cmd/com): Comments/uncomments piped text
* [cmd/f+](./cmd/f+): Increase font size
* [cmd/f-](./cmd/f-): Decrease font size
* [cmd/font](./cmd/font): Wrapper around f+ or f- intended to be invoked from a tool like skhd
* [cmd/md](./cmd/md): Shortcuts for working with markdown
* [cmd/move](./cmd/move): Shortcuts for moving the cursor
* [cmd/nyne](./cmd/nyne): The core autoformatting engine that is run from within acme
* [cmd/nynetab](./cmd/nynetab): Implements tab expansion and indentation
* [cmd/save](./cmd/save): Utility to execute Put via keyboard bindings
* [cmd/xcom](./cmd/xcom): Wrapper around `com` intended to be invoked from a tool like skhd
* [cmd/xec](./cmd/xec): Execute a command in the focused window as if it had been clicked with B2

## Configuration

Nyne and the bundled utilities are configured in
[config.go](https://github.com/dnjp/nyne/blob/master/config.go)
which defineshow nyne handles different file types, what to write
to the menu, etc.  Alter this file to your liking before building
and installing nyne.

Several of the included tools are intended to be called from a tool
like [skhd](https://github.com/koekeishiya/skhd) which allows for
overriding the application handlers for particular key bindings.
See [skhdrc](./skhdrc) for an example of how to use nyne tooling
with skhd.

## Install

To install nyne, first make sure that you have properly installed
[Go](https://go.dev/learn/) and then execute the following commands:

```
% git clone https://github.com/dnjp/nyne
% cd nyne
% go install ./...
```

This will build and install the included commands and the nyne library itself.

## Bugs or Feature Requests

Please feel free to file an [issue](https://github.com/dnjp/nyne/issues) if you
run into any bugs or problems during normal usage. Should you have any questions
about how to use or setup nyne, you can start a new
[discussion](https://github.com/dnjp/nyne/discussions) thread.

## Contributing

Please do! If you have a fix for a bug or a new feature you'd like added to
nyne, please fork this repository, commit your changes to a new branch, and
submit a [PR](https://github.com/dnjp/nyne/pulls) with your changes.
