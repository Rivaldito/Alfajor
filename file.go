package alfajor

import (
	"log"
	"os"
)

func (Alfajor Alfajor) createfile(fileName string) {

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
}

func (Alfajor Alfajor) appendText(text string, file string) {

	//Opend the file
	f, err := os.OpenFile(file,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	//Write the file
	if _, err := f.WriteString(text + lineBreak); err != nil {
		log.Println(err)
	}
}
