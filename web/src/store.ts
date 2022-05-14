import create from 'zustand';
import { GuildFilters, GuildInfo, Sound } from './api';
import { SnackBarModel } from './components/SnackBar';

type Store = {
  wsDisconnected: boolean;
  setWsDisconnected: (wsDisconnected: boolean) => void;

  loggedIn: boolean;
  setLoggedIn: (loggedIn: boolean) => void;

  snackBar: SnackBarModel;
  setSnackBar: (snackBar: Partial<SnackBarModel>) => void;

  sounds: Sound[];
  setSounds: (sounds: Sound[]) => void;
  addSound: (sound: Sound) => void;
  removeSound: (sound: Sound) => void;
  updateSound: (sound: Sound) => void;

  order: string;
  setOrder: (order: string) => void;

  connected: boolean;
  setConnected: (connected: boolean) => void;

  guild: GuildInfo | undefined;
  setGuild: (guild: GuildInfo | undefined) => void;

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
  wsDisconnected: false,
  setWsDisconnected: (wsDisconnected) => set({ wsDisconnected }),

  loggedIn: false,
  setLoggedIn: (loggedIn) => set({ loggedIn }),

  snackBar: { show: false } as SnackBarModel,
  setSnackBar: (snackBar) => set({ snackBar: { ...get().snackBar, ...snackBar } }),

  sounds: [],
  setSounds: (sounds) => set({ sounds }),
  addSound: (sound) => set({ sounds: [sound, ...get().sounds] }),
  removeSound: (sound) => set({ sounds: [...get().sounds.filter((s) => s.uid !== sound.uid)] }),
  updateSound: (sound) =>
    set({ sounds: [sound, ...get().sounds.filter((s) => s.uid !== sound.uid)] }),

  order: 'created',
  setOrder: (order) => set({ order }),

  joined: false,
  setJoined: (joined) => set({ joined }),

  connected: false,
  setConnected: (connected) => set({ connected }),

  guild: undefined,
  setGuild: (guild) => set({ guild }),

  playing: undefined,
  setPlaying: (playing) => set({ playing }),

  volume: 0,
  setVolume: (volume) => set({ volume }),

  filters: undefined,
  setFilters: (filters) => set({ filters }),
}));
