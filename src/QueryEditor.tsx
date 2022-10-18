import { QueryEditorProps } from '@grafana/data';
import { InlineFieldRow, InlineLabel, Input, Select, Switch } from '@grafana/ui';
import { defaults, get as lodashGet, setWith as lodashSetWith } from 'lodash';
import React, { PureComponent } from 'react';
import { DataSource } from './datasource';
import { DataSourceOptions, GrafanaQuery, defaultQuery } from './types';
// import { ObjectSchema } from 'react-schema-based-json-editor';
import { JSONSchema7 } from 'json-schema'

type Props = QueryEditorProps<DataSource, GrafanaQuery, DataSourceOptions>;

export class QueryEditor extends PureComponent<Props>{
  onPayloadChange = (value: string) => {
    const { onChange, query } = this.props;
    onChange({ ...query, payload: value });
  };
  onAliasChange: React.FormEventHandler<HTMLInputElement> = (ev) => {
    const { onChange, query } = this.props;
    onChange({ ...query, alias: ev.currentTarget.value });
  };

  iterateSchema = (s: JSONSchema7, propTree: string[] = [], titleTree: string[], compList: JSX.Element[] = []) => {
    let payloadStr = this.props.query.payload;
    // console.log(this.props.datasource.schema)
    const payload: object = (() => {
      try {
        return JSON.parse(payloadStr)
      } catch (_) { return {} }
    })()

    for (const objKey in s.properties) {
      // rejection criteria start
      if (!s.properties.hasOwnProperty(objKey)) { continue }

      const sProperty = s.properties[objKey];
      if (typeof sProperty === "boolean") { continue }

      const propType = sProperty.type
      if (typeof propType !== "string") { continue }
      // rejection criteria end

      const propVal = lodashGet(payload, [...propTree, objKey], sProperty.default)
      const propLabel = sProperty.title ?? objKey
      const onPayloadChange = this.onPayloadChange
      const isPropEnum: boolean = sProperty.hasOwnProperty("enum")
      if (isPropEnum) {
        if (["number", "string", "integer"].indexOf(propType) === -1) { continue }
        const propOpts = sProperty.enum;
        if (!Array.isArray(propOpts)) { continue }
        let selOpts = propOpts.map((v) => {
          let val: string | number = (propType === "string") ? `${v}` : parseFloat(`${v}`)
          return { label: "" + v, value: val }
        });
        const el = <InlineFieldRow>
          <InlineLabel tooltip={sProperty.description}>{[...titleTree, propLabel].join('.')}</InlineLabel>
          <Select
            options={selOpts}
            value={propVal}
            onChange={(v) => {
              const newObj = lodashSetWith(payload, [...propTree, objKey], v.value, Object)
              onPayloadChange(JSON.stringify(newObj))
            }}
          />
        </InlineFieldRow>
        compList.push(el)
      }
      else if (propType === "number" || propType === "integer" || propType === "string") {
        const el = <InlineFieldRow>
          <InlineLabel tooltip={sProperty.description}>{[...titleTree, propLabel].join('.')}</InlineLabel>
          <Input width={12} value={propVal}
            onChange={function (e) {
              let valueOk = true
              // console.log(payload)
              let newVal: string | number = e.currentTarget.value
              if (propType === "number") {
                newVal = parseFloat(newVal)
                valueOk = !isNaN(newVal)
              }
              else if (propType === "integer") {
                newVal = parseInt(newVal, 10)
                valueOk = !isNaN(newVal)
              }
              // console.log(newObj)
              if (valueOk) {
                const newObj = lodashSetWith(payload, [...propTree, objKey], newVal, Object)
                onPayloadChange(JSON.stringify(newObj))
              }
            }} />
        </InlineFieldRow>
        compList.push(el)
      }
      else if (propType === "boolean") {
        const el = <InlineFieldRow>
          <InlineLabel tooltip={sProperty.description}>{[...titleTree, propLabel].join('.')}</InlineLabel>
          <Switch value={propVal}
            onChange={(e) => {
              const newObj = lodashSetWith(payload, [...propTree, objKey], e.currentTarget.checked, Object)
              onPayloadChange(JSON.stringify(newObj))
            }} />
        </InlineFieldRow>
        compList.push(el)
      }
    }
    return compList;
  }

  render() {
    const query = defaults(this.props.query, defaultQuery);
    const { alias } = query;
    const schema = this.props.datasource.schema ?? {}
    return (
      <>
        <InlineFieldRow>
          <InlineLabel tooltip="If left blank, the field uses the name of the queried element.">Alias</InlineLabel>
          <Input width={12} value={alias} onChange={this.onAliasChange} />
        </InlineFieldRow>
        <>{this.iterateSchema(schema, [], [], [])}</>
      </>
    )
  }
}
