package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/redmejia/internal/models"
)

func (app *App) AITextHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		var prompt models.TextPrompt
		err := json.NewDecoder(r.Body).Decode(&prompt)

		if err != nil {
			app.ErrorLog.Println("error decoding requested body", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Invalid request body",
			})
		}

		// get prompt by conversation ID
		textPrompt, conversationID, err := app.DB.GetPromptByConversationID(prompt.ConversationID)
		if err != nil {
			app.ErrorLog.Println("error getting prompt ", err)
		}

		app.InfoLog.Println("Prompt: I ASK :  ", textPrompt)

		// ctx := context.Background()
		// client, err := text.Client(ctx, app.GeminiKey)
		// if err != nil {
		// 	app.ErrorLog.Println("error creating genai client", err)
		// }
		// generatedText, err := text.GenTextContent(ctx, client, textPrompt)
		// if err != nil {
		// 	app.ErrorLog.Println("error generating text", err)
		// }

		// Simulate a delay to mimic the time it would take to generate a response
		time.Sleep(2 * time.Second)

		// Example result of prompt "what is a go routine?" use for testing
		// generatedText := "A Go routine is a lightweight, concurrent function that executes independently alongside other Go routines. Think of them as lightweight threads managed by the Go runtime. They are a core part of Go's concurrency model, enabling you to write programs that can perform multiple tasks simultaneously and efficiently.\n\nHere's a breakdown of what makes Go routines special:\n\n**Key Features:**\n\n* **Lightweight:**  Go routines are much cheaper to create and manage than traditional operating system threads. They consume far less memory and have lower overhead. You can easily spawn thousands (or even millions) of Go routines without overwhelming the system.\n* **Concurrent, not necessarily parallel:** Concurrency means that multiple tasks *can* make progress seemingly at the same time. Parallelism means that multiple tasks are *actually* running at the same time, typically on different cores. Go routines support concurrency naturally.  Whether they run in parallel depends on the number of available cores and the Go runtime's scheduler.\n* **Managed by the Go Runtime:** Go's runtime handles the complexities of scheduling, context switching, and managing Go routines. This frees the programmer from dealing with low-level thread management details.  The Go scheduler multiplexes Go routines onto a smaller number of OS threads.\n* **Communication via Channels:** Go encourages communication between Go routines using channels. Channels provide a safe and reliable way to exchange data and synchronize operations, preventing race conditions and other concurrency issues.\n* **Easy to Create:** Launching a Go routine is as simple as adding the `go` keyword before a function call.\n\n**How Go Routines Work:**\n\n1. **Scheduling:** The Go runtime has its own scheduler, called the \"M:N\" scheduler. This means that it multiplexes `M` Go routines onto `N` operating system threads.  Often, `N` will be equal to the number of available CPU cores, allowing for parallel execution.\n2. **Stack Growth:**  Go routines start with a small stack (e.g., a few kilobytes). As the routine needs more memory, the stack automatically grows and shrinks dynamically. This helps conserve memory.\n3. **Context Switching:** The Go runtime efficiently switches between Go routines.  Context switching happens more quickly than switching between OS threads.  The runtime scheduler can make scheduling decisions based on various factors, such as I/O operations or communication via channels.\n\n**Example:**\n\n```go\npackage main\n\nimport (\n\t\"fmt\"\n\t\"time\"\n)\n\nfunc sayHello(name string) {\n\tfor i := 0; i < 5; i++ {\n\t\tfmt.Println(\"Hello,\", name)\n\t\ttime.Sleep(100 * time.Millisecond) // Simulate some work\n\t}\n}\n\nfunc main() {\n\t// Start a Go routine to say hello to \"Alice\"\n\tgo sayHello(\"Alice\")\n\n\t// Start another Go routine to say hello to \"Bob\"\n\tgo sayHello(\"Bob\")\n\n\t// The main function also does some work\n\tfor i := 0; i < 3; i++ {\n\t\tfmt.Println(\"Main routine\")\n\t\ttime.Sleep(200 * time.Millisecond)\n\t}\n\n\t// Allow the Go routines some time to complete\n\ttime.Sleep(1 * time.Second)\n}\n```\n\n**Explanation:**\n\n* We define a function `sayHello` that prints a greeting with a given name.\n* In `main`, we use the `go` keyword to launch two instances of `sayHello` as separate Go routines.\n* The `main` function also performs its own work.\n* `time.Sleep` is added to allow the Go routines to run and print their greetings before the `main` function exits.  Without the `Sleep` in `main`, the program might terminate before the launched routines complete.\n\n**Output (order may vary):**\n\n```\nMain routine\nHello, Alice\nHello, Bob\nMain routine\nHello, Alice\nHello, Bob\nMain routine\nHello, Alice\nHello, Bob\nHello, Alice\nHello, Bob\nHello, Alice\nHello, Bob\n```\n\n**Benefits of Using Go Routines:**\n\n* **Improved Performance:**  Concurrency allows your program to utilize multiple CPU cores and execute tasks in parallel (when possible), leading to faster overall execution.\n* **Responsiveness:**  Go routines can prevent a single long-running task from blocking the entire program, keeping it responsive to user input and other events.\n* **Simplified Concurrency:** Go routines and channels provide a simple and elegant way to handle concurrent operations, reducing the risk of common concurrency bugs like race conditions and deadlocks.\n* **Scalability:** Go routines are easy to scale up as your application's workload increases. You can simply launch more Go routines to handle the additional tasks.\n\n**In summary, a Go routine is a lightweight, concurrent function managed by the Go runtime that enables efficient and safe concurrent programming.  It's a key building block for writing high-performance and scalable Go applications.**\n"

		generatedText := "This is a simulated response from the Gemini 2.0-flash model."
		// app.InfoLog.Println("Generated Text: ", generatedText)
		generatedResponse := app.DB.InsertGeneratedText("Gemini 2.0-flash", "Gemini_AI", conversationID, generatedText)
		// jsonByte, _ := json.Marshal(&prompt)
		// app.InfoLog.Println("PROMPT [>] ", string(jsonByte))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(generatedResponse)
	}
}
