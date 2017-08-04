package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//var numberofContracts int = 5

const nc = 5

const contracts string = "Contract"
const lc string = "LC"
const shipment string = "Shipment"
const payment string = "Payment"
const completed string = "Completed"

const dateFormat string = "2006-01-02"

type Sorted []contract

func (slice Sorted) Len() int {
	return len(slice)
}

func (slice Sorted) Less(i, j int) bool {
	return slice[i].ContractCreateDate.After(slice[j].ContractCreateDate)
}

func (slice Sorted) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func initializeChaincode(stub shim.ChaincodeStubInterface, args []string) error {
	var ok bool
	var err error
	ok, err = createDatabase(stub, args)
	if !ok {
		return err
	}
	return nil
}

func initializeUser(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 arguments")
	}
	userId := args[0]
	ok, err := insertUserBlankRecord(stub, userId)
	if !ok {
		return nil, err
	}

	return nil, nil
}

func saveContractDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var contractDetails contract
	var err error
	var ok bool

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 arguments")
	}

	json.Unmarshal([]byte(args[0]), &contractDetails)

	/* Commented becouse v0.6 does not support it
		//Datatype checking
		ok, err = dataTypeCheck(contractDetails)
		if !ok {
			return nil, err
		}

		//Mandatory Field checking
		ok, err = mandatoryFieldCheck(contractDetails)
		if !ok {
			return nil, err
		}

		//Delivary Date Duration checking
		DeliveryDate, _ := time.Parse(time.RFC3339, contractDetails.DeliveryDetails.DeliveryDate)
		CurrentDate := time.Now()
		Days := DiffDays(int(DeliveryDate.Year()), int(DeliveryDate.Month()), int(DeliveryDate.Day()), int(CurrentDate.Year()), int(CurrentDate.Month()), int(CurrentDate.Day()))

		if Days < Min_Days_DeliveryDuration {
			return nil, errors.New("Delivery Duration must be greater than " + string(Min_Days_DeliveryDuration) + " days")
		} else if Days > Max_Days_DeliveryDuration {
			return nil, errors.New("Payment Duration must be less than " + string(Max_Days_DeliveryDuration) + " days")
		}
	  //duration := time.Since(deliveryDate)
		//if int(duration.Hours()) < (Min_Days_DeliveryDuration * 24) {
			//return nil, errors.New("Delivery Duration must be greater than " + string(Min_Days_DeliveryDuration) + " days")
		//} else if int(duration.Hours()) > (Max_Days_DeliveryDuration * 24) {
			//return nil, errors.New("Payment Duration must be less than " + string(Max_Days_DeliveryDuration) + " days")
		//}

		// Payment duartion checking
		PaymentDuration, _ := strconv.Atoi(contractDetails.TradeConditions.PaymentDuration)
		if PaymentDuration < Min_Days_PaymentDuration {
			return nil, errors.New("Payment Duration must be gereater than " + string(Min_Days_PaymentDuration) + " days")
		} else if PaymentDuration > Max_Days_PaymentDuration {
			return nil, errors.New("Payment Duration must be less than " + string(Max_Days_PaymentDuration) + " days")
		}

		// Transport duartion checking
		TransportDuration, _ := strconv.Atoi(contractDetails.TradeConditions.TransportDuration)
		if TransportDuration < Min_Days_TransportDuration {
			return nil, errors.New("Payment Duration must be gereater than " + string(Min_Days_TransportDuration) + " days")
		} else if TransportDuration > Max_Days_TransportDuration {
			return nil, errors.New("Payment Duration must be less than " + string(Max_Days_TransportDuration) + " days")
		}
	  comment ending */

	contractDetails = addContractInformation(contractDetails)

	ok, err = insertContractDetails(stub, contractDetails)
	if !ok && err == nil {
		return nil, errors.New("Error in adding OrderDetails record")
	}

	ok, err = updateUsersContractList(stub, contractDetails)
	if !ok && err == nil {
		return nil, errors.New("Error in adding OrderDetails record")
	}

	return nil, nil
}

func getContractDetailsByContractId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	contractId := args[0]
	contractDetails, _ := getContractDetails(stub, contractId)

	jsonAsBytes, _ := json.Marshal(contractDetails)
	return jsonAsBytes, nil

}

func saveAttachmentDetails(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	var ok bool

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
	}

	contractId := args[0]
	attachmentName := args[1]
	documentBlob := args[2]

	ok, err = insertAttachmentDetails(stub, contractId, attachmentName, documentBlob)
	if !ok && err == nil {
		return nil, errors.New("Error in inserting attachment")
	}

	return nil, err
}

func getAttachment(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Need 2 arguments")
	}

	contractId := args[0]
	attachmentName := args[1]

	jsonAsBytes, err := getAttachmentDetails(stub, contractId, attachmentName)
	if err != nil {
		return nil, errors.New("Error in downloading the attachment")
	}

	return jsonAsBytes, nil
}

