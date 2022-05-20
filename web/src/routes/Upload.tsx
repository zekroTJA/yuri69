import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router';
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

export const UploadRoute: React.FC<Props> = ({}) => {
  const [file, setFile] = useState<File>();
  const [sound, setSound] = useState({ normalize: true } as CreateSoundRequest);
  const fetch = useApi();
  const nav = useNavigate();
  const { show } = useSnackBar();
  const [state, setState] = useState(0);
  const randomImageRef = useRef<string>(randomFrom(LOADINGS));

  const _create = async () => {
    if (!file || !sound) return;
    try {
      const req = { ...sound };
      setState(1);
      const res = await fetch((c) => c.soundsUpload(file));
      req.upload_id = res.upload_id;
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
    if (!!file && !sound.uid) setSound({ ...sound, uid: fileName(file.name) });
  }, [file]);

  return (
    <RouteContainer maxWidth="50em">
      <h1>Upload</h1>
      {(state === 0 && (
        <>
          <StyledFileDrop file={file} onFileInput={setFile} />
          <SoundEditor
            isNew
            disabled={!file}
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
  return name;
}
