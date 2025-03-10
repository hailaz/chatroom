class Chat {
    constructor() {
        this.currentUser = null;
        this.currentRoom = null;
        this.ws = new ChatWebSocket();
        this.ui = new ChatUI();
        this.setupEventListeners();
        this.setupWebSocketHandlers();
    }

    async init() {
        // 检查登录状态
        const token = Api.getToken();
        if (!token) {
            window.location.href = '/html/login.html';
            return;
        }

        // 获取用户信息并全局设置
        this.currentUser = JSON.parse(localStorage.getItem('user'));
        window.currentUser = this.currentUser; // 设置全局currentUser
        document.getElementById('userInfo').textContent = this.currentUser.nickname;

        // 加载聊天室列表
        await this.loadRoomList();

        // 绑定消息输入事件
        const messageInput = document.getElementById('messageInput');
        messageInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage();
            }
        });
    }

    setupEventListeners() {
        // 文件上传事件
        document.getElementById('imageInput').addEventListener('change', (e) => this.handleImageUpload(e));
        document.getElementById('fileInput').addEventListener('change', (e) => this.handleFileUpload(e));

        // WebSocket错误处理
        window.addEventListener('ws:error', (e) => {
            const errorMessage = document.createElement('div');
            errorMessage.className = 'alert alert-danger';
            errorMessage.innerHTML = `
                <strong>连接错误</strong><br>
                ${e.detail.description}<br>
                <small>请检查网络连接</small>
            `;
            this.ui.messageList.insertBefore(errorMessage, this.ui.messageList.firstChild);
        });

        window.addEventListener('ws:reconnectFailed', () => {
            const finalError = document.createElement('div');
            finalError.className = 'alert alert-warning';
            finalError.textContent = '多次重连失败，请刷新页面重试';
            this.ui.messageList.insertBefore(finalError, this.ui.messageList.firstChild);
        });
    }

    setupWebSocketHandlers() {
        // 处理不同类型的消息
        this.ws.on(0, (message) => this.ui.appendMessage(message)); // 文本消息
        this.ws.on(1, (message) => this.ui.appendMessage(message)); // 图片消息
        this.ws.on(2, (message) => this.ui.appendMessage(message)); // 文件消息
        this.ws.on(3, (message) => this.ui.appendMessage(message)); // 系统消息
        this.ws.on(4, (message) => {                                // 用户列表更新
            this.ui.updateUserList(message.data);
            this.loadRoomList(); // 刷新聊天室列表以更新在线人数
        });
    }

    async loadRoomList() {
        try {
            const data = await Api.getRoomList();
            this.ui.updateRoomList(data.list, this.currentRoom);
        } catch (err) {
            console.error('加载聊天室列表失败:', err);
        }
    }

    async joinRoom(roomId) {
        try {
            // 如果已在房间中，先离开当前房间
            if (this.currentRoom) {
                await this.leaveRoom(this.currentRoom, false);
            }
            
            await Api.joinRoom(roomId);
            this.currentRoom = roomId;
            this.ws.connect(roomId);

            // 加载历史消息
            const history = await Api.getChatHistory(roomId);
            this.ui.clearChatArea();
            history.messages.reverse().forEach(msg => {
                this.ui.appendMessage(msg);
            });
            
            // 更新UI状态
            this.loadRoomList();
        } catch (err) {
            console.error('加入聊天室失败:', err);
            alert(err.message || '加入聊天室失败');
        }
    }

    async leaveRoom(roomId, clearUI = true) {
        if (!roomId) return;
        
        try {
            await Api.leaveRoom(roomId);
            
            if (this.currentRoom === roomId) {
                this.currentRoom = null;
                this.ws.close();
                if (clearUI) {
                    this.ui.clearChatArea();
                }
            }
            
            await this.loadRoomList();
        } catch (err) {
            console.error('离开聊天室失败:', err);
            alert(err.message || '离开聊天室失败');
        }
    }

    async createRoom() {
        const name = document.getElementById('roomName').value;
        const description = document.getElementById('roomDescription').value;
        const isPrivate = document.getElementById('roomPrivate').checked;

        try {
            const room = await Api.createRoom({ name, description, isPrivate });
            this.ui.createRoomModal.hide();
            await this.loadRoomList();
            this.joinRoom(room.id);
        } catch (err) {
            alert(err.message);
        }
    }

    async deleteRoom(roomId) {
        if (!confirm('确定要删除这个聊天室吗？所有聊天记录将被永久删除。')) {
            return;
        }

        try {
            await Api.deleteRoom(roomId);
            
            if (this.currentRoom === roomId) {
                this.ui.clearChatArea();
                this.currentRoom = null;
                this.ws.close();
            }
            
            await this.loadRoomList();
        } catch (err) {
            alert(err.message);
        }
    }

    sendMessage() {
        const input = document.getElementById('messageInput');
        const content = input.value.trim();
        if (!content) return;

        this.ws.sendTextMessage(content);
        input.value = '';
    }

    async handleImageUpload(e) {
        const file = e.target.files[0];
        if (!file) return;

        try {
            const base64 = await this.readFileAsBase64(file);
            this.ws.sendImageMessage(base64);
        } catch (err) {
            console.error('图片上传失败:', err);
            alert('图片上传失败');
        }

        e.target.value = '';
    }

    async handleFileUpload(e) {
        const file = e.target.files[0];
        if (!file) return;

        try {
            const base64 = await this.readFileAsBase64(file);
            this.ws.sendFileMessage({
                name: file.name,
                size: file.size,
                type: file.type,
                data: base64
            });
        } catch (err) {
            console.error('文件上传失败:', err);
            alert('文件上传失败');
        }

        e.target.value = '';
    }

    readFileAsBase64(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = () => resolve(reader.result);
            reader.onerror = reject;
            reader.readAsDataURL(file);
        });
    }

    logout() {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/html/login.html';
    }
}

// 等待DOM和其他脚本加载完成后再初始化
document.addEventListener('DOMContentLoaded', () => {
    // 创建全局实例
    window.chat = new Chat();
    window.currentUser = null;
    // 初始化
    chat.init();
});