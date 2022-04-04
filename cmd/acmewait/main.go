package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	var out []byte
	var err error
	sock := fmt.Sprintf("/tmp/ns.%s.:0/acme\n", os.Getenv("USER"))
	for !bytes.Contains(out, []byte(sock)) {
		time.Sleep(time.Second)
		cmd := exec.Command("lsof", "-U")
		out, err = cmd.CombinedOutput()
		if err != nil {
			panic(err)
		}
	}
}
