package tools

// reade sheet
type ReadSheetArgs struct {
	SheetName     string `json:"sheetName" jsonschema:"Name of the sheet to read from"`
	RangeNotation string `json:"rangeNotation" jsonschema:"Range in A1 notation (e.g., 'A1:D10')"`
}

type ReadSheetResult struct {
	Status string          `json:"status"`
	Data   [][]interface{} `json:"data,omitempty"`
	Error  string          `json:"error,omitempty"`
}

// write sheet
type WriteSheetArgs struct {
	SheetName     string          `json:"sheetName" jsonschema:"Name of the sheet to write to"`
	RangeNotation string          `json:"rangeNotation" jsonschema:"Range in A1 notation"`
	Values        [][]interface{} `json:"values" jsonschema:"2D array of values to write"`
}

type WriteSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// append sheet
type AppendSheetArgs struct {
	SheetName string          `json:"sheetName" jsonschema:"Name of the sheet to append to"`
	Values    [][]interface{} `json:"values" jsonschema:"2D array of rows to append"`
}

type AppendSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// create sheet
type CreateSheetArgs struct {
	SheetTitle string `json:"sheetTitle" jsonschema:"Title for the new sheet"`
}

type CreateSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// list sheet
type ListSheetsResult struct {
	Status string                   `json:"status"`
	Sheets []map[string]interface{} `json:"sheets,omitempty"`
	Error  string                   `json:"error,omitempty"`
}

// check empty sheet
type CheckSheetEmptyArgs struct {
	SheetName string `json:"sheetName" jsonschema:"Name of the sheet to check"`
}

type CheckSheetEmptyResult struct {
	Status  string `json:"status"`
	IsEmpty bool   `json:"isEmpty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}
