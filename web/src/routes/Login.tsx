import styled from 'styled-components';
import ImgAvatar from '../../assets/avatar.jpg';
import { Button } from '../components/Button';
import { ReactComponent as DcLogo } from '../assets/dc-logo.svg';
import { ReactComponent as TwitchLogo } from '../assets/twitch-logo.svg';
import { Smol } from '../components/Smol';
import { ApiClientInstance } from '../instances';
import { useApi } from '../hooks/useApi';
import { useEffect, useState } from 'react';
import { Spinner } from '../components/Spinner';

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
  const fetch = useApi();
  const [capabilities, setCapabilities] = useState<string[]>();

  const _onLoginDiscord = () => {
    window.location.assign(ApiClientInstance.loginUrl('discord'));
  };

  const _onLoginTwitch = () => {
    window.location.assign(ApiClientInstance.loginUrl('twitch'));
  };

  useEffect(() => {
    fetch((c) => c.loginCapabilities())
      .then(setCapabilities)
      .catch();
  }, []);

  return (
    <LoginContainer>
      <Avatar src={ImgAvatar} />
      <ButtonsContainer>
        {capabilities === undefined ? (
          <Spinner />
        ) : (
          <>
            {capabilities.includes('oauth2:discord') && (
              <Button color="#404eed" onClick={_onLoginDiscord}>
                <DcLogo /> Login with Discord
              </Button>
            )}
            {capabilities.includes('oauth2:twitch') && (
              <Button color="#9146FF" onClick={_onLoginTwitch}>
                <TwitchLogo /> Login with Twitch
              </Button>
            )}
          </>
        )}
      </ButtonsContainer>
      <Smol>
        Yuri only uses OAuth2 to authenticate you and to get a unique ID for you (from the User IDs
        of the OAuth providers).
      </Smol>
    </LoginContainer>
  );
};
