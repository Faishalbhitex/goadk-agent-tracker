package agent

import (
	"context"
	"fmt"
	"os"

	"finagent/internal/tools"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type readSheetArgs struct {
	SheetName     string `json:"sheetName" jsonschema:"Name of the sheet to read from"`
	RangeNotation string `json:"rangeNotation" jsonschema:"Range in A1 notation (e.g., 'A1:D10')"`
}

type readSheetResult struct {
	Status string          `json:"status"`
	Data   [][]interface{} `json:"data,omitempty"`
	Error  string          `json:"error,omitempty"`
}

type writeSheetArgs struct {
	SheetName     string          `json:"sheetName" jsonschema:"Name of the sheet to write to"`
	RangeNotation string          `json:"rangeNotation" jsonschema:"Range in A1 notation"`
	Values        [][]interface{} `json:"values" jsonschema:"2D array of values to write"`
}

type writeSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type appendSheetArgs struct {
	SheetName string          `json:"sheetName" jsonschema:"Name of the sheet to append to"`
	Values    [][]interface{} `json:"values" jsonschema:"2D array of rows to append"`
}

type appendSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type createSheetArgs struct {
	SheetTitle string `json:"sheetTitle" jsonschema:"Title for the new sheet"`
}

type listSheetsResult struct {
	Status string                   `json:"status"`
	Sheets []map[string]interface{} `json:"sheets,omitempty"`
	Error  string                   `json:"error,omitempty"`
}

type checkSheetEmptyArgs struct {
	SheetName string `json:"sheetName" jsonschema:"Name of the sheet to check"`
}

type checkSheetEmptyResult struct {
	Status  string `json:"status"`
	IsEmpty bool   `json:"isEmpty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type createSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func readFromSheet(ctx tool.Context, args readSheetArgs) (readSheetResult, error) {
	data, err := tools.ReadFromSheet(context.Background(), args.SheetName, args.RangeNotation)
	if err != nil {
		return readSheetResult{Status: "error", Error: err.Error()}, nil
	}
	return readSheetResult{Status: "success", Data: data}, nil
}

func writeToSheet(ctx tool.Context, args writeSheetArgs) (writeSheetResult, error) {
	err := tools.WriteToSheet(context.Background(), args.SheetName, args.RangeNotation, args.Values)
	if err != nil {
		return writeSheetResult{Status: "error", Error: err.Error()}, nil
	}
	msg := fmt.Sprintf("Successfully wrote to %s!%s", args.SheetName, args.RangeNotation)
	return writeSheetResult{Status: "success", Message: msg}, nil
}

func appendToSheet(ctx tool.Context, args appendSheetArgs) (appendSheetResult, error) {
	err := tools.AppendToSheet(context.Background(), args.SheetName, args.Values)
	if err != nil {
		return appendSheetResult{Status: "error", Error: err.Error()}, nil
	}
	msg := fmt.Sprintf("Successfully appended %d rows to %s", len(args.Values), args.SheetName)
	return appendSheetResult{Status: "success", Message: msg}, nil
}

func createNewSheet(ctx tool.Context, args createSheetArgs) (createSheetResult, error) {
	err := tools.CreateNewSheet(context.Background(), args.SheetTitle)
	if err != nil {
		return createSheetResult{Status: "error", Error: err.Error()}, nil
	}
	msg := fmt.Sprintf("Successfully created sheet '%s'", args.SheetTitle)
	return createSheetResult{Status: "success", Message: msg}, nil
}

func listSheets(ctx tool.Context, args struct{}) (listSheetsResult, error) {
	sheets, err := tools.ListSheets(context.Background())
	if err != nil {
		return listSheetsResult{Status: "error", Error: err.Error()}, nil
	}
	return listSheetsResult{Status: "success", Sheets: sheets}, nil
}

func checkSheetEmpty(ctx tool.Context, args checkSheetEmptyArgs) (checkSheetEmptyResult, error) {
	isEmpty, err := tools.CheckSheetEmpty(context.Background(), args.SheetName)
	if err != nil {
		return checkSheetEmptyResult{Status: "error", Error: err.Error()}, nil
	}

	msg := fmt.Sprintf("Sheet '%s' is ", args.SheetName)
	if isEmpty {
		msg += "empty"
	} else {
		msg += "not empty (contains data)"
	}

	return checkSheetEmptyResult{
		Status:  "success",
		IsEmpty: isEmpty,
		Message: msg,
	}, nil
}

func NewTrackerAgent(ctx context.Context) (agent.Agent, error) {
	if err := tools.InitSheetClient(ctx); err != nil {
		return nil, err
	}

	model, err := gemini.NewModel(ctx, "gemini-2.0-flash-lite", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		return nil, err
	}

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
			Description: "List all sheets in the spreadsheet with their info",
		},
		listSheets,
	)
	if err != nil {
		return nil, err
	}

	checkEmptyTool, err := functiontool.New(
		functiontool.Config{
			Name:        "check_sheet_empty",
			Description: "Check if a specific sheet is empty or contains data",
		},
		checkSheetEmpty,
	)
	if err != nil {
		return nil, err
	}

	trackerAgent, err := llmagent.New(llmagent.Config{
		Name:        "financial_tracker",
		Model:       model,
		Description: "A financial transaction tracker that manages data in Google Sheets",
		Instruction: SystemPrompt,
		Tools:       []tool.Tool{readTool, writeTool, appendTool, createSheetTool, listSheetsTool, checkEmptyTool},
	})
	if err != nil {
		return nil, err
	}

	return trackerAgent, nil
}
