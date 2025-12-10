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
      const { accessToken, refreshToken, username: realUsername, role } =
        await loginApi({ loginId: username, password });

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
    try {
      await logoutApi();
    } catch {
      // 不做任何处理
    }
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
      const payload = token.split('.')[1];
      const decodedPayload = JSON.parse(atob(payload));

      // 从 JWT payload 中提取用户信息
      const { login_id, role } = decodedPayload;

      // 构造用户信息对象
      userInfo = {
        userId: login_id,
        username: login_id,
        realName: login_id,
        avatar: '',
        roles: [role],
      };

      userStore.setUserInfo(userInfo);
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
