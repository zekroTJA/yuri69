import { DiscordImage } from '../DiscordImage';
import { Embed } from '../Embed';
import { Smol } from '../Smol';
import { User } from '../../api';
import styled from 'styled-components';

type Props = {
  user: User;
};

const UserTileContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 1em;
  background-color: ${(p) => p.theme.background3};
  padding: 0.5em;
  width: 100%;

  > img {
    width: 3em;
  }

  > div {
    display: flex;
    flex-direction: column;
    gap: 0.2em;
  }
`;

export const UserTile: React.FC<Props> = ({ user }) => {
  return (
    <UserTileContainer>
      <DiscordImage src={user.avatar_url} round />
      <div>
        <strong>{user.username}</strong>
        <Smol>{user.id}</Smol>
      </div>
    </UserTileContainer>
  );
};

