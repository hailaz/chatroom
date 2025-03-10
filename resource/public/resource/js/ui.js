class ChatUI {
    constructor() {
        this.messageList = document.getElementById('messageList');
        this.userList = document.getElementById('userList');
        this.roomList = document.getElementById('roomList');
        this.imageModal = new bootstrap.Modal(document.getElementById('imageModal'));
        this.createRoomModal = new bootstrap.Modal(document.getElementById('createRoomModal'));
    }

    // 消息渲染
    appendMessage(message) {
        const div = document.createElement('div');
        const isCurrentUser = window.currentUser && message.userId === window.currentUser.id;
        div.className = `message d-flex align-items-start ${isCurrentUser ? 'self' : ''}`;
        
        const timestamp = new Date(message.timestamp).toLocaleTimeString();
        let content = '';

        switch (message.type) {
            case 0: // 文本消息
                content = `<div class="content">${this.escapeHtml(message.content)}</div>`;
                break;
            case 1: // 图片消息
                content = `
                    <div class="content p-0">
                        <img src="${message.content}" class="image-content" onclick="chatUI.showImagePreview('${message.content}')">
                    </div>`;
                break;
            case 2: // 文件消息
                const fileInfo = JSON.parse(message.content);
                content = `
                    <div class="content file-content" onclick="chatUI.downloadFile('${fileInfo.data}', '${fileInfo.name}')">
                        <i class="fas fa-file me-2"></i>
                        ${fileInfo.name}<br>
                        <small class="text-muted">${this.formatFileSize(fileInfo.size)}</small>
                    </div>`;
                break;
            case 3: // 系统消息
                div.className = 'message text-center text-muted small py-2';
                content = message.content;
                break;
        }
        
        if (message.type === 3) {
            div.innerHTML = content;
        } else {
            div.innerHTML = `
                <img src="${message.avatar || '/resource/image/avatar/default.png'}" class="avatar me-2">
                <div>
                    <div class="small text-muted mb-1">${message.nickname || message.username} - ${timestamp}</div>
                    ${content}
                </div>
            `;
        }
        
        this.messageList.appendChild(div);
        this.scrollToBottom();
    }

    // 更新用户列表
    updateUserList(users) {
        const header = this.userList.querySelector('.bg-light');
        this.userList.innerHTML = '';
        this.userList.appendChild(header);

        users.forEach(user => {
            const div = document.createElement('div');
            div.className = 'user-item d-flex align-items-center';
            div.innerHTML = `
                <img src="${user.avatar || '/resource/image/avatar/default.png'}" class="avatar me-2">
                <div>
                    <div>${user.nickname}</div>
                    <small class="text-muted">${user.username}</small>
                </div>
                <i class="fas fa-circle status ms-auto ${user.status === 1 ? 'text-success' : 'text-secondary'}"></i>
            `;
            this.userList.appendChild(div);
        });
    }

    // 更新聊天室列表
    updateRoomList(rooms, currentRoomId) {
        this.roomList.innerHTML = '';
        rooms.forEach(room => {
            const div = document.createElement('div');
            div.className = `list-group-item room-item ${currentRoomId === room.id ? 'active' : ''}`;
            div.onclick = () => chat.joinRoom(room.id);
            div.innerHTML = `
                <div class="room-header">
                    <h6 class="room-title">${room.name}</h6>
                    <div class="d-flex align-items-center">
                        <span class="online-count">
                            <i class="fas fa-users"></i>${room.userCount}
                        </span>
                        <div class="room-actions">
                            ${this.renderRoomActions(room)}
                        </div>
                    </div>
                </div>
                <p class="mb-0 small ${currentRoomId === room.id ? 'text-white-50' : 'text-muted'}">${room.description || '无描述'}</p>
            `;
            this.roomList.appendChild(div);
        });
    }

    // 渲染聊天室操作按钮
    renderRoomActions(room) {
        if (!window.currentUser) return '';
        
        let actions = '';
        if (room.creatorId === window.currentUser.id) {
            actions = `
                <li><a class="dropdown-item text-danger" href="#" onclick="event.stopPropagation(); chat.deleteRoom(${room.id})">
                    <i class="fas fa-trash"></i> 删除聊天室
                </a></li>`;
        }
        
        // 添加离开房间选项
        actions += `
            <li><a class="dropdown-item" href="#" onclick="event.stopPropagation(); chat.leaveRoom(${room.id})">
                <i class="fas fa-sign-out-alt"></i> 离开房间
            </a></li>`;

        return `
            <div class="dropdown ms-2">
                <button class="btn btn-sm btn-outline-secondary dropdown-toggle" type="button" data-bs-toggle="dropdown" onclick="event.stopPropagation()">
                    <i class="fas fa-ellipsis-v"></i>
                </button>
                <ul class="dropdown-menu dropdown-menu-end">
                    ${actions}
                </ul>
            </div>`;
    }

    // 清空聊天区域
    clearChatArea() {
        this.messageList.innerHTML = '';
        this.userList.innerHTML = `
            <div class="p-3 bg-light border-bottom">
                <h6 class="mb-0">在线用户</h6>
            </div>`;
    }

    // 显示创建聊天室模态框
    showCreateRoomModal() {
        document.getElementById('createRoomForm').reset();
        this.createRoomModal.show();
    }

    // 显示图片预览
    showImagePreview(src) {
        document.getElementById('imagePreview').src = src;
        this.imageModal.show();
    }

    // 下载文件
    downloadFile(base64Data, fileName) {
        const link = document.createElement('a');
        link.href = base64Data;
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    }

    // 工具方法
    scrollToBottom() {
        this.messageList.scrollTop = this.messageList.scrollHeight;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    formatFileSize(bytes) {
        const units = ['B', 'KB', 'MB', 'GB'];
        let size = bytes;
        let unit = 0;
        while (size >= 1024 && unit < units.length - 1) {
            size /= 1024;
            unit++;
        }
        return `${size.toFixed(2)} ${units[unit]}`;
    }
}

// 导出UI类
window.ChatUI = ChatUI;