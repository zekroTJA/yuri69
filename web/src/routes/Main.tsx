import { Outlet } from 'react-router-dom';
import styled from 'styled-components';
import { NavBar } from '../components/NavBar';

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
  return (
    <MainContainer>
      <NavBar />
      <main>
        <Outlet />
      </main>
    </MainContainer>
  );
};
