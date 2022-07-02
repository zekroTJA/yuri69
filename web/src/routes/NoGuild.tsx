import styled from 'styled-components';
import { Button } from '../components/Button';
import { RouteContainer } from '../components/RouteContainer';
import { ApiClientInstance } from '../instances';
import { ReactComponent as IconLogout } from '../assets/logout.svg';

type Props = {};

const StyledButton = styled(Button)`
  margin-top: 2em;
`;

export const NoGuildRoute: React.FC<Props> = ({}) => {
  const _onLogout = () => {
    window.location.assign(ApiClientInstance.logoutUrl());
  };

  return (
    <RouteContainer center>
      <h3>Oh no. 😲</h3>
      <span>Looks like you are not sharing any guild with Yuri.</span>
      <StyledButton variant="red" onClick={_onLogout}>
        <IconLogout />
        Logout
      </StyledButton>
    </RouteContainer>
  );
};

