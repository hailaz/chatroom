class Api {
    static getToken() {
        return localStorage.getItem('token');
    }

    static async request(url, options = {}) {
        const token = this.getToken();
        const defaultOptions = {
            headers: {
                'Authorization': token ? `Bearer ${token}` : '',
                'Content-Type': 'application/json'
            }
        };

        const finalOptions = {
            ...defaultOptions,
            ...options,
            headers: {
                ...defaultOptions.headers,
                ...options.headers
            }
        };

        try {
            const response = await fetch(url, finalOptions);
            const data = await response.json();

            if (data.code === 401) {
                // Token过期或无效，跳转到登录页
                localStorage.removeItem('token');
                localStorage.removeItem('user');
                window.location.href = '/html/login.html';
                return null;
            }

            if (data.code !== 0) {
                throw new Error(data.message || '请求失败');
            }

            return data.data;
        } catch (error) {
            console.error('API请求错误:', error);
            throw error;
        }
    }

    // 聊天室相关接口
    static async getRoomList(page = 1, size = 50) {
        return this.request(`/api/chatroom/list?page=${page}&size=${size}`);
    }

    static async createRoom(roomData) {
        return this.request('/api/chatroom/create', {
            method: 'POST',
            body: JSON.stringify(roomData)
        });
    }

    static async joinRoom(roomId) {
        return this.request(`/api/chatroom/join/${roomId}`, {
            method: 'POST'
        });
    }

    static async leaveRoom(roomId) {
        return this.request(`/api/chatroom/leave/${roomId}`, {
            method: 'POST'
        });
    }

    static async deleteRoom(roomId) {
        return this.request(`/api/chatroom/delete/${roomId}`, {
            method: 'POST'
        });
    }

    // 聊天相关接口
    static async getChatHistory(roomId, page = 1, size = 50) {
        return this.request(`/api/chat/history/${roomId}?page=${page}&size=${size}`);
    }
}

// 导出API类
window.Api = Api;