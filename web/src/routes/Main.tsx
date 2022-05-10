import { Outlet } from 'react-router-dom';
import styled from 'styled-components';
import { Sidebar } from '../components/Sidebar';

type Props = {};

const MainContainer = styled.div`
  height: 100%;
  display: flex;

  > main {
    margin-left: 5.5em;
    width: 100%;
  }
`;

export const MainRoute: React.FC<Props> = ({}) => {
  return (
    <MainContainer>
      <Sidebar />
      <main>
        <Outlet />
      </main>
    </MainContainer>
  );
};

