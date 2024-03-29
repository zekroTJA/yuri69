import { useEffect, useReducer, useRef, useState } from 'react';
import { uid } from 'react-uid';
import styled, { useTheme } from 'styled-components';
import { APIError, Sound, TwitchPageState } from '../../api';
import { Embed } from '../../components/Embed';
import { InfoPanel } from '../../components/InfoPanel';
import { RouteContainer } from '../../components/RouteContainer';
import { SearchBar } from '../../components/SearchBar';
import { Smol } from '../../components/Smol';
import { SoundButton } from '../../components/SoundButton';
import { useApi } from '../../hooks/useApi';
import { useSounds } from '../../hooks/useSounds';
import { useStore } from '../../store';
import { possessive } from '../../util/format';

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

const StyledInfoPanel = styled(InfoPanel)`
  margin-bottom: 1em;

  p {
    margin: 0 0 0.4em 0;
  }
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
  const theme = useTheme();
  const [playing] = useStore((s) => [s.playing]);
  const [twitchState, setTwitchState] = useState<TwitchPageState>();
  const [limited, setLimited] = useState(false);
  const limitedTimerRef = useRef<ReturnType<typeof setTimeout>>();
  const [search, dispatchSearch] = useReducer(searchReducer, { show: false, filter: '' });
  const { filteredSounds: sounds, refetchSounds } = useSounds(
    search.filter,
    !!twitchState?.channel
      ? (order) => (c) => c.twitchPageSounds(order)
      : (_) => (_) => Promise.resolve([]),
  );

  const _activateSound = (s: Sound) => {
    fetch((c) => c.twitchPagePlay(s.uid))
      .then((res) => {
        if (res.ratelimit.remaining === 0) {
          setLimited(true);
          const reset = new Date(res.ratelimit.reset!).getTime() - Date.now();
          if (reset > 0) {
            if (limitedTimerRef.current) clearTimeout(limitedTimerRef.current);
            limitedTimerRef.current = setTimeout(() => setLimited(false), reset);
            console.log('timeout', reset);
          } else {
            setLimited(false);
          }
        }
      })
      .catch();
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
      .catch((err) => {
        if (err instanceof APIError && err.message === 'invalid scopes')
          window.location.assign('/');
      });
  }, []);

  useEffect(() => {
    if (!twitchState?.channel) return;
    refetchSounds();
  }, [twitchState]);

  useEffect(() => {
    window.addEventListener('keydown', _onKeyDown);
    return () => window.removeEventListener('keydown', _onKeyDown);
  }, []);

  return twitchState?.channel ? (
    <RouteContainer>
      <SearchBar
        show={search.show}
        value={search.filter}
        onInput={(v) => dispatchSearch({ type: 'setFilter', payload: v })}
      />

      {(limited && (
        <StyledInfoPanel color={theme.red}>
          Oh snap, you have been rate limited! Please wait a short while until you can play sounds
          again.
        </StyledInfoPanel>
      )) || (
        <StyledInfoPanel>
          <p>Click a button to play a sound in {possessive(twitchState.channel)} stream.</p>
          <Smol>
            💡 Pro Tip: Press <Embed>CTRL</Embed> + <Embed>F</Embed> to pop up a search!
          </Smol>
        </StyledInfoPanel>
      )}

      <ButtonsContainer>
        {sounds?.map((s) => (
          <SoundButton
            key={uid(s)}
            sound={s}
            activate={_activateSound}
            active={s.uid === playing}
            playable={!limited}
          />
        ))}
      </ButtonsContainer>
    </RouteContainer>
  ) : (
    <RouteContainer center>
      <h3>You are nor connected to any twitch chat.</h3>
      <p>Please connect to the desired twitch chat and then reload the page.</p>
      <Smol>
        If you are indeed connected to the stream, simply write something in the chat so that Yuri
        can find you more easily. 😉
      </Smol>
    </RouteContainer>
  );
};

