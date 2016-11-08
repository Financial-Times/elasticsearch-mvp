package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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

	if elasticClient != nil {
		output, err := elasticClient.ClusterHealth().Do()
		if err != nil {
			return "Cluster is not healthy: ", err
		} else if output.Status != "green" {
			return fmt.Sprintf("Cluster is %v", output.Status), nil
		}
		return "Cluster is healthy", nil
	}

	return "Couldn't check the cluster's health.", errors.New("Couldn't establish connectivity.")
}

func connectivityHealthyCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Could not connect to Elasticsearch",
		Name:             "Check connectivity to the Elasticsearch cluster",
		PanicGuide:       "todo",
		Severity:         1,
		TechnicalSummary: "Connection to Elasticsearch cluster could not be created. Please check your AWS credentials.",
		Checker:          connectivityChecker,
	}
}

func connectivityChecker() (string, error) {
	if elasticClient == nil {
		return "", errors.New("Could not connect to elasticsearch, please check the application parameters/env variables, and restart the service.")
	}

	return "Successfully connected to the cluster", nil
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := healthChecker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

//HealthDetails returns the response from elasticsearch service /__health endpoint - describing the cluster health
func HealthDetails(writer http.ResponseWriter, req *http.Request) {

	if elasticClient == nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	output, err := elasticClient.ClusterHealth().Do()
	if err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var response []byte
	response, err = json.Marshal(*output)
	if err != nil {
		response = []byte(err.Error())
	}

	_, err = writer.Write(response)
	if err != nil {
		log.Errorf(err.Error())
	}
}
