package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	Atomic_Counters "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Atomic_Counters"
	Buffered_Channels "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Buffered_Channels"
	Channels_unbuffered "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Channels_unbuffered"
	Cond_syncCond "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Cond_syncCond"
	Context_Cancellation "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Context_Cancellation"
	Goroutines_WaitGroup "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Goroutines_WaitGroup"
	Mutex "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Mutex"
	Pool_Once_Map "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Pool_Once_Map"
	RWMutex "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/RWMutex"
	Worker_Pool "github.com/shahinzaman102/Go_JumpStart/internal/concurrency/Worker_Pool"
)

// captureOutput redirects stdout to a buffer and returns it as string
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// helper to set header and write captured output to HTTP response
func writeDemoOutput(w http.ResponseWriter, runFunc func()) {
	// Output plain text safely
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// It tells the client (browser, API caller, etc.) that the server’s response body is plain text (UTF-8 encoded).
	fmt.Fprintln(w, captureOutput(runFunc))
}

// ==== Handlers ====
func GoroutinesWaitGroupHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Goroutines_WaitGroup.Run)
}

func ChannelsUnbufferedHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Channels_unbuffered.Run)
}

func BufferedChannelsHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Buffered_Channels.Run)
}

func MutexHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Mutex.Run)
}

func RWMutexHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, RWMutex.Run)
}

func WorkerPoolHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Worker_Pool.Run)
}

func AtomicCountersHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Atomic_Counters.Run)
}

func CondSyncCondHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Cond_syncCond.Run)
}

func PoolOnceMapHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Pool_Once_Map.Run)
}

func ContextCancellationHandler(w http.ResponseWriter, r *http.Request) {
	writeDemoOutput(w, Context_Cancellation.Run)
}
