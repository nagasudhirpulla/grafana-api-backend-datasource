import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { DataSourceOptions, GrafanaQuery } from './types';

export class DataSource extends DataSourceWithBackend<GrafanaQuery, DataSourceOptions> {
  url?: string;
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);
    this.url = instanceSettings.url;
  }
}
