import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router';
import styled from 'styled-components';
import { CreateSoundRequest, Sound } from '../api';
import { Embed } from '../components/Embed';
import { FileDrop } from '../components/FileDrop';
import { RouteContainer } from '../components/RouteContainer';
import { SoundEditor } from '../components/SoundEditor';
import { useApi } from '../hooks/useApi';
import { useSnackBar } from '../hooks/useSnackBar';

type Props = {};

const StyledFileDrop = styled(FileDrop)`
  margin-bottom: 1.5em;
  width: 100%;
`;

export const UploadRoute: React.FC<Props> = ({}) => {
  const [file, setFile] = useState<File>();
  const [sound, setSound] = useState({ normalize: true } as CreateSoundRequest);
  const fetch = useApi();
  const nav = useNavigate();
  const { show } = useSnackBar();

  const _create = async () => {
    if (!file || !sound) return;
    try {
      const req = { ...sound };
      const res = await fetch((c) => c.soundsUpload(file));
      req.upload_id = res.upload_id;
      await fetch((c) => c.soundsCreate(req));
      nav(-1);
      show(
        <span>
          Sound <Embed>{req.uid}</Embed> has successfully been created.
        </span>,
        'success',
      );
    } catch {}
  };

  useEffect(() => {
    if (!!file && !sound.uid) setSound({ ...sound, uid: fileName(file.name) });
  }, [file]);

  return (
    <RouteContainer maxWidth="50em">
      <h1>Upload</h1>
      <StyledFileDrop file={file} onFileInput={setFile} />
      <SoundEditor
        isNew
        disabled={!file}
        sound={sound}
        updateSound={setSound}
        onCancel={() => nav(-1)}
        onSave={_create}
      />
    </RouteContainer>
  );
};

function fileName(name: string): string {
  const i = name.lastIndexOf('.');
  if (i != -1) name = name.substring(0, i);
  return name;
}
