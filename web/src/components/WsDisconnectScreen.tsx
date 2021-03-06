import styled from 'styled-components';
import PepeHands from '../../assets/pepehands.png';
import { useStore } from '../store';

type Props = {};

const WsDisconnectScreenContainer = styled.div<{ show: boolean }>`
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 15;
  background-color: rgba(0 0 0 / 75%);
  opacity: ${(p) => (p.show ? 1 : 0)};
  pointer-events: ${(p) => (p.show ? 'all' : 'none')};
  transition: all 0.2s ease;

  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  line-height: 1.8em;
  justify-content: center;
  gap: 3em;

  > img {
    width: 10em;
  }
`;

export const WsDisconnectScreen: React.FC<Props> = ({}) => {
  const [wsDisconnected, loggedIn] = useStore((s) => [s.wsDisconnected, s.loggedIn]);
  return (
    <WsDisconnectScreenContainer show={wsDisconnected && loggedIn}>
      <img src={PepeHands} />
      <span>
        The web socket connection has been disconnected.
        <br />I try my best to reconnect you as soon as possible.
        <br />
        Please stand by ...
      </span>
    </WsDisconnectScreenContainer>
  );
};
