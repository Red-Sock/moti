/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";
import * as BasicExampleApiMessages from "./messages.pb";


export class BasicExampleApi {
  static Hello(this:void, req: BasicExampleApiMessages.HelloRequest, initReq?: fm.InitReq): Promise<BasicExampleApiMessages.HelloResponse> {
    return fm.fetchRequest<BasicExampleApiMessages.HelloResponse>(`/api/hello?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
}