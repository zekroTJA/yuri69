import { useSounds } from "../hooks/useSounds";

type Props = {};

export const SoundsRoute: React.FC<Props> = ({}) => {
  const { sounds } = useSounds();

  return (
    <>
      {sounds?.map((s) => (
        <p>{s.uid}</p>
      ))}
    </>
  );
};
