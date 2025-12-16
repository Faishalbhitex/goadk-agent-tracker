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
	msg := fmt.Sprintf("Successfully created sheet '%s'", args.SheetTitle)
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
			Name:        "read_from_sheet",
			Description: "Read data from a Google Sheet range",
		},
		readFromSheet,
	)
	if err != nil {
		return nil, err
	}

	writeTool, err := functiontool.New(
		functiontool.Config{
			Name:        "write_to_sheet",
			Description: "Write data to a specific Google Sheet range",
		},
		writeToSheet,
	)
	if err != nil {
		return nil, err
	}

	appendTool, err := functiontool.New(
		functiontool.Config{
			Name:        "append_to_sheet",
			Description: "Append new rows to the end of a Google Sheet",
		},
		appendToSheet,
	)
	if err != nil {
		return nil, err
	}

	createSheetTool, err := functiontool.New(
		functiontool.Config{
			Name:        "create_new_sheet",
			Description: "Create a new sheet with the given title",
		},
		createNewSheet,
	)
	if err != nil {
		return nil, err
	}

	listSheetsTool, err := functiontool.New(
		functiontool.Config{
			Name:        "list_sheets",
			Description: "List all sheets with complete info (title, isEmpty, rowCount, colCount)",
		},
		listSheets,
	)
	if err != nil {
		return nil, err
	}

	return []tool.Tool{
		readTool,
		writeTool,
		appendTool,
		createSheetTool,
		listSheetsTool,
	}, nil
}
