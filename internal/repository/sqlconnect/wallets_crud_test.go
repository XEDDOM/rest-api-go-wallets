package sqlconnect

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func mockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cannot create sqlmock: %v", err)
	}
	ConnectDb = func() (*sql.DB, error) {
		return db, nil
	}
	return db, mock
}

func TestGetWalletByID_OK(t *testing.T) {
	db, mock := mockDB(t)
	defer db.Close()
	rows := sqlmock.NewRows([]string{"valletid", "amount"}).AddRow(1, 300)
	mock.ExpectQuery(`SELECT valletid, amount FROM wallets WHERE valletid = \$1`).WithArgs(1).WillReturnRows(rows)
	w, err := GetWalletByID(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w.Amount != 300 {
		t.Fatalf("expected amount=300, got %d", w.Amount)
	}
}

func TestGetWalletByID_NotFound(t *testing.T) {
	db, mock := mockDB(t)
	defer db.Close()
	mock.ExpectQuery(`SELECT .* FROM wallets WHERE valletid = \$1`).WithArgs(1).WillReturnError(sql.ErrNoRows)
	_, err := GetWalletByID(1)
	if err == nil || err.Error() != "wallet not found" {
		t.Fatalf("expected wallet not found error")
	}
}

func TestUpdateWallet_Deposit_OK(t *testing.T) {
	db, mock := mockDB(t)
	defer db.Close()
	input := map[string]any{
		"valletId":      1,
		"operationType": "DEPOSIT",
		"amount":        100,
	}
	mock.ExpectQuery(`SELECT valletid, amount FROM wallets WHERE valletid = \$1`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"valletid", "amount"}).AddRow(1, 200))
	mock.ExpectExec(regexp.QuoteMeta("UPDATE wallets SET amount = $1 WHERE valletid = $2")).WithArgs(300, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	err := UpdateWallet(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateWallet_InsufficientFunds(t *testing.T) {
	db, mock := mockDB(t)
	defer db.Close()
	input := map[string]any{
		"valletId":      1,
		"operationType": "WITHDRAW",
		"amount":        500,
	}
	mock.ExpectQuery(`SELECT valletid, amount FROM wallets WHERE valletid = \$1`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"valletid", "amount"}).AddRow(1, 100))
	err := UpdateWallet(input)
	if err == nil || err.Error() != "insufficient funds" {
		t.Fatalf("expected insufficient funds, got %#v", err)
	}
}

func TestUpdateWallet_NotFound(t *testing.T) {
	db, mock := mockDB(t)
	defer db.Close()
	input := map[string]any{
		"valletId":      2,
		"operationType": "DEPOSIT",
		"amount":        50,
	}
	mock.ExpectQuery(`SELECT .* FROM wallets WHERE valletid = \$1`).WithArgs(2).WillReturnError(sql.ErrNoRows)
	err := UpdateWallet(input)
	if err == nil || err.Error() != "wallet not found" {
		t.Fatalf("expected wallet not found error")
	}
}

func TestUpdateWallet_InvalidOperation(t *testing.T) {
	db, mock := mockDB(t)
	defer db.Close()
	input := map[string]any{
		"valletId":      1,
		"operationType": "SOMETHING",
		"amount":        10,
	}
	mock.ExpectQuery(`SELECT .* FROM wallets WHERE valletid = \$1`).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"valletid", "amount"}).AddRow(1, 500))
	err := UpdateWallet(input)
	if err == nil || err.Error() != "unknown operationType. Must be DEPOSIT or WITHDRAW" {
		t.Fatalf("expected operationType error")
	}
}

func TestUpdateWallet_NegativeAmount(t *testing.T) {
	input := map[string]any{
		"valletId":      1,
		"operationType": "DEPOSIT",
		"amount":        -10,
	}
	_, mock := mockDB(t)
	mock.ExpectationsWereMet()
	err := UpdateWallet(input)
	if err == nil || err.Error() != "amount must be positive" {
		t.Fatalf("expected amount positive error")
	}
}
