package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// sample implementation
// https://pkg.go.dev/github.com/grafana/grafana-plugin-sdk-go@v0.139.0/backend/datasource#example-package

var (
	_ backend.QueryDataHandler      = (*TestDataSource)(nil)
	_ backend.CheckHealthHandler    = (*TestDataSource)(nil)
	_ instancemgmt.InstanceDisposer = (*TestDataSource)(nil)
)

func NewDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	opts, err := setting.HTTPClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := httpclient.New(opts)
	if err != nil {
		return nil, err
	}

	return &TestDataSource{
		httpClient: client,
		baseUrl:    setting.URL,
	}, nil
}

func (s *TestDataSource) Dispose() {
	// Cleanup
}

type TestDataSource struct {
	httpClient *http.Client
	baseUrl    string
}

func (ds *TestDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// call the /test URL of the datasource
	log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	resp, err := ds.httpClient.Get(ds.baseUrl)
	if err != nil {
		status = backend.HealthStatusError
		message = "Unable to connect to datasource via get request"
	} else if resp.StatusCode != 200 {
		status = backend.HealthStatusError
		message = "datasource responded with status code " + strconv.Itoa(resp.StatusCode) + "instead of 200"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (ds *TestDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	resp := backend.NewQueryDataResponse()
	log.DefaultLogger.Info("Query Data called", "Request", req)
	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := ds.query(ctx, req.PluginContext, q)

		// save the response in a hashmap
		// based on with RefID as identifier
		resp.Responses[q.RefID] = res
	}
	// log.DefaultLogger.Info("Response generated", "Response", resp)
	return resp, nil
}

func (ds *TestDataSource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// create data frame response.
	postBody, err := json.Marshal(query)
	if err != nil {
		return response
	}
	_, e := ds.httpClient.Post(ds.baseUrl, "application/json", bytes.NewBuffer(postBody))
	if e != nil {
		return response
	}
	frame := data.NewFrame("response")

	// convert time to unix epoch and vice-versa - https://stackoverflow.com/questions/43915900/how-to-convert-unix-time-to-time-time-in-golang

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
		data.NewField("values", nil, []int64{10, 20}),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)

	return response
}
