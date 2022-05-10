import styled, { useTheme } from 'styled-components';
import { Entry } from './Entry';
import { useStore } from '../../store';
import { useApi } from '../../hooks/useApi';
import ImgAvatar from '../../../assets/avatar.jpg';
import { ReactComponent as IconOrder } from '../../../assets/order.svg';
import { ReactComponent as IconJoin } from '../../../assets/join.svg';
import { ReactComponent as IconLeave } from '../../../assets/leave.svg';
import { ReactComponent as IconStop } from '../../../assets/stop.svg';
import { ReactComponent as IconVolume } from '../../../assets/volume.svg';
import { ReactComponent as IconSettings } from '../../../assets/settings.svg';
import { ReactComponent as IconUpload } from '../../../assets/upload.svg';
import { Slider } from '../Slider';
import { debounce } from 'debounce';
import { useCallback } from 'react';

type Props = {};

const SidebarContainer = styled.div`
  position: fixed;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 10;
`;

const SidebarBackground = styled.div`
  position: fixed;
  z-index: -1;
  width: 100%;
  height: 100%;
  background-color: rgba(0 0 0 / 0);
  transition: all 0.2s ease 0.15s;
`;

const EntryContainer = styled.nav`
  position: fixed;
  height: 100%;
  background-color: ${(p) => p.theme.background2};
  width: 4em;
  overflow: hidden;
  pointer-events: all;
  transition: all 0.2s ease 0.15s;

  &:hover {
    width: 20em;

    ~ ${SidebarBackground} {
      background-color: rgba(0 0 0 / 40%);
    }
  }
`;

const Avatar = styled.img`
  width: 4em;
`;

const Spacer = styled.div`
  width: 100%;
  height: 0.6em;
`;

export const Sidebar: React.FC<Props> = ({}) => {
  const fetch = useApi();
  const [order, setOrder, connected, joined, playing, volume, setVolume] = useStore((s) => [
    s.order,
    s.setOrder,
    s.connected,
    s.joined,
    s.playing,
    s.volume,
    s.setVolume,
  ]);
  const theme = useTheme();

  const _setVolume = (v: number) => {
    setVolume(v);
    _updateVolume(v);
  };

  const _updateVolume = useCallback(
    debounce((v: number) => {
      fetch((c) => c.playersVolume(v)).catch();
    }, 250),
    [],
  );

  return (
    <SidebarContainer>
      <EntryContainer>
        <Entry to="/" icon={<Avatar src={ImgAvatar} />} label="Yuri" />
        <Spacer />
        <Entry
          action={() => setOrder(order === 'created' ? 'name' : 'created')}
          icon={<IconOrder />}
          label={`Order by ${order === 'created' ? 'Name' : 'Date'}`}
          color={theme.green}
        />
        <Entry
          action={() => fetch((c) => (joined ? c.playersLeave() : c.playersJoin())).catch()}
          icon={joined ? <IconLeave /> : <IconJoin />}
          label={joined ? 'Leave' : 'Join'}
          disabled={!connected}
          color={theme.orange}
        />
        <Entry
          action={() => fetch((c) => c.playersStop()).catch()}
          icon={<IconStop />}
          label="Stop"
          disabled={!connected || !playing}
          color={theme.pink}
        />
        <Entry
          icon={<IconVolume />}
          label={
            <Slider min={1} max={200} value={volume} onChange={_setVolume} disabled={!connected} />
          }
          disabled={!connected}
          color={theme.cyan}
        />
        <Spacer />
        <Entry icon={<IconUpload />} label="Upload" to="upload" color={theme.gray} />
        <Entry icon={<IconSettings />} label="Settings" to="settings" color={theme.gray} />
      </EntryContainer>
      <SidebarBackground />
    </SidebarContainer>
  );
};
