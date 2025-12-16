package tools

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}
	return &SheetClient{
		service:       srv,
		spreadsheetID: spreadsheetID,
	}, nil
}

// === Basic CRUD Operations ===

func (s *SheetClient) Read(ctx context.Context, sheetName, rangeNotation string) ([][]interface{}, error) {
	// Tambahkan single quotes di '%s'
	fullRange := fmt.Sprintf("'%s'!%s", sheetName, rangeNotation)

	resp, err := s.service.Spreadsheets.Values.Get(s.spreadsheetID, fullRange).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("read failed: %w", err)
	}
	return resp.Values, nil
}

func (s *SheetClient) Write(ctx context.Context, sheetName, rangeNotation string, values [][]interface{}) error {
	// Tambahkan single quotes di '%s'
	fullRange := fmt.Sprintf("'%s'!%s", sheetName, rangeNotation)

	valueRange := &sheets.ValueRange{Values: values}

	_, err := s.service.Spreadsheets.Values.Update(
		s.spreadsheetID,
		fullRange,
		valueRange,
	).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("write failed: %w", err)
	}
	return nil
}

// === Sheet Management ===

func (s *SheetClient) Append(ctx context.Context, sheetName string, values [][]interface{}) error {
	valueRange := &sheets.ValueRange{Values: values}

	// Bungkus sheetName dengan single quotes ('') untuk menangani spasi
	safeRange := fmt.Sprintf("'%s'", sheetName)

	_, err := s.service.Spreadsheets.Values.Append(
		s.spreadsheetID,
		safeRange, // Gunakan safeRange, bukan sheetName
		valueRange,
	).ValueInputOption("USER_ENTERED").Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("append failed: %w", err)
	}
	return nil
}

func (s *SheetClient) Create(ctx context.Context, title string) (int64, error) {
	req := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: title,
				GridProperties: &sheets.GridProperties{
					FrozenRowCount: 1,
				},
			},
		},
	}

	batchReq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{req},
	}

	resp, err := s.service.Spreadsheets.BatchUpdate(s.spreadsheetID, batchReq).Context(ctx).Do()
	if err != nil {
		return 0, fmt.Errorf("create sheet failed: %w", err)
	}

	if len(resp.Replies) == 0 || resp.Replies[0].AddSheet == nil {
		return 0, fmt.Errorf("invalid create response")
	}

	return resp.Replies[0].AddSheet.Properties.SheetId, nil
}

func (s *SheetClient) ListSheets(ctx context.Context) ([]SheetInfo, error) {
	resp, err := s.service.Spreadsheets.Get(s.spreadsheetID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("list sheets failed: %w", err)
	}

	var sheets []SheetInfo
	for _, sheet := range resp.Sheets {
		isEmpty := false
		data, err := s.Read(ctx, sheet.Properties.Title, "A1")
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

func (s *SheetClient) GetSheetInfo(ctx context.Context, sheetName string) (*SheetInfo, error) {
	sheets, err := s.ListSheets(ctx)
	if err != nil {
		return nil, err
	}

	for _, sheet := range sheets {
		if sheet.Title == sheetName {
			return &sheet, nil
		}
	}
	return nil, fmt.Errorf("sheet '%s' not found", sheetName)
}

// === Formatting ===

func (s *SheetClient) FormatHeader(ctx context.Context, sheetID int64, colCount int) error {
	formatReq := &sheets.Request{
		RepeatCell: &sheets.RepeatCellRequest{
			Range: &sheets.GridRange{
				SheetId:          sheetID,
				StartRowIndex:    0,
				EndRowIndex:      1,
				StartColumnIndex: 0,
				EndColumnIndex:   int64(colCount),
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

	batchReq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{formatReq},
	}

	_, err := s.service.Spreadsheets.BatchUpdate(s.spreadsheetID, batchReq).Context(ctx).Do()
	return err
}

// === Helpers ===

func (s *SheetClient) GetLastRowNumber(ctx context.Context, sheetName string) (int, error) {
	info, err := s.GetSheetInfo(ctx, sheetName)
	if err != nil {
		return 0, err
	}

	if info.RowCount <= 1 {
		return 0, nil
	}

	lastRow := int(info.RowCount)
	rangeNotation := fmt.Sprintf("A%d", lastRow)
	data, err := s.Read(ctx, sheetName, rangeNotation)

	if err != nil || len(data) == 0 || len(data[0]) == 0 {
		return 0, nil
	}

	var lastNum int
	fmt.Sscanf(fmt.Sprintf("%v", data[0][0]), "%d", &lastNum)
	return lastNum, nil
}

// === Global singleton ===

var globalClient *SheetClient

func InitSheetClient(ctx context.Context) error {
	credPath := os.Getenv("GOOGLE_SA_PATH")
	spreadsheetID := os.Getenv("SPREADSHEET_ID")

	var err error
	globalClient, err = NewSheetClient(ctx, credPath, spreadsheetID)
	return err
}
