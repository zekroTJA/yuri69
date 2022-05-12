import styled from 'styled-components';

export const LinkButton = styled.button`
  background: none;
  border: none;
  color: ${(p) => p.theme.accent};
  text-decoration: underline;
  cursor: pointer;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  &:enabled:hover {
    filter: brightness(1.2);
  }
`;
