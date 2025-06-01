class WsConnection {
  private socket: WebSocket | null = null;
  private url: string;
  private reconnectAttempts: number = 0
  private maxReconnectAttempts: number = 6;
  private reconnectDelay: number = 2000;
  private onOpenCallback: (() => void) | null = null
  private onMessageCallback: ((message: any) => void) | null = null

  constructor(url: string) {
    this.url = url
    this.connect();
  }

  private connect(): void {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error("Max reconnect attempts reached. Stopping WebSocket.");
      return;
    } 
  
    this.socket = new WebSocket(this.url);

    this.socket.onopen = () => {
      this.reconnectAttempts = 0
      if (this.onOpenCallback) {
        this.onOpenCallback();
      }
    }

    this.socket.onclose = () => {
      console.warn("WebSocket closed. reconnecting in", this.reconnectDelay)
      this.reconnectAttempts++;
      setTimeout(() => this.connect(), this.reconnectDelay)
      this.reconnectDelay *= 2;
    }


    this.socket.onmessage = (event) => {
      if (this.onMessageCallback) {
        try {
          const data = JSON.parse(event.data);
          this.onMessageCallback(data)
        } catch (error) {
          console.error("Error parsing the WebSocket event")
        }
      }
    }

    this.socket.onerror = (error) => {
      console.error("websocket error", error)
    }
  }

  send(message: object): void {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message))
    } else {
      console.error("websocket is not in open state")
    }
  }

  onMessage(callback: (message: any) => void): void {
    this.onMessageCallback = callback;
  }

  onOpen(callback: () => void): void {
    this.onOpenCallback = callback;
  }

  close(): void {
    if (this.socket) {
      this.socket.close();
    }
  }
};

export default WsConnection;
