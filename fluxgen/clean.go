package fluxgen

import (
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

	printMsg("_site/ cleaned", "tick")

}
