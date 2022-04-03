/*
The core autoformatting engine that is run from within acme.

	Usage of nyne:
		nyne

Once you have built and installed nyne, simply execute `nyne` in
acme by middle clicking on the text "nyne" typed in the upper most
window tag. Nyne will watch for windows to be opened that match any
of the extensions you have configured.  If it finds a match, it
will write the menu options you've configured to the scratch area
and begin listening for file save events received when you middle
click `Put`. When this event is received, it will format the buffer
using your configured external formatting programs. If the program
does not print to stdout, a new file will be written to `/tmp`,
formatted using youc configured commands, and the output applied
to your active buffer in acme. If `tabexpand` is enabled for a given
file extension, `nynetab` will be used to convert tabs to spaces
when you enter `tab` with your keyboard.

*/
package main

import (
	"log"

	"github.com/dnjp/nyne"
)

func main() {
	f, err := nyne.NewFormatter(nyne.Filetypes, nyne.Menu)
	if err != nil {
		log.Fatal(err)
	}
	err = f.Run()
	if err != nil {
		log.Fatal(err)
	}
}
