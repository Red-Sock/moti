/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as GoogleProtobufTimestamp from "./google/protobuf/timestamp.pb";

type Absent<T, K extends keyof T> = { [k in Exclude<keyof T, K>]?: undefined };

type OneOf<T> =
  | { [k in keyof T]?: undefined }
  | (keyof T extends infer K
      ? K extends string & keyof T
        ? { [k in K]: T[K] } & Absent<T, K>
        : never
      : never);

export enum ConfigType {
  plain = "plain",
  verv = "verv",
  minio = "minio",
  pg = "pg",
  nginx = "nginx",
  kv = "kv",
}

export enum Format {
  yaml = "yaml",
  env = "env",
}

export enum SortType {
  default = "default",
  by_name = "by_name",
  by_updated_at = "by_updated_at",
}

export type Paging = {
  limit?: string;
  offset?: string;
};

export type Sort = {
  type?: SortType;
  desc?: boolean;
};

export type Node = {
  name?: string;
  value?: string;
  innerNodes?: Node[];
};

export type ConfigBase = {
  id?: number;
  name?: string;
  createdAt?: GoogleProtobufTimestamp.Timestamp;
  updatedAt?: GoogleProtobufTimestamp.Timestamp;
  configType?: ConfigType;
};

export type ConfigInfo = {
  configBase?: ConfigBase;
  versions?: string[];
};

type BasePatch = {
  fieldName?: string;
};

export type Patch = BasePatch &
  OneOf<{
    rename: string;
    updateValue: string;
    delete: boolean;
  }>;