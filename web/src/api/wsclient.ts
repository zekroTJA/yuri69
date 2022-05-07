import { Event, EventAuthRequest } from "./models";
import { WS_ENDPOINT } from "./static";

export class WSClient {
  private ws: WebSocket | undefined;

  constructor(
    private authTokenGetter: () => Promise<string>,
    private onEvent: (e: Event<any>) => void
  ) {}

  connect() {
    console.log(WS_ENDPOINT);
    this.ws = new WebSocket(WS_ENDPOINT);
    this.ws.addEventListener("open", async () => {
      const token = await this.authTokenGetter();
      this.send<EventAuthRequest>({
        type: "auth",
        payload: {
          token,
        },
      });
    });
    this.ws.addEventListener("message", (m) => this.onMessage(m));
  }

  send<T>(e: Event<T>) {
    this.ws?.send(JSON.stringify(e));
  }

  close() {
    this.ws?.close();
  }

  private onMessage(message: MessageEvent<any>) {
    const data = JSON.parse(message.data) as Event<any>;

    if (data.type === "authpromptfailed") {
    }

    this.onEvent(data);
  }
}
