import type { Recordable, UserInfo } from '@vben/types';

import { ref } from 'vue';
import { useRouter } from 'vue-router';

import { LOGIN_PATH } from '@vben/constants';
import { preferences } from '@vben/preferences';
import { resetAllStores, useAccessStore, useUserStore } from '@vben/stores';

import { ElNotification } from 'element-plus';
import { defineStore } from 'pinia';

import { loginApi, logoutApi } from '#/api';
import { $t } from '#/locales';

export const useAuthStore = defineStore('auth', () => {
  const accessStore = useAccessStore();
  const userStore = useUserStore();
  const router = useRouter();

  const loginLoading = ref(false);

  /**
   * 异步处理登录操作
   * Asynchronously handle the login process
   * @param params 登录表单数据
   */
  async function authLogin(
    params: Recordable<any>,
    onSuccess?: () => Promise<void> | void,
  ) {
    // 异步处理用户登录操作并获取 accessToken
    let userInfo: null | UserInfo = null;
    try {
      loginLoading.value = true;

      // 将表单的 username 映射为后端需要的 loginId
      const { username, password } = params;
      const {
        accessToken,
        refreshToken,
        username: realUsername,
        role,
        department,
        requirePasswordChange
      } = await loginApi({ loginId: username, password });

      // 如果成功获取到 accessToken
      if (accessToken) {
        // 将 accessToken 和 refreshToken 存储到 accessStore 中
        accessStore.setAccessToken(accessToken);
        accessStore.setRefreshToken(refreshToken);

        // 构造用户信息
        userInfo = {
          userId: username,
          username: realUsername,
          realName: realUsername,
          avatar: '',
          roles: [role],
        };

        // 存储用户信息
        userStore.setUserInfo(userInfo);

        // 根据role设置权限码（简单映射，后续可以扩展）
        accessStore.setAccessCodes([role]);

        // 如果需要修改密码，可以在这里添加逻辑
        if (requirePasswordChange) {
          console.warn('User needs to change password');
          // TODO: 跳转到修改密码页面或显示提示
        }

        if (accessStore.loginExpired) {
          accessStore.setLoginExpired(false);
        } else {
          onSuccess
            ? await onSuccess?.()
            : await router.push(preferences.app.defaultHomePath);
        }

        ElNotification({
          message: `${$t('authentication.loginSuccessDesc')}:${realUsername}`,
          title: $t('authentication.loginSuccess'),
          type: 'success',
        });
      }
    } finally {
      loginLoading.value = false;
    }

    return {
      userInfo,
    };
  }

  async function logout(redirect: boolean = true) {
    // 发送 logout 请求到后端（不等待结果，不关心成功失败）
    logoutApi().catch(() => {
      // 忽略所有错误，包括 401
    });

    // 立即清理前端状态
    resetAllStores();
    accessStore.setLoginExpired(false);

    // 回登录页带上当前路由地址
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
   * 从 accessToken 中解析用户信息
   */
  async function fetchUserInfo() {
    let userInfo: null | UserInfo = null;

    // 从 accessStore 获取 accessToken
    const token = accessStore.accessToken;
    if (!token) {
      return null;
    }

    try {
      // 解析 JWT token (payload 是 base64 编码的第二部分)
      const parts = token.split('.');
      if (parts.length < 2 || !parts[1]) {
        return null;
      }
      const decodedPayload = JSON.parse(atob(parts[1]));

      // 从 JWT payload 中提取用户信息
      const { login_id, user_name, role, department } = decodedPayload;

      // 构造用户信息对象
      userInfo = {
        userId: login_id,
        username: user_name || login_id,
        realName: user_name || login_id,
        avatar: '',
        roles: [role],
      };

      userStore.setUserInfo(userInfo);
      accessStore.setAccessCodes([role]);
    } catch (error) {
      console.error('Failed to parse JWT token:', error);
      return null;
    }

    return userInfo;
  }

  function $reset() {
    loginLoading.value = false;
  }

  return {
    $reset,
    authLogin,
    fetchUserInfo,
    loginLoading,
    logout,
  };
});
