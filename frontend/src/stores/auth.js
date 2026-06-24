import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '');
  const username = ref(localStorage.getItem('username') || '');
  const isAuthenticated = ref(!!token.value);

  const login = async (user, password) => {
    try {
      const response = await fetch('/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: user,
          password: password
        })
      });

      const text = await response.text();
      let data;

      try {
        data = JSON.parse(text);
      } catch (e) {
        console.error('JSON parse error:', text);
        throw new Error('服务器响应格式错误');
      }

      if (!response.ok) {
        throw new Error(data?.message || '登录请求失败');
      }

      if (!data.success) {
        throw new Error(data.message || '登录失败');
      }

      if (!data.token) {
        throw new Error('未获取到认证令牌');
      }

      token.value = data.token;
      username.value = user;
      isAuthenticated.value = true;

      localStorage.setItem('token', data.token);
      localStorage.setItem('username', user);

      return data;
    } catch (error) {
      throw error;
    }
  };

  const logout = async () => {
    try {
      await fetch('/api/auth/logout', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      token.value = '';
      username.value = '';
      isAuthenticated.value = false;
      localStorage.removeItem('token');
      localStorage.removeItem('username');
    }
  };

  const checkAuth = async () => {
    if (!token.value) {
      return false;
    }

    try {
      const response = await fetch('/api/auth/check', {
        headers: {
          'Authorization': `Bearer ${token.value}`
        }
      });

      return response.ok;
    } catch (error) {
      return false;
    }
  };

  return {
    token,
    username,
    isAuthenticated,
    login,
    logout,
    checkAuth
  };
});
