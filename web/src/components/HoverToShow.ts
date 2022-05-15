import styled from 'styled-components';

export const HoverToShow = styled.p`
  position: relative;
  width: fit-content;
  padding: 0.4em 0.6em;
  background-color: ${(p) => p.theme.gray};
  font-family: 'Roboto Mono', monospace;

  &::before {
    content: 'Hover to show';
    font-family: 'Rubik', sans-serif;
    position: absolute;
    top: 0;
    left: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    background-color: ${(p) => p.theme.gray};
  }

  &:hover::before {
    display: none;
  }
`;
