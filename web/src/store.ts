import create from 'zustand';

type Store = {
  order: string;
  setOrder: (order: string) => void;
};

export const useStore = create<Store>((set, get) => ({
  order: '',
  setOrder: (order) => set({ order }),
}));
