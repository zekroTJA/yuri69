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

    const _filteredSounds = sounds.filter(soundsFilter(filter));
    setFilteredSounds(_filteredSounds);
  }, [sounds, filter]);

  return { sounds, filteredSounds: filteredSounds ?? sounds };
};

const soundsFilter = (filter: string) => {
  filter = filter.toLowerCase();

  let check: (v: string) => boolean;
  if (filter.includes('*')) {
    const rx = new RegExp('^' + filter.replaceAll('*', '.*') + '$');
    check = (v) => rx.test(v);
  } else {
    check = (v) => v.includes(filter);
  }

  return (sound: Sound) =>
    (!!sound.display_name && check(sound.display_name.toLowerCase())) ||
    check(sound.uid) ||
    sound.tags?.find((t) => check(t.toLowerCase()));
};
