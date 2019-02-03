package srv

import (
	"net/http"

	"github.com/davyj0nes/worker-example/metrics"
	"github.com/justinas/alice"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/v1/hidden", apiHandler)
	mux.Handle("/metrics", metrics.Handler())

	chain := alice.New(metrics.Measure).Then(mux)

	return chain
}
