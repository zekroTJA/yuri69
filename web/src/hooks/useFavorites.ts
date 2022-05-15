import { useEffect } from 'react';
import { useStore } from '../store';
import { useApi } from './useApi';

export const useFavorites = () => {
  const fetch = useApi();
  const [favorites, set, add, remove] = useStore((s) => [
    s.favorites,
    s.setFavorites,
    s.addFavorite,
    s.removeFavorite,
  ]);

  const addFavorite = (ident: string) => {
    fetch((c) => c.addFavorite(ident))
      .then(() => add(ident))
      .catch();
  };

  const removeFavorite = (ident: string) => {
    fetch((c) => c.removeFavorite(ident))
      .then(() => remove(ident))
      .catch();
  };

  useEffect(() => {
    if (!favorites || favorites.length === 0) {
      fetch((c) => c.favorites())
        .then(set)
        .catch();
    }
  }, []);

  return { favorites, addFavorite, removeFavorite };
};
