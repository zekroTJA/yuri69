import { Sound } from '../../api';
import { Button } from '../Button';
import { Disableable } from '../props';

type Props = Disableable & {
  sound: Sound;
  active: boolean;
  activate: (sound: Sound) => void;
  openContext?: (e: React.MouseEvent<HTMLButtonElement, MouseEvent>, sound: Sound) => void;
};

export const SoundButton: React.FC<Props> = ({
  sound,
  disabled = false,
  active,
  activate,
  openContext = () => {},
}) => {
  return (
    <Button
      disabled={disabled}
      onClick={() => activate(sound)}
      onContextMenu={(e) => openContext(e, sound)}
      variant={active ? 'pink' : 'default'}>
      {sound.display_name || sound.uid}
    </Button>
  );
};
