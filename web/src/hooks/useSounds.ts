import { useEffect, useState } from "react";
import { Sound } from "../api";
import { useApi } from "./useApi";

export const useSounds = () => {
  const fetch = useApi();

  const [sounds, setSounds] = useState<Sound[]>();

  useEffect(() => {
    fetch((c) => c.sounds())
      .then((sounds) => setSounds(sounds))
      .catch();
  }, []);

  return { sounds };
};
