/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";
import * as MatreshkaApiMatreshkaCommon from "./matreshka_common.pb";


export type VersionRequest = Record<string, never>;

export type VersionResponse = {
  version?: string;
};

export type Version = Record<string, never>;

export type ListConfigsRequest = {
  paging?: MatreshkaApiMatreshkaCommon.Paging;
  sort?: MatreshkaApiMatreshkaCommon.Sort;
  searchPattern?: string;
};

export type ListConfigsResponse = {
  configs?: MatreshkaApiMatreshkaCommon.ConfigBase[];
  totalRecords?: string;
};

export type ListConfigs = Record<string, never>;

export type CreateConfigRequest = {
  configName?: string;
  configType?: MatreshkaApiMatreshkaCommon.ConfigType;
};

export type CreateConfigResponse = Record<string, never>;

export type CreateConfig = Record<string, never>;

export type GetConfigRequest = {
  configName?: string;
  version?: string;
  format?: MatreshkaApiMatreshkaCommon.Format;
};

export type GetConfigResponse = {
  config?: Uint8Array;
  info?: MatreshkaApiMatreshkaCommon.ConfigInfo;
};

export type GetConfig = Record<string, never>;

export type PatchConfigRequest = {
  configName?: string;
  version?: string;
  patches?: MatreshkaApiMatreshkaCommon.Patch[];
};

export type PatchConfigResponse = Record<string, never>;

export type PatchConfig = Record<string, never>;

export type SaveConfigRequest = {
  format?: MatreshkaApiMatreshkaCommon.Format;
  configName?: string;
  version?: string;
  config?: Uint8Array;
};

export type SaveConfigResponse = Record<string, never>;

export type SaveConfig = Record<string, never>;

export type GetConfigNodeRequest = {
  configName?: string;
  version?: string;
};

export type GetConfigNodeResponse = {
  root?: MatreshkaApiMatreshkaCommon.Node;
  versions?: string[];
};

export type GetConfigNode = Record<string, never>;

export type RenameConfigRequest = {
  configName?: string;
  newName?: string;
};

export type RenameConfigResponse = {
  newName?: string;
};

export type RenameConfig = Record<string, never>;

export type SubscribeOnChangesRequest = {
  subscribeConfigNames?: string[];
  unsubscribeConfigNames?: string[];
};

export type SubscribeOnChangesResponse = {
  configName?: string;
  timestamp?: number;
  patches?: MatreshkaApiMatreshkaCommon.Patch[];
};

export type SubscribeOnChanges = Record<string, never>;

export type DeleteConfigRequest = {
  configName?: string;
  configVersion?: string;
};

export type DeleteConfigResponse = Record<string, never>;

export type DeleteConfig = Record<string, never>;

export class MatreshkaApi {
  static Version(this:void, req: VersionRequest, initReq?: fm.InitReq): Promise<VersionResponse> {
    return fm.fetchRequest<VersionResponse>(`/api/version?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
  static ListConfigs(this:void, req: ListConfigsRequest, initReq?: fm.InitReq): Promise<ListConfigsResponse> {
    return fm.fetchRequest<ListConfigsResponse>(`/api/config/list`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static CreateConfig(this:void, req: CreateConfigRequest, initReq?: fm.InitReq): Promise<CreateConfigResponse> {
    return fm.fetchRequest<CreateConfigResponse>(`/api/config/create`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static SaveConfig(this:void, req: SaveConfigRequest, initReq?: fm.InitReq): Promise<SaveConfigResponse> {
    return fm.fetchRequest<SaveConfigResponse>(`/api/config/${req.configName}/Save`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static GetConfig(this:void, req: GetConfigRequest, initReq?: fm.InitReq): Promise<GetConfigResponse> {
    return fm.fetchRequest<GetConfigResponse>(`/api/config/${req.configName}?${fm.renderURLSearchParams(req, ["configName"])}`, {...initReq, method: "GET"});
  }
  static GetConfigNodes(this:void, req: GetConfigNodeRequest, initReq?: fm.InitReq): Promise<GetConfigNodeResponse> {
    return fm.fetchRequest<GetConfigNodeResponse>(`/api/config/nodes`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static PatchConfig(this:void, req: PatchConfigRequest, initReq?: fm.InitReq): Promise<PatchConfigResponse> {
    return fm.fetchRequest<PatchConfigResponse>(`/api/config/${req.configName}/patch`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
  static RenameConfig(this:void, req: RenameConfigRequest, initReq?: fm.InitReq): Promise<RenameConfigResponse> {
    return fm.fetchRequest<RenameConfigResponse>(`/api/config/${req.configName}/rename/${req.newName}`, {...initReq, method: "POST"});
  }
  static DeleteConfig(this:void, req: DeleteConfigRequest, initReq?: fm.InitReq): Promise<DeleteConfigResponse> {
    return fm.fetchRequest<DeleteConfigResponse>(`/api/config/${req.configName}/delete`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}