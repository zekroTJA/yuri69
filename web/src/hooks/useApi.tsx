import { useNavigate } from 'react-router';
import { APIClient, APIError } from '../api';
import { Embed } from '../components/Embed';
import { ApiClientInstance } from '../instances';
import { useStore } from '../store';
import { useSnackBar } from './useSnackBar';

export const useApi = () => {
  const [setLoggedIn] = useStore((s) => [s.setLoggedIn]);
  const nav = useNavigate();
  const { show } = useSnackBar();

  async function fetch<T>(
    req: (c: APIClient) => Promise<T>,
    silenceErrors: boolean = false,
  ): Promise<T> {
    try {
      return await req(ApiClientInstance);
    } catch (e) {
      if (!silenceErrors) {
        if (e instanceof APIError) {
          if (e.code === 401) {
            nav('/login');
            setLoggedIn(false);
          } else {
            show(
              <span>
                <strong>API Error:</strong>&nbsp;{e.message} <Embed>({e.code})</Embed>
              </span>,
              'error',
              6000,
            );
          }
        } else {
          show(
            <span>
              <strong>Error:</strong>&nbsp;Unknown Request Error: <Embed>{`${e}`}</Embed>
            </span>,
            'error',
            6000,
          );
        }
      }
      throw e;
    }
  }

  return fetch;
};
