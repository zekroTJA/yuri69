import { useEffect, useState } from 'react';
import { Sound } from '../api';
import { useStore } from '../store';
import { useApi } from './useApi';

export const useSounds = (filter?: string) => {
  const fetch = useApi();

  const [sounds, setSounds, order] = useStore((s) => [s.sounds, s.setSounds, s.order]);
  const [filteredSounds, setFilteredSounds] = useState<Sound[]>();

  useEffect(() => {
    fetch((c) => c.sounds(order))
      .then((sounds) => setSounds(sounds))
      .catch();
  }, [order]);

  useEffect(() => {
    if (!filter) {
      setFilteredSounds(undefined);
      return;
    }

    const _filter = filter.toLowerCase();
    const _filteredSounds = sounds.filter(
      (s) =>
        s.uid.includes(_filter) ||
        s.display_name?.toLowerCase().includes(_filter) ||
        s.tags?.find((t) => t.toLowerCase().includes(_filter)),
    );
    setFilteredSounds(_filteredSounds);
  }, [sounds, filter]);

  return { sounds, filteredSounds: filteredSounds ?? sounds };
};