func addContractInformation(contractDetails contract) contract {
	//contractDetails.ContractId = strconv.Itoa(rand.Intn(1000000000) + 1)

	contractDetails.ContractId = time.Now().Format("0102150405")
	contractDetails.ContractCreateDate = time.Now().Local()
	contractDetails.LastUpdatedDate = contractDetails.ContractCreateDate.Format(dateFormat)
	contractDetails.IsLCAttached = false
	contractDetails.IsPOAttached = true
	contractDetails.IsBillOfLedingAttached = false
	contractDetails.IsInvoiceListAttached = false
	contractDetails.ActionPendingOn = "buyer"
	contractDetails.ContractStatus = "Contract Created"

	//calculate TotalTradeAmount
	var TotalTradeAmount float64
	TotalTradeAmount = 0
	for _, element := range contractDetails.TradeDetails {
		Amount, _ := strconv.ParseFloat(element.TotalAmount, 64)
		TotalTradeAmount = TotalTradeAmount + Amount
	}
	contractDetails.TotalTradeAmount = TotalTradeAmount

	return contractDetails
}

func updateUsersContractList(stub shim.ChaincodeStubInterface, contractDetails contract) (bool, error) {
	var ok bool
	var userContractList []string

	//Update Seller's Contract List
	userContractList, ok = getUserContractList(stub, contractDetails.SellerDetails.Seller.UserId)
	if !ok {
		return ok, errors.New("Error in geting Seller's contract list")
	}
	userContractList = append(userContractList, contractDetails.ContractId)
	ok = updateUserContractList(stub, contractDetails.SellerDetails.Seller.UserId, userContractList)
	if !ok {
		return ok, errors.New("Error in updating Seller's contract list")
	}

	//Update SellerBank's Contract List
	userContractList, ok = getUserContractList(stub, contractDetails.SellerDetails.SellerBank.UserId)
	if !ok {
		return ok, errors.New("Error in geting SellerBank's contract list")
	}
	userContractList = append(userContractList, contractDetails.ContractId)
	ok = updateUserContractList(stub, contractDetails.SellerDetails.SellerBank.UserId, userContractList)
	if !ok {
		return ok, errors.New("Error in updating SellerBank's contract list")
	}

	//Update Buyer's Contract List
	userContractList, ok = getUserContractList(stub, contractDetails.BuyerDetails.Buyer.UserId)
	if !ok {
		return ok, errors.New("Error in geting Buyer's contract list")
	}
	userContractList = append(userContractList, contractDetails.ContractId)
	ok = updateUserContractList(stub, contractDetails.BuyerDetails.Buyer.UserId, userContractList)
	if !ok {
		return ok, errors.New("Error in updating Buyer's contract list")
	}

	//Update BuyerBank's Contract List
	userContractList, ok = getUserContractList(stub, contractDetails.BuyerDetails.BuyerBank.UserId)
	if !ok {
		return ok, errors.New("Error in geting BuyerBank's contract list")
	}
	userContractList = append(userContractList, contractDetails.ContractId)
	ok = updateUserContractList(stub, contractDetails.BuyerDetails.BuyerBank.UserId, userContractList)
	if !ok {
		return ok, errors.New("Error in updating BuyerBank's contract list")
	}

	//Update Transporter's Contract List
	userContractList, ok = getUserContractList(stub, contractDetails.DeliveryDetails.TransporterDetails.UserId)
	if !ok {
		return ok, errors.New("Error in geting Transporter's contract list")
	}
	userContractList = append(userContractList, contractDetails.ContractId)
	ok = updateUserContractList(stub, contractDetails.DeliveryDetails.TransporterDetails.UserId, userContractList)
	if !ok {
		return ok, errors.New("Error in updating Transporter's contract list")
	}

	return true, nil
}

