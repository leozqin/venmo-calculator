package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Rhymond/go-money"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type LineItem struct {
	Item   string  `json:"item,omitempty"`
	Cost   float32 `json:"cost,omitempty"`
	Person string  `json:"person,omitempty"`
	Shared bool    `json:"shared,omitempty"`
}

type Receipt struct {
	LineItems []LineItem `json:"line_items,omitempty"`
	TotalCost float32    `json:"total_cost,omitempty"`
	Subtotal  float32    `json:"subtotal,omitempty"`
}

type SplitCheckPerson struct {
	Person string  `json:"person,omitempty"`
	Cost   float32 `json:"cost,omitempty"`
}

type SplitCheck struct {
	People []SplitCheckPerson `json:"people,omitempty"`
}

type ValidationError struct {
	Receipt Receipt
	Msg     string
}

func (e *ValidationError) Error() string {
	return e.Msg
}

func (r Receipt) split() SplitCheck {
	costMap := map[string]float32{}
	var shared float32 = 0.0

	for _, i := range r.LineItems {
		if i.Shared {
			shared += i.Cost
		} else {
			costMap[i.Person] += i.Cost
		}
	}

	var checkPeople []SplitCheckPerson
	var sharedSplit float32 = shared / float32(len(costMap))

	for person, cost := range costMap {
		proportion := (cost + sharedSplit) / r.Subtotal

		split := SplitCheckPerson{
			Person: person,
			Cost:   proportion * r.TotalCost,
		}

		checkPeople = append(checkPeople, split)
	}

	return SplitCheck{
		People: checkPeople,
	}
}

func (r Receipt) validate() error {
	if r.Subtotal > r.TotalCost {
		log.Printf("Subtotal: %v, Total: %v", r.Subtotal, r.TotalCost)
		return &ValidationError{
			Receipt: r,
			Msg:     "Subtotal cannot be more than TotalCost",
		}
	}

	if r.Subtotal <= 0 || r.TotalCost <= 0 {
		return &ValidationError{
			Receipt: r,
			Msg:     "Both Subtotal and TotalCost must be greater than zero",
		}
	}

	if len(r.LineItems) == 0 {
		return &ValidationError{
			Receipt: r,
			Msg:     "There must be at least one line item",
		}
	}

	var runningSubtotal float32

	for _, m := range r.LineItems {
		runningSubtotal += m.Cost
	}

	runningSubtotalMoney := money.NewFromFloat(float64(runningSubtotal), money.USD)
	subTotalMoney := money.NewFromFloat(float64(r.Subtotal), money.USD)

	matches, _ := runningSubtotalMoney.Equals(subTotalMoney)

	if !(matches) {
		return &ValidationError{
			Receipt: r,
			Msg:     fmt.Sprintf("The sum of line items (%v) does not match the Subtotal (%v)", runningSubtotalMoney.Display(), subTotalMoney.Display()),
		}
	}

	return nil

}

func hello(w http.ResponseWriter, r *http.Request) {
	var rc Receipt

	err := json.NewDecoder(r.Body).Decode(&rc)

	if err != nil {
		panic(err)
	}

	err = rc.validate()

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		log.Println(err.Error())
		return
	}

	response := json.NewEncoder(w)

	response.Encode(rc.split())
}

func main() {
	baseUrl, ok := os.LookupEnv("BASE_URL")
	if !ok {
		baseUrl = "http://localhost:8080"
	}
	log.Println(baseUrl)
	r := chi.NewRouter()

	r.Use(
		cors.Handler(
			cors.Options{
				AllowedOrigins:   []string{baseUrl},
				AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
				AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
				ExposedHeaders:   []string{"Link"},
				AllowCredentials: false,
				MaxAge:           300,
			},
		),
	)

	r.Route("/split", func(r chi.Router) {
		r.Post("/", hello)
	})

	http.ListenAndServe(":3333", r)
}
