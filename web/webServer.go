package web

import (
	"awesomeProject/cache"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type JSONResponse struct {
	Status string `json:"status"`
}

type WebServer struct {
	CacheStore *cache.InMemoryStore
}

func (ws *WebServer) getAllHandler(w http.ResponseWriter, r *http.Request) {
	orders := ws.CacheStore.GetAllOrders()

	ordersJSON, err := json.Marshal(orders)
	if err != nil {
		log.Printf("Error marshaling orders to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersJSON)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ws *WebServer) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/get_order/"):]

	order, _ := ws.CacheStore.GetOrder(id)

	ordersJSON, err := json.Marshal(order)
	if err != nil {
		log.Printf("Error marshaling orders to JSON: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ordersJSON)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (ws *WebServer) Run(wg *sync.WaitGroup) {
	defer wg.Done()

	staticDir := "./static"

	http.Handle("/", http.FileServer(http.Dir(staticDir)))

	http.HandleFunc("/get_all", ws.getAllHandler)
	http.HandleFunc("/get_order/", ws.getOrderHandler)

	port := ":7000"
	log.Printf("Server is running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