func getContractDetailsByUserId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var contractDetails []contract
	var contract contract
	var sortedDetails Sorted

	if len(args) > 1 && len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Need 1 or 3 argument")
	}

	if len(args) == 1 {
		userId := args[0]

		contractIdList, ok := getUserContractList(stub, userId)
		if !ok {
			return nil, errors.New("Error in geting user specific contract list")
		}

		for _, element := range contractIdList {
			contractId := element
			contract, _ = getContractDetails(stub, contractId)
			contractDetails = append(contractDetails, contract)
		}

		sortedDetails = contractDetails
		sort.Sort(sortedDetails)
		contractAsBytes, _ := json.Marshal(sortedDetails)
		return contractAsBytes, nil
	}

	if len(args) == 3 {

		userId := args[0]
		chartName := args[1]
		chartStatus := args[2]

		contractIdList, ok := getUserContractList(stub, userId)
		if !ok {
			return nil, errors.New("Error in geting user specific contract list")
		}

		if chartName == "CountStatus" {
			// Count status calculation
			if chartStatus == "Contract" {

				for _, element := range contractIdList {
					contractId := element
					contract, _ = getContractDetails(stub, contractId)
					if contract.ContractStatus == Contract_Created || contract.ContractStatus == Contract_Accepted {
						contractDetails = append(contractDetails, contract)
					}
				}
				sortedDetails = contractDetails
				sort.Sort(sortedDetails)
				contractAsBytes, _ := json.Marshal(sortedDetails)
				return contractAsBytes, nil
			}

			if chartStatus == "LC" {

				for _, element := range contractIdList {
					contractId := element
					contract, _ = getContractDetails(stub, contractId)
					if contract.ContractStatus == LC_Created || contract.ContractStatus == LC_Approved {
						contractDetails = append(contractDetails, contract)
					}
				}
				sortedDetails = contractDetails
				sort.Sort(sortedDetails)
				contractAsBytes, _ := json.Marshal(sortedDetails)
				return contractAsBytes, nil
			}

			if chartStatus == "Shipment" {

				for _, element := range contractIdList {
					contractId := element
					contract, _ = getContractDetails(stub, contractId)
					if contract.ContractStatus == Ready_For_Shipment || contract.ContractStatus == Shipment_Inprogress || contract.ContractStatus == Shipment_Delivered {
						contractDetails = append(contractDetails, contract)
					}
				}
				sortedDetails = contractDetails
				sort.Sort(sortedDetails)
				contractAsBytes, _ := json.Marshal(sortedDetails)
				return contractAsBytes, nil
			}

			if chartStatus == "Payment" {

				for _, element := range contractIdList {
					contractId := element
					contract, _ = getContractDetails(stub, contractId)
					if contract.ContractStatus == Invoice_Created || contract.ContractStatus == Payment_Completed_to_Seller_Bank || contract.ContractStatus == Payment_Completed_to_Seller {
						contractDetails = append(contractDetails, contract)
					}
				}
				sortedDetails = contractDetails
				sort.Sort(sortedDetails)
				contractAsBytes, _ := json.Marshal(sortedDetails)
				return contractAsBytes, nil
			}

			if chartStatus == "Completed" {

				for _, element := range contractIdList {
					contractId := element
					contract, _ = getContractDetails(stub, contractId)
					if contract.ContractStatus == Contract_Completed {
						contractDetails = append(contractDetails, contract)
					}
				}
				sortedDetails = contractDetails
				sort.Sort(sortedDetails)
				contractAsBytes, _ := json.Marshal(sortedDetails)
				return contractAsBytes, nil
			}

		}

		/*
			if chartName == "ProgressStatus" {

				CurrentDate := time.Now().Local()

				if chartStatus == "Ontime" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)

						paymentDuration, _ := strconv.Atoi(contract.TradeConditions.PaymentDuration)
						expectedDeliveryDate := contract.ContractCreateDate.AddDate(0, 0, paymentDuration)

						if contract.ContractStatus != Contract_Completed {

							if inTimeSpan(contract.ContractCreateDate, expectedDeliveryDate, CurrentDate) ||
								CurrentDate.Equal(expectedDeliveryDate) ||
								CurrentDate.Equal(contract.ContractCreateDate) {

								contractDetails = append(contractDetails, contract)
							}
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}
				if chartStatus == "Delayed" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)

						paymentDuration, _ := strconv.Atoi(contract.TradeConditions.PaymentDuration)
						expectedDeliveryDate := contract.ContractCreateDate.AddDate(0, 0, paymentDuration)

						if CurrentDate.After(expectedDeliveryDate) {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

				if chartStatus == "Completed" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Contract_Completed {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

			}

			if chartName == "PaymentStatus" {
				if chartStatus == "PendingSellerBank" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Invoice_Created {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

				if chartStatus == "PendingBuyerBank" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Payment_Completed_to_Seller {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

				if chartStatus == "PendingBuyer" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Payment_Completed_to_Seller_Bank {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

				if chartStatus == "CompletedBuyer" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Contract_Completed {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

			}
			if chartName == "ShipmentStatus" {
				if chartStatus == "Pending" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Ready_For_Shipment {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}
				if chartStatus == "InProgress" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Shipment_Inprogress {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}
				if chartStatus == "Completed" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Shipment_Delivered {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

			}
			if chartName == "DeliveryStatus" {

				if chartStatus == "NeedToStart" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Ready_For_Shipment {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

				if chartStatus == "OnTimeDelivery" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						if contract.ContractStatus == Shipment_Inprogress {
							contractDetails = append(contractDetails, contract)
						}
					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

				if chartStatus == "Delayed" {

					for _, element := range contractIdList {
						contractId := element
						contract, _ = getContractDetails(stub, contractId)
						deliveryDate, _ := time.Parse(time.RFC3339, contract.DeliveryDetails.DeliveryDate)

						if time.Now().Local().After(deliveryDate) == true {
							contractDetails = append(contractDetails, contract)
						}

					}
					sortedDetails = contractDetails
					sort.Sort(sortedDetails)
					contractAsBytes, _ := json.Marshal(sortedDetails)
					return contractAsBytes, nil
				}

			}*/

	}

	return nil, nil

}

func getCountStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Need 1 argument")
	}
	userId := args[0]

	var countStatus countStatus

	var contractCount int
	var lcCount int
	var shipmentCount int
	var paymnetCount int
	var completedCount int

	var contractVar contract

	contractIdList := []string{}
	//contractDetails := []contract{}

	contractIdList, ok := getUserContractList(stub, userId)

	if !ok {
		return nil, errors.New("Error in geting user specific contract list")
	}

	for _, element := range contractIdList {
		contractId := element
		contractVar, _ = getContractDetails(stub, contractId)

		status := mapping_status(contractVar.ContractStatus)

		// Counts Check

		if status == contracts {
			contractCount++
		} else if status == lc {
			lcCount++
		} else if status == shipment {
			shipmentCount++
		} else if status == payment {
			paymnetCount++
		} else if status == completed {
			completedCount++
		}
	}

	countStatus.ContractCount = contractCount
	countStatus.LCCount = lcCount
	countStatus.PaymentCount = paymnetCount
	countStatus.ShipmentCount = shipmentCount
	countStatus.CompletedCount = completedCount

	countStatusAsBytes, _ := json.Marshal(countStatus)
	return countStatusAsBytes, nil

}

func getStaticDetailsByUserId(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var staticDetails staticData

	var latestContracts []contract
	var sortedDetails Sorted

	var contractDetails []contract
	var contractVar contract

	var notificationCount int
	//var totalContracts int
	//var thisMonth int
	//var lastMonth int

	var contractCount int
	var lcCount int
	var shipmentCount int
	var paymnetCount int
	var completedCount int

	var ontimeOrder int
	var delayedOrder int

	var pending int
	var inprogress int
	var delivered int

	var needtostart int
	var ontimedelivery int
	var delayed int

	var pendingfrombuyer int
	var pendingfromsellerbank int
	var pendingfrombuyerbank int
	var completedbuyer int

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Need 2 argument")
	}
	userId := args[0]
	userRole := args[1]

	contractIdList := []string{}
	//contractDetails := []contract{}

	contractIdList, ok := getUserContractList(stub, userId)

	if !ok {
		return nil, errors.New("Error in geting user specific contract list")
	}

	for _, element := range contractIdList {
		contractId := element
		contractVar, _ = getContractDetails(stub, contractId)

		contractDetails = append(contractDetails, contractVar)

		CurrentDate := time.Now().Local()

		if CurrentDate.Month() == contractVar.ContractCreateDate.Month() && CurrentDate.Year() == contractVar.ContractCreateDate.Year() {
			staticDetails.CurrentMonthContracts++
		}

		lastMonthDate := CurrentDate.AddDate(0, -1, 0)
		if lastMonthDate.Month() == contractVar.ContractCreateDate.Month() && lastMonthDate.Year() == contractVar.ContractCreateDate.Year() {
			staticDetails.LastMonthContracts++
		}

		status := mapping_status(contractVar.ContractStatus)
		fmt.Print("Staus is", status)

		// NotificationCount Check

		if contractVar.ActionPendingOn == userRole {
			notificationCount++
		}

		// Counts Check

		if status == contracts {
			contractCount++
		} else if status == lc {
			lcCount++
		} else if status == shipment {
			shipmentCount++
		} else if status == payment {
			paymnetCount++
		} else if status == completed {
			completedCount++
		}

		// Progress Check

		paymentDuration, _ := strconv.Atoi(contractVar.TradeConditions.PaymentDuration)
		expectedDeliveryDate := contractVar.ContractCreateDate.AddDate(0, 0, paymentDuration)
		fmt.Println("Delivery Date", expectedDeliveryDate)
		if contractVar.ContractStatus != Contract_Completed {

			if inTimeSpan(contractVar.ContractCreateDate, expectedDeliveryDate, CurrentDate) ||
				CurrentDate.Equal(expectedDeliveryDate) ||
				CurrentDate.Equal(contractVar.ContractCreateDate) {

				ontimeOrder++
			}

		}

		if CurrentDate.After(expectedDeliveryDate) {
			delayedOrder++
		}

		// Payment Staus Check
		if contractVar.ContractStatus == Invoice_Created {
			pendingfromsellerbank++
		}
		if contractVar.ContractStatus == Payment_Completed_to_Seller {
			pendingfrombuyerbank++
		}
		if contractVar.ContractStatus == Payment_Completed_to_Seller_Bank {
			pendingfrombuyer++
		}
		if contractVar.ContractStatus == Contract_Completed {
			completedbuyer++
		}

		// Shipment, Delivery Status Check
		if contractVar.ContractStatus == Ready_For_Shipment {
			pending++
		} else if contractVar.ContractStatus == Shipment_Inprogress {
			inprogress++
		} else if contractVar.ContractStatus == Shipment_Delivered {
			delivered++
		}

		deliveryDate, _ := time.Parse(time.RFC3339, contractVar.DeliveryDetails.DeliveryDate)

		if time.Now().Local().After(deliveryDate) == true {
			delayed++
		}

	}

	staticDetails.TotalContracts = len(contractIdList)

	if staticDetails.TotalContracts == 0 {
		//staticDetails.ContractList = latestContracts
		//staticDataAsBytes, _ := json.Marshal(staticDetails)

		staticDataAsBytes, _ := json.Marshal(staticDetails)
		return staticDataAsBytes, nil
		//return nil
	}

	needtostart = pending
	ontimedelivery = inprogress

	//staticDetails.TotalContracts = totalContracts
	//staticDetails.CurrentMonthContracts = thisMonth
	//staticDetails.LastMonthContracts = lastMonth

	staticDetails.NotificationCount = notificationCount
	staticDetails.CountStatus.ContractCount = contractCount
	staticDetails.CountStatus.LCCount = lcCount
	staticDetails.CountStatus.PaymentCount = paymnetCount
	staticDetails.CountStatus.ShipmentCount = shipmentCount
	staticDetails.CountStatus.CompletedCount = completedCount

	staticDetails.ProgressStatus.Ontime = ontimeOrder
	staticDetails.ProgressStatus.Delayed = delayedOrder
	staticDetails.ProgressStatus.Completed = completedCount

	staticDetails.ShipmentStatus.Pending = pending
	staticDetails.ShipmentStatus.InProgress = inprogress
	staticDetails.ShipmentStatus.Delivered = delivered

	staticDetails.PaymentStatus.PendingSellerBank = pendingfromsellerbank
	staticDetails.PaymentStatus.PendingBuyerBank = pendingfrombuyerbank
	staticDetails.PaymentStatus.CompletedBuyer = completedbuyer
	staticDetails.PaymentStatus.PendingBuyer = pendingfrombuyer

	staticDetails.DeliveryStatus.NeedToStarted = needtostart
	staticDetails.DeliveryStatus.OnTimeDelivery = ontimedelivery
	staticDetails.DeliveryStatus.Delayed = delayed

	staticDetails.ContractList = contractDetails

	sortedDetails = contractDetails
	sort.Sort(sortedDetails)

	if sortedDetails.Len() < 5 {
		numberofContracts := sortedDetails.Len()
		for j := 0; j < numberofContracts; j++ {
			latestContracts = append(latestContracts, sortedDetails[j])
		}
	} else {
		for j := 0; j < nc; j++ {
			latestContracts = append(latestContracts, sortedDetails[j])
		}
	}

	staticDetails.ContractList = latestContracts

	staticDataAsBytes, _ := json.Marshal(staticDetails)
	return staticDataAsBytes, nil

}

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func getNotificationStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var contractDetails []contract
	var contract contract
	var sortedDetails Sorted

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Need 2 argument")
	}
	userId := args[0]
	userRole := args[1]

	contractIdList, ok := getUserContractList(stub, userId)
	if !ok {
		return nil, errors.New("Error in geting user specific contract list")
	}

	for _, element := range contractIdList {
		contractId := element
		contract, _ = getContractDetails(stub, contractId)

		if contract.ActionPendingOn == userRole {
			contractDetails = append(contractDetails, contract)
		}

	}
	sortedDetails = contractDetails
	sort.Sort(sortedDetails)
	contractAsBytes, _ := json.Marshal(sortedDetails)
	return contractAsBytes, nil

}

func getNotificationCountStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var contract contract

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Need 2 argument")
	}

	userId := args[0]
	userRole := args[1]

	var countStatus countStatus

	var contractCount int
	var lcCount int
	var shipmentCount int
	var paymnetCount int
	var completedCount int

	contractIdList := []string{}
	//contractDetails := []contract{}

	contractIdList, ok := getUserContractList(stub, userId)

	if !ok {
		return nil, errors.New("Error in geting user specific contract list")
	}

	for _, element := range contractIdList {
		contractId := element
		contract, _ = getContractDetails(stub, contractId)
		if contract.ActionPendingOn == userRole {

			status := mapping_status(contract.ContractStatus)

			// Counts Check

			if status == contracts {
				contractCount++
			} else if status == lc {
				lcCount++
			} else if status == shipment {
				shipmentCount++
			} else if status == payment {
				paymnetCount++
			} else if status == completed {
				completedCount++
			}
		}
	}

	countStatus.ContractCount = contractCount
	countStatus.LCCount = lcCount
	countStatus.PaymentCount = paymnetCount
	countStatus.ShipmentCount = shipmentCount
	countStatus.CompletedCount = completedCount

	countStatusAsBytes, _ := json.Marshal(countStatus)
	return countStatusAsBytes, nil

}

func UpdateContractStatus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var ok bool
	var err error
	//var status statusMaintained
	//var contractLists contract

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Need 2 arguments")
	}

	userID := args[0]
	contractID := args[1]
	current_time := time.Now().Local()
	contractList, _ := getContractDetails(stub, contractID)

	contractStatus := contractList.ContractStatus
	//for seller
	if contractList.SellerDetails.Seller.UserId == userID {
		if contractStatus == LC_Approved {
			contractList.ContractStatus = Ready_For_Shipment
			contractList.ActionPendingOn = "transporter"
			contractList.ReadyForShipmentBySellerDate = current_time.Format("2006-01-02")
		} else if contractStatus == Shipment_Delivered {
			contractList.ContractStatus = Invoice_Created
			contractList.ActionPendingOn = "sellerbank"
			contractList.InvoiceCreatedBySellerDate = current_time.Format("2006-01-02")
		}
	}

	//for buyer
	if contractList.BuyerDetails.Buyer.UserId == userID {
		if contractStatus == Contract_Created {
			contractList.ContractStatus = Contract_Accepted
			contractList.ActionPendingOn = "buyerbank"
			contractList.ApprovedContractByBuyerDate = current_time.Format("2006-01-02")
		} else if contractStatus == Payment_Completed_to_Seller_Bank {
			contractList.ContractStatus = Contract_Completed
			contractList.ActionPendingOn = Contract_Completed
			contractList.ContractCompletedByBuyerDate = current_time.Format("2006-01-02")
		} else if contractStatus == Shipment_Inprogress {
			contractList.ContractStatus = Shipment_Delivered
			contractList.ActionPendingOn = "seller"
			contractList.ShipmentDeliveredByBuyerDate = current_time.Format("2006-01-02")
		}
	}

	//for sellerBank
	if contractList.SellerDetails.SellerBank.UserId == userID {
		if contractStatus == LC_Created {
			contractList.ContractStatus = LC_Approved
			contractList.ActionPendingOn = "seller"
			contractList.LCApprovedBySellerBankDate = current_time.Format("2006-01-02")
		} else if contractStatus == Invoice_Created {
			contractList.ContractStatus = Payment_Completed_to_Seller
			contractList.ActionPendingOn = "buyerbank"
			contractList.PaymentCompletedToSellerBySellerBankDate = current_time.Format("2006-01-02")
		}
	}

	//for buyerBank
	if contractList.BuyerDetails.BuyerBank.UserId == userID {
		if contractStatus == Contract_Accepted {
			contractList.ContractStatus = LC_Created
			contractList.ActionPendingOn = "sellerbank"
			contractList.LCCreatedByBuyerBankDate = current_time.Format("2006-01-02")
		} else if contractStatus == Payment_Completed_to_Seller {
			contractList.ContractStatus = Payment_Completed_to_Seller_Bank
			contractList.ActionPendingOn = "buyer"
			contractList.PaymentCompletedToSellerBankByBuyerBankDate = current_time.Format("2006-01-02")
		}
	}

	//for transporter
	if contractList.DeliveryDetails.TransporterDetails.UserId == userID {
		if contractStatus == Ready_For_Shipment {
			contractList.ContractStatus = Shipment_Inprogress
			contractList.ActionPendingOn = "buyer"
			contractList.ShipmentInProgressByTransDate = current_time.Format("2006-01-02")
		}
	}

	contractList.LastUpdatedDate = current_time.Format("2006-01-02")
	//status = setStructStatus(stub, status, userID, contractStatus)

	//Discount logic - Ready For Shipment
	DeliveryDate, _ := time.Parse(time.RFC3339, contractList.DeliveryDetails.DeliveryDate)
	CurrentDate := time.Now().Local()
	if contractList.SellerDetails.Seller.UserId == userID {
		if contractStatus == LC_Approved {
			if CurrentDate.After(DeliveryDate) {
				Days := DiffDays(int(CurrentDate.Year()), int(CurrentDate.Month()), int(CurrentDate.Day()), int(DeliveryDate.Year()), int(DeliveryDate.Month()), int(DeliveryDate.Day()))
				if (Days > 0) && (Days <= 5) {
					contractList.DiscountPercentage = 5
					contractList.DiscountedAmount = contractList.TotalTradeAmount - (contractList.TotalTradeAmount * 0.5)
					//return []byte("Disscount 5%"), nil //errors.New("Disscount 5%")
				} else if (Days >= 6) && (Days <= 15) {
					contractList.DiscountPercentage = 10
					contractList.DiscountedAmount = contractList.TotalTradeAmount - (contractList.TotalTradeAmount * 0.10)
					//return []byte("Disscount 10%"), nil //errors.New("Disscount 10%")
				} else if Days >= 16 {
					contractList.DiscountPercentage = 20
					contractList.DiscountedAmount = contractList.TotalTradeAmount - (contractList.TotalTradeAmount * 0.20)
					//return []byte("Disscount 20%"), nil //errors.New("Disscount 20%")
				}

			}
		}
	}

	ok = updateContractListByContractID(stub, contractID, contractList)
	if !ok {
		return nil, errors.New("Error in updating contract list")
	}

	return nil, err
}

