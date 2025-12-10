import { baseRequestClient, requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    loginId: string;
    password: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    accessToken: string;
    refreshToken: string;
    username: string;
    role: string;
    requirePasswordChange: boolean;
  }

  /** 刷新Token参数 */
  export interface RefreshTokenParams {
    refreshToken: string;
  }

  /** 刷新Token返回值 */
  export interface RefreshTokenResult {
    accessToken: string;
  }
}

/**
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  return requestClient.post<AuthApi.LoginResult>('/auth/login', data);
}

/**
 * 刷新accessToken
 */
export async function refreshTokenApi(refreshToken: string) {
  return baseRequestClient.post<AuthApi.RefreshTokenResult>(
    '/auth/refresh',
    { refreshToken },
  );
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return baseRequestClient.post('/auth/logout');
}
