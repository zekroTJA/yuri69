import { useSounds } from '../hooks/useSounds';
import { uid } from 'react-uid';

type Props = {};

export const SoundsRoute: React.FC<Props> = ({}) => {
  const { sounds } = useSounds();

  return (
    <>
      {sounds?.map((s) => (
        <p key={uid(s)}>{s.uid}</p>
      ))}
    </>
  );
};
