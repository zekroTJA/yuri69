import styled from 'styled-components';

type Props = {
  wrap?: boolean;
  gap?: string;
  direction?: 'row' | 'column';
  vCenter?: boolean;
  margin?: string;
};

export const Flex = styled.div<Props>`
  display: flex;
  flex-wrap: ${(p) => (p.wrap ? 'wrap' : 'nowrap')};
  gap: ${(p) => p.gap};
  flex-direction: ${(p) => p.direction};
  margin: ${(p) => p.margin};

  ${(p) => p.vCenter && 'align-items: center;'}
`;

Flex.defaultProps = {
  direction: 'row',
};
