import { useSounds } from '../hooks/useSounds';
import { uid } from 'react-uid';
import { SoundButton } from '../components/SoundButton';
import styled from 'styled-components';
import { useStore } from '../store';
import { Sound } from '../api/models';
import { useApi } from '../hooks/useApi';
import { RouteContainer } from '../components/RouteContainer';

type Props = {};

const ButtonsContainer = styled(RouteContainer)`
  display: flex;
  gap: 0.7em;
`;

export const SoundsRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const { sounds } = useSounds();
  const [connected, playing] = useStore((s) => [s.connected, s.playing]);

  const _activateSound = (s: Sound) => {
    fetch((c) => c.playersPlay(s.uid)).catch();
  };

  return (
    <ButtonsContainer>
      {sounds?.map((s) => (
        <SoundButton
          key={uid(s)}
          sound={s}
          activate={_activateSound}
          active={s.uid === playing}
          disabled={!connected}
        />
      ))}
    </ButtonsContainer>
  );
};
