package alfajor

import (
	"log"
	"os"
)

func (Alfajor Alfajor) mkdir(folderName string) {

	if err := os.Mkdir(folderName, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
