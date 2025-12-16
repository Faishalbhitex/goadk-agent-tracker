package tools

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SheetClient struct {
	service       *sheets.Service
	spreadsheetID string
}

func NewSheetClient(ctx context.Context, credPath, spreadsheetID string) (*SheetClient, error) {
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(credPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %v", err)
	}
	return &SheetClient{
		service:       srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

func (s *SheetClient) ReadSheet(ctx context.Context, sheetName, rangeNotation string) ([][]interface{}, error) {
	readRange := fmt.Sprintf("%s!%s", sheetName, rangeNotation)
	resp, err := s.service.Spreadsheets.Values.Get(s.spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet: %v", err)
	}
	return resp.Values, nil
}

func (s *SheetClient) WriteSheet(ctx context.Context, sheetName, rangeNotation string, values [][]interface{}) error {
	writeRange := fmt.Sprintf("%s!%s", sheetName, rangeNotation)
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	_, err := s.service.Spreadsheets.Values.Update(
		s.spreadsheetID,
		writeRange,
		valueRange,
	).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to write sheet: %v", err)
	}
	return nil
}

func (s *SheetClient) AppendSheet(ctx context.Context, sheetName string, values [][]interface{}) error {
	valueRange := &sheets.ValueRange{
		Values: values,
	}
	_, err := s.service.Spreadsheets.Values.Append(
		s.spreadsheetID,
		sheetName,
		valueRange,
	).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to append sheet: %v", err)
	}
	return nil
}

var DefaultHeaders = []string{
	"no",
	"item_name",
	"qty",
	"unit",
	"unit_price",
	"amount",
	"category",
	"merchant",
	"receipt_date",
	"input_source",
	"receipt_id",
}

func (s *SheetClient) CreateSheet(ctx context.Context, sheetTitle string) error {
	// 1. Buat sheet baru
	addSheetReq := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetTitle,
			},
		},
	}

	batchReq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{addSheetReq},
	}

	resp, err := s.service.Spreadsheets.BatchUpdate(s.spreadsheetID, batchReq).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create sheet: %v", err)
	}

	if len(resp.Replies) == 0 || resp.Replies[0].AddSheet == nil {
		return fmt.Errorf("sheet creation response invalid")
	}
	sheetID := resp.Replies[0].AddSheet.Properties.SheetId

	// 2. Tulis header ke baris pertama
	headerRange := fmt.Sprintf("%s!A1:%s1", sheetTitle, columnLetter(len(DefaultHeaders)))
	_, err = s.service.Spreadsheets.Values.Update(s.spreadsheetID, headerRange, &sheets.ValueRange{
		Values: [][]interface{}{stringSliceToInterfaceSlice(DefaultHeaders)},
	}).ValueInputOption("RAW").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	// 3. Format header: bold + background hijau muda
	formatReq := &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    0,
				EndRowIndex:      1,
				StartColumnIndex: 0,
				EndColumnIndex:   int64(len(DefaultHeaders)),
			},
			Cell: &sheets.CellData{
				UserEnteredFormat: &sheets.CellFormat{
					TextFormat: &sheets.TextFormat{Bold: true},
					BackgroundColorStyle: &sheets.ColorStyle{
						RgbColor: &sheets.Color{Red: 0.85, Green: 0.95, Blue: 0.85},
					},
				},
			},
			Fields: "userEnteredFormat(textFormat,backgroundColorStyle)",
		},
	}

	// 4. Auto-resize kolom
	var resizeRequests []*sheets.Request
	for i := 0; i < len(DefaultHeaders); i++ {
		resizeRequests = append(resizeRequests, &sheets.Request{
			AutoResizeDimensions: &sheets.AutoResizeDimensionsRequest{
				Dimensions: &sheets.DimensionRange{
					SheetId:    sheetID,
					Dimension:  "COLUMNS",
					StartIndex: int64(i),
					EndIndex:   int64(i + 1),
				},
			},
		})
	}

	// 5. Kirim semua format + resize
	formatBatch := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: append([]*sheets.Request{formatReq}, resizeRequests...),
	}
	_, err = s.service.Spreadsheets.BatchUpdate(s.spreadsheetID, formatBatch).Context(ctx).Do()
	if err != nil {
		log.Printf("Warning: failed to format sheet (headers may be unstyled): %v", err)
		// Tetap lanjut, jangan fail total
	}

	return nil
}

// Helper: konversi []string → []interface{}
func stringSliceToInterfaceSlice(s []string) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		res[i] = v
	}
	return res
}

// Helper: angka kolom → huruf (A, B, ..., Z, AA, AB, ...)
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

func (s *SheetClient) ListSheetsWithInfo(ctx context.Context) ([]SheetInfo, error) {
	resp, err := s.service.Spreadsheets.Get(s.spreadsheetID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get spreadsheet info: %v", err)
	}

	var sheets []SheetInfo
	for _, sheet := range resp.Sheets {
		// Check if empty by reading first cell
		isEmpty := false
		data, err := s.ReadSheet(ctx, sheet.Properties.Title, "A1")
		if err == nil && len(data) == 0 {
			isEmpty = true
		}

		sheets = append(sheets, SheetInfo{
			Title:    sheet.Properties.Title,
			SheetID:  sheet.Properties.SheetId,
			RowCount: sheet.Properties.GridProperties.RowCount,
			ColCount: sheet.Properties.GridProperties.ColumnCount,
			IsEmpty:  isEmpty,
		})
	}
	return sheets, nil
}

var sheetClient *SheetClient

func InitSheetClient(ctx context.Context) error {
	credPath := os.Getenv("GOOGLE_SA_PATH")
	spreadsheetID := os.Getenv("SPREADSHEET_ID")

	var err error
	sheetClient, err = NewSheetClient(ctx, credPath, spreadsheetID)
	return err
}

func ReadFromSheet(ctx context.Context, sheetName, rangeNotation string) ([][]interface{}, error) {
	return sheetClient.ReadSheet(ctx, sheetName, rangeNotation)
}

func WriteToSheet(ctx context.Context, sheetName, rangeNotation string, values [][]interface{}) error {
	return sheetClient.WriteSheet(ctx, sheetName, rangeNotation, values)
}

func AppendToSheet(ctx context.Context, sheetName string, values [][]interface{}) error {
	return sheetClient.AppendSheet(ctx, sheetName, values)
}

func CreateNewSheet(ctx context.Context, sheetTitle string) error {
	return sheetClient.CreateSheet(ctx, sheetTitle)
}

func ListSheetsWithInfo(ctx context.Context) ([]SheetInfo, error) {
	return sheetClient.ListSheetsWithInfo(ctx)
}
