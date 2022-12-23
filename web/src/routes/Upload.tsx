import { useEffect, useRef, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';
import styled from 'styled-components';
import { CreateSoundRequest, Sound } from '../api';
import { Embed } from '../components/Embed';
import { FileDrop } from '../components/FileDrop';
import { RouteContainer } from '../components/RouteContainer';
import { SoundEditor } from '../components/SoundEditor';
import { useApi } from '../hooks/useApi';
import { useSnackBar } from '../hooks/useSnackBar';

import Loading0 from '../../assets/loading/0.gif';
import Loading1 from '../../assets/loading/1.gif';
import Loading2 from '../../assets/loading/2.gif';
import Loading3 from '../../assets/loading/3.gif';
import Loading4 from '../../assets/loading/4.gif';
import Loading5 from '../../assets/loading/5.gif';
import { randomFrom } from '../util/rand';
import { useSearchParams } from 'react-router-dom';
import { YouTubePreview } from '../components/YouTubePreview';
import { sanitizeUid } from '../util/uid';

const LOADINGS = [Loading0, Loading1, Loading2, Loading3, Loading4, Loading5];

type Props = {};

const StyledFileDrop = styled(FileDrop)`
  margin-bottom: 1.5em;
  width: 100%;
`;

const LoadingContainer = styled.div`
  > img {
    height: 14em;
  }
`;

const getSeconds = (v: string) => {
  const split = v.split(':').reverse();
  var secs = parseFloat(split[0]);
  if (split.length > 1) secs += parseInt(split[1]) * 60;
  if (split.length > 2) secs += parseInt(split[2]) * 3600;
  return secs;
};

export const UploadRoute: React.FC<Props> = ({}) => {
  const [searchParams] = useSearchParams();
  const [file, setFile] = useState<File>();
  const [sound, setSound] = useState({ normalize: true } as CreateSoundRequest);
  const fetch = useApi();
  const nav = useNavigate();
  const [youtubeUrl, setYoutubeUrl] = useState<string | null>();
  const { show } = useSnackBar();
  const [state, setState] = useState(0);
  const randomImageRef = useRef<string>(randomFrom(LOADINGS));

  const _create = async () => {
    if ((!youtubeUrl && !file) || !sound) return;
    try {
      const req = { ...sound };
      setState(1);
      if (youtubeUrl) {
        req.youtube = {
          url: youtubeUrl,
          start_time_seconds: getSeconds(sound._start_time_str ?? '0'),
          end_time_seconds: getSeconds(sound._end_time_str ?? '0'),
        };
      } else {
        const res = await fetch((c) => c.soundsUpload(file!));
        req.upload_id = res.upload_id;
      }
      setState(2);
      await fetch((c) => c.soundsCreate(req));
      nav(-1);
      show(
        <span>
          Sound <Embed>{req.uid}</Embed> has successfully been created.
        </span>,
        'success',
      );
    } catch {
    } finally {
      setState(0);
    }
  };

  useEffect(() => {
    setYoutubeUrl(searchParams.get('youtube_url'));
  }, [searchParams]);

  useEffect(() => {
    if (!!file && !sound.uid) setSound({ ...sound, uid: fileName(file.name) });
  }, [file]);

  return (
    <RouteContainer maxWidth="50em">
      {(youtubeUrl && <h1>Import from YouTube</h1>) || <h1>Upload</h1>}
      {(state === 0 && (
        <>
          {(youtubeUrl && <YouTubePreview url={youtubeUrl} />) || (
            <StyledFileDrop file={file} onFileInput={setFile} />
          )}
          <SoundEditor
            isNew
            isYouTube={!!youtubeUrl}
            disabled={!youtubeUrl && !file}
            sound={sound}
            updateSound={setSound}
            onCancel={() => nav(-1)}
            onSave={_create}
          />
        </>
      )) || (
        <LoadingContainer>
          <img src={randomImageRef.current} />
          <p>
            {state === 1 && 'Uploading sound ...'}
            {state === 2 && 'Processing sound ...'}
          </p>
        </LoadingContainer>
      )}
    </RouteContainer>
  );
};

function fileName(name: string): string {
  const i = name.lastIndexOf('.');
  if (i != -1) name = name.substring(0, i);
  return sanitizeUid(name);
}
