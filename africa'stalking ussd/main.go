package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	log.Printf("Server is running at port %s\n", port)

	// Setup routes
	http.HandleFunc("/status", serverStatus) // check on the server status of our ussd callback
	http.HandleFunc("/ussd", ussdCallBack)

	//run
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serverStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := map[string]string{
		"status_code": "200",
		"reason":      "up",
	}
	bytesRes, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "An error occurred", http.StatusInternalServerError)
	}
	w.Write(bytesRes)
}

func ussdCallBack(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")

	// receive formValue from AT
	sessionId := r.FormValue("sessionId")
	serviceCode := r.FormValue("serviceCode")
	phoneNumber := r.FormValue("phoneNumber")
	text := r.FormValue("text")

	fmt.Printf("%s,%s", sessionId, serviceCode)

	if len(text) == 0 { // this indicates the beginning of the project
		w.Write([]byte("CON What would you want to check \n1. My Account \n2. My Phone Number"))
		return
	} else {
		switch text {
		case "1":
			w.Write([]byte("CON Choose account information you want to view \n1. Account Number\n2. Account Balance"))
			return
		case "2":
			w.Write([]byte(fmt.Sprintf("END Your Phone Number is %s", phoneNumber)))
			return
		case "1*1":
			w.Write([]byte("END Your Account Number is 00001"))
			return
		case "1*2":
			w.Write([]byte("END Your Balance is ksh 20"))
			return
		default:
			w.Write([]byte("END Invalid input"))
			return
		}
	}
}
