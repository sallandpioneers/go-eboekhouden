package model

type (
	CustomerType string
)

const (
	Business CustomerType = "b"
	Private  CustomerType = "p"
)

type Customer struct {
	ID         int64
	Code       string
	Business   string
	Addresses  CustomerAddresses
	Phone      string
	Email      string
	Website    string
	Notition   string
	VAT        string
	COC        string
	Salutation string
	IBAN       string
	BIC        string
	Type       CustomerType
}

type CustomerAddresses struct {
	Business Address
	Mailing  Address
}

type Address struct {
	Address string
	ZipCode string
	City    string
	Country string
}
