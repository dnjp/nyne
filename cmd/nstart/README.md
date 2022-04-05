# nstart

nstart is used for launching acme along with all of its
dependencies and helpers. An example of its usage can
be found in the (acme start script) .[./../mac/Acme.app/Contents/MacOS/acme](./../mac/Acme.app/Contents/MacOS/acme).

The programs that are launched can be configured in
(config.go) .[./../config.go](./../config.go). Stderr and Stdout are
grouped together and written to the default log
file which can be found at $HOME/.config/acme/acme.log.
