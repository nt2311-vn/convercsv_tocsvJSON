package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Invoice struct {
	InternalID string  `json:"internal_id"`
	AppliedAmt float64 `json:"applied_amt"`
}

type Journal struct {
	InternalID string  `json:"internal_id"`
	AppliedAmt float64 `json:"applied_amt"`
}

type PaymentInfo struct {
	PaymentRef    string  `json:"payment_ref"`
	PaymentDate   string  `json:"payment_date"`
	CustomerCode  string  `json:"customer_code"`
	Memo          string  `json:"memo"`
	BankName      string  `json:"bank_name"`
	LocationCode  string  `json:"location_code"`
	PaymentAmount float64 `json:"payment_amount"`
}

type GroupRecords struct {
	PaymentObj PaymentInfo `json:"payment_info"`
	Invoices   []Invoice   `json:"invoices"`
	Journals   []Journal   `json:"journals"`
}

func main() {
	dir := "./data"

	csvFile, err := findCSVFile(dir)
	if err != nil {
		log.Fatalf("Error finding CSV file: %v", err)
	}

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error on read all records: %s", err)
	}

	data := records[1:]

	groupedData := make(map[string]*GroupRecords)

	for _, row := range data {
		paymentID, customerCode, paymentDate := row[0], row[1], row[2]
		memo, bankName, locationCode := row[3], row[4], row[5]
		paymentAmt, internalID, recordType, appliedAmt := row[6], row[7], row[8], row[9]

		if groupedData[paymentID] == nil {

			paidAmt, err := strconv.ParseFloat(paymentAmt, 64)
			if err != nil {
				log.Fatalf("Error parsing applied amount: %v", err)
			}
			groupedData[paymentID] = &GroupRecords{
				PaymentObj: PaymentInfo{
					PaymentRef:    paymentID,
					CustomerCode:  customerCode,
					PaymentDate:   paymentDate,
					Memo:          memo,
					BankName:      bankName,
					LocationCode:  locationCode,
					PaymentAmount: paidAmt,
				},
			}
		}

		applyingAmt, err := strconv.ParseFloat(appliedAmt, 64)
		if err != nil {
			log.Fatalf("Error parsing applied amount: %v", err)
		}

		switch recordType {
		case "invoice":
			groupedData[paymentID].Invoices = append(groupedData[paymentID].Invoices, Invoice{
				InternalID: internalID,
				AppliedAmt: applyingAmt,
			})

		case "journal":
			groupedData[paymentID].Journals = append(groupedData[paymentID].Journals, Journal{
				InternalID: internalID,
				AppliedAmt: applyingAmt,
			})

		default:
			log.Fatalf("Unknown record type: %s", recordType)
		}
	}

	if err := os.MkdirAll("./result", os.ModePerm); err != nil {
		log.Fatalf("Error creating result directory: %v", err)
	}

	fileName := fmt.Sprintf("%s_output.csv", time.Now().Format("20060102_150405"))
	outputFile, err := os.Create("./result/" + fileName)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}

	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	defer writer.Flush()

	header := []string{"File Name", "Payment Info", "Invoice", "Journal"}
	writer.Write(header)

	for _, group := range groupedData {
		paymentInfo, err := json.Marshal(group.PaymentObj)
		if err != nil {
			log.Fatalf("Error marshalling payment info: %v", err)
		}

		invoicesJSON, err := json.Marshal(group.Invoices)
		if err != nil {
			log.Fatalf("Error marshalling invoices: %v", err)
		}

		journalsJSON, err := json.Marshal(group.Journals)
		if err != nil {
			log.Fatalf("Error marshalling journals: %v", err)
		}

		writer.Write(
			[]string{fileName, string(paymentInfo), string(invoicesJSON), string(journalsJSON)},
		)

	}
}

func findCSVFile(dir string) (string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".csv" {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no CSV file found in %s", dir)
}
