import { QueryEditorProps } from '@grafana/data';
import { CodeEditor, InlineFieldRow, InlineLabel } from '@grafana/ui';
import { defaults } from 'lodash';
import React, { PureComponent } from 'react';
import AutoSizer from 'react-virtualized-auto-sizer';
import { DataSource } from './datasource';

import { DataSourceOptions, GrafanaQuery, defaultQuery } from './types';

type Props = QueryEditorProps<DataSource, GrafanaQuery, DataSourceOptions>;


export class QueryEditor extends PureComponent<Props>{
  onPayloadChange = (value: string) => {
    const { onChange, query } = this.props;
    onChange({ ...query, payload: value });
  };

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { payload } = query;
    return (
      <>
        <InlineFieldRow>
          <AutoSizer disableHeight>
            {({ width }) => (
              <div style={{ width: width + 'px' }}>
                <InlineLabel>Payload</InlineLabel>
                <CodeEditor
                  width="100%"
                  height="200px"
                  language="json"
                  showLineNumbers={true}
                  showMiniMap={payload.length > 100}
                  value={payload}
                  onBlur={(value: string) => this.onPayloadChange(value)}
                />
              </div>
            )}
          </AutoSizer>
        </InlineFieldRow>
      </>
    )
  }
}
