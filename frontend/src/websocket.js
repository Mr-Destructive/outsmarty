import { EventEmitter } from 'events';

class WebSocketService extends EventEmitter {
  constructor() {
    super();
    this.socket = null;
  }

  connect(slug, uid) {
    this.socket = new WebSocket(`ws://localhost:8080/ws/room/${slug}?uid=${uid}`);
    
    this.socket.onopen = () => {
      console.log('WebSocket connection established');
    };

    this.socket.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data);
        this.emit('message', message);
      } catch (error) {
        console.error('Error parsing JSON:', error);
      }
    };

    this.socket.onclose = () => {
      console.log('WebSocket connection closed');
    };

    this.socket.onerror = (error) => {
      console.error('WebSocket error', error);
    };
  }

  send(event, payload) {
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify({ event, ...payload }));
    }
  }

  disconnect() {
    if (this.socket) {
      this.socket.close();
    }
  }
}

const webSocketService = new WebSocketService();
export default webSocketService;
