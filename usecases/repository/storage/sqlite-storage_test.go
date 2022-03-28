package storage

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/norby7/shortening-service/entities"
	"testing"
)

var (
	dbMock sqlmock.Sqlmock
)

var MockOpener = func(string, string) (*sql.DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return nil, fmt.Errorf("unable to create mock driver: %s", err.Error())
	}

	dbMock = mock

	return db, nil
}

var MockErrOpener = func(d string, p string) (*sql.DB, error) {
	if p == "errPath" {
		return nil, fmt.Errorf("unable to connect to database")
	}

	return nil, nil
}

func TestNewRepository(t *testing.T) {
	SqlOpen = MockErrOpener
	testCases := []struct {
		name    string
		input   string
		isError bool
	}{{
		name:    "valid path",
		input:   "./database/sqlite/urls.db",
		isError: false,
	}, {
		name:    "invalid path",
		input:   "errPath",
		isError: true,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewSqliteStorage(tc.input, 0)

			if (err != nil) != tc.isError {
				t.Errorf("expected error (%v), got error (%v)", tc.isError, err.Error())
			}
		})
	}
}

func TestValidAdd(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	u := entities.Url{
		Id:       1,
		Code:     "84gfj4i9",
		Url:      "https://google.com",
		ShortUrl: "http://localhost/84gfj4i9",
		Domain:   "http://localhost",
		Counter:  1,
	}

	dbMock.ExpectExec(`INSERT INTO urls`).WithArgs(u.Code, u.Url, u.Counter, u.ShortUrl, u.Domain).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Add(&u)
	if err != nil {
		t.Fatalf("unable to execute add call: %s", err.Error())
	}
}

func TestInsertErrorAdd(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	u := entities.Url{
		Id:       1,
		Code:     "84gfj4i9",
		Url:      "https://google.com",
		ShortUrl: "http://localhost/84gfj4i9",
		Domain:   "http://localhost",
		Counter:  1,
	}

	insertErr := fmt.Errorf("error executing insert query")

	dbMock.ExpectExec(`INSERT INTO urls`).WithArgs(u.Code, u.Url, u.Counter, u.ShortUrl, u.Domain).WillReturnError(insertErr)

	err = repo.Add(&u)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", insertErr)
	}
}

func TestValidDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	dbMock.ExpectExec(`DELETE FROM urls`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(1)
	if err != nil {
		t.Fatalf("unable to execute delete call: %s", err.Error())
	}
}

func TestErrorDelete(t *testing.T) {
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	deleteErr := fmt.Errorf("erorr executing delete query")
	dbMock.ExpectExec(`DELETE FROM urls`).WithArgs(1).WillReturnError(deleteErr)

	err = repo.Delete(1)
	if err == nil {
		t.Errorf("expected error (%v), got error nil", deleteErr)
	}
}

func TestValidGetUrlByCode(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	rows := sqlmock.NewRows([]string{"url"})
	rows.AddRow("84gfj4i9")

	dbMock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	_, err = repo.GetUrlByCode("84gfj4i9")
	if err != nil{
		t.Fatalf("unable to execute get by code call: %s", err.Error())
	}
}

func TestNoRowsGetUrlByCode(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	dbMock.ExpectQuery(`SELECT`).WillReturnError(sql.ErrNoRows)

	_, err = repo.GetUrlByCode("84gfj4i9")
	if err != nil{
		t.Fatalf("expected no error, got: %s", err.Error())
	}
}

func TestErrorGetUrlByCode(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	queryErr := fmt.Errorf("error fetching data")
	dbMock.ExpectQuery(`SELECT`).WillReturnError(queryErr)

	_, err = repo.GetUrlByCode("84gfj4i9")
	if err == nil{
		t.Errorf("expected error (%v), got error nil", queryErr)
	}
}

func TestValidGetById(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	rows := sqlmock.NewRows([]string{"id", "code", "url", "shortUrl", "domain", "counter"})
	rows.AddRow("1", "84gfj4i9", "https://google.com", "http://localhost/84gfj4i9", "http://localhost", "0")

	dbMock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	_, err = repo.GetById(1)
	if err != nil{
		t.Fatalf("unable to execute get by code call: %s", err.Error())
	}
}

func TestNoRowsGetById(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	dbMock.ExpectQuery(`SELECT`).WillReturnError(sql.ErrNoRows)

	_, err = repo.GetById(0)
	if err != nil{
		t.Fatalf("expected no error, got: %s", err.Error())
	}
}

func TestErrorGetById(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	queryErr := fmt.Errorf("error fetching data")
	dbMock.ExpectQuery(`SELECT`).WillReturnError(queryErr)

	_, err = repo.GetById(1)
	if err == nil{
		t.Errorf("expected error (%v), got error nil", queryErr)
	}
}

func TestValidGetByUrl(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	rows := sqlmock.NewRows([]string{"id", "code", "url", "shortUrl", "domain", "counter"})
	rows.AddRow("1", "84gfj4i9", "https://google.com", "http://localhost/84gfj4i9", "http://localhost", "0")

	dbMock.ExpectQuery(`SELECT`).WillReturnRows(rows)

	_, err = repo.GetByUrl("https://google.com")
	if err != nil{
		t.Fatalf("unable to execute get by code call: %s", err.Error())
	}
}

func TestNoRowsGetByUrl(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	dbMock.ExpectQuery(`SELECT`).WillReturnError(sql.ErrNoRows)

	_, err = repo.GetByUrl("https://google1.com")
	if err != nil{
		t.Fatalf("expected no error, got: %s", err.Error())
	}
}

func TestErrorGetByUrl(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	queryErr := fmt.Errorf("error fetching data")
	dbMock.ExpectQuery(`SELECT`).WillReturnError(queryErr)

	_, err = repo.GetByUrl("https://google.com")
	if err == nil{
		t.Errorf("expected error (%v), got error nil", queryErr)
	}
}

func TestValidIncrementCounter(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	dbMock.ExpectExec(`UPDATE urls`).WithArgs("www.test.com").WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.IncrementCounter("www.test.com")
	if err != nil{
		t.Errorf("expected no error, got: %s", err.Error())
	}
}

func TestErrorIncrementCounter(t *testing.T){
	SqlOpen = MockOpener
	repo, err := NewSqliteStorage("./database/sqlite/test.db", 0)
	if err != nil {
		t.Fatalf("unable to create mock repository: %s", err.Error())
	}

	updateErr := fmt.Errorf("error executing update query")
	dbMock.ExpectExec(`UPDATE urls`).WithArgs("www.test.com").WillReturnError(updateErr)

	err = repo.IncrementCounter("www.test.com")
	if err == nil{
		t.Errorf("expected error (%v), got error nil", updateErr)
	}
}