import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router';
import { CreateSoundRequest } from '../api';
import { Embed } from '../components/Embed';
import { RouteContainer } from '../components/RouteContainer';
import { SoundEditor } from '../components/SoundEditor';
import { useApi } from '../hooks/useApi';
import { useSnackBar } from '../hooks/useSnackBar';
import { useStore } from '../store';

type Props = {};

export const EditRoute: React.FC<Props> = ({}) => {
  const { uid } = useParams();
  const [sound, setSound] = useState<CreateSoundRequest>();
  const fetch = useApi();
  const nav = useNavigate();
  const { show } = useSnackBar();

  const _update = async () => {
    if (!sound) return;
    try {
      const res = await fetch((c) => c.soundsUpdate(sound));
      nav(-1);
      show(
        <span>
          Sound <Embed>{res.uid}</Embed> has successfully been updated.
        </span>,
        'success',
      );
    } catch {}
  };

  useEffect(() => {
    if (uid) {
      fetch((c) => c.sound(uid))
        .then((s) => {
          setSound(s as CreateSoundRequest);
        })
        .catch();
    }
  }, [uid]);

  return (
    <RouteContainer maxWidth="50em">
      <h1>Edit Sound</h1>
      {sound && (
        <SoundEditor
          sound={sound}
          updateSound={setSound}
          onCancel={() => nav(-1)}
          onSave={_update}
        />
      )}
    </RouteContainer>
  );
};
