import { useEffect, useRef, useState } from 'react';
import styled from 'styled-components';
import { Input } from '../Input';

type Props = {
  show?: boolean;
  value?: string;
  onInput?: (v: string) => void;
};

const SearchBarContainer = styled.div<{ show: boolean }>`
  height: ${(p) => (p.show ? '4em' : '0')};
  transition: all 0.25s ease;

  > ${Input} {
    width: 100%;
    transform: translateY(${(p) => (p.show ? '0' : '-4em')});
    transition: all 0.25s ease;
  }
`;

export const SearchBar: React.FC<Props> = ({ show = false, value, onInput = () => {} }) => {
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    console.log(inputRef.current);
    if (show) inputRef.current?.focus();
  }, [show]);

  return (
    <SearchBarContainer show={show}>
      <Input
        ref={inputRef}
        value={value}
        onInput={(e) => onInput(e.currentTarget.value)}
        placeholder="Search for IDs, display name or tags"
      />
    </SearchBarContainer>
  );
};
