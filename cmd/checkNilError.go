package cmd

import "log"

func CheckNilError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}