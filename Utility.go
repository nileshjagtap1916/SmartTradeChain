package main

var Contract_Created = "Contract Created"
var Contract_Accepted = "Contract Accepted"
var LC_Created = "LC Created"
var LC_Approved = "LC Approved"
var Ready_For_Shipment = "Ready For Shipment"
var Shipment_Inprogress = "Shipment Inprogress"
var Shipment_Delivered = "Shipment Delivered"
var Invoice_Created = "Invoice Created"
var Payment_Completed_to_Seller_Bank = "Payment Completed to Seller Bank"
var Payment_Completed_to_Seller = "Payment Completed to Seller"
var Contract_Completed = "Contract Completed"

func mapping_status(contract_status string) string {
	category := map[string]string{
		"Contract Created":                 "Contract",
		"Contract Accepted":                "Contract",
		"LC Created":                       "LC",
		"LC Approved":                      "LC",
		"Ready For Shipment":               "Shipment",
		"Shipment Inprogress":              "Shipment",
		"Shipment Delivered":               "Shipment",
		"Invoice Created":                  "Payment",
		"Payment Completed to Seller":      "Payment",
		"Payment Completed to Seller Bank": "Payment",
		"Contract Completed":               "Completed",
	}
	category_status := category[contract_status]
	return category_status
}
