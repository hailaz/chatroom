import { MessageType } from './constants.js';

class ChatWebSocket {
    constructor() {
        this.ws = null;
        this.currentRoom = null;
        this.messageHandlers = new Map();
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 3000;
    }

    connect(roomId) {
        if (this.ws) {
            this.ws.close();
        }

        this.currentRoom = roomId;
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const token = Api.getToken();
        const wsUrl = `${protocol}//${window.location.host}/ws/chat?roomId=${roomId}&token=${token}`;

        this.ws = new WebSocket(wsUrl);
        this.setupWebSocketHandlers();
    }

    setupWebSocketHandlers() {
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            const handler = this.messageHandlers.get(message.type);
            if (handler) {
                handler(message);
            }
        };

        this.ws.onclose = () => {
            console.log('WebSocket连接已关闭');
            this.handleReconnect();
        };

        this.ws.onerror = (error) => {
            this.handleError(error);
        };
    }

    handleReconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            this.emit('reconnectFailed');
            return;
        }

        if (this.currentRoom) {
            setTimeout(() => {
                console.log(`尝试重新连接WebSocket... (${this.reconnectAttempts + 1}/${this.maxReconnectAttempts})`);
                this.connect(this.currentRoom);
                this.reconnectAttempts++;
            }, this.reconnectDelay);
        }
    }

    handleError(error) {
        let errorDesc = '';
        switch (this.ws.readyState) {
            case WebSocket.CONNECTING:
                errorDesc = '正在连接中...';
                break;
            case WebSocket.OPEN:
                errorDesc = '连接已建立';
                break;
            case WebSocket.CLOSING:
                errorDesc = '连接正在关闭';
                break;
            case WebSocket.CLOSED:
                errorDesc = '连接已关闭或无法建立连接';
                break;
        }

        console.error('WebSocket错误:', {
            error,
            readyState: this.ws.readyState,
            stateDesc: errorDesc
        });

        this.emit('error', {
            error,
            description: errorDesc,
            readyState: this.ws.readyState
        });
    }

    on(messageType, handler) {
        this.messageHandlers.set(messageType, handler);
    }

    emit(event, data) {
        const customEvent = new CustomEvent(`ws:${event}`, { detail: data });
        window.dispatchEvent(customEvent);
    }

    send(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket未连接');
        }
    }

    sendTextMessage(content) {
        this.send({
            type: MessageType.TEXT,
            content,
            roomId: this.currentRoom
        });
    }

    sendImageMessage(base64Data) {
        this.send({
            type: MessageType.IMAGE,
            content: base64Data,
            roomId: this.currentRoom
        });
    }

    sendFileMessage(fileInfo) {
        this.send({
            type: MessageType.FILE,
            content: JSON.stringify(fileInfo),
            roomId: this.currentRoom
        });
    }

    close() {
        if (this.ws) {
            this.ws.close();
            this.ws = null;
            this.currentRoom = null;
        }
    }
}

// 导出WebSocket类
window.ChatWebSocket = ChatWebSocket;