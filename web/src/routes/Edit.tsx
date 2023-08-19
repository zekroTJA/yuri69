import { CreateSoundRequest, Sound } from '../api';
import { useEffect, useReducer, useState } from 'react';
import { useNavigate, useParams } from 'react-router';

import { Button } from '../components/Button';
import { DiscordImage } from '../components/DiscordImage';
import { Embed } from '../components/Embed';
import { ReactComponent as IconDelete } from '..//assets/delete.svg';
import { Modal } from '../components/Modal';
import { RouteContainer } from '../components/RouteContainer';
import { Smol } from '../components/Smol';
import { SoundEditor } from '../components/SoundEditor';
import styled from 'styled-components';
import { useApi } from '../hooks/useApi';
import { useSnackBar } from '../hooks/useSnackBar';

type Props = {};

const Heading = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const CreatorContainer = styled.div`
  margin-bottom: 1.5em;

  > h4 {
    margin-bottom: 1em;
  }

  > div {
    display: flex;
    gap: 1em;

    > img {
      height: 3em;
    }

    > div {
      display: flex;
      flex-direction: column;
      gap: 0.2em;
      justify-content: center;
    }
  }
`;

const deleteReducer = (
  state: { show: boolean; sound?: Sound },
  action: { type: 'show' | 'hide'; payload?: Sound },
) => {
  switch (action.type) {
    case 'show':
      return { show: true, sound: action.payload };
    case 'hide':
      return { ...state, show: false };
    default:
      return state;
  }
};

export const EditRoute: React.FC<Props> = ({}) => {
  const { uid } = useParams();
  const [sound, setSound] = useState<CreateSoundRequest>();
  const fetch = useApi();
  const nav = useNavigate();
  const { show } = useSnackBar();
  const [remove, dispatchRemove] = useReducer(deleteReducer, { show: false, sound: undefined });

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

  const _deleteSound = () => {
    if (!remove.sound) return;
    fetch((c) => c.soundsDelete(remove.sound!)).then(() => {
      show(
        <span>
          Sound <Embed>{remove.sound!.uid}</Embed> has successfully been deleted.
        </span>,
        'success',
      );
      nav(-1);
    });
    dispatchRemove({ type: 'hide' });
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
      <Heading>
        <h1>Edit Sound</h1>
        <Button variant="red" onClick={() => dispatchRemove({ type: 'show', payload: sound })}>
          <IconDelete />
        </Button>
      </Heading>

      {sound && (
        <>
          {sound.creator && (
            <CreatorContainer>
              <h4>Original Creator</h4>
              <div>
                <DiscordImage src={sound.creator.avatar_url} round />
                <div>
                  <span>{sound.creator.username}</span>
                  <Smol>{sound.creator.id}</Smol>
                </div>
              </div>
            </CreatorContainer>
          )}
          <SoundEditor
            sound={sound}
            updateSound={setSound}
            onCancel={() => nav(-1)}
            onSave={_update}
          />
        </>
      )}

      <Modal
        show={remove.show}
        onClose={() => dispatchRemove({ type: 'hide' })}
        heading="Delete Sound"
        controls={
          <>
            <Button variant="red" onClick={_deleteSound}>
              Delete
            </Button>
            <Button variant="gray" onClick={() => dispatchRemove({ type: 'hide' })}>
              Cancel
            </Button>
          </>
        }>
        <span>
          Do you really want to delete the sound <Embed>{remove.sound?.uid}</Embed>?
        </span>
      </Modal>
    </RouteContainer>
  );
};
