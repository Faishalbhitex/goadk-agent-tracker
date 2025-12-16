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

func (s *SheetClient) CreateSheet(ctx context.Context, sheetTitle string) error {
	req := &sheets.Request{
		AddSheet: &sheets.AddSheetRequest{
			Properties: &sheets.SheetProperties{
				Title: sheetTitle,
			},
		},
	}

	batchUpdate := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{req},
	}

	_, err := s.service.Spreadsheets.BatchUpdate(s.spreadsheetID, batchUpdate).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to create sheet: %v", err)
	}
	return nil
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
