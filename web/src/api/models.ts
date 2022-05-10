export enum EventType {
  AuthOK = 'authok',
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
}

export type Sound = {
  uid: string;
  display_name: string;
  created: string;
  creator_id: string;
  tags: string[];
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

export type EventVoiceJoinPayload = {
  volume: number;
  filters: GuildFilters;
};

export type EventStatePayload = EventVoiceJoinPayload & {
  connected: boolean;
  joined: boolean;
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
