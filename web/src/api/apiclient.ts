import { HttpClient } from './httpclient';
import {
  CreateSoundRequest,
  Event,
  FastTrigger,
  GuildFilters,
  Sound,
  Status,
  UploadSoundResponse,
} from './models';

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

  sound(id: string): Promise<Sound> {
    return this.req('GET', `sounds/${id}`);
  }

  sounds(order: string = 'created'): Promise<Sound[]> {
    return this.req('GET', `sounds?order=${order}`);
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

  usersSetFasttrigger(fasttrigger: FastTrigger): Promise<Status> {
    return this.req('POST', 'users/settings/fasttrigger', fasttrigger);
  }
}
