import { useEffect, useState } from 'react';
import styled from 'styled-components';
import { CreateSoundRequest, Sound } from '../../api';
import { Button } from '../Button';
import { Input } from '../Input';
import { Smol } from '../Smol';

type Props = {
  sound: CreateSoundRequest;
  updateSound?: (s: CreateSoundRequest) => void;
  isNew?: boolean;
  disabled?: boolean;
  onSave?: () => void;
  onCancel?: () => void;
};

const SoundEditorContainer = styled.div`
  > label {
    display: block;
    margin-bottom: 0.2em;
  }

  > ${Input} {
    width: 100%;
    margin-bottom: 1em;
  }
`;

const Controls = styled.div`
  display: flex;
  gap: 1em;
  margin-top: 0.5em;

  > * {
    width: 100%;
  }
`;

const CheckboxControls = styled.div`
  display: flex;
  gap: 1em;
  margin: 0 0 1em 0;

  > span {
    width: 100%;
    > input {
      margin-right: 0.5em;
    }
  }
`;

export const SoundEditor: React.FC<Props> = ({
  sound,
  updateSound = () => {},
  isNew = false,
  disabled = false,
  onCancel = () => {},
  onSave = () => {},
}) => {
  const [tagsValue, setTagsValue] = useState('');

  const _update = (s: Partial<CreateSoundRequest>) => {
    updateSound({ ...sound, ...s });
  };

  const _valueToTags = (v: string) =>
    v
      .split(',')
      .map((t) => t.trim())
      .filter((t) => !!t);

  const _tagsToValue = (t: string[]) => t.join(', ');

  useEffect(() => {
    setTagsValue(_tagsToValue(sound.tags ?? []));
  }, [sound]);

  return (
    <SoundEditorContainer>
      <label htmlFor="uid">
        Unique ID <Smol>(required)</Smol>
      </label>
      <Input
        id="uid"
        value={sound.uid ?? ''}
        onInput={(e) => _update({ uid: e.currentTarget.value })}
      />
      <label htmlFor="displayname">Display Name</label>
      <Input
        id="displayname"
        value={sound.display_name ?? ''}
        placeholder={sound.uid ?? ''}
        onInput={(e) => _update({ display_name: e.currentTarget.value })}
      />
      <label htmlFor="tags">
        Tags <Smol>(comma separated)</Smol>
      </label>
      <Input
        id="tags"
        value={tagsValue}
        onInput={(e) => setTagsValue(e.currentTarget.value)}
        onBlur={(e) => _update({ tags: _valueToTags(e.currentTarget.value) })}
      />
      {isNew && (
        <CheckboxControls>
          <span>
            <Input
              id="cb-normalize"
              type="checkbox"
              checked={sound.normalize}
              onChange={(e) => _update({ normalize: e.currentTarget.checked })}
            />
            <label htmlFor="cb-normalize">Normalize Volume</label>
          </span>
          <span>
            <Input id="cb-overdrive" type="checkbox" disabled />
            <label htmlFor="cb-overdrive">Create Overdrive Version</label>
          </span>
        </CheckboxControls>
      )}
      <Controls>
        <Button variant="green" disabled={disabled || !sound.uid} onClick={() => onSave()}>
          {isNew ? 'Create' : 'Save'}
        </Button>
        <Button variant="gray" onClick={() => onCancel()}>
          Cancel
        </Button>
      </Controls>
    </SoundEditorContainer>
  );
};
