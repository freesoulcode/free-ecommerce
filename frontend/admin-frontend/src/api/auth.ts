import { authClient } from './client'
import type { ApiResponse, LoginRequest, LoginResponse } from '@/types'

export async function login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  const res = await authClient.post<ApiResponse<LoginResponse>>('/login', {
    email: data.email,
    password: data.password,
  })
  return res.data
}

export async function refreshToken(refreshToken: string): Promise<ApiResponse<LoginResponse>> {
  const res = await authClient.post<ApiResponse<LoginResponse>>('/token/refresh', { refresh_token: refreshToken })
  return res.data
}

export async function logout(): Promise<void> {
  await authClient.post('/logout')
}
