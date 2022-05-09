export enum EventType {
  EventSoundCreated = 'soundcreated',
  EventSoundUpdated = 'soundupdated',
  EventSoundDeleted = 'sounddeleted',
  EventVolumeUpdated = 'volumeupdated',
  EventGuildFilterUpdated = 'guildfilterupdated',
  EventPlayStart = 'playstart',
  EventPlayEnd = 'playend',
  EventPlayStuck = 'playstuck',
  EventPlayException = 'playexception',
  EventVoiceJoin = 'voicejoin',
  EventVoiceLeave = 'voiceleave',
  EventVoiceInit = 'voiceinit',
  EventVoiceDeinit = 'voicedeinit',
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
  normalize: string;
  overdrive: string;
};

export type UpdateSoundRequest = Sound;

export type UploadSoundResponse = {
  upload_id: string;
  deadline: string;
};

export type SetSoundRequest = {
  volume: string;
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
};

export type EventErrorPlayload = {
  code: number;
  message: string;
};
