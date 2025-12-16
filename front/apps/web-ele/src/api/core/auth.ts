import { requestClient } from '#/api/request';

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
export async function loginApi(data: AuthApi.LoginParams): Promise<AuthApi.LoginResult> {
  const response = await requestClient.post<AuthApi.LoginResult>('/auth/login', data);
  // 从 axios 响应中提取数据
  return response.data;
}

/**
 * 刷新accessToken
 */
export async function refreshTokenApi(refreshToken: string): Promise<AuthApi.RefreshTokenResult> {
  const response = await requestClient.post<AuthApi.RefreshTokenResult>(
    '/auth/refresh',
    { refreshToken },
  );
  // 从 axios 响应中提取数据
  return response.data;
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return requestClient.post('/auth/logout');
}
