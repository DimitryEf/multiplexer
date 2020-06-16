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
	fmt.Fprint(w, body)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{id:[0-9]+}", get)
	http.Handle("/", r)
	fmt.Println("Server is listening...")
	http.ListenAndServe(":8081", nil)
}
