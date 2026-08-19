package main

import (
	"backend-test/cmd"
	"backend-test/helpers"
)

func main() {
	// load config
	helpers.SetupConfig()

	// load log
	helpers.SetupLogger()

	// load db
	helpers.SetupPostgreSQL()

	// run http
	cmd.ServeHTTP()
}
