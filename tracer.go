package bop

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/opentracing/opentracing-go"
	"sourcegraph.com/sourcegraph/appdash"
	appdashot "sourcegraph.com/sourcegraph/appdash/opentracing"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

// Trace applies the given tracer to a handler in order to add tracing meta data to the request
func Trace(tracer opentracing.Tracer, subject string, h HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req, resp *Message) error {
		var sp opentracing.Span
		wireContext, err := tracer.Extract(opentracing.TextMap, opentracing.TextMapCarrier(req.Values))
		if err != nil {
			sp = opentracing.StartSpan(subject)
		} else {
			sp = opentracing.StartSpan(subject, opentracing.ChildOf(wireContext))
		}
		defer sp.Finish()
		if err := sp.Tracer().Inject(sp.Context(), opentracing.TextMap, opentracing.TextMapCarrier(req.Values)); err != nil {
			return err
		}
		return h(ctx, req, resp)
	}
}

// StartAppDash is a helper function to start an instance of appdash.
// It returns a tracer which can be used in a server
func StartAppDash() opentracing.Tracer {
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 0})
	if err != nil {
		log.Fatal(err)
	}
	collectorPort := l.Addr().(*net.TCPAddr).Port
	collectorAdd := fmt.Sprintf(":%d", collectorPort)

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashPort := 8700
	appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
	}
	fmt.Printf("To see your traces, go to %s/traces\n", appdashURL)

	// Start the web UI in a separate goroutine.
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatal(err)
	}
	tapp.Store = store
	tapp.Queryer = store
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
	}()

	return appdashot.NewTracer(appdash.NewRemoteCollector(collectorAdd))
}
