import styled from 'styled-components';
import { useStore } from '../../store';
import { LinearGradient } from '../styleParts';
import { SnackBarModel } from './model';

type Props = {};

const SnackBarContaienr = styled.div<SnackBarModel>`
  position: fixed;
  bottom: 0;
  left: 0;
  width: 100%;
  padding: 1em;
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 20;
  transition: all 0.3s ease;

  transform: translateY(${(p) => (p.show ? 0 : '3em')});
  opacity: ${(p) => (p.show ? 1 : 0)};
  pointer-events: ${(p) => (p.show ? 'all' : 'none')};

  ${(p) => {
    switch (p.type) {
      case 'error':
        return LinearGradient(p.theme.red);
      case 'warning':
        return LinearGradient(p.theme.orange);
      case 'success':
        return LinearGradient(p.theme.green);
      default:
        return LinearGradient(p.theme.accent);
    }
  }}
`;

export const SnackBar: React.FC<Props> = ({}) => {
  const [snackBar] = useStore((s) => [s.snackBar]);

  return <SnackBarContaienr {...snackBar}>{snackBar.content}</SnackBarContaienr>;
};
