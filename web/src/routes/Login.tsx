import styled from 'styled-components';
import ImgAvatar from '../../assets/avatar.jpg';
import { Button } from '../components/Button';
import { ReactComponent as DcLogo } from '../assets/dc-logo.svg';
import { ReactComponent as TwitchLogo } from '../assets/twitch-logo.svg';
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

const ButtonsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 2em;
  padding: 2em;
  border: solid 0.2em ${(p) => p.theme.gray};
`;

const Avatar = styled.img`
  border-radius: 100%;
  width: 10em;
`;

export const LoginRoute: React.FC<Props> = ({}) => {
  const _onLoginDiscord = () => {
    window.location.assign(ApiClientInstance.loginUrl('discord'));
  };

  const _onLoginTwitch = () => {
    window.location.assign(ApiClientInstance.loginUrl('twitch'));
  };

  return (
    <LoginContainer>
      <Avatar src={ImgAvatar} />
      <ButtonsContainer>
        <Button color="#404eed" onClick={_onLoginDiscord}>
          <DcLogo /> Login with Discord
        </Button>
        <Button color="#9146FF" onClick={_onLoginTwitch}>
          <TwitchLogo /> Login with Twitch
        </Button>
      </ButtonsContainer>
      <Smol>
        Yuri only uses OAuth2 to authenticate you and to get a unique ID for you (from the User IDs
        of the OAuth providers).
      </Smol>
    </LoginContainer>
  );
};
