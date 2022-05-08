import { useNavigate } from "react-router";
import { APIClient, APIError } from "../api";
import { ApiClientInstance } from "../instances";

export const useApi = () => {
  const nav = useNavigate();
  // const { pushNotification } = useNotifications();

  async function fetch<T>(
    req: (c: APIClient) => Promise<T>,
    silenceErrors: boolean = false
  ): Promise<T> {
    try {
      return await req(ApiClientInstance);
    } catch (e) {
      if (!silenceErrors) {
        if (e instanceof APIError) {
          if (e.code === 401) {
            nav("/login");
          } else {
            // pushNotification({
            //   type: NotificationType.ERROR,
            //   delay: 8000,
            //   heading: "API Error",
            //   message: `${e.message} (${e.code})`,
            // });
          }
        } else {
          // pushNotification({
          //   type: NotificationType.ERROR,
          //   delay: 8000,
          //   heading: "Error",
          //   message: `Unknown Request Error: ${e}`,
          // });
        }
      }
      throw e;
    }
  }

  return fetch;
};
