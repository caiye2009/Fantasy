import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';
import { ElNotification } from 'element-plus';
import { defineStore } from 'pinia';
import { loginApi, logoutApi, refreshTokenApi } from '#/api';
import { $t } from '#/locales';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  /**
   * 异步处理登录操作
   */
  async function authLogin(
    params: Recordable<any>,
    onSuccess?: () => Promise<void> | void,
  ) {
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;
      
      // 调用登录 API，返回扁平化数据
      const response = await loginApi(params);
      
      if (response && response.accessToken) {
        const { accessToken, refreshToken, username, role, requirePasswordChange } = response;
        
        // 存储 accessToken 到内存（Pinia）
        accessStore.setAccessToken(accessToken);
        
        // 存储 refreshToken 到 localStorage
        if (refreshToken) {
          localStorage.setItem('refreshToken', refreshToken);
        }

        // 构造用户信息（用于显示和前端权限判断）
        userInfo = {
          userId: '', // 后端没返回 ID，可以不需要
          username: username,
          realName: username, // 使用 username 作为显示名称
          role: role,
          homePath: '/dashboard',
        };

        // 存储用户信息到 localStorage + Pinia（持久化）
        userStore.setUserInfo(userInfo);
        localStorage.setItem('userInfo', JSON.stringify(userInfo));

        // 处理登录过期状态
        if (accessStore.loginExpired) {
          accessStore.setLoginExpired(false);
        }

        // 显示登录成功通知
        ElNotification({
          message: `欢迎回来，${username}`,
          title: '登录成功',
          type: 'success',
        });

        // 执行登录成功后的跳转
        if (onSuccess) {
          await onSuccess();
        } else {
          await router.push(
            userInfo.homePath || preferences.app.defaultHomePath,
          );
        }

        // 如果需要修改密码，跳转到修改密码页面
        // if (requirePasswordChange) {
        //   await router.push('/change-password');
        // }
      }
    } catch (error) {
      console.error('Login failed:', error);
      ElNotification({
        message: '登录失败，请检查用户名和密码',
        title: '错误',
        type: 'error',
      });
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  /**
   * 登出
   */
  async function logout(redirect: boolean = true) {
    try {
      await logoutApi();
    } catch {
      // 忽略错误
    }
    
    // 清除所有存储
    resetAllStores();
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('userInfo');
    accessStore.setLoginExpired(false);

    // 跳转到登录页
    await router.replace({
      path: LOGIN_PATH,
      query: redirect
        ? {
            redirect: encodeURIComponent(router.currentRoute.value.fullPath),
          }
        : {},
    });
  }

  /**
   * 初始化：从 localStorage 恢复用户状态
   */
  async function initUserState() {
    const storedUserInfo = localStorage.getItem('userInfo');
    const refreshToken = localStorage.getItem('refreshToken');
    
    if (storedUserInfo && refreshToken) {
      try {
        const userInfo = JSON.parse(storedUserInfo);
        userStore.setUserInfo(userInfo);
        
        // 尝试刷新 token
        const { accessToken } = await refreshTokenApi({ refreshToken });
        accessStore.setAccessToken(accessToken);
        
        return true;
      } catch (error) {
        // refresh token 过期，清除状态
        await logout(false);
        return false;
      }
    }
    
    return false;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    initUserState,
    loginLoading,
    logout,
  };
});