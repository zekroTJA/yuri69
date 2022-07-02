import { RouteContainer } from '../components/RouteContainer';

type Props = {};

export const NoGuildRoute: React.FC<Props> = ({}) => {
  return (
    <RouteContainer center>
      <h3>Oh no. 😲</h3>
      <span>Looks like you are not sharing any guild with Yuri.</span>
    </RouteContainer>
  );
};

