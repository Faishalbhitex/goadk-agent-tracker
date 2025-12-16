package tools

// Standard column indices
const (
	ColNo          = 0
	ColItemName    = 1
	ColQty         = 2
	ColUnit        = 3
	ColUnitPrice   = 4
	ColAmount      = 5
	ColCategory    = 6
	ColMerchant    = 7
	ColReceiptDate = 8
	ColInputSource = 9
	ColReceiptID   = 10
)

var DefaultHeaders = []string{
	"no", "item_name", "qty", "unit", "unit_price",
	"amount", "category", "merchant", "receipt_date",
	"input_source", "receipt_id",
}

// Sheet info
type SheetInfo struct {
	Title    string `json:"title"`
	SheetID  int64  `json:"sheetId"`
	RowCount int64  `json:"rowCount"`
	ColCount int64  `json:"colCount"`
	IsEmpty  bool   `json:"isEmpty"`
}

// Tool args & results
type ReadSheetArgs struct {
	SheetName     string `json:"sheetName"`
	RangeNotation string `json:"rangeNotation"`
}

type ReadSheetResult struct {
	Status string          `json:"status"`
	Data   [][]interface{} `json:"data,omitempty"`
	Error  string          `json:"error,omitempty"`
}

type WriteSheetArgs struct {
	SheetName     string          `json:"sheetName"`
	RangeNotation string          `json:"rangeNotation"`
	Values        [][]interface{} `json:"values"`
}

type WriteSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type AppendSheetArgs struct {
	SheetName string          `json:"sheetName"`
	Values    [][]interface{} `json:"values"`
}

type AppendSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type CreateSheetArgs struct {
	SheetTitle string `json:"sheetTitle"`
}

type CreateSheetResult struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ListSheetsResult struct {
	Status      string      `json:"status"`
	TotalSheets int         `json:"totalSheets,omitempty"`
	Sheets      []SheetInfo `json:"sheets,omitempty"`
	Error       string      `json:"error,omitempty"`
}
