package main

import (
	"database/sql"
	"fmt"
	"testing"
)

var (
	email        = "test1@test.com"
	emailInvalid = "bad!email@@test.com"
)

func TestInit(t *testing.T) {
	initialize("meiru_test")
	db.Exec("DELETE FROM emails")
}

func TestInsert(t *testing.T) {
	fmt.Print("Running insert test ... ")
	err := insertEmail(email)
	if err != nil {
		t.Fatal("TestInsert failed!")
	}
	fmt.Printf("OK: %s inserted.\n", email)
}

func TestDuplicateInsert(t *testing.T) {
	fmt.Print("Running duplicate insert test ... ")
	var numResults string
	insertEmail(email)
	insertEmail(email)
	db.QueryRow("SELECT COUNT(email) FROM emails").Scan(&numResults)
	if numResults != "1" {
		t.Fatalf("numResults %v NOT 1", numResults)
	}
	fmt.Printf("OK: %s result after duplicate insert.\n", numResults)
}

func TestInvalidInsert(t *testing.T) {
	fmt.Print("Running insert test ... ")
	var res string
	insertEmail(emailInvalid)
	err := db.QueryRow("SELECT email FROM emails WHERE email=$1", emailInvalid).Scan(&res)
	if err != sql.ErrNoRows {
		t.Fatalf("Invalid email was inserted")
	}
	fmt.Printf("OK: %s not inserted.\n", emailInvalid)
}

func TestGetMails(t *testing.T) {
	fmt.Print("Running get mails test ... ")
	mails := getMails()
	if len(mails) != 1 {
		t.Fatalf("number of results %v NOT 1", len(mails))
	}
	fmt.Printf("OK: %s\n", mails[0])
}

func TestTeardown(t *testing.T) {
	db.Close()
}
