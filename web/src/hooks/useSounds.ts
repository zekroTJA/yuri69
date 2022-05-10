import { useEffect } from 'react';
import { useStore } from '../store';
import { useApi } from './useApi';

export const useSounds = () => {
  const fetch = useApi();

  const [sounds, setSounds, order] = useStore((s) => [s.sounds, s.setSounds, s.order]);

  useEffect(() => {
    fetch((c) => c.sounds(order))
      .then((sounds) => setSounds(sounds))
      .catch();
  }, [order]);

  return { sounds };
};
