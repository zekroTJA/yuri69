import { HttpClient } from './httpclient';
import {
  CreateSoundRequest,
  Event,
  FastTrigger,
  GuildFilters,
  OTAToken,
  PlaybackLogEntry,
  PlaybackStats,
  Sound,
  StateStats,
  Status,
  UploadSoundResponse,
  User,
} from './models';
import { buildQueryParams } from './util';

export class APIClient extends HttpClient {
  private _onWsEvent: (e: Event<any>) => void = () => {};

  constructor() {
    super();
    this.connectWS((e) => {
      this._onWsEvent(e);
    });
  }

  set onWsEvent(handler: (e: Event<any>) => void) {
    this._onWsEvent = handler;
  }

  loginUrl(): string {
    return this.basePath('auth/login');
  }

  logoutUrl(): string {
    return this.basePath('auth/logout');
  }

  checkAuth(): Promise<Status> {
    return this.req('GET', 'auth/check');
  }

  getOTAToken(): Promise<OTAToken> {
    return this.req('GET', 'auth/ota/token');
  }

  sound(id: string): Promise<Sound> {
    return this.req('GET', `sounds/${id}`);
  }

  sounds(order: string = 'created'): Promise<Sound[]> {
    return this.req('GET', `sounds${buildQueryParams({ order })}`);
  }

  soundsUpload(file: File): Promise<UploadSoundResponse> {
    return this.req('PUT', 'sounds/upload', file);
  }

  soundsCreate(sound: CreateSoundRequest): Promise<Sound> {
    return this.req('POST', 'sounds/create', sound);
  }

  soundsUpdate(sound: Sound): Promise<Sound> {
    return this.req('POST', `sounds/${sound.uid}`, sound);
  }

  soundsDelete(sound: Sound): Promise<Status> {
    return this.req('DELETE', `sounds/${sound.uid}`);
  }

  playersJoin(): Promise<Status> {
    return this.req('POST', 'players/join');
  }

  playersLeave(): Promise<Status> {
    return this.req('POST', 'players/leave');
  }

  playersPlay(ident: string): Promise<Status> {
    return this.req('POST', `players/play/${ident}`);
  }

  playersPlayExternal(url: string): Promise<Status> {
    return this.req('POST', `players/play/external${buildQueryParams({ url })}`);
  }

  playersStop(): Promise<Status> {
    return this.req('POST', 'players/stop');
  }

  playersVolume(volume: number): Promise<Status> {
    return this.req('POST', 'players/volume', { volume });
  }

  guildsGetFilters(): Promise<GuildFilters> {
    return this.req('GET', 'guilds/filters');
  }

  guildsSetFilters(filters: GuildFilters): Promise<GuildFilters> {
    return this.req('POST', 'guilds/filters', filters);
  }

  usersGetFasttrigger(): Promise<FastTrigger> {
    return this.req('GET', 'users/settings/fasttrigger');
  }

  usersSetFasttrigger(fast_trigger: string): Promise<Status> {
    return this.req('POST', 'users/settings/fasttrigger', { fast_trigger });
  }

  statsLog(
    guildid: string = '',
    userid: string = '',
    ident: string = '',
    limit: number = 50,
    offset: number = 0,
  ): Promise<PlaybackLogEntry[]> {
    return this.req(
      'GET',
      `stats/log${buildQueryParams({ guildid, userid, ident, limit, offset })}`,
    );
  }

  statsCount(guildid: string = '', userid: string = ''): Promise<PlaybackStats[]> {
    return this.req('GET', `stats/count${buildQueryParams({ guildid, userid })}`);
  }

  statsState(): Promise<StateStats> {
    return this.req('GET', 'stats/state');
  }

  admins(): Promise<User[]> {
    return this.req('GET', 'admins');
  }

  setAdmin(id: string): Promise<User> {
    return this.req('PUT', `admins/${id}`);
  }

  removeAdmin(id: string): Promise<User[]> {
    return this.req('DELETE', `admins/${id}`);
  }
}
