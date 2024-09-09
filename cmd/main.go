package main

import (
	"log"

	"github.com/PNYwise/graha-bank-transaction-data-fixer-tool/internal"
)

const (
	purchasePrefix       string = "PB"
	purchaseReturnPrefix string = "RB"
	salePrefix           string = "PJ"
	saleReturnPrefix     string = "RJ"
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
	purchaseReturnRepo := internal.NewPurchaseReturnRepository(internal.DB.Db)
	saleRepo := internal.NewSaleRepository(internal.DB.Db)
	saleReturnRepo := internal.NewSaleReturnRepository(internal.DB.Db)

	bankTransactionRepo := internal.NewBankTransactionRepository(internal.DB.Db)
	masterBankTransactionRepo := internal.NewMasterBankTransactionRepository(internal.DB.Db)

	purchases, err := purchaseRepo.FindPurchaseWithOutBankTransaction()
	if err != nil {
		log.Fatalf("%v", err)
	}

	purchaseReturns, err := purchaseReturnRepo.FindPurchaseReturnWithOutBankTransaction()
	if err != nil {
		log.Fatalf("%v", err)
	}

	sales, err := saleRepo.FindSaleWithOutBankTransaction()
	if err != nil {
		log.Fatalf("%v", err)
	}

	saleReturns, err := saleReturnRepo.FindSaleReturnWithOutBankTransaction()
	if err != nil {
		log.Fatalf("%v", err)
	}

	purchaseBankTransactions, err := bankTransactionRepo.FindBankTransactionWithDateNotEqualPurchaseDate()
	if err != nil {
		log.Fatalf("%v", err)
	}
	purchaseReturnBankTransactions, err := bankTransactionRepo.FindBankTransactionWithDateNotEqualPurchaseReturnDate()
	if err != nil {
		log.Fatalf("%v", err)
	}

	saleBankTransactions, err := bankTransactionRepo.FindBankTransactionWithDateNotEqualSaleDate()
	if err != nil {
		log.Fatalf("%v", err)
	}
	saleReturnBankTransactions, err := bankTransactionRepo.FindBankTransactionWithDateNotEqualSaleReturnDate()
	if err != nil {
		log.Fatalf("%v", err)
	}

	masterBankTransactions, err := masterBankTransactionRepo.Find([]string{"PURCHPAY", "PPNIN", "PURCHRET", "PPNINRET", "PPNOUT", "SALESPAY", "INCDELIV", "PPNOUTRET", "SALESRET"})
	if err != nil {
		log.Fatalf("%v", err)
	}

	PURCHPAY := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "PURCHPAY"
	})
	PPNIN := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "PPNIN"
	})
	PURCHRET := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "PURCHRET"
	})
	PPNINRET := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "PPNINRET"
	})
	SALESPAY := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "SALESPAY"
	})
	PPNOUT := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "PPNOUT"
	})

	INCDELIV := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "INCDELIV"
	})
	PPNOUTRET := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "PPNOUTRET"
	})
	SALESRET := internal.Find(masterBankTransactions, func(v internal.MasterBankTransactionEntity) bool {
		return v.Code == "SALESRET"
	})

	if PURCHPAY == nil || PPNIN == nil || PURCHRET == nil || PPNINRET == nil || PPNOUT == nil || SALESPAY == nil || INCDELIV == nil || PPNOUTRET == nil || SALESRET == nil {
		log.Fatalln("some masterBankTransactions not found")
	}

	// Update purchase bank transaction
	for _, bankTransaction := range *purchaseBankTransactions {
		purchase := bankTransaction.Purchase
		date := purchase.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, purchasePrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(purchasePrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}
		bankTransaction.BankNumber = newCode
		bankTransaction.Date = date
		bankTransactionRepo.Update(&bankTransaction)

	}

	// Update purchase return bank transaction
	for _, bankTransaction := range *purchaseReturnBankTransactions {
		purchaseReturn := bankTransaction.PurchaseReturn
		date := purchaseReturn.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, purchaseReturnPrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(purchaseReturnPrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}
		bankTransaction.BankNumber = newCode
		bankTransaction.Date = date
		bankTransactionRepo.Update(&bankTransaction)

	}

	// Update sale bank transaction
	for _, bankTransaction := range *saleBankTransactions {
		sale := bankTransaction.Sale
		date := sale.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, salePrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(salePrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}
		bankTransaction.BankNumber = newCode
		bankTransaction.Date = date
		bankTransactionRepo.Update(&bankTransaction)

	}

	// Update sale return bank transaction
	for _, bankTransaction := range *saleReturnBankTransactions {
		saleReturn := bankTransaction.SaleReturn
		date := saleReturn.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, saleReturnPrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(saleReturnPrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}
		bankTransaction.BankNumber = newCode
		bankTransaction.Date = date
		bankTransactionRepo.Update(&bankTransaction)

	}

	// Create purchase bank transaction
	for _, purchase := range *purchases {
		date := purchase.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, purchasePrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(purchasePrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}

		purchaseId := purchase.ID
		bankTransaction := internal.BankTransactionEntity{
			BankNumber: newCode,
			Date:       date,
			Note:       "-",
			Amount:     float64(purchase.Total),
			IsMain:     true,
			PurchaseID: &purchaseId,
			BankTransactionItems: []internal.BankTransactionItemEntity{
				{Note: purchase.Code, SubAmount: float64(purchase.TotalNotIncludingPPN), MasterBankTransactionID: PURCHPAY.ID},
				{Note: purchase.Code, SubAmount: purchase.PPNInValue, MasterBankTransactionID: PPNIN.ID},
			},
		}
		bankTransactionRepo.Create(&bankTransaction)
	}

	// Create purchase return bank transaction
	for _, purchaseReturn := range *purchaseReturns {
		date := purchaseReturn.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, purchaseReturnPrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(purchaseReturnPrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}

		purchaseReturnId := purchaseReturn.ID
		bankTransaction := internal.BankTransactionEntity{
			BankNumber:       newCode,
			Date:             date,
			Note:             "-",
			Amount:           float64(purchaseReturn.Total),
			IsMain:           true,
			PurchaseReturnID: &purchaseReturnId,
			BankTransactionItems: []internal.BankTransactionItemEntity{
				{Note: purchaseReturn.Code, SubAmount: float64(purchaseReturn.TotalNotIncludingPPN), MasterBankTransactionID: PURCHRET.ID},
				{Note: purchaseReturn.Code, SubAmount: purchaseReturn.PPNInValue, MasterBankTransactionID: PPNINRET.ID},
			},
		}
		bankTransactionRepo.Create(&bankTransaction)
	}

	// Create sale bank transaction
	for _, sale := range *sales {
		date := sale.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, salePrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(salePrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}

		saleId := sale.ID
		bankTransaction := internal.BankTransactionEntity{
			BankNumber: newCode,
			Date:       date,
			Note:       "-",
			Amount:     float64(sale.Total) + float64(sale.DeliveryFee) - float64(sale.DownPayment),
			IsMain:     true,
			SaleID:     &saleId,
			BankTransactionItems: []internal.BankTransactionItemEntity{
				{Note: sale.SaleCode, SubAmount: float64(sale.TotalNotIncludingPpn), MasterBankTransactionID: SALESPAY.ID},
				{Note: sale.SaleCode, SubAmount: float64(sale.PpnInValue), MasterBankTransactionID: PPNOUT.ID},
			},
		}
		if sale.DeliveryFee > 0 {
			bankTransaction.BankTransactionItems = append(bankTransaction.BankTransactionItems, internal.BankTransactionItemEntity{
				Note: sale.SaleCode, SubAmount: float64(sale.DeliveryFee), MasterBankTransactionID: INCDELIV.ID,
			})
		}
		bankTransactionRepo.Create(&bankTransaction)
	}

	// Create sale return bank transaction
	for _, saleReturn := range *saleReturns {
		date := saleReturn.Date
		lastCode, err := bankTransactionRepo.FindLastCode(date, saleReturnPrefix)
		if err != nil {
			log.Fatalf("Error FIndLastCode %v", err)
		}

		newCode, err := internal.GetNewCode(saleReturnPrefix, date, lastCode)
		if err != nil {
			log.Fatalf("%v", err)
		}

		saleReturnId := saleReturn.ID
		bankTransaction := internal.BankTransactionEntity{
			BankNumber: newCode,
			Date:       date,
			Note:       "-",
			Amount:     float64(saleReturn.Total),
			IsMain:     true,
			SaleID:     &saleReturnId,
			BankTransactionItems: []internal.BankTransactionItemEntity{
				{Note: saleReturn.Code, SubAmount: float64(saleReturn.TotalNotIncludingPpn), MasterBankTransactionID: PPNOUTRET.ID},
				{Note: saleReturn.Code, SubAmount: float64(saleReturn.PpnInValue), MasterBankTransactionID: SALESRET.ID},
			},
		}
		bankTransactionRepo.Create(&bankTransaction)
	}

}
