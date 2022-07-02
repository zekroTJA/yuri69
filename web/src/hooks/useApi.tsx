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
    silenceErrors?: boolean | number | number[],
  ): Promise<T> {
    if (typeof silenceErrors === 'number') silenceErrors = [silenceErrors];
    try {
      return await req(ApiClientInstance);
    } catch (e) {
      if (typeof silenceErrors === 'boolean' && silenceErrors) throw e;
      if (e instanceof APIError) {
        if (silenceErrors && silenceErrors.includes(e.code)) throw e;
        if (e.code === 401) {
          nav('/login');
          setLoggedIn(false);
        } else if (
          e.code === 403 &&
          e.message.toLowerCase() === 'you need to share a guild with yuri to access this resource'
        ) {
          nav('/noguild');
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
      throw e;
    }
  }

  return fetch;
};
