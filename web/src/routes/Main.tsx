import { useApi } from "../hooks/useApi";
import { useSounds } from "../hooks/useSounds";

type Props = {};

export const MainRoute: React.FC<Props> = ({}) => {
  const { sounds } = useSounds();

  return (
    <>
      {sounds?.map((s) => (
        <p>{s.uid}</p>
      ))}
    </>
  );
};
