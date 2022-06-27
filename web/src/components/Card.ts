import styled from 'styled-components';

type Props = {
  margin?: string;
};

export const Card = styled.section<Props>`
  border: solid 3px ${(p) => p.theme.background3};
  padding: 1em;
  height: fit-content;
  margin: ${(p) => p.margin ?? '0'};
  min-width: 20em;
`;
