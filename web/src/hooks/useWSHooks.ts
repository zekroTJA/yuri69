import { useEffect } from 'react';
import { ApiClientInstance } from '../instances';
import { useStore } from '../store';
import {
  Event,
  EventStatePayload,
  EventType,
  EventPlayer,
  EventVoiceJoinPayload,
  GuildFilters,
  SetVolumeRequest,
} from '../api';

export const useWSHooks = () => {
  const [setConnected, setJoined, setPlaying, setVolume, setFilters] = useStore((s) => [
    s.setConnected,
    s.setJoined,
    s.setPlaying,
    s.setVolume,
    s.setFilters,
  ]);

  const _eventHandler = (e: Event<any>) => {
    console.log('WS Event', e);

    switch (e.type) {
      case EventType.AuthOK: {
        const pl = e.payload as EventStatePayload;
        setConnected(pl.connected);
        setJoined(pl.joined);
        setVolume(pl.volume);
        setFilters(pl.filters);
        break;
      }

      case EventType.VoiceJoin: {
        const pl = e.payload as EventVoiceJoinPayload;
        setJoined(true);
        setVolume(pl.volume);
        setFilters(pl.filters);
        break;
      }

      case EventType.VoiceLeave:
        setJoined(false);
        break;

      case EventType.VoiceInit: {
        const pl = e.payload as EventVoiceJoinPayload;
        setConnected(true);
        setVolume(pl.volume);
        setFilters(pl.filters);
        break;
      }

      case EventType.VoiceDeinit:
        setConnected(false);
        break;

      case EventType.PlayStart: {
        const pl = e.payload as EventPlayer;
        setPlaying(pl.ident);
        break;
      }

      case EventType.PlayException:
      case EventType.PlayEnd:
        setPlaying(undefined);
        break;

      case EventType.GuildFilterUpdated: {
        const pl = e.payload as GuildFilters;
        setFilters(pl);
        break;
      }

      case EventType.VolumeUpdated: {
        const pl = e.payload as SetVolumeRequest;
        setVolume(pl.volume);
        break;
      }
    }
  };

  useEffect(() => {
    ApiClientInstance.onWsEvent = _eventHandler;
  });
};
