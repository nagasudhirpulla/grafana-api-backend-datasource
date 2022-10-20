# <img src="src/img/logo.svg" alt="logo.svg" width="50px" height="50px"/> Grafana API Backend Datasource Plugin

* This is a Grafana Data Source Plugin for getting data from HTTP API
* User inputs in the query editor can be customized by setting JSON schema in the datasource configuration page
* Now end user can create queries in query editor without understanding JSON
* Alerts are also supported since it is a backend datasource plugin 

# Plugin installation
* Download the plugin from https://github.com/nagasudhirpulla/grafana-api-backend-datasource/releases/
* unzip and paste the folder in the data/plugins folder of Grafana installation
* Restart Grafana  
* A data source by the name `api-backend` should appear

## Create JSON schema for query payload
* You can generate JSON schema from example query payload JSON [here](https://easy-json-schema.github.io/) 

## Example JSON schema
* Paste the following in the Query schema text box of the json-backend datasource settings screen
```json
{
	"description": "Query Schema",
	"properties": {
		"fetchCached": {
			"default": true,
			"description": "Whether cached data fetching is ok",
			"title": "Fetch Cached?",
			"type": "boolean"
		},
		"measId": {
			"default": "",
			"description": "Measurement ID",
			"maxLength": 20,
			"minLength": 15,
			"title": "Measurement ID",
			"type": "string"
		},
		"samplingSecs": {
			"default": 60,
			"description": "Sampling Freq. in secs",
			"multipleOf": 2,
			"title": "Sampling Freq (in secs)",
			"type": "integer"
		},
		"samplingType": {
			"default": "snap",
			"enum": [
				"snap",
				"average",
				"max",
				"min",
				"sum"
			],
			"title": "Sampling Type",
			"type": "string"
		}
	},
	"title": "Query Editor GUI",
	"type": "object"
}
```

![plugin-settings-example.png](https://github.com/nagasudhirpulla/grafana-api-backend-datasource/raw/main/readme-img/plugin-settings-example.png)

* The grafana query editor for this datasource will look like the following

![query-inspector-example.png](https://github.com/nagasudhirpulla/grafana-api-backend-datasource/raw/main/readme-img/query-editor-example.png)

* It can be seen in the query inspector that the JSON payload is generated from user inputs and request is sent

![query-inspector-payload-example.png](https://github.com/nagasudhirpulla/grafana-api-backend-datasource/raw/main/readme-img/query-inspector-payload-example.png)

## End points to be implemented by the API
* The plugin requires these endpoints to be implemented by the API server
* API server should listen for GET requests at `/` and return 200 OK status code as a health check endpoint implementation
* API server should listen for POST requests at `/` to respond to queries with data. The POST request body will contain the query information like series name, start and end timestamps, query payload etc. The response should have timestamps and values for a single series.
* A sample API request that will sent to API server by grafana can be like
```json
{
    "RefID": "A",
    "QueryType": "",
    "MaxDataPoints": 974,
    "Interval": 60000000000,
    "TimeRange": {
        "From": "2022-10-08T00:00:00+05:30",
        "To": "2022-10-08T14:44:09.586+05:30"
    },
    "JSON": {
        "refId": "A",
        "alias": "abcd",
        "payload": "{\"fetchCached\":true,\"measId\":\"abcd\",\"samplingType\":\"snap\",\"samplingSecs\":60}",
        "bucketAggs": [
            {
                "field": "@timestamp",
                "id": "2",
                "settings": {
                    "interval": "auto"
                },
                "type": "date_histogram"
            }
        ],
        "datasource": "api-backend",
        "datasourceId": 10,
        "intervalMs": 60000,
        "maxDataPoints": 974,
        "metrics": [
            {
                "id": "1",
                "type": "count"
            }
        ],
        "query": "",
        "timeField": "@timestamp"
    }
}
```
* A sample API response can be like 
```json
{
	"frames":[
		{
			"columns":[
				{"name": "@timestamp", "values": [1665216413748, 1665219977028], "labels": null},
				{"name": "abcd", "values": [5, 10], "labels": null}
			]
		}
	]
}
```

The times in response should be UNIX epoch timestamps in seconds 

## Data flow
![json-api-backend.png](https://github.com/nagasudhirpulla/grafana-api-backend-datasource/raw/main/readme-img/json-api-backend.png)

## Developer Documentation
Developer documentation can be found [here](https://github.com/nagasudhirpulla/grafana-api-backend-datasource/blob/main/devDocs.md)