import { HttpClient } from './httpclient';
import {
  ApiKey,
  CreateSoundRequest,
  Event,
  FastTrigger,
  GuildFilters,
  GuildInfo,
  ImportSoundsResult,
  OTAToken,
  PlaybackLogEntry,
  PlaybackStats,
  Sound,
  StateStats,
  Status,
  StatusWithReservation,
  TwitchPageState,
  TwitchSettings,
  TwitchState,
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

  loginCapabilities(): Promise<string[]> {
    return this.req('GET', 'auth/logincapabilities');
  }

  loginUrl(provider: 'discord' | 'twitch'): string {
    return this.basePath(`auth/oauth2/${provider}/login`);
  }

  logoutUrl(): string {
    return this.basePath('auth/logout');
  }

  soundDownloadUrl(uid: string): string {
    return this.basePath(`sounds/${uid}/download?accessToken=${this.accessToken}`);
  }

  allSoundsDownloadUrl(): string {
    return this.basePath(`sounds/downloadall?accessToken=${this.accessToken}`);
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

  favorites(): Promise<string[]> {
    return this.req('GET', 'users/settings/favorites');
  }

  addFavorite(ident: string): Promise<Status> {
    return this.req('PUT', `users/settings/favorites/${ident}`);
  }

  removeFavorite(ident: string): Promise<Status> {
    return this.req('DELETE', `users/settings/favorites/${ident}`);
  }

  apiKey(): Promise<ApiKey> {
    return this.req('GET', 'users/settings/apikey');
  }

  generateApiKey(): Promise<ApiKey> {
    return this.req('POST', 'users/settings/apikey');
  }

  removeApiKey(): Promise<ApiKey> {
    return this.req('DELETE', 'users/settings/apikey');
  }

  twitchState(): Promise<TwitchState> {
    return this.req('GET', 'users/settings/twitch/state');
  }

  setTwitchSettings(settings: TwitchSettings): Promise<Status> {
    return this.req('POST', 'users/settings/twitch/settings', settings);
  }

  joinTwitch(settings?: TwitchSettings): Promise<Status> {
    return this.req('POST', 'users/settings/twitch/join', settings);
  }

  leaveTwitch(): Promise<Status> {
    return this.req('POST', 'users/settings/twitch/leave');
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

  guilds(): Promise<GuildInfo[]> {
    return this.req('GET', `admins/guilds`);
  }

  removeGuild(id: string): Promise<any> {
    return this.req('DELETE', `admins/guilds/${id}`);
  }

  soundsImport(file: File): Promise<ImportSoundsResult> {
    return this.req('POST', 'sounds/import', file);
  }

  twitchPageState(): Promise<TwitchPageState> {
    return this.req('GET', 'twitch/state');
  }

  twitchPageSounds(order: string): Promise<Sound[]> {
    return this.req('GET', `twitch/sounds?order=${order}`);
  }

  twitchPagePlay(ident: string): Promise<StatusWithReservation> {
    return this.req('POST', `twitch/play/${ident}`);
  }

  twitchPagePlayRandom(): Promise<StatusWithReservation> {
    return this.req('POST', 'twitch/play/random');
  }
}
