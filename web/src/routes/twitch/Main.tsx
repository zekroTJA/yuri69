import { useEffect } from 'react';
import { Outlet } from 'react-router';
import styled from 'styled-components';
import { useApi } from '../../hooks/useApi';
import { useStore } from '../../store';

type Props = {};

const MainContainer = styled.div`
  height: 100%;
  display: flex;

  > main {
    width: 100%;

    @media screen and (orientation: portrait) {
      margin-left: 1.5em;
    }
  }
`;

export const TwitchMainRoute: React.FC<Props> = ({}) => {
  const [setLoggedIn] = useStore((s) => [s.setLoggedIn]);
  const fetch = useApi();

  useEffect(() => {
    fetch((c) => c.checkAuth()).then(() => setLoggedIn(true));
  }, []);

  return (
    <MainContainer>
      <main>
        <Outlet />
      </main>
    </MainContainer>
  );
};

