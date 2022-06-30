import { useEffect, useReducer, useState } from 'react';
import { uid } from 'react-uid';
import styled from 'styled-components';
import { Sound, TwitchPageState } from '../../api';
import { RouteContainer } from '../../components/RouteContainer';
import { SearchBar } from '../../components/SearchBar';
import { SoundButton } from '../../components/SoundButton';
import { useApi } from '../../hooks/useApi';
import { useSounds } from '../../hooks/useSounds';

type Props = {};

const ButtonsContainer = styled.div`
  display: flex;
  position: relative;
  flex-wrap: wrap;
  gap: 0.7em;
  height: fit-content;
  overflow-y: auto;
  padding-bottom: 1.5em;
`;

const searchReducer = (
  state: { show: boolean; filter: string },
  action: { type: 'show' | 'hide' | 'toggle' } | { type: 'setFilter'; payload: string },
) => {
  switch (action.type) {
    case 'show':
      return { filter: '', show: true };
    case 'hide':
      return { filter: '', show: false };
    case 'toggle':
      return { filter: '', show: !state.show };
    case 'setFilter':
      return { show: true, filter: action.payload };
    default:
      return state;
  }
};

export const TwitchSoundsRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const [twitchState, setTwitchState] = useState<TwitchPageState>();
  const [search, dispatchSearch] = useReducer(searchReducer, { show: false, filter: '' });
  const { filteredSounds: sounds, refetchSounds } = useSounds(
    search.filter,
    !!twitchState?.channel
      ? (order) => (c) => c.twitchPageSounds(order)
      : (_) => (_) => Promise.resolve([]),
  );

  const _activateSound = (s: Sound) => {
    fetch((c) => c.twitchPagePlay(s.uid)).catch();
  };

  const _onKeyDown = (e: KeyboardEvent) => {
    if (e.ctrlKey && e.key === 'f') {
      e.preventDefault();
      dispatchSearch({ type: 'toggle' });
    } else if (e.key === 'Escape') {
      dispatchSearch({ type: 'hide' });
    }
  };

  useEffect(() => {
    fetch((c) => c.twitchPageState(), true)
      .then(setTwitchState)
      .catch();
  }, []);

  useEffect(() => {
    if (!twitchState?.channel) return;
    refetchSounds();
  }, [twitchState]);

  useEffect(() => {
    window.addEventListener('keydown', _onKeyDown);
    return () => window.removeEventListener('keydown', _onKeyDown);
  }, []);

  return (
    <RouteContainer>
      {twitchState?.channel ? (
        <div>
          <SearchBar
            show={search.show}
            value={search.filter}
            onInput={(v) => dispatchSearch({ type: 'setFilter', payload: v })}
          />

          <ButtonsContainer>
            {sounds?.map((s) => (
              <SoundButton
                key={uid(s)}
                sound={s}
                activate={_activateSound}
                active={false}
                playable={true}
              />
            ))}
          </ButtonsContainer>
        </div>
      ) : (
        <div>
          <p>You are nor connected to any twitch chat.</p>
          <p>Please connect to the desired twitch chat and then reload the page.</p>
        </div>
      )}
    </RouteContainer>
  );
};

