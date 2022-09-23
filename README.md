# Grafana API Backend Datasource Plugin

This is a Grafana Data Source Backend Plugin for getting data from HTTP API

Documentation on [Backend plugins](https://grafana.com/docs/grafana/latest/developers/plugins/backend/).

A data source backend plugin consists of both frontend and backend components.

## Pre-requisites for plugin development

### For Back-end
* Install Go
* Install Mage
* Install Go extension in VS code (optional)

### For Front-end
* Install node
* Install yarn

## Setting up the development environment

### Frontend

1. Install dependencies with ```yarn install```

2. Build plugin in development mode or run in watch mode with ```yarn dev``` or ```yarn watch```

3. Build plugin in production mode with ```yarn build```

### Backend

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

   ```bash
   go get -u github.com/grafana/grafana-plugin-sdk-go
   go mod tidy
   ```

2. Build backend plugin binaries for Linux, Windows and Darwin with ```mage -v``` or ```mage -v build:windows``` for building only for windows

3. List all available Mage targets for additional commands using ```mage -l```

## Learn more

- [Build a data source backend plugin tutorial](https://grafana.com/tutorials/build-a-data-source-backend-plugin)
- [Grafana documentation](https://grafana.com/docs/)
- [Grafana Tutorials](https://grafana.com/tutorials/) - Grafana Tutorials are step-by-step guides that help you make the most of Grafana
- [Grafana UI Library](https://developers.grafana.com/ui) - UI components to help you build interfaces using Grafana Design System
- [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/)
