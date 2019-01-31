// Package helloworld provides a set of Cloud Function samples.
package helloworld

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func registerStackdriverExporter() {
	sd, err := stackdriver.NewExporter(stackdriver.Options{})
	if err != nil {
		log.Fatalf("Failed to create the Stackdriver exporter: %v", err)
	}
	// It is imperative to invoke flush before your main function exits
	defer sd.Flush()

	trace.RegisterExporter(sd)

	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// Register it as a metrics exporter
	view.RegisterExporter(sd)
	view.SetReportingPeriod(60 * time.Second)
}

// initialize function
func init() {
	registerStackdriverExporter()

}

// WhatIsMyIp is an HTTP Cloud Function returns client's ip address.
func WhatIsMyIp(w http.ResponseWriter, r *http.Request) {
	ctx, span := trace.StartSpan(r.Context(), "WhatIsMyIp")
	defer span.End()

	switch r.Method {
	case http.MethodGet:
		userIp := getUserIp(ctx, *r)
		_, span := trace.StartSpan(ctx, "write response")
		_, err := fmt.Fprintf(w, "%s\n", userIp)
		span.End()
		if err != nil {
			log.Fatalf("can't writer response to %+v\n", w)
		}
	case http.MethodPut:
		http.Error(w, "403 - Forbidden", http.StatusForbidden)
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}

}

func getUserIp(ctx context.Context, request http.Request) string {
	ctx, span := trace.StartSpan(ctx, "getUserIp")
	defer span.End()

	// X-Appengine-User-Ip used to pass user ip address.
	ips, ok := request.Header["X-Appengine-User-Ip"]
	if !ok {
		return "mystery"
	}

	if len(ips) > 1 {
		return fmt.Sprintf("many of them! %q", strings.Join(ips, ","))
	}

	return ips[0]
}
