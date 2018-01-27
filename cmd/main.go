package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	var httpAddr = flag.String("listen", ":9000", "HTTP listen and serve address for service")
	http.HandleFunc("/", handler)
	http.ListenAndServe(*httpAddr, nil)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "https://thechive.files.wordpress.com/2017/10/dog-memes-36-photos-21.jpg?quality=85&strip=info&w=600")
}
