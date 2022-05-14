import { useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import styled from 'styled-components';
import { NavBar } from '../components/NavBar';
import { useApi } from '../hooks/useApi';
import { useStore } from '../store';

type Props = {};

const MainContainer = styled.div`
  height: 100%;
  display: flex;

  > main {
    margin-left: 5.5em;
    width: 100%;

    @media screen and (orientation: portrait) {
      margin-left: 1.5em;
    }
  }
`;

export const MainRoute: React.FC<Props> = ({}) => {
  const [setLoggedIn] = useStore((s) => [s.setLoggedIn]);
  const fetch = useApi();

  useEffect(() => {
    fetch((c) => c.checkAuth())
      .then(() => setLoggedIn(true))
      .catch();
  }, []);

  return (
    <MainContainer>
      <NavBar />
      <main>
        <Outlet />
      </main>
    </MainContainer>
  );
};
