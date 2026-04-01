/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "../fetch.pb";


export type HelloRequest = Record<string, never>;

export type HelloResponse = Record<string, never>;

export type Hello = Record<string, never>;

export class BasicExampleApi {
  static Hello(this:void, req: HelloRequest, initReq?: fm.InitReq): Promise<HelloResponse> {
    return fm.fetchRequest<HelloResponse>(`/api/hello?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
}