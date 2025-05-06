package pkg

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/joho/godotenv"
)

func WithMetrics(mainFunc func()) {
	start := time.Now()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	startHeapAlloc := mem.HeapAlloc

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	mainFunc()

	elapsed := time.Since(start)
	runtime.ReadMemStats(&mem)
	usedHeap := mem.HeapAlloc - startHeapAlloc

	fmt.Println("--- Resource Usage ---")
	fmt.Printf("Execution Time: %s\nHeap Memory Used: %.2f MiB\n",
		elapsed, float64(usedHeap)/1024/1024)
}
