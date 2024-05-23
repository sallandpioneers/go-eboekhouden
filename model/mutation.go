package model

import (
	"time"
)

type (
	MutationType string
	TaxType      string
	VATCode      string
)

const (
	InvoiceReceived        MutationType = "FactuurOntvangen"
	InvoiceSend            MutationType = "FactuurVerstuurd"
	InvoicePaymentReceived MutationType = "FactuurbetalingOntvangen"
	InvoicePaymentSend     MutationType = "FactuurbetalingVerstuurd"
	MoneyReceived          MutationType = "GeldOntvangen"
	MoneySpend             MutationType = "GeldUitgegeven"
	Memorial               MutationType = "Memoriaal"

	Inclusive TaxType = "IN"
	Exclusive TaxType = "EX"

	VATHighSales           VATCode = "HOOG_VERK"    // VAT High, sales 19%
	VATHighSales21         VATCode = "HOOG_VERK_21" // VAT High, sales 21%
	VATHighReverseCharge21 VATCode = "VERL_VERK"    // VAT High, reverse charge 21%
	VATLowSales            VATCode = "LAAG_VERK"    // VAT Low, sales, on transaction before 1-1-19 6%, after 9%
	VATLowSales9           VATCode = "LAAG_VERK_9"  // VAT Low, sales 9%
	VATLowReverseCharge9   VATCode = "LAAG_VERK_L9" // VAT Low, reverse charge 9%
	VATVariable            VATCode = "AFW"          // Variable VAT Sales
	VATOutsideEU0          VATCode = "BU_EU_VERK"   // Sales outside of EU 0%
	VatInsideEU0           VATCode = "BI_EU_VERK"   // Sales inside of EU 0%
	VATNo                  VATCode = "GEEN"         // No VAT
)

type Mutation struct {
	BoekhoudenCustomerID string
	Number               int
	Type                 MutationType
	Date                 time.Time
	LedgerAccountCode    string
	InvoiceNumber        string
	InvoiceURL           string
	Description          string
	PaymentTerm          string
	PaymentFeature       string
	InExBtw              TaxType
	Items                []MutationItem
}

type MutationItem struct {
	Amount            float64
	AmountExVAT       float64
	AmountVAT         float64
	AmountInVAT       float64
	VATCode           VATCode
	VATPercentage     float64
	LedgerAccountCode string
	KostenplaatsID    int64
}
