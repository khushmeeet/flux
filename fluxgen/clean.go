package fluxgen

import (
	"fmt"
	"log"
	"os"
)

func FluxClean() {
	err := os.RemoveAll(SiteDir)
	if err != nil {
		log.Fatal("Unable to delete _site folder")
	}

	err = os.Mkdir(SiteDir, 0777)
	if err != nil {
		log.Fatal("Unable to create _site folder")
	}

	fmt.Printf("Flux %v cleaned!\n", SiteDir)

}
