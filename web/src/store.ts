import create from 'zustand';
import { GuildFilters, Sound } from './api';

type Store = {
  sounds: Sound[];
  setSounds: (sounds: Sound[]) => void;
  addSound: (sound: Sound) => void;
  removeSound: (sound: Sound) => void;

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
  sounds: [],
  setSounds: (sounds) => set({ sounds }),
  addSound: (sound) => set({ sounds: [sound, ...get().sounds] }),
  removeSound: (sound) => set({ sounds: [...get().sounds.filter((s) => s.uid !== sound.uid)] }),

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
