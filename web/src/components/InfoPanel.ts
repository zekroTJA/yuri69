import styled from 'styled-components';
import { LinearGradient } from './styleParts';

export const InfoPanel = styled.div<{ color?: string }>`
  position: relative;
  padding: 0.6em;

  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    opacity: 0.3;
    z-index: -1;
    ${(p) => LinearGradient(p.color ?? p.theme.accent)};
  }
`;
