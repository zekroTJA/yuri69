import { useEffect, useState } from 'react';
import { Sound } from '../api';
import { useStore } from '../store';
import { useApi } from './useApi';

export const useSounds = () => {
  const fetch = useApi();

  const [order] = useStore((s) => [s.order]);
  const [sounds, setSounds] = useState<Sound[]>();

  useEffect(() => {
    fetch((c) => c.sounds(order))
      .then((sounds) => setSounds(sounds))
      .catch();
  }, [order]);

  return { sounds };
};
