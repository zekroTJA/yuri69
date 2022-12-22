import styled from 'styled-components';
import { GuildInfo } from '../../api';
import { DiscordImage } from '../DiscordImage';
import { Embed } from '../Embed';
import { Flex } from '../Flex';
import { Styled } from '../props';

type Props = Styled & {
  guild: GuildInfo;
};

const StyledDiscordImage = styled(DiscordImage)`
  width: 3em;
`;

const GuildTileContainer = styled.div`
  display: flex;
  gap: 0.5em;
  align-items: center;
  background-color: rgba(0 0 0 / 10%);
  width: fit-content;
  padding: 0.5em;
`;

export const GuildTile: React.FC<Props> = ({ guild, ...props }) => {
  return (
    <GuildTileContainer {...props}>
      <StyledDiscordImage src={guild.icon_url} />
      <Flex direction="column" gap="0.2em">
        <strong>{guild.name}</strong>
        <Embed>{guild.id}</Embed>
      </Flex>
    </GuildTileContainer>
  );
};
