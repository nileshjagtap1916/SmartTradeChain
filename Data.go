package main

import "time"

type contract struct {
	ContractId                                  string          `json:"contractId"`
	SellerDetails                               sellerDetails   `json:"sellerDetails"`
	BuyerDetails                                buyerDetails    `json:"buyerDetails"`
	TradeDetails                                []product       `json:"tradeDetails"`
	TradeConditions                             tradeConditions `json:"tradeConditions"`
	DeliveryDetails                             deliveryDetails `json:"deliveryDetails"`
	ContractCreateDate                          time.Time       `json:"createDate"`
	IsLCAttached                                bool            `json:"isLCAttached"`
	IsPOAttached                                bool            `json:"isPOAttached"`
	IsInvoiceListAttached                       bool            `json:"isInvoiceListAttached"`
	IsBillOfLedingAttached                      bool            `json:"isBillOfLedingAttached"`
	ActionPendingOn                             string          `json:"actionPendingOn"`
	ContractStatus                              string          `json:"contractStatus"`
	LastUpdatedDate                             string          `json:"LastUpdatedDate"`
	ApprovedContractByBuyerDate                 string          `json:"ApprovedContractByBuyerDate"`
	LCCreatedByBuyerBankDate                    string          `json:"LCCreatedByBuyerBankDate"`
	LCApprovedBySellerBankDate                  string          `json:"LCApprovedBySellerBankDate"`
	ReadyForShipmentBySellerDate                string          `json:"ReadyForShipmentBySellerDate"`
	ShipmentInProgressByTransDate               string          `json:"ShipmentInProgressByTransDate"`
	ShipmentDeliveredByBuyerDate                string          `json:"ShipmentDeliveredByBuyerDate"`
	InvoiceCreatedBySellerDate                  string          `json:"InvoiceCreatedBySellerDate"`
	PaymentCompletedToSellerBySellerBankDate    string          `json:"PaymentCompletedToSellerBySellerBankDate"`
	PaymentCompletedToSellerBankByBuyerBankDate string          `json:"PaymentCompletedToSellerBankByBuyerBankDate"`
	ContractCompletedByBuyerDate                string          `json:"ContractCompletedByBuyerDate"`
}

type tradeConditions struct {
	PaymentDuration   string `"json:"paymentDuration"`
	TransportDuration string `"json:"transportDuration"`
	Currency          string `"json:"currency"`
	PaymentTerms      string `"json:"paymentTerms"`
}

type product struct {
	ProductName     string `json:"productName"`
	ProductDesc     string `json:"productDesc"`
	ProductPrice    string `json:"productPrice"`
	ProductQuantity string `json:"productQuantity"`
	TotalAmount     string `json:"totalAmount"`
}
type sellerDetails struct {
	Seller     user `json:"seller"`
	SellerBank user `json:"sellerBank"`
}
type buyerDetails struct {
	Buyer     user `json:"buyer"`
	BuyerBank user `json:"buyerBank"`
}

type deliveryDetails struct {
	PickupAddress      string `json:"pickupAddress"`
	DeliveryAddress    string `json:"deliveryAddress"`
	DeliveryDate       string `json:"deliveryDate"`
	TransporterDetails user   `json:"transporterDetails"`
}

type user struct {
	UserId    string `json:"userId"`
	UserName  string `json:"userName"`
	ContactNo string `json:"contactNo"`
	Address   string `json:"address"`
}

type staticData struct {
	TotalContracts        int            `json:"totalContracts"`
	CurrentMonthContracts int            `josn:"currentMonthContracts"`
	LastMonthContracts    int            `json:"lastMonthContracts"`
	NotificationCount     int            `json:"notificationCount"`
	CountStatus           countStatus    `json:"countStatus"`
	ProgressStatus        progressStatus `json:"progressStatus"`
	PaymentStatus         paymentStatus  `json:"paymentStatus"`
	ShipmentStatus        shipmentStatus `json:"shipmentStatus"`
	DeliveryStatus        deliveryStatus `json:"deliveryStatus"`
	ContractList          []contract     `json:"contractList"`
}

type countStatus struct {
	ContractCount  int `json:"contractCount"`
	LCCount        int `json:"lcCount"`
	ShipmentCount  int `json:"shipmentCount"`
	PaymentCount   int `json:"paymentCount"`
	CompletedCount int `json:"completedCount"`
}

type progressStatus struct {
	Ontime    int `json:"ontime"`
	Delayed   int `json:"delayed"`
	Completed int `json:"completed"`
}

type paymentStatus struct {
	PendingBuyerBank  int `json:"pendingBuyerBank"`
	PendingSellerBank int `json:"pendingSellerBank"`
	CompletedBuyer    int `json:"completedBuyer"`
	PendingBuyer      int `json:"pendingBuyer"`
}

type shipmentStatus struct {
	Pending    int `json:"pending"`
	InProgress int `json:"inProgress"`
	Delivered  int `json:"delivered"`
}

type deliveryStatus struct {
	NeedToStarted  int `json:"needToStarted"`
	OnTimeDelivery int `json:"onTimeDelivery"`
	Delayed        int `json:"delayed"`
}
