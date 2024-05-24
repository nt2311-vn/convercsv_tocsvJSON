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

			paidAmt, err := strconv.ParseFloat(appliedAmt, 64)
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
				break

			case "journal":
				groupedData[paymentID].Journals = append(groupedData[paymentID].Journals, Journal{
					InternalID: internalID,
					AppliedAmt: applyingAmt,
				})
				break

			default:
				log.Fatalf("Unknown record type: %s", recordType)
			}

			var groupRecordsList []GroupRecords

			for _, groupRecords := range groupedData {
				groupRecordsList = append(groupRecordsList, *groupRecords)
			}

			jsonData, err := json.Marshal(groupRecordsList)
			if err != nil {
				log.Fatalf("Error marshaling JSON: %v", err)
			}

			fileName := fmt.Sprintf("%s_output.json", time.Now().Format("2006-01-02T15:04:05"))
			outputFile, err := os.Create()

		}
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
