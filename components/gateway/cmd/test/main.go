//nolint:mnd,nlreturn,exhaustruct,godot,errorlint // Ok for test.
package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// PdfJob represents a PDF job
type PdfJob struct {
	ID     string `json:"id"`
	Status string `json:"status"` // queued, processing, completed
	Path   string `json:"path,omitempty"`
}

// backendSubscriber subscribes to pdf.* and prints job updates
func backendSubscriber(ctx context.Context, i int) error {
	nc, err := nats.Connect("nats://127.0.0.1:3020")
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	// Ensure stream exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "PDF_JOBS",
		Subjects: []string{"pdf.*"},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return err
	}

	sub, err := js.QueueSubscribeSync(
		"pdf.*",
		"pdf-workers", // same group name for all workers
	)
	if err != nil {
		return err
	}

	log.Printf("Backend %v subscribed to pdf.*", i)

	for {
		select {
		case <-ctx.Done():
			log.Println(i, "Backend subscriber stopping")
			return nil
		default:
			msg, e := sub.NextMsg(500 * time.Millisecond)
			if e == nats.ErrTimeout {
				continue
			}
			if e != nil {
				return e
			}

			var job PdfJob
			if e = json.Unmarshal(msg.Data, &job); e != nil {
				log.Println("Invalid job:", e)
				continue
			}

			log.Printf("Backend %v received job: %+v\n", i, job)
			e = msg.Ack() // important for JetStream queue
			if e != nil {
				return e
			}
		}
	}
}

// simulateJobPublisher creates a job and updates it
func simulateJobPublisher(ctx context.Context) error {
	nc, err := nats.Connect("nats://127.0.0.1:3020")
	if err != nil {
		return err
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		return err
	}

	job := PdfJob{
		ID:     uuid.NewString(),
		Status: "queued",
	}

	data, _ := json.Marshal(job)

	_, err = js.Publish("pdf."+job.ID, data)
	if err != nil {
		return err
	}

	log.Println("Published queued job", job.ID)

	time.Sleep(2 * time.Second)

	job.ID = uuid.NewString()
	job.Status = "processing"
	data, _ = json.Marshal(job)

	_, err = js.Publish("pdf."+job.ID, data)
	if err != nil {
		return err
	}

	log.Println("Published processing job", job.ID)

	time.Sleep(2 * time.Second)

	job.ID = uuid.NewString()
	job.Status = "completed"
	job.Path = "s3://bucket/" + job.ID + ".pdf"

	data, _ = json.Marshal(job)

	_, err = js.Publish("pdf."+job.ID, data)
	if err != nil {
		return err
	}

	log.Println("Published completed job", job.ID)

	<-ctx.Done()

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// Give the server a moment to start
	time.Sleep(500 * time.Millisecond)

	// 2️⃣ Start backend subscriber
	for i := range 2 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := backendSubscriber(ctx, i); err != nil {
				log.Println("Backend subscriber error:", err)
			}
		}()
	}

	// 3️⃣ Simulate job creation
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := simulateJobPublisher(ctx); err != nil {
			log.Println("Job publisher error:", err)
		}
	}()

	// Wait for all goroutines to finish (or context to timeout)
	wg.Wait()
	log.Println("All goroutines finished, exiting")
}
