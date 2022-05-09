import styled, { useTheme } from 'styled-components';
import { Entry } from './Entry';
import ImgAvatar from '../../../assets/avatar.jpg';
import { ReactComponent as IconOrder } from '../../../assets/order.svg';
import { useStore } from '../../store';

type Props = {};

const SidebarContainer = styled.div`
  position: fixed;
  width: 100%;
  height: 100%;
  pointer-events: none;
`;

const SidebarBackground = styled.div`
  position: fixed;
  z-index: -1;
  width: 100%;
  height: 100%;
  background-color: rgba(0 0 0 / 0);
  transition: all 0.2s ease;
`;

const EntryContainer = styled.nav`
  position: fixed;
  height: 100%;
  background-color: ${(p) => p.theme.background2};
  width: 4em;
  overflow: hidden;
  pointer-events: all;
  transition: all 0.2s ease;

  &:hover {
    width: 18em;

    ~ ${SidebarBackground} {
      background-color: rgba(0 0 0 / 40%);
    }
  }
`;

const Avatar = styled.img`
  width: 4em;
`;

export const Sidebar: React.FC<Props> = ({}) => {
  const [order, setOrder] = useStore((s) => [s.order, s.setOrder]);
  const theme = useTheme();

  return (
    <SidebarContainer>
      <EntryContainer>
        <Entry to="/" icon={<Avatar src={ImgAvatar} />} label="Yuri" />
        <Entry
          action={() => setOrder(order === 'created' ? 'name' : 'created')}
          icon={<IconOrder />}
          label={`Order by ${order === 'created' ? 'Name' : 'Date'}`}
          color={theme.green}
        />
        {/* <Entry /> */}
      </EntryContainer>
      <SidebarBackground />
    </SidebarContainer>
  );
};
