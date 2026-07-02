import { adminClient as client } from './client'
import type { ApiResponse, Order, OrderStatus, PaginationParams } from '@/types'

interface OrdersResponse {
  order_groups: Order[]
  total: number
  page: number
  page_size: number
}

export async function getOrders(params: PaginationParams & { status?: OrderStatus }): Promise<ApiResponse<OrdersResponse>> {
  const res = await client.get<ApiResponse<OrdersResponse>>('/orders', { params })
  return res.data
}
