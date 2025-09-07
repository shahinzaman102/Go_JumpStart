package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// GoBasics demonstrates core Go features through HTTP output.
func GoBasics(w http.ResponseWriter, r *http.Request) {
	// --- Type Conversion ---
	priceStr := "199"
	price, err := strconv.Atoi(priceStr) // string → int
	if err != nil {
		fmt.Fprintln(w, "Error converting price:", err)
	} else {
		fmt.Fprintf(w, "Converted string '%s' → int: %d\n", priceStr, price)
	}

	discount := 15
	finalPrice := float64(price) * (1 - float64(discount)/100)
	fmt.Fprintf(w, "Final Price after %d%% discount: %.2f\n", discount, finalPrice)

	rating := 4.8
	fmt.Fprintf(w, "User rating %.1f stored as int stars: %d\n", rating, int(rating))

	arr := []string{"Go", "Python", "Rust"}
	uidx := uint(2)
	fmt.Fprintf(w, "Accessing array[%d] → %s\n", uidx, arr[uidx])

	fmt.Fprintln(w)

	// --- Switch ---
	day := 2
	switch day {
	case 1:
		fmt.Fprintln(w, "Switch Case: Friday")
	case 2:
		fmt.Fprintln(w, "Switch Case: Saturday")
	case 3:
		fmt.Fprintln(w, "Switch Case: Sunday")
	default:
		fmt.Fprintln(w, "Switch Case: Unknown Day")
	}

	fmt.Fprintln(w)

	// --- iota ---
	const (
		First = iota
		Second
		Third
	)
	fmt.Fprintf(w, "Basic iota: First=%d, Second=%d, Third=%d\n", First, Second, Third)

	const (
		FlagRead = 1 << iota
		FlagWrite
		FlagExecute
	)
	fmt.Fprintf(w, "Bit Flags: Read=%d, Write=%d, Execute=%d\n", FlagRead, FlagWrite, FlagExecute)

	fmt.Fprintln(w)

	// --- Command-line arguments ---
	fmt.Fprintf(w, "Command-line arguments:\n")
	for i, arg := range os.Args {
		fmt.Fprintf(w, "Arg %d: %s\n", i, arg)
	}
	if len(os.Args) <= 1 {
		fmt.Fprintln(w, "\nTip: run with args, e.g.: go run ./cmd/server arg1 arg2 arg3")
	}

	fmt.Fprintln(w)

	// --- Panic and Recover ---
	processPayment := func(amount float64) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(w, "Recovered from failure: %v\n", r)
			}
		}()

		fmt.Fprintln(w, "Connecting to payment gateway...")
		if amount <= 0 {
			panic("invalid payment amount")
		}

		_, err := os.Open("credentials.txt")
		if err != nil {
			panic("missing payment credentials file")
		}

		fmt.Fprintf(w, "Payment processed successfully: %.2f\n", amount)
	}

	fmt.Fprintln(w, "\n=== Payment Demo ===")
	processPayment(0)   // triggers panic, recovered
	processPayment(100) // may panic if file missing

	fmt.Fprintln(w)

	// --- Exit Codes (simulated for safety in HTTP) ---
	if exitVal := r.URL.Query().Get("exit"); exitVal != "" {
		code, err := strconv.Atoi(exitVal)
		if err != nil {
			http.Error(w, "Invalid exit code", http.StatusBadRequest)
			return
		}

		var msg string
		switch code {
		case 0:
			msg = "Exit code 0: Success"
		case 1:
			msg = "Exit code 1: Generic error"
		case 2:
			msg = "Exit code 2: Invalid input"
		case 5:
			msg = "Exit code 5: Config missing"
		default:
			msg = fmt.Sprintf("Exit code %d: Unknown reason", code)
		}

		fmt.Fprintln(w, msg)
		fmt.Fprintf(w, "Server would exit with code %d (skipped for safety)\n", code)
	}

	// Example test URLs
	fmt.Fprintln(w, "\n💡 Safe URLs to test exit codes:")
	fmt.Fprintln(w, "http://localhost:8080/go-basics?exit=0")
	fmt.Fprintln(w, "http://localhost:8080/go-basics?exit=1")
	fmt.Fprintln(w, "http://localhost:8080/go-basics?exit=2")
	fmt.Fprintln(w, "http://localhost:8080/go-basics?exit=5")
	fmt.Fprintln(w, "http://localhost:8080/go-basics?exit=999")

	fmt.Fprintln(w, "\nGo Basics Demo Complete.")
}
