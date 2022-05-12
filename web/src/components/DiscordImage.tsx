import styled from 'styled-components';
import DCLogoURL from '../../assets/dc-logo.svg';

type ImgProps = {
  round?: boolean;
};

type Props = React.ImgHTMLAttributes<any> & ImgProps;

const StyledImg = styled.img<ImgProps>`
  border-radius: ${(p) => (p.round ? '100%' : '0')};
`;

export const DiscordImage: React.FC<Props> = ({ src, ...props }) => {
  return <StyledImg src={!!src ? src : DCLogoURL} {...props} />;
};
