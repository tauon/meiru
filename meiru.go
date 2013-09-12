// Meiru: collect 'them mails
// By jbidzos
// GPLv3: http://www.gnu.org/licenses/gpl-3.0.html

package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"regexp"
	"strings"
)

var (
	db          *sql.DB
	emailregexp *regexp.Regexp
	address     *string = flag.String("address", "127.0.0.1", "Address server listens on")
	port        *int    = flag.Int("port", 10025, "TCP port server listens on")
	dbName      *string = flag.String("dbname", "meiru", "Name of database")
	dbUser      *string = flag.String("dbuser", "postgres", "Database user name")
	useTLS      *bool   = flag.Bool("tls", false, "Use TLS; requires certificate (cert.pem) and private key (key.pem)")
	verbosity   *int    = flag.Int("v", 0, "Output verbosity level (0=error, 1=info, 2=debug)")
	listen      *bool   = flag.Bool("l", false, "Listen for POST requests")
	separator   *string = flag.String("sep", ",", "Separator to use when dumping email list")
)

func initialize(dbname string) {
	emailregexp, _ = regexp.Compile(`^(|(([A-Za-z0-9]+_+)|([A-Za-z0-9]+\-+)|([A-Za-z0-9]+\.+)|([A-Za-z0-9]+\++))*[A-Za-z0-9]+@((\w+\-+)|(\w+\.))*\w{1,63}\.[a-zA-Z]{2,6})$`)
	var dberr error
	db, dberr = sql.Open("postgres", fmt.Sprintf("dbname=%v sslmode=disable", dbname))
	if dberr != nil {
		log.Fatal(dberr)
	}
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS emails (email varchar(255) PRIMARY KEY)")
	if err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch {
	case *verbosity >= 2:
		log.Printf("%v", r)
	case *verbosity == 1:
		log.Printf("%v %v", r.Method, r.RequestURI)
	}
	email := r.PostFormValue("email")
	_ = insertEmail(email)
}

func getMails() []string {
	var mails []string
	rows, err := db.Query("SELECT email FROM emails")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}
		mails = append(mails, email)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return mails
}

func dumpMails(emails []string, sep string) {
	fmt.Print(strings.Join(emails, sep))
}

func insertEmail(email string) error {
	if *verbosity > 0 {
		log.Printf("inserting: %s", email)
	}
	email = strings.Trim(email, " ")
	email = strings.ToLower(email)
	if email != "" && emailregexp.MatchString(email) {
		_, err := db.Exec("INSERT INTO emails (email) VALUES ($1)", email)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid email address")
	}
	return nil
}

func main() {
	flag.Parse()
	initialize(*dbName)
	if *listen {
		http.HandleFunc("/", handleRequest)
		if *useTLS {
			log.Printf("Meiru listening on https://%v:%v [TLS]...\n", *address, *port)
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%v:%v", *address, *port), "cert.pem", "key.pem", nil))
		} else {
			log.Printf("Meiru listening on http://%v:%v...\n", *address, *port)
			log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", *address, *port), nil))
		}
	} else {
		dumpMails(getMails(), *separator)
	}
}
