import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { CodeEditor, DataSourceHttpSettings, InlineFieldRow, InlineLabel } from '@grafana/ui';
import AutoSizer from 'react-virtualized-auto-sizer';
import React, { ComponentType } from 'react';
import { DataSourceOptions } from './types';
import { JSONSchema7 } from 'json-schema';

type Props = DataSourcePluginOptionsEditorProps<DataSourceOptions>;

export const ConfigEditor: ComponentType<Props> = ({ options, onOptionsChange }) => {
  const schemaStr: string = (() => {
    try {
      return JSON.stringify(options.jsonData.schema, null, "\t")
    } catch (_) {
      return "";
    }
  })()
  return (
    <>
      <DataSourceHttpSettings
        defaultUrl={'http://localhost:8080'}
        dataSourceConfig={options}
        showAccessOptions={true}
        onChange={onOptionsChange}
      />
      <InlineFieldRow>
        <AutoSizer disableHeight>
          {({ width }) => (
            <div style={{ width: width + 'px' }}>
              <InlineLabel>Query Schema</InlineLabel>
              <CodeEditor
                width="100%"
                height="200px"
                language="json"
                showLineNumbers={true}
                showMiniMap={(schemaStr??"").length > 100}
                value={schemaStr}
                onBlur={(txt: string) => {
                  const schema: JSONSchema7 = (() => {
                    try {
                      return JSON.parse(txt)
                    } catch (_) {
                      return {}
                    }
                  })()
                  onOptionsChange({ ...options, jsonData: { ...options.jsonData, schema: schema } })
                }}
              />
            </div>
          )}
        </AutoSizer>
      </InlineFieldRow>
    </>
  );
};
