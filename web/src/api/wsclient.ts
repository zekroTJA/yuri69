import { randomNumber } from '../util/rand';
import { Event, EventAuthRequest } from './models';
import { WS_ENDPOINT } from './static';

export class WSClient {
  private ws: WebSocket | undefined;
  private reconnectTries = 1;
  private connected: boolean = false;

  constructor(
    private authTokenGetter: () => Promise<string>,
    private onEvent: (e: Event<any>) => void,
  ) {}

  connect() {
    if (this.connected) return;

    this.ws = new WebSocket(WS_ENDPOINT);

    this.ws.addEventListener('open', async () => {
      const token = await this.authTokenGetter();
      this.send<EventAuthRequest>({
        type: 'auth',
        payload: {
          token,
        },
      });
    });

    this.ws.addEventListener('message', (m) => this.onMessage(m));

    this.ws.addEventListener('close', (e) => {
      if (e.wasClean && e.code === 1001) return;
      this.connected = false;
      if (this.reconnectTries + 1 < 15) {
        this.reconnectTries++;
      }
      const to = randomNumber(2000, 200) * this.reconnectTries;
      console.log(`Web socket disconnected. Trying to reconnect in ${to}ms ...`, to);
      setTimeout(() => this.connect(), to);
    });
  }

  send<T>(e: Event<T>) {
    this.ws?.send(JSON.stringify(e));
  }

  close() {
    this.ws?.close();
  }

  private onMessage(message: MessageEvent<any>) {
    const data = JSON.parse(message.data) as Event<any>;

    if (data.type === 'authpromptfailed') {
    }

    this.onEvent(data);
  }
}
