import { useSounds } from '../hooks/useSounds';
import { uid } from 'react-uid';
import { SoundButton } from '../components/SoundButton';
import styled, { css } from 'styled-components';
import { useStore } from '../store';
import { Sound } from '../api/models';
import { useApi } from '../hooks/useApi';
import { RouteContainer } from '../components/RouteContainer';
import { animation, Item, ItemParams, Menu, theme, useContextMenu } from 'react-contexify';
import { ReactComponent as IconDelete } from '../../assets/delete.svg';
import { ReactComponent as IconEdit } from '../../assets/edit.svg';
import { useNavigate } from 'react-router';
import { Modal } from '../components/Modal';
import { useState } from 'react';
import { Embed } from '../components/Embed';
import { Button } from '../components/Button';

const SOUNDS_MENU_ID = 'sounds-menu';

type Props = {};

const ButtonsContainer = styled(RouteContainer)`
  display: flex;
  gap: 0.7em;
`;

const StyledItem = styled(Item)<{ delete?: boolean }>`
  > div {
    display: flex;
    align-items: center;
    gap: 0.5em;

    > svg {
      width: 1.1em;
      height: 1.1em;
    }

    ${(p) =>
      p.delete &&
      css`
        color: ${p.theme.red} !important;
      `}
  }
`;

export const SoundsRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const { sounds } = useSounds();
  const [connected, playing] = useStore((s) => [s.connected, s.playing]);
  const { show } = useContextMenu({ id: SOUNDS_MENU_ID });
  const nav = useNavigate();
  const [showDeleteModal, setShowDeleteModal] = useState(false);
  const [toDelete, setToDelete] = useState<Sound>();

  const _activateSound = (s: Sound) => {
    fetch((c) => c.playersPlay(s.uid)).catch();
  };

  const _openSoundOptions = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>, sound: Sound) => {
    e.preventDefault();
    show(e, {
      props: { sound },
    });
  };

  const _onSoundEdit = ({ props }: ItemParams<{ sound: Sound }, any>) => {
    nav(`sounds/${props!.sound.uid}`);
  };

  const _onSoundDelete = ({ props }: ItemParams<{ sound: Sound }, any>) => {
    setToDelete(props!.sound);
    setShowDeleteModal(true);
  };

  const _deleteSound = () => {
    if (!toDelete) return;
    fetch((c) => c.soundsDelete(toDelete));
    setShowDeleteModal(false);
  };

  return (
    <>
      <ButtonsContainer>
        {sounds?.map((s) => (
          <SoundButton
            key={uid(s)}
            sound={s}
            activate={_activateSound}
            active={s.uid === playing}
            playable={connected}
            openContext={_openSoundOptions}
          />
        ))}
      </ButtonsContainer>

      <Menu id={SOUNDS_MENU_ID} theme={theme.dark} animation={animation.fade}>
        <StyledItem onClick={_onSoundEdit}>
          <IconEdit /> <span>Edit</span>
        </StyledItem>
        <StyledItem delete onClick={_onSoundDelete}>
          <IconDelete />
          <span>Delete</span>
        </StyledItem>
      </Menu>

      <Modal
        show={showDeleteModal}
        onClose={() => setShowDeleteModal(false)}
        heading="Delete Sound"
        controls={
          <>
            <Button variant="red" onClick={_deleteSound}>
              Delete
            </Button>
            <Button variant="gray" onClick={() => setShowDeleteModal(false)}>
              Cancel
            </Button>
          </>
        }>
        <span>
          Do you really want to delete the sound <Embed>{toDelete?.uid}</Embed>?
        </span>
      </Modal>
    </>
  );
};
