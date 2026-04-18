package main

import (
	"database/sql"
	"fmt"
	"html/template"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"time"
)

type Msg struct {
	Name         string
	Content      string
	Timestamp    time.Time
	FormattedTime string
}

type TemplateData struct {
	Messages    []Msg
	Nickname   string
	NicknameAttr string
}

func init_db() {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	fmt.Println("Connected to the SQLite database successfully.")

	// Get the version of SQLite
	sql := `CREATE TABLE IF NOT EXISTS messages(
		nick VARCHAR(50),
		contents VARCHAR(1024),
		time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	result, err := db.Exec(sql)
	_ = result

	if err != nil {
		fmt.Println(err)
		return
	}
}

func HtmlSpecialchars(html string) string {
	return template.HTMLEscapeString(html)
}



func send(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()

	nickname := query.Get("nickname")
	if nickname == "" {
		nickname = "anon"
	}
	if len(nickname) > 50 {
		nickname = nickname[:50]
	}

	if query.Get("text") == "" {
		http.Redirect(w, req, fmt.Sprintf("/?nickname=%s", HtmlSpecialchars(nickname)), http.StatusSeeOther);
		return
	}

	message := query.Get("text")
	if len(message) > 1024 {
		message = message[:1024]
	}

	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error while opening database connection.")
		return
	}

	defer db.Close()
	fmt.Println("Connected to the SQLite database successfully.")

	stmt, err := db.Prepare("INSERT INTO messages (nick, contents, time) VALUES (?, ?, CURRENT_TIMESTAMP)")
	if err != nil {
		fmt.Println("Prepare error:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(nickname, message)
	if err != nil {
		fmt.Println("Error SQLITE:", err)
	}
	fmt.Println("sql executed")
	http.Redirect(w, req, fmt.Sprintf("/?nickname=%s", HtmlSpecialchars(nickname)), http.StatusSeeOther);
}

func getMessages() []Msg {
	db, err := sql.Open("sqlite3", "./chat.db")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error while opening database connection.")
		return nil
	}

	defer db.Close()
	fmt.Println("Connected to the SQLite database successfully.")
	sql := "SELECT * FROM messages"
	rows, err := db.Query(sql)

	if err != nil {
		fmt.Println("Query error")
		return nil
	}

	var messages []Msg
	for rows.Next() {
		var msg Msg
		if err := rows.Scan(&msg.Name, &msg.Content, &msg.Timestamp); err != nil {
			fmt.Println("Scan error:", err)
		} else {
			msg.FormattedTime = msg.Timestamp.Format(time.DateTime)
			msg.Content = msg.Content
			messages = append(messages, msg)
		}
	}
	return messages
}
func root(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	nickname := query.Get("nickname")
	nicknameAttr := ""
	if nickname != "" {
		nicknameAttr = "class=\"dim\" readonly"
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := TemplateData{
		Messages:    getMessages(),
		Nickname:   nickname,
		NicknameAttr: nicknameAttr,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	init_db()
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/send", send)
	http.HandleFunc("/", root)
	log.Println("Listening on 0.0.0.0:7866")
	err := http.ListenAndServe(":7866", nil)
	if err != nil {
		log.Println("Error: ", err)
	}
}
