import { useState } from 'react';
import styled from 'styled-components';
import { NoEmbedModel, useNoEmbed } from '../../hooks/useNoEmbed';
import { useClipboardEvent } from '../../hooks/useClipboardEvent';
import { Button } from '../Button';
import { Flex } from '../Flex';
import { ReactComponent as IconLink } from '../../../assets/link.svg';
import { useApi } from '../../hooks/useApi';
import { useStore } from '../../store';

type Props = {};

const UrlImportContainer = styled.div<{ show: boolean }>`
  position: fixed;
  top: 0;
  left: 4em;
  z-index: 5;
  width: 100%;
  height: fit-content;
  background-color: ${(p) => p.theme.background3};
  display: flex;
  transform: translateY(${(p) => (p.show ? '0' : '-15em')});
  opacity: ${(p) => (p.show ? '1' : '0')};
  transition: all 0.25s ease;

  img {
    height: 50vh;
    max-height: 15em;
  }

  svg {
    width: 4em;
    height: 4em;
    margin: auto 2em auto 3em;
  }

  > div {
    padding: 1.5em;
    display: flex;
    flex-direction: column;
  }
`;

const Controls = styled.div`
  display: flex;
  gap: 1em;
  margin-top: auto;
`;

export const UrlImport: React.FC<Props> = ({}) => {
  const [url, setUrl] = useState<string>();
  const embed = useNoEmbed(url);
  const fetch = useApi();
  const [connected] = useStore((s) => [s.connected]);

  // https://youtu.be/DczHdNwooY4
  useClipboardEvent((e) => {
    const data = e.clipboardData?.items[0];
    if (!data) return;
    if (data.kind !== 'string') return;
    data.getAsString((d) => {
      if (!d.startsWith('https://')) return;
      setUrl(d);
    });
  });

  const _close = () => setUrl(undefined);

  const _play = () => {
    if (!url) return;
    fetch((c) => c.playersPlayExternal(url)).catch();
    _close();
  };

  return (
    <UrlImportContainer show={!!url && !!embed}>
      {(embed && <img src={embed.thumbnail_url} />) || <IconLink />}
      <div>
        <h3>Do you want to play this medias audio?</h3>
        {(embed && (
          <>
            <h2>{embed.title}</h2>
            <span>
              by {embed.author_name} via {embed.provider_name}
            </span>
          </>
        )) || <h2>{url}</h2>}
        <Controls>
          <Button disabled={!connected} onClick={_play}>
            Play
          </Button>
          <Button variant="gray" onClick={_close}>
            Cancel
          </Button>
        </Controls>
      </div>
    </UrlImportContainer>
  );
};
