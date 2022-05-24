import styled from 'styled-components';

export const SplitContainer = styled.div<{ margin?: string }>`
  display: flex;
  gap: 1em;
  > * {
    width: 100%;
  }

  margin: ${(p) => p.margin};

  @media screen and (max-width: 70em) {
    flex-direction: column;
  }
`;
