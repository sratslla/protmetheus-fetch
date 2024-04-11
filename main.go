package main

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

func main() {
	prometheusURL := "http://192.168.88.129:30000"
	promClient, err := NewPrometheusClient(prometheusURL)
	if err != nil {
		fmt.Println("Error creating Prometheus client:", err)
		return
	}

	query := `sum(rate(container_cpu_usage_seconds_total{pod=~"frontend-.*", container = "", namespace="default"}[1m]))`
	result, err := QueryPrometheus(promClient, query)
	if err != nil {
		fmt.Println("Error querying Prometheus:", err)
		return
	}

	fmt.Println("Query result:", result)
}

func NewPrometheusClient(prometheusURL string) (v1.API, error) {
	client, err := api.NewClient(api.Config{
		Address: prometheusURL,
	})
	if err != nil {
		return nil, err
	}
	return v1.NewAPI(client), nil
}

func QueryPrometheus(api v1.API, query string) (model.Value, error) {
	result, warnings, err := api.Query(context.Background(), query, time.Now())
	if err != nil {
		return nil, err
	}
	if len(warnings) > 0 {
		fmt.Println("Warnings received during query execution:", warnings)
	}
	return result, nil
}
