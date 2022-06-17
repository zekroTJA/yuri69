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
  const filters = filter
    .toLowerCase()
    .split(',')
    .map((f) => f.trim())
    .filter((f) => !!f);

  const checkFuncs = filters.map((filter) => {
    if (filter.includes('*')) {
      const rx = new RegExp('^' + filter.replaceAll('*', '.*') + '$');
      return (v: string) => rx.test(v);
    }
    return (v: string) => v.includes(filter);
  });

  console.log(checkFuncs);

  const check = (v: string) => !!checkFuncs.find((filter) => filter(v));

  return (sound: Sound) =>
    (!!sound.display_name && check(sound.display_name.toLowerCase())) ||
    check(sound.uid) ||
    sound.tags?.find((t) => check(t.toLowerCase()));
};
