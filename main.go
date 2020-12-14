package main

import (
	"encoding/json"
	"fmt"
	"github.com/jordan-wright/email"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

type Person struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Err struct {
	E string `json:"error"`
}

var (
	emailPassword = os.Getenv("EMAIL_PASS")
	emailAdr      = os.Getenv("EMAIL")
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", mainHandler)
	fmt.Printf("server port: %v \n", port)
	err := http.ListenAndServe(":"+port, nil) // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	persons := new([]Person)
	err := json.NewDecoder(r.Body).Decode(persons)
	if err != nil {
		fmt.Fprintln(w, JSONError(err))
		log.Println("messages has been not sent: " + err.Error())
		return
	}
	err = SendMessages(*persons)
	if err != nil {
		fmt.Fprintln(w, JSONError(err))
		log.Println("messages has been not sent: " + err.Error())
		return
	}
	log.Printf("messages has been sent: %v \n", *persons)
}

func SendMessages(persons []Person) error {
	arr := make([]int, len(persons))
	arrShfl := make([]int, len(persons))
	for i := 0; i < len(persons); i++ {
		arr[i] = i
		arrShfl[i] = i
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(arrShfl), func(i, j int) {
		if (arr[i] == arrShfl[i]) || (arr[j] == arrShfl[j]) {
			arrShfl[i], arrShfl[j] = arrShfl[j], arrShfl[i]
		}
	})

	if arr[len(persons)-1] == arrShfl[len(persons)-1] {
		arrShfl[0], arrShfl[len(persons)-1] = arrShfl[len(persons)-1], arrShfl[0]
	}
	for i := 0; i < len(persons)-1; i++ {
		if arr[i] == arrShfl[i] {
			arrShfl[i], arrShfl[i+1] = arrShfl[i+1], arrShfl[i]
		}
	}

	for i := 0; i < len(persons); i++ {
		txt := persons[arrShfl[i]].Name + " (" + persons[arrShfl[i]].Email + ")"
		err := SendEmail(persons[arr[i]].Email, txt)
		if err != nil {
			return err
		}
	}

	return nil
}

func SendEmail(to, txt string) error {
	e := email.NewEmail()
	e.From = "Secret Santa <your.anonymous.sender.from.russia@gmail.com>"
	e.To = []string{to}
	e.Subject = "Secret Santa"
	e.Text = []byte(txt)
	//e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("",
		emailAdr,
		emailPassword,
		"smtp.gmail.com"))
	return err
}

func JSONError(err error) string {
	jsn, _ := json.Marshal(Err{err.Error()})
	return string(jsn)
}
