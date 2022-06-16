import styled from 'styled-components';
import { useNoEmbed } from '../../hooks/useNoEmbed';

type Props = {
  url: string;
};

const PreviewContainer = styled.div`
  display: flex;
  gap: 2em;
  margin-bottom: 1.5em;

  > img {
    height: 6em;
  }

  > div {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
`;

export const YouTubePreview: React.FC<Props> = ({ url }) => {
  const embed = useNoEmbed(url);
  return (
    (embed && (
      <PreviewContainer>
        <img src={embed.thumbnail_url} />
        <div>
          <h3>{embed.title}</h3>
          <span>
            by {embed.author_name} via {embed.provider_name}
          </span>
        </div>
      </PreviewContainer>
    )) || <></>
  );
};
