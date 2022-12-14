package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
		HttpClient: client,
		BaseUrl:    setting.URL,
	}, nil
}

func (s *TestDataSource) Dispose() {
	// Cleanup
}

type TestDataSource struct {
	HttpClient *http.Client
	BaseUrl    string
}

func (ds *TestDataSource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// log.DefaultLogger.Info("CheckHealth called", "request", req)

	var status = backend.HealthStatusOk
	var message = "Data source is working"

	resp, err := ds.HttpClient.Get(ds.BaseUrl)
	if err != nil {
		status = backend.HealthStatusError
		message = "Unable to connect to datasource via get request"
	} else if resp.StatusCode != 200 && resp.StatusCode != 202 && resp.StatusCode != 204 {
		status = backend.HealthStatusError
		message = "datasource responded with status code " + strconv.Itoa(resp.StatusCode) + "instead of 200/202/204"
	}
	// closing the response after function exit
	defer resp.Body.Close()
	// TODO check if JSON schema format is ok
	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func (ds *TestDataSource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	resp := backend.NewQueryDataResponse()
	// log.DefaultLogger.Info("Query Data called", "Request", req)
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

type ApiRespColumn struct {
	Name   string
	Labels map[string]string
	Values []float64
}

type ApiRespFrame struct {
	Columns []ApiRespColumn
}

type ApiResponse struct {
	Frames []ApiRespFrame
}

func (ds *TestDataSource) query(_ context.Context, pCtx backend.PluginContext, query backend.DataQuery) backend.DataResponse {
	response := backend.DataResponse{}

	// create data frame response
	postBody, err := json.Marshal(query)
	if err != nil {
		log.DefaultLogger.Error("post request json marshal error", "error", err)
		return response
	}

	resp, e := ds.HttpClient.Post(ds.BaseUrl, "application/json", bytes.NewBuffer(postBody))
	if e != nil {
		log.DefaultLogger.Error("post request error", "error", err)
		defer resp.Body.Close()
		return response
	} else if resp.StatusCode != 200 {
		defer resp.Body.Close()
		return response
	}

	defer resp.Body.Close()
	var respData ApiResponse
	er := json.NewDecoder(resp.Body).Decode(&respData)
	if er != nil {
		log.DefaultLogger.Error("api response decoding error", "error", er)
		return response
	}
	// log.DefaultLogger.Info("api response received", "api response", respData)

	for fItr, respFrame := range respData.Frames {
		frame := data.NewFrame(fmt.Sprintf("frame_%d", fItr))
		for _, respCol := range respFrame.Columns {
			if respCol.Name == "@timestamp" {
				var timeVals []time.Time
				for _, t := range respCol.Values {
					timeVals = append(timeVals, time.UnixMilli(int64(t)))
				}
				frame.Fields = append(frame.Fields,
					data.NewField(respCol.Name, nil, timeVals),
				)
			} else {
				frame.Fields = append(frame.Fields,
					data.NewField(respCol.Name, nil, respCol.Values),
				)
			}
		}
		response.Frames = append(response.Frames, frame)
	}

	// add the frames to the response.
	return response
}
