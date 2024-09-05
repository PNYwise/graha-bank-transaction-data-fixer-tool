package main

import (
	"log"

	"github.com/PNYwise/graha-bank-transaction-data-fixer-tool/internal"
)

func main()  {
	
	/**
	 * Open DB connection
	 *
	**/
	internal.ConnectDb()
	defer func() {
		if err := internal.CloseDb(); err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}()

	if err := internal.Ping(); err != nil {
		log.Fatalf("Error ping database connection: %v", err)}
}