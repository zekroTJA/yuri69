import styled from 'styled-components';
import { Sound } from '../../api';
import { Button, ButtonVariant } from '../Button';

type Props = {
  sound: Sound;
  active: boolean;
  playable?: boolean;
  activate: (sound: Sound) => void;
  openContext?: (e: React.MouseEvent<HTMLButtonElement, MouseEvent>, sound: Sound) => void;
};

const StyledButton = styled(Button)<{ fav?: boolean }>`
  border-bottom: solid 0 ${(p) => p.theme.orange};
  ${(p) => p.fav && 'border-bottom-width: 0.3em'};
`;

export const SoundButton: React.FC<Props> = ({
  sound,
  playable = false,
  active,
  activate,
  openContext = () => {},
}) => {
  const _variant: ButtonVariant = active ? 'pink' : playable ? 'default' : 'gray';

  return (
    <StyledButton
      fav={sound._favorite}
      onClick={() => playable && activate(sound)}
      onContextMenu={(e) => openContext(e, sound)}
      variant={_variant}>
      {sound.display_name || sound.uid}
    </StyledButton>
  );
};
