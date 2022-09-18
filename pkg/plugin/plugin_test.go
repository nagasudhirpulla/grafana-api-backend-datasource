package plugin_test

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/nagasudhirpulla/grafana-api-backend-datasource/pkg/plugin"
)

// "testing"

// This is where the tests for the datasource backend live.
func TestQueryData(t *testing.T) {
	ds := plugin.TestDataSource{}
	ds.BaseUrl = "http://localhost:8080"
	var err error
	ds.HttpClient, err = httpclient.New()
	if err != nil {
		t.Error(err)
	}

	resp, err := ds.QueryData(
		context.Background(),
		&backend.QueryDataRequest{
			PluginContext: backend.PluginContext{},
			Headers:       map[string]string{},
			Queries: []backend.DataQuery{{
				RefID:         "A",
				QueryType:     "",
				MaxDataPoints: 0,
				Interval:      0,
				TimeRange: backend.TimeRange{
					From: time.Now().Add(time.Duration(time.Duration.Minutes(-30))),
					To:   time.Now(),
				},
				JSON: []byte("{}"),
			}},
		},
	)
	if err != nil {
		t.Error(err)
	}

	if len(resp.Responses) != 1 {
		t.Fatal("QueryData must return a response")
	}
}
