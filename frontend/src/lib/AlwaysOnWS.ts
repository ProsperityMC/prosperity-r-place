export class AlwaysOnWS {
  ws: WebSocket;
  url: string;
  onopen;
  onmessage;
  onerror;
  onclose;
  grow = 0;
  running: boolean;
  isOpen: boolean;
  catchUp: (string | ArrayBufferLike | Blob | ArrayBufferView)[] = [];

  constructor(url) {
    let that = this;
    this.running = true;
    this.url = url;
    this.ws = new WebSocket(url);
    this.ws.onopen = () => that._triggerOpen();
    this.ws.onmessage = x => that._triggerMessage(x);
    this.ws.onerror = x => that._triggerError(x);
    this.ws.onclose = x => that._triggerClose(x);
  }

  send(data: string | ArrayBufferLike | Blob | ArrayBufferView) {
    if (this.isOpen) this.ws.send(data);
    else this.catchUp.push(data);
  }

  close() {
    this.running = false;
    console.log("Closing connection to WS");
    this.isOpen = false;
    this.ws.close();
    if (this.onclose) this.onclose();
  }

  _triggerOpen() {
    this.isOpen = true;
    this._triggerReconnect();
    if (this.onopen) this.onopen();
  }

  _triggerReconnect() {
    this.grow = 0;
    while (this.catchUp.length > 0) this.send(this.catchUp.shift());
  }

  _triggerMessage(x) {
    if (this.onmessage) this.onmessage(x);
  }

  _triggerError(x) {
    this.isOpen = false;
    if (this.onerror) this.onerror(x);
    console.log("Socket error occured.", x.message, "Closing socket.");
    this.ws.close();
  }

  _triggerClose(x) {
    console.error("Closed", x);
    this.isOpen = false;
    if (!this.running) return;

    let n = this.grow;
    this.grow *= 2;
    if (this.grow == 0) this.grow = 1;
    if (this.grow > 15) this.grow = 15;

    console.log("Socket has closed.", `A reconnect will be atempted in ${n} second${n === 1 ? "s" : ""}.`, x.reason);
    let that = this;
    setTimeout(() => {
      let ws = new WebSocket(that.url);
      ws.onopen = () => {
        that.ws = ws;
      };
      ws.onmessage = x => that._triggerMessage(x);
      ws.onclose = x => that._triggerClose(x);
      ws.onerror = x => that._triggerError(x);
    }, n * 1000);
  }
}
