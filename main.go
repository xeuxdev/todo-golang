package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Todo struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type Response struct {
	Message string `json:"message"`
	Status string 	`json:"status"`
}

var todos []Todo


func main (){

	todos = []Todo{
		{Id: "1", Title: "Finish coding task"},
		{Id: "2", Title: "Go for a run"},
	}

	http.HandleFunc("/todos", getTodos)
	http.HandleFunc("/todos/add", addTodos)
	http.HandleFunc("/todos/delete", deleteTodo)
	http.HandleFunc("/todos/edit", editTodo)

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func addTodos(w http.ResponseWriter, r *http.Request){
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var response = Response {
		Message: "Todo with id " + todo.Id + "  added successfully",
		Status: http.StatusText(http.StatusCreated),
	}

	todos = append(todos, todo)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func deleteTodo(w http.ResponseWriter, r *http.Request){
	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}
	var response = Response {
		Message: "Todo with id " + id + " deleted successfully",
		Status: http.StatusText(http.StatusOK),
	}

	for i, todo := range todos {

		if todo.Id == id {
			todos = append(todos[:i], todos[i+1: ]...)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	
	http.Error(w, "Todo not found", http.StatusNotFound)
}


func editTodo(w http.ResponseWriter, r *http.Request){
	var updatedTodo Todo
	err := json.NewDecoder(r.Body).Decode(&updatedTodo)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if updatedTodo.Id == ""{
		http.Error(w, "Id is required as a param", http.StatusBadRequest)
		return
	}

	var response = Response {
		Message: "Todo with id " + updatedTodo.Id + " updated successfully",
		Status: http.StatusText(http.StatusOK),
	}

	for i, todo := range todos {
		if todo.Id == updatedTodo.Id {
			todos[i] = updatedTodo
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	http.Error(w, "Todo not found", http.StatusNotFound)

}
