import create from 'zustand';
import { GuildFilters } from './api';

type Store = {
  order: string;
  setOrder: (order: string) => void;

  connected: boolean;
  setConnected: (connected: boolean) => void;

  joined: boolean;
  setJoined: (joined: boolean) => void;

  playing: string | undefined;
  setPlaying: (playing: string | undefined) => void;

  volume: number;
  setVolume: (volume: number) => void;

  filters: GuildFilters | undefined;
  setFilters: (filters: GuildFilters | undefined) => void;
};

export const useStore = create<Store>((set, get) => ({
  order: 'created',
  setOrder: (order) => set({ order }),

  joined: false,
  setJoined: (joined) => set({ joined }),

  connected: false,
  setConnected: (connected) => set({ connected }),

  playing: undefined,
  setPlaying: (playing) => set({ playing }),

  volume: 0,
  setVolume: (volume) => set({ volume }),

  filters: undefined,
  setFilters: (filters) => set({ filters }),
}));
