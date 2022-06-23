import styled from 'styled-components';

export const Spinner = styled.div`
  @keyframes spinner_rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }

  width: 1.5em;
  height: 1.5em;
  border-radius: 100%;
  border: solid 0.3em ${(p) => p.theme.text};
  border-bottom: solid 4px transparent;
  animation: spinner_rotate 1s infinite linear;
`;
