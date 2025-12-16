package tools

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// === Public API untuk ADK Tools ===

func ReadFromSheet(ctx context.Context, sheetName, rangeNotation string) ([][]interface{}, error) {
	return globalClient.Read(ctx, sheetName, rangeNotation)
}

func WriteToSheet(ctx context.Context, sheetName, rangeNotation string, values [][]interface{}) error {
	return globalClient.Write(ctx, sheetName, rangeNotation, values)
}

func AppendToSheet(ctx context.Context, sheetName string, values [][]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no data to append")
	}

	// Normalize & validate rows
	normalized, err := normalizeRows(ctx, sheetName, values)
	if err != nil {
		return err
	}

	return globalClient.Append(ctx, sheetName, normalized)
}

func CreateNewSheet(ctx context.Context, sheetTitle string) error {
	// Format: Transaction_{title}_{YYYYMMDD}
	timestamp := time.Now().Format("20060102")
	formattedTitle := fmt.Sprintf("Transaction_%s_%s", sheetTitle, timestamp)

	// Create sheet
	sheetID, err := globalClient.Create(ctx, formattedTitle)
	if err != nil {
		return err
	}

	// Write header
	headerRange := fmt.Sprintf("A1:%s1", columnLetter(len(DefaultHeaders)))
	headerValues := [][]interface{}{toInterfaceSlice(DefaultHeaders)}

	if err := globalClient.Write(ctx, formattedTitle, headerRange, headerValues); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Format header (non-critical, don't fail)
	if err := globalClient.FormatHeader(ctx, sheetID, len(DefaultHeaders)); err != nil {
		log.Printf("⚠ Warning: failed to format header: %v", err)
	}

	log.Printf("✓ Created sheet: %s", formattedTitle)
	return nil
}

func ListSheetsWithInfo(ctx context.Context) ([]SheetInfo, error) {
	return globalClient.ListSheets(ctx)
}

// === Internal helpers ===

func normalizeRows(ctx context.Context, sheetName string, rows [][]interface{}) ([][]interface{}, error) {
	lastNo, _ := globalClient.GetLastRowNumber(ctx, sheetName)
	nextNo := lastNo + 1

	normalized := make([][]interface{}, 0, len(rows))

	for i, row := range rows {
		normalizedRow, err := normalizeRow(row, nextNo, i+1)
		if err != nil {
			return nil, err
		}
		normalized = append(normalized, normalizedRow)
		nextNo++
	}

	return normalized, nil
}

func normalizeRow(row []interface{}, nextNo, rowIndex int) ([]interface{}, error) {
	// Ensure 11 columns
	normalized := make([]interface{}, 11)
	copy(normalized, row)

	// Auto-fill defaults
	if isEmpty(normalized[ColNo]) {
		normalized[ColNo] = nextNo
	}
	if isEmpty(normalized[ColQty]) {
		normalized[ColQty] = 1
	}
	if isEmpty(normalized[ColReceiptDate]) {
		normalized[ColReceiptDate] = time.Now().Format(time.RFC3339)
	}
	if isEmpty(normalized[ColInputSource]) {
		normalized[ColInputSource] = "manual"
	}

	// Validate required fields
	required := map[int]string{
		ColItemName:  "item_name",
		ColAmount:    "amount",
		ColMerchant:  "merchant",
		ColReceiptID: "receipt_id",
	}

	for col, name := range required {
		if isEmpty(normalized[col]) {
			return nil, fmt.Errorf("row %d: missing required field '%s'", rowIndex, name)
		}
	}

	return normalized, nil
}

func isEmpty(val interface{}) bool {
	if val == nil {
		return true
	}
	return strings.TrimSpace(fmt.Sprintf("%v", val)) == ""
}

func toInterfaceSlice(s []string) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}

func columnLetter(col int) string {
	col--
	if col < 0 {
		return "A"
	}
	letters := ""
	for col >= 0 {
		letters = string(rune('A'+(col%26))) + letters
		col = col/26 - 1
	}
	return letters
}
