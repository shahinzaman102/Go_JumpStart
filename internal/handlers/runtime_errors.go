package handlers

import (
	"fmt"
	"net/http"
)

// RuntimeErrorsHandler provides examples of common Go runtime errors.
// Trigger using query params, e.g. /runtime-errors?type=divide
func RuntimeErrorsHandler(w http.ResponseWriter, r *http.Request) {
	errType := r.URL.Query().Get("type")

	defer func() {
		if rec := recover(); rec != nil {
			http.Error(w, fmt.Sprintf("Recovered from panic: %v", rec), http.StatusInternalServerError)
		}
	}()

	switch errType {
	case "divide":
		// Division by zero
		x := 10
		y := 0
		_ = x / y
	case "nilptr":
		// Nil pointer dereference
		var p *int
		_ = *p // Keep it as-is. It intentionally triggers the nil pointer error.
	case "outofbounds":
		// Indexing outside array/slice
		arr := []int{1, 2, 3}
		_ = arr[10]
	case "typeassert":
		// Invalid type assertion
		var i any = "hello"
		_ = i.(int)
	default:
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid type. Use one of: divide, nilptr, outofbounds, typeassert\n")
		return
	}

	// If panic didn’t happen (should not reach here normally)
	fmt.Fprintln(w, "No error occurred (unexpected).")
}
