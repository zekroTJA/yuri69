import { useNavigate } from 'react-router';
import styled from 'styled-components';
import { Styled } from '../props';
import { LinearGradient } from '../styleParts';

type Props = Styled & {
  icon: JSX.Element;
  label: string;
  to?: string;
  action?: () => void;
  color?: string;
};

const EntryContainer = styled.div<{ color?: string }>`
  cursor: pointer;
  display: flex;
  align-items: center;

  ${(p) => LinearGradient(p.color ?? p.theme.background3)}
`;

const Icon = styled.div`
  min-width: 4em;
  width: 4em;
  height: 4em;
  display: flex;
  align-items: center;
  justify-content: center;

  > svg {
    width: 50%;
    height: 50%;
  }
`;

const Label = styled.span`
  font-weight: 600;
  text-transform: uppercase;
  font-size: 1.3rem;
  padding: 0 1em;
  white-space: nowrap;
`;

export const Entry: React.FC<Props> = ({ icon, label, to, action = () => {}, color, ...props }) => {
  const nav = useNavigate();

  const _onClick = () => {
    if (!!to) nav(to);
    else action();
  };

  return (
    <EntryContainer onClick={_onClick} color={color} {...props}>
      <Icon>{icon}</Icon>
      <Label>{label}</Label>
    </EntryContainer>
  );
};
