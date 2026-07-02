import { adminClient as client } from './client'
import type { ApiResponse, UserInfo } from '@/types'

export async function getUser(id: string): Promise<ApiResponse<UserInfo>> {
  const res = await client.get<ApiResponse<UserInfo>>(`/users/${id}`)
  return res.data
}
