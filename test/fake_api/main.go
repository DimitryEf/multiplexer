package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	body := fmt.Sprintf("api %s", id)
	fmt.Printf("get request for %s\n", body)
	//if i, _ := strconv.Atoi(id); i == 5 {
	//	time.Sleep(500 * time.Millisecond)
	//}
	time.Sleep(500 * time.Millisecond)
	if _, err := fmt.Fprint(w, body); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{id:[0-9]+}", get)
	http.Handle("/", r)
	fmt.Println("Server is listening...")
	if err := http.ListenAndServe(":8081", nil); err != nil && err != http.ErrServerClosed {
		fmt.Println(err.Error())
	}
}
