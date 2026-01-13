package daemon

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ok")
	})

	fmt.Println("simpleDeploypod running on :7777")
	http.ListenAndServe(":7777", nil)
}
