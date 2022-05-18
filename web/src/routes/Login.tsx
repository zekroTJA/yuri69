import styled from 'styled-components';
import ImgAvatar from '../../assets/avatar.jpg';
import { Button } from '../components/Button';
import { ReactComponent as DcLogo } from '..//assets/dc-logo.svg';
import { Smol } from '../components/Smol';
import { ApiClientInstance } from '../instances';

type Props = {};

const LoginContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2em;
  align-items: center;
  justify-content: center;
  padding: 5em;
  text-align: center;

  > * {
    max-width: 18em;
  }
`;

const Avatar = styled.img`
  border-radius: 100%;
  width: 10em;
`;

export const LoginRoute: React.FC<Props> = ({}) => {
  const _onLogin = () => {
    window.location.assign(ApiClientInstance.loginUrl());
  };

  return (
    <LoginContainer>
      <Avatar src={ImgAvatar} />
      <Button onClick={_onLogin}>
        <DcLogo /> Login with Discord
      </Button>
      <Smol>
        No worries! Yuri will only rip off all your data from your Discord account by logging in. 😉
      </Smol>
    </LoginContainer>
  );
};
