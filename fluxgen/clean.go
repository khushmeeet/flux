package fluxgen

import (
	"log"
	"os"
)

func FluxClean() {
	err := os.RemoveAll(SiteFolder)
	if err != nil {
		log.Fatal("Unable to delete _site folder")
	}

	err = os.Mkdir(SiteFolder, 0777)
	if err != nil {
		log.Fatal("Unable to create _site folder")
	}

}
