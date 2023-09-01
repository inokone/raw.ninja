package web

import (
	"log"
	"os"
)

func main() {
	logfile, err := os.Create("app.log")

	if err != nil {
		log.Fatal(err)
	}

	defer logfile.Close()
	log.SetOutput(logfile)
}
