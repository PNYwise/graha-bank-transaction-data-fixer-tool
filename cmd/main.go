package main

import (
	"fmt"
	"log"

	"github.com/PNYwise/graha-bank-transaction-data-fixer-tool/internal"
)

func main() {

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
		log.Fatalf("Error ping database connection: %v", err)
	}

	purchaseRepo := internal.NewPurchaseRepository(internal.DB.Db)
	bankTransactionRepo := internal.NewBankTransactionRepository(internal.DB.Db)

	_ = purchaseRepo
	bankTransactions, err := bankTransactionRepo.FindBankTransactionWithDateNotEqualPurchaseDate()
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, bankTransaction := range *bankTransactions {
		purchase := bankTransaction.Purchase
		date := purchase.Date
		prefix := "PB"
		lastCode, err := bankTransactionRepo.FindLastCode(date, prefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(prefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}

		fmt.Println(lastCode)
		fmt.Printf("NEW CODE %s \n", newCode)

		bankTransaction.BankNumber = newCode
		bankTransaction.Date = date
		bankTransactionRepo.Update(&bankTransaction)

	}
}
