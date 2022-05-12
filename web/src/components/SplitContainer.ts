import styled from 'styled-components';

export const SplitContainer = styled.div`
  display: flex;
  gap: 1em;
  > * {
    width: 100%;
  }

  @media screen and (max-width: 70em) {
    flex-direction: column;
  }
`;
