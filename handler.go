package main

import (
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	"net/http"
	"encoding/json"
)

func clusterIsHealthyCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Full or partial degradation in serving requests from Elasticsearch",
		Name:             "Check Elasticsearch cluster health",
		PanicGuide:       "todo",
		Severity:         1,
		TechnicalSummary: "Elasticsearch cluster is not healthy. Details on __elasticsearch-mvp/__health-details",
		Checker:          healthChecker,
	}
}

func healthChecker() (string, error) {
	output, err := elasticClient.ClusterHealth().Do()
	if err != nil {
		return "Cluster is not healthy: ", err
	} else if output.Status != "green" {
		return fmt.Sprintf("Cluster is %v", output.Status), nil
	}

	return "Cluster is healthy", nil
}

func GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := healthChecker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func HealthDetails(writer http.ResponseWriter, req *http.Request) {
	output, err := elasticClient.ClusterHealth().Do()
	if err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
	out, err := json.Marshal(*output)
	if (err!=nil) {
		writer.Write([]byte(err.Error()))
	}
	writer.Write(out)
}
