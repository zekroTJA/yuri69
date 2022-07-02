import styled, { css } from 'styled-components';

type Props = {
  maxWidth?: string;
  center?: boolean;
};

export const RouteContainer = styled.div<Props>`
  padding: 1.5em;
  width: 100%;
  height: 100%;

  ${(p) =>
    p.maxWidth &&
    css`
      max-width: ${p.maxWidth};
    `}

  ${(p) =>
    p.center &&
    css`
      display: flex;
      flex-direction: column;
      justify-content: center;
      align-items: center;
      text-align: center;
    `}
`;
