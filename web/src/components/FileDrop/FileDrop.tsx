import { useEffect, useRef, useState } from 'react';
import styled, { css } from 'styled-components';
import { ReactComponent as IconUpload } from '../..//assets/upload.svg';
import { ReactComponent as IconFile } from '../..//assets/file.svg';
import { ReactComponent as IconError } from '../..//assets/fileerror.svg';
import { Flex } from '../Flex';
import { byteFormatter } from 'byte-formatter';
import { Smol } from '../Smol';
import { Styled } from '../props';
import { useClipboardEvent } from '../../hooks/useClipboardEvent';

type Props = Styled & {
  file?: File;
  onFileInput?: (file: File) => void;
};

const FileDropContainer = styled.div<{ isError: boolean; isDragging: boolean }>`
  width: fit-content;
  height: fit-content;
  border: dashed 3px currentColor;
  padding: 1em;
  display: flex;
  align-items: center;
  gap: 1em;
  cursor: pointer;
  transition: all 0.2s ease;

  > svg {
    width: 3em;
    height: 3em;
  }

  ${(p) =>
    p.isDragging &&
    css`
      color: ${p.theme.green};
    `}

  ${(p) =>
    p.isError &&
    css`
      color: ${p.theme.red};
    `}
`;

export const FileDrop: React.FC<Props> = ({ file, onFileInput = () => {}, ...props }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string>();
  const [dragging, setDragging] = useState(false);

  const _fileInputChange: React.ChangeEventHandler<HTMLInputElement> = (e) =>
    _setFile((e.currentTarget.files ?? [])[0]);

  const _onDragOver: React.DragEventHandler<HTMLDivElement> = (e) => {
    e.stopPropagation();
    e.preventDefault();
    e.dataTransfer.dropEffect = 'copy';
    setDragging(true);
    setError(undefined);
  };

  const _onDragEnd: React.DragEventHandler<HTMLDivElement> = (e) => {
    setDragging(false);
  };

  const _onDrop: React.DragEventHandler<HTMLDivElement> = (e) => {
    e.preventDefault();
    setDragging(false);
    _setFile((e.dataTransfer.files ?? [])[0]);
  };

  const _onPaste = (e: ClipboardEvent) => {
    const item = (e.clipboardData?.items ?? [])[0];
    if (!item || item.kind !== 'file') return;
    const file = item.getAsFile();
    if (file) _setFile(file);
  };

  const _setFile = (file?: File) => {
    if (!!file) {
      try {
        onFileInput(file);
      } catch (e) {
        setError(e instanceof Error ? e.message : (e as string));
      }
    }
  };

  useEffect(() => {
    setError(undefined);
    setDragging(false);
  }, [file]);

  useClipboardEvent(_onPaste);

  const _info = !!error ? (
    <span>{error}</span>
  ) : !!file ? (
    <Flex direction="column" gap="0.5em">
      <strong>{file.name}</strong>
      <Smol>{byteFormatter(file.size)}</Smol>
    </Flex>
  ) : (
    <span>
      Either drop the file to be uploaded here, paste it from your clipboard or just click to select
      a file.
    </span>
  );

  return (
    <FileDropContainer
      onClick={() => fileInputRef.current?.click()}
      onDrop={_onDrop}
      onDragOver={_onDragOver}
      onDragExit={_onDragEnd}
      isError={!!error}
      isDragging={dragging}
      {...props}>
      {(!!error && <IconError />) || (!!file && <IconFile />) || <IconUpload />}
      {_info}
      <input ref={fileInputRef} type="file" hidden onChange={_fileInputChange} />
    </FileDropContainer>
  );
};
