import create from 'zustand';

type Store = {
  order: string;
  setOrder: (order: string) => void;

  joined: boolean;
  setJoined: (joined: boolean) => void;
};

export const useStore = create<Store>((set, get) => ({
  order: '',
  setOrder: (order) => set({ order }),

  joined: false,
  setJoined: (joined) => set({ joined }),
}));
