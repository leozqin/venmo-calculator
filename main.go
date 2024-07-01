package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type LineItem struct {
	Item   string  `json:"item,omitempty"`
	Cost   float32 `json:"cost,omitempty"`
	Person string  `json:"person,omitempty"`
}

type Receipt struct {
	LineItems []LineItem `json:"line_items,omitempty"`
	TotalCost float32    `json:"total_cost,omitempty"`
	Subtotal  float32    `json:"subtotal,omitempty"`
}

type SplitCheckPerson struct {
	Person string `json:"person,omitempty"`
	Cost   string `json:"cost,omitempty"`
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

}

func (r Receipt) validate() (bool, error) {
	if r.Subtotal > r.TotalCost {
		return false, &ValidationError{
			Receipt: r,
			Msg:     "Subtotal cannot be more than TotalCost",
		}
	}

	if r.Subtotal <= 0 || r.TotalCost <= 0 {
		return false, &ValidationError{
			Receipt: r,
			Msg:     "Both Subtotal and TotalCost must be greater than zero",
		}
	}

	if len(r.LineItems) == 0 {
		return false, &ValidationError{
			Receipt: r,
			Msg:     "There must be at least one line item",
		}
	}

	var runningSubtotal float32

	for _, m := range r.LineItems {
		runningSubtotal += m.Cost
		log.Println(runningSubtotal)
	}

	if runningSubtotal-r.Subtotal != 0 {
		return false, &ValidationError{
			Receipt: r,
			Msg:     "The sum of line items does not match the Subtotal",
		}
	}

	return true, nil

}

func hello(w http.ResponseWriter, r *http.Request) {
	var rc Receipt

	err := json.NewDecoder(r.Body).Decode(&rc)

	if err != nil {
		panic(err)
	}

	status, err := rc.validate()

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	log.Println(status)
}

func main() {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/", hello)
	})

	http.ListenAndServe(":3333", r)
}
