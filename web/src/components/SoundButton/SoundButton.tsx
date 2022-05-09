import { useTheme } from 'styled-components';
import { Sound } from '../../api';
import { Button } from '../Button';
import { Disableable } from '../props';

type Props = Disableable & {
  sound: Sound;
  active: boolean;
  activate: (sound: Sound) => void;
};

export const SoundButton: React.FC<Props> = ({ sound, disabled = false, active, activate }) => {
  return (
    <Button
      disabled={disabled}
      onClick={() => activate(sound)}
      variant={active ? 'pink' : 'default'}>
      {sound.display_name || sound.uid}
    </Button>
  );
};
