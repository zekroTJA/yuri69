import { useEffect, useState } from 'react';
import { useSnackBar } from './useSnackBar';

const apiUrl = (url: string) => `https://noembed.com/embed?url=${url}`;

export type NoEmbedModel = {
  width: number;
  author_name: string;
  author_url: string;
  version: string;
  provider_url: string;
  provider_name: string;
  thumbnail_width: number;
  thumbnail_url: string;
  height: number;
  thumbnail_height: number;
  html: string;
  url: string;
  type: string;
  title: string;
};

export const useNoEmbed = (url?: string) => {
  const [embed, setEmbed] = useState<NoEmbedModel>();
  const { show } = useSnackBar();

  const _clearEmbed = () => {
    setEmbed(undefined);
  };

  useEffect(() => {
    if (!url) return;
    _clearEmbed();
    window
      .fetch(apiUrl(url))
      .then((res) => res.json())
      .then((data) => {
        if (data.error) return;
        setEmbed(data);
      })
      .catch((err) => {
        show(`Failed getting URL embed: ${err}`);
      });
  }, [url]);

  return embed;
};