/* Commented becouse v0.6 does not support it
func mandatoryFieldCheck(contractDetails contract) (bool, error) {

	if contractDetails.SellerDetails.Seller.UserId == "" {
		return false, errors.New("Seller UserId field is mandatory")
	} else if contractDetails.SellerDetails.Seller.UserName == "" {
		return false, errors.New("Seller UserName field is mandatory")
	} else if contractDetails.SellerDetails.Seller.ContactNo == "" {
		return false, errors.New("Seller ContactNo field is mandatory")
	} else if contractDetails.SellerDetails.Seller.Address == "" {
		return false, errors.New("Seller Address field is mandatory")
	} else if contractDetails.SellerDetails.SellerBank.UserId == "" {
		return false, errors.New("SellerBank UserId field is mandatory")
	} else if contractDetails.SellerDetails.SellerBank.UserName == "" {
		return false, errors.New("SellerBank UserName field is mandatory")
	} else if contractDetails.SellerDetails.SellerBank.ContactNo == "" {
		return false, errors.New("SellerBank ContactNo field is mandatory")
	} else if contractDetails.SellerDetails.SellerBank.Address == "" {
		return false, errors.New("SellerBank Address field is mandatory")
	} else if contractDetails.BuyerDetails.Buyer.UserId == "" {
		return false, errors.New("Buyer UserId field is mandatory")
	} else if contractDetails.BuyerDetails.Buyer.UserName == "" {
		return false, errors.New("Buyer UserName field is mandatory")
	} else if contractDetails.BuyerDetails.Buyer.ContactNo == "" {
		return false, errors.New("Buyer ContactNo field is mandatory")
	} else if contractDetails.BuyerDetails.Buyer.Address == "" {
		return false, errors.New("Buyer Address field is mandatory")
	} else if contractDetails.BuyerDetails.BuyerBank.UserId == "" {
		return false, errors.New("BuyerBank UserId field is mandatory")
	} else if contractDetails.BuyerDetails.BuyerBank.UserName == "" {
		return false, errors.New("BuyerBank UserName field is mandatory")
	} else if contractDetails.BuyerDetails.BuyerBank.ContactNo == "" {
		return false, errors.New("BuyerBank ContactNo field is mandatory")
	} else if contractDetails.BuyerDetails.BuyerBank.Address == "" {
		return false, errors.New("BuyerBank Address field is mandatory")
	} else if contractDetails.TradeConditions.PaymentDuration == "" {
		return false, errors.New("PaymentDuration field is mandatory")
	} else if contractDetails.TradeConditions.TransportDuration == "" {
		return false, errors.New("TransportDuration field is mandatory")
	} else if contractDetails.TradeConditions.Currency == "" {
		return false, errors.New("Currency field is mandatory")
	} else if contractDetails.TradeConditions.PaymentTerms == "" {
		return false, errors.New("PaymentTerms field is mandatory")
	} else if contractDetails.DeliveryDetails.PickupAddress == "" {
		return false, errors.New("PickupAddress field is mandatory")
	} else if contractDetails.DeliveryDetails.DeliveryAddress == "" {
		return false, errors.New("DeliveryAddress field is mandatory")
	} else if contractDetails.DeliveryDetails.DeliveryDate == "" {
		return false, errors.New("DeliveryDate field is mandatory")
	} else if contractDetails.DeliveryDetails.Incoterm == "" {
		return false, errors.New("Incoterm field is mandatory")
	} else if contractDetails.DeliveryDetails.TransporterDetails.UserId == "" {
		return false, errors.New("TransporterDetails UserId field is mandatory")
	} else if contractDetails.DeliveryDetails.TransporterDetails.UserName == "" {
		return false, errors.New("TransporterDetails UserName field is mandatory")
	} else if contractDetails.DeliveryDetails.TransporterDetails.ContactNo == "" {
		return false, errors.New("TransporterDetails ContactNo field is mandatory")
	} else if contractDetails.DeliveryDetails.TransporterDetails.Address == "" {
		return false, errors.New("TransporterDetails Address field is mandatory")
	}

	for _, element := range contractDetails.TradeDetails {
		if element.ProductName == "" {
			return false, errors.New("ProductName field is mandatory")
		} else if element.ProductDesc == "" {
			return false, errors.New("ProductDesc field is mandatory")
		} else if element.ProductPrice == "" {
			return false, errors.New("ProductPrice field is mandatory")
		} else if element.ProductQuantity == "" {
			return false, errors.New("ProductQuantity field is mandatory")
		} else if element.TotalAmount == "" {
			return false, errors.New("TotalAmount field is mandatory")
		}
	}

	return true, nil
}

func dataTypeCheck(contractDetails contract) (bool, error) {

	if reflect.TypeOf(contractDetails.SellerDetails.Seller.UserId).Kind() != reflect.String {
		return false, errors.New("String required for Seller UserId field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.Seller.UserName).Kind() != reflect.String {
		return false, errors.New("String required for Seller UserName field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.Seller.ContactNo).Kind() != reflect.String {
		return false, errors.New("String required for Seller ContactNo field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.Seller.Address).Kind() != reflect.String {
		return false, errors.New("String required for Seller Address field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.SellerBank.UserId).Kind() != reflect.String {
		return false, errors.New("String required for SellerBank UserId field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.SellerBank.UserName).Kind() != reflect.String {
		return false, errors.New("String required for SellerBank UserName field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.SellerBank.ContactNo).Kind() != reflect.String {
		return false, errors.New("String required for SellerBank ContactNo field ")
	} else if reflect.TypeOf(contractDetails.SellerDetails.SellerBank.Address).Kind() != reflect.String {
		return false, errors.New("String required for SellerBank Address field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.Buyer.UserId).Kind() != reflect.String {
		return false, errors.New("String required for Buyer UserId field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.Buyer.UserName).Kind() != reflect.String {
		return false, errors.New("String required for Buyer UserName field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.Buyer.ContactNo).Kind() != reflect.String {
		return false, errors.New("String required for Buyer ContactNo field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.Buyer.Address).Kind() != reflect.String {
		return false, errors.New("String required for Buyer Address field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.BuyerBank.UserId).Kind() != reflect.String {
		return false, errors.New("String required for BuyerBank UserId field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.BuyerBank.UserName).Kind() != reflect.String {
		return false, errors.New("String required for BuyerBank UserName field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.BuyerBank.ContactNo).Kind() != reflect.String {
		return false, errors.New("String required for BuyerBank ContactNo field ")
	} else if reflect.TypeOf(contractDetails.BuyerDetails.BuyerBank.Address).Kind() != reflect.String {
		return false, errors.New("String required for BuyerBank Address field ")
	} else if reflect.TypeOf(contractDetails.TradeConditions.PaymentDuration).Kind() != reflect.String {
		return false, errors.New("String required for PaymentDuration field ")
	} else if reflect.TypeOf(contractDetails.TradeConditions.TransportDuration).Kind() != reflect.String {
		return false, errors.New("String required for TransportDuration field ")
	} else if reflect.TypeOf(contractDetails.TradeConditions.Currency).Kind() != reflect.String {
		return false, errors.New("String required for Currency field ")
	} else if reflect.TypeOf(contractDetails.TradeConditions.PaymentTerms).Kind() != reflect.String {
		return false, errors.New("String required for PaymentTerms field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.PickupAddress).Kind() != reflect.String {
		return false, errors.New("String required for PickupAddress field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.DeliveryAddress).Kind() != reflect.String {
		return false, errors.New("String required for DeliveryAddress field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.DeliveryDate).Kind() != reflect.String {
		return false, errors.New("String required for DeliveryDate field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.Incoterm).Kind() != reflect.String {
		return false, errors.New("String required for Incoterm field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.TransporterDetails.UserId).Kind() != reflect.String {
		return false, errors.New("String required for TransporterDetails UserId field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.TransporterDetails.UserName).Kind() != reflect.String {
		return false, errors.New("String required for TransporterDetails UserName field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.TransporterDetails.ContactNo).Kind() != reflect.String {
		return false, errors.New("String required for TransporterDetails ContactNo field ")
	} else if reflect.TypeOf(contractDetails.DeliveryDetails.TransporterDetails.Address).Kind() != reflect.String {
		return false, errors.New("String required for TransporterDetails Address field ")
	}

	for _, element := range contractDetails.TradeDetails {
		if reflect.TypeOf(element.ProductName).Kind() != reflect.String {
			return false, errors.New("String required for ProductName field ")
		} else if reflect.TypeOf(element.ProductDesc).Kind() != reflect.String {
			return false, errors.New("String required for ProductDesc field ")
		} else if reflect.TypeOf(element.ProductPrice).Kind() != reflect.String {
			return false, errors.New("String required for ProductPrice field ")
		} else if reflect.TypeOf(element.ProductQuantity).Kind() != reflect.String {
			return false, errors.New("String required for ProductQuantity field ")
		} else if reflect.TypeOf(element.TotalAmount).Kind() != reflect.String {
			return false, errors.New("String required for TotalAmount field ")
		}
	}
	return true, nil
}*/
