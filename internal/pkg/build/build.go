package build

import (
	"flag"
	"log"
)

func RunCmd() {

	if flag.NArg() == 0 {
		log.Fatal("No command specified")
	}

	for _, cmd := range flag.Args() {
		switch cmd {
		case "docs":
			Docs()
		default:
			log.Fatalf("Unknown command %s", cmd)
		}
	}

}
