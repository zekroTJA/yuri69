import { Status } from "./models";

export class APIError extends Error {
  constructor(private _res: Response, private _body?: Status) {
    super(_body?.message ?? "unknown");
  }

  get response() {
    return this._res;
  }

  get status() {
    return this._res.status;
  }

  get code() {
    return this._body?.status ?? 0;
  }
}
