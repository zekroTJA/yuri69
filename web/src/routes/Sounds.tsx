import { useSounds } from '../hooks/useSounds';
import { uid } from 'react-uid';
import { SoundButton } from '../components/SoundButton';
import styled, { css } from 'styled-components';
import { useStore } from '../store';
import { Sound } from '../api/models';
import { useApi } from '../hooks/useApi';
import { RouteContainer } from '../components/RouteContainer';
import {
  animation,
  Item,
  ItemParams,
  Menu,
  PredicateParams,
  Separator,
  theme,
  useContextMenu,
} from 'react-contexify';
import { ReactComponent as IconDelete } from '..//assets/delete.svg';
import { ReactComponent as IconEdit } from '..//assets/edit.svg';
import { ReactComponent as IconStar } from '..//assets/star.svg';
import { ReactComponent as IconUnstar } from '..//assets/unstar.svg';
import { useNavigate } from 'react-router';
import { Modal } from '../components/Modal';
import { useEffect, useReducer, useState } from 'react';
import { Embed } from '../components/Embed';
import { Button } from '../components/Button';
import { useSnackBar } from '../hooks/useSnackBar';
import { UrlImport } from '../components/UrlImport';
import { SearchBar } from '../components/SearchBar';
import { useFavorites } from '../hooks/useFavorites';

const SOUNDS_MENU_ID = 'sounds-menu';

type Props = {};

const SoundsRouteContainer = styled(RouteContainer)``;

const ButtonsContainer = styled.div`
  display: flex;
  position: relative;
  flex-wrap: wrap;
  gap: 0.7em;
  height: fit-content;
  overflow-y: auto;
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

export const SoundsRoute: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const [connected, playing, filters] = useStore((s) => [s.connected, s.playing, s.filters]);
  const { show: showCtx } = useContextMenu({ id: SOUNDS_MENU_ID });
  const nav = useNavigate();
  const { show } = useSnackBar();
  const [remove, dispatchRemove] = useReducer(deleteReducer, { show: false, sound: undefined });
  const [showSearch, setShowSearch] = useState(false);
  const [searchFilter, setSearchFilter] = useState('');
  const { filteredSounds: sounds } = useSounds(searchFilter);
  const { favorites, addFavorite, removeFavorite } = useFavorites();

  const _activateSound = (s: Sound) => {
    fetch((c) => c.playersPlay(s.uid)).catch();
  };

  const _openSoundOptions = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>, sound: Sound) => {
    e.preventDefault();
    showCtx(e, {
      props: { sound },
    });
  };

  const _onFavorize = ({ props }: ItemParams<{ sound: Sound }, any>) =>
    addFavorite(props!.sound.uid);

  const _onUnfavorize = ({ props }: ItemParams<{ sound: Sound }, any>) =>
    removeFavorite(props!.sound.uid);

  const _onSoundEdit = ({ props }: ItemParams<{ sound: Sound }, any>) => {
    nav(`sounds/${props!.sound.uid}`);
  };

  const _onSoundDelete = ({ props }: ItemParams<{ sound: Sound }, any>) => {
    dispatchRemove({ type: 'show', payload: props!.sound });
  };

  const _deleteSound = () => {
    if (!remove.sound) return;
    fetch((c) => c.soundsDelete(remove.sound!)).then(() =>
      show(
        <span>
          Sound <Embed>{remove.sound!.uid}</Embed> has successfully been deleted.
        </span>,
        'success',
      ),
    );
    dispatchRemove({ type: 'hide' });
  };

  const _hideSearch = () => {
    setShowSearch(false);
    setSearchFilter('');
  };

  const _onKeyDown = (e: KeyboardEvent) => {
    if (e.ctrlKey && e.key === 'f') {
      e.preventDefault();
      setShowSearch(true);
    } else if (e.key === 'Escape') {
      _hideSearch();
    }
  };

  const _isExcluded = (s: Sound) => {
    if (!s.tags || !filters?.exclude) return false;
    return !!filters.exclude.find((e) => s.tags!.includes(e));
  };

  const _hideCtxFavorize = ({ props }: PredicateParams<{ sound: Sound }, any>) =>
    !!props?.sound._favorite;
  const _hideCtxUnavorize = ({ props }: PredicateParams<{ sound: Sound }, any>) =>
    !props?.sound._favorite;

  useEffect(() => {
    window.addEventListener('keydown', _onKeyDown);
    return () => window.removeEventListener('keydown', _onKeyDown);
  }, []);

  const _sounds = sounds?.map((s) => ({
    ...s,
    _favorite: favorites.includes(s.uid),
    _exclude: _isExcluded(s),
  }));
  const _favs = _sounds.filter((s) => s._favorite);
  const _nonfavs = _sounds.filter((s) => !s._favorite);
  const _sortedSounds = [..._favs, ..._nonfavs];

  return (
    <>
      <SoundsRouteContainer>
        <UrlImport />
        <SearchBar show={showSearch} value={searchFilter} onInput={setSearchFilter} />

        <ButtonsContainer>
          {_sortedSounds?.map((s) => (
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
      </SoundsRouteContainer>

      <Menu id={SOUNDS_MENU_ID} theme={theme.dark} animation={animation.fade}>
        <StyledItem onClick={_onFavorize} hidden={_hideCtxFavorize}>
          <IconStar /> <span>Favorize</span>
        </StyledItem>
        <StyledItem onClick={_onUnfavorize} hidden={_hideCtxUnavorize}>
          <IconUnstar /> <span>Unfavorize</span>
        </StyledItem>
        <StyledItem onClick={_onSoundEdit}>
          <IconEdit /> <span>Edit</span>
        </StyledItem>
        <Separator />
        <StyledItem delete onClick={_onSoundDelete}>
          <IconDelete />
          <span>Delete</span>
        </StyledItem>
      </Menu>

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
    </>
  );
};
