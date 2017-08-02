package main

import "time"

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

//Payment Condotions
var Max_Days_PaymentDuration = 30
var Min_Days_PaymentDuration = 15
var Min_Days_TransportDuration = 10
var Max_Days_TransportDuration = 20
var Min_Days_DeliveryDuration = 15
var Max_Days_DeliveryDuration = 30

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

func DiffDays(year2, month2, day2, year1, month1, day1 int) int {
	if year2 < year1 {
		return -DiffDays(year1, month1, day1, year2, month2, day2)
	}
	d2 := time.Date(year2, time.Month(month2), day2, 0, 0, 0, 0, time.UTC)
	d1 := time.Date(year1, time.Month(month1), day1, 0, 0, 0, 0, time.UTC)
	diff := d2.YearDay() - d1.YearDay()

	for y := year1; y < year2; y++ {
		diff += time.Date(y, time.December, 31, 0, 0, 0, 0, time.UTC).YearDay()
	}
	/* if debug && !d1.AddDate(0, 0, diff).Equal(d2) {
	    panic("invalid diff")
	} */
	return diff
}
