package redirect

import (
	"github.com/spf13/cobra"
	"net/http"
)

// Serve starts the HTTP server which will do the redirect.
// shortUrl redirected to -> LongUrl (Original Destination)
var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the server that handles redirects",
	RunE:  RunServe,
}

func RunServe(cmd *cobra.Command, args []string) error {
	// Stub implementation to validate runtime constraints.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// FixMe:- Unhandled error
		w.Write([]byte("ok"))
	})

	return http.ListenAndServe("localhost:8080", http.DefaultServeMux)
}

func init() {
}
