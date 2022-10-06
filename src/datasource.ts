import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';
import { JSONSchema7 } from 'json-schema';
import { DataSourceOptions, GrafanaQuery } from './types';

export class DataSource extends DataSourceWithBackend<GrafanaQuery, DataSourceOptions> {
  url?: string;
  schema?: JSONSchema7;
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);
    this.url = instanceSettings.url;
    this.schema = instanceSettings.jsonData.schema;
  }
}
