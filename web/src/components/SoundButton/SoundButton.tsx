import { Sound } from '../../api';
import { Button, ButtonVariant } from '../Button';

type Props = {
  sound: Sound;
  active: boolean;
  playable?: boolean;
  activate: (sound: Sound) => void;
  openContext?: (e: React.MouseEvent<HTMLButtonElement, MouseEvent>, sound: Sound) => void;
};

export const SoundButton: React.FC<Props> = ({
  sound,
  playable = false,
  active,
  activate,
  openContext = () => {},
}) => {
  const _variant: ButtonVariant = active ? 'pink' : playable ? 'default' : 'gray';

  return (
    <Button
      onClick={() => playable && activate(sound)}
      onContextMenu={(e) => openContext(e, sound)}
      variant={_variant}>
      {sound.display_name || sound.uid}
    </Button>
  );
};
