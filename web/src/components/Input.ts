import Color from 'color';
import styled, { css } from 'styled-components';

export const Input = styled.input<{ invalid?: boolean }>`
  background-color: ${(p) => p.theme.background3};
  border: none;
  font-size: 1rem;
  color: ${(p) => p.theme.text};
  padding: 0.5em;
  transition: outline 0.2s ease;
  outline: solid 2px ${(p) => new Color(p.theme.accent).fade(1).hexa()};

  ${(p) =>
    p.invalid &&
    css`
      color: ${p.theme.red};
      outline: solid 2px ${p.theme.red} !important;
    `}

  &:enabled:focus {
    outline: solid 2px ${(p) => p.theme.accent};
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;
