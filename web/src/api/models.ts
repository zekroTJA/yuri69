export enum EventType {
  AuthOK = 'authok',
  AuthFailed = 'authpromptfailed',
  SoundCreated = 'soundcreated',
  SoundUpdated = 'soundupdated',
  SoundDeleted = 'sounddeleted',
  VolumeUpdated = 'volumeupdated',
  GuildFilterUpdated = 'guildfilterupdated',
  PlayStart = 'playstart',
  PlayEnd = 'playend',
  PlayStuck = 'playstuck',
  PlayException = 'playexception',
  VoiceJoin = 'voicejoin',
  VoiceLeave = 'voiceleave',
  VoiceInit = 'voiceinit',
  VoiceDeinit = 'voicedeinit',
  _Disconnected = '_disconnected',
  _Reconnected = '_reconnected',
}

export type Sound = {
  uid: string;
  display_name?: string;
  created?: string;
  creator_id?: string;
  tags?: string[];
};

export type GuildFilters = {
  include: string[];
  exclude: string[];
};

export type Status = {
  status: number;
  message: string;
};

export type AccessToken = {
  access_token: string;
  expires: string;

  expiresDate: Date;
};

export type CreateSoundRequest = Sound & {
  upload_id: string;
  normalize: boolean;
  overdrive: boolean;
};

export type UpdateSoundRequest = Sound;

export type UploadSoundResponse = {
  upload_id: string;
  deadline: string;
};

export type SetVolumeRequest = {
  volume: number;
};

export type FastTrigger = {
  fast_trigger: string;
};

export type Event<T> = {
  type: string;
  origin?: string;
  payload?: T;
};

export type EventAuthPromptPayload = {
  deadline: string;
  token_type: string;
};

export type EventAuthRequest = {
  token: string;
};

export type GuildInfo = {
  id: string;
  name: string;
  icon_url: string;
};

export type EventVoiceJoinPayload = {
  volume: number;
  filters: GuildFilters;
  guild: GuildInfo;
};

export type EventStatePayload = EventVoiceJoinPayload & {
  connected: boolean;
  joined: boolean;
  is_admin: boolean;
};

export type EventErrorPlayload = {
  code: number;
  message: string;
};

export type EventPlayer = {
  ident?: string;
  guild_id?: string;
  user_id?: string;
  error?: string;
};

export type PlaybackLogEntry = {
  id: string;
  ident: string;
  guild_id: string;
  user_id: string;
  timestamp: string;
};

export type PlaybackStats = {
  ident: string;
  count: number;
};

export type StateStats = {
  n_sounds: number;
  n_plays: number;
};

export type OTAToken = {
  deadline: string;
  token: string;
  qrcode_data: string;
};

export type User = {
  id: string;
  username: string;
  discriminator: string;
  avatar_url: string;
  is_owner: boolean;
};
