export type Handler = (payload?: object) => void;

class EventBus {
  private handlers: { [key: string]: Handler[] } = {};

  publish(event: string, payload?: object) {
    this.handlers[event]?.forEach((h) => h(payload));
  }

  subscribe(event: string, handler: Handler): () => void {
    if (!this.handlers[event]) this.handlers[event] = [handler];
    else this.handlers[event].push(handler);

    return () => {
      const i = this.handlers[event].indexOf(handler);
      if (i != -1) this.handlers[event].splice(i, 1);
    };
  }
}

export const EVENT_BUS = new EventBus();

