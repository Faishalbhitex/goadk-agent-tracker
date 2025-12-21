package tools

import (
	"context"
	"fmt"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

func readFromSheet(ctx tool.Context, args ReadSheetArgs) (ReadSheetResult, error) {
	data, err := ReadFromSheet(context.Background(), args.SheetName, args.RangeNotation)
	if err != nil {
		return ReadSheetResult{Status: "error", Error: err.Error()}, nil
	}
	return ReadSheetResult{Status: "success", Data: data}, nil
}

func writeToSheet(ctx tool.Context, args WriteSheetArgs) (WriteSheetResult, error) {
	err := WriteToSheet(context.Background(), args.SheetName, args.RangeNotation, args.Values)
	if err != nil {
		return WriteSheetResult{Status: "error", Error: err.Error()}, nil
	}
	msg := fmt.Sprintf("Successfully wrote to %s!%s", args.SheetName, args.RangeNotation)
	return WriteSheetResult{Status: "success", Message: msg}, nil
}

func appendToSheet(ctx tool.Context, args AppendSheetArgs) (AppendSheetResult, error) {
	err := AppendToSheet(context.Background(), args.SheetName, args.Values)
	if err != nil {
		return AppendSheetResult{Status: "error", Error: err.Error()}, nil
	}
	msg := fmt.Sprintf("Successfully appended %d rows to %s", len(args.Values), args.SheetName)
	return AppendSheetResult{Status: "success", Message: msg}, nil
}

func createNewSheet(ctx tool.Context, args CreateSheetArgs) (CreateSheetResult, error) {
	err := CreateNewSheet(context.Background(), args.SheetTitle)
	if err != nil {
		return CreateSheetResult{Status: "error", Error: err.Error()}, nil
	}
	msg := fmt.Sprintf("Successfully created sheet 'Transaction_%s_{date}'", args.SheetTitle)
	return CreateSheetResult{Status: "success", Message: msg}, nil
}

func listSheets(ctx tool.Context, args struct{}) (ListSheetsResult, error) {
	sheets, err := ListSheetsWithInfo(context.Background())
	if err != nil {
		return ListSheetsResult{Status: "error", Error: err.Error()}, nil
	}
	return ListSheetsResult{
		Status:      "success",
		TotalSheets: len(sheets),
		Sheets:      sheets,
	}, nil
}

func NewAdkToolSheets() ([]tool.Tool, error) {
	readTool, err := functiontool.New(
		functiontool.Config{
			Name: "read_from_sheet",
			Description: `Read data from a specific Google Sheet range.
Usage: Get existing transaction data or check sheet contents.
Args: sheetName (string), rangeNotation (e.g., "A1:K10")
Returns: 2D array of cell values`,
		},
		readFromSheet,
	)
	if err != nil {
		return nil, err
	}

	writeTool, err := functiontool.New(
		functiontool.Config{
			Name: "write_to_sheet",
			Description: `Write/overwrite data to a specific Google Sheet range.
Usage: Update existing cells or modify specific ranges.
Args: sheetName, rangeNotation, values (2D array)
WARNING: This overwrites existing data. Use append_to_sheet for adding new rows.`,
		},
		writeToSheet,
	)
	if err != nil {
		return nil, err
	}

	appendTool, err := functiontool.New(
		functiontool.Config{
			Name: "append_to_sheet",
			Description: `Append new transaction rows to the end of a Google Sheet.
Usage: Add new transactions extracted from receipts or user input.
Args: 
  - sheetName: Target sheet name
  - values: Array of rows, each row MUST have 11 columns:
    [no, item_name, qty, unit, unit_price, amount, category, merchant, receipt_date, input_source, receipt_id]
    
IMPORTANT:
  - Leave 'no' empty ("") for auto-increment
  - Required fields: item_name, amount, merchant, receipt_id
  - Backend auto-fills: no, qty (default=1), receipt_date (if empty), input_source
  
Example:
  values: [
    ["", "Nasi Goreng", "1", "", "25000", "25000", "Food", "Warung Pak Budi", "2025-01-15T12:00:00", "image", "REC_20250115_001"]
  ]`,
		},
		appendToSheet,
	)
	if err != nil {
		return nil, err
	}

	createSheetTool, err := functiontool.New(
		functiontool.Config{
			Name: "create_new_sheet",
			Description: `Create a new Google Sheet with auto-naming and standard transaction headers.
Usage: Only when user explicitly requests a new sheet or no empty sheets exist.
Args: sheetTitle (string) - Short descriptive name (e.g., "Groceries", "Restaurant")

Auto-naming format: Transaction_{sheetTitle}_{YYYYMMDD}
Example: "Groceries" → "Transaction_Groceries_20251217"

The sheet will be created with:
  - Standard 11-column header (no, item_name, qty, unit, unit_price, amount, category, merchant, receipt_date, input_source, receipt_id)
  - Frozen header row
  - Formatted header (bold, light green background)
  
IMPORTANT: Always check existing sheets with list_sheets first. Prefer using empty sheets.`,
		},
		createNewSheet,
	)
	if err != nil {
		return nil, err
	}

	listSheetsTool, err := functiontool.New(
		functiontool.Config{
			Name: "list_sheets",
			Description: `List all available Google Sheets with metadata.
Usage: Check which sheets exist, find empty sheets, or get sheet info before operations.
Returns: {
  status: "success",
  totalSheets: number,
  sheets: [
    {
      title: string,
      sheetId: number,
      rowCount: number,
      colCount: number,
      isEmpty: boolean (true if only header exists or completely empty)
    }
  ]
}

Use this BEFORE creating new sheets to find available empty sheets.`,
		},
		listSheets,
	)
	if err != nil {
		return nil, err
	}

	return []tool.Tool{
		listSheetsTool, // List first (untuk discovery)
		readTool,
		appendTool, // Most used for transactions
		createSheetTool,
		writeTool, // Least used (careful operation)
	}, nil
}
