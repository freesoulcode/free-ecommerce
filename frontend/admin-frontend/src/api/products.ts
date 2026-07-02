import { adminClient as client } from './client'
import type { ApiResponse, PaginationParams, Product, ProductStatus } from '@/types'

interface ProductsResponse {
  products: Product[]
  total: number
  page: number
  page_size: number
}

export async function getProducts(params: PaginationParams & { status?: ProductStatus }): Promise<ApiResponse<ProductsResponse>> {
  const res = await client.get<ApiResponse<ProductsResponse>>('/products', { params })
  return res.data
}

export async function reviewProduct(id: string, action: 'approve' | 'reject'): Promise<ApiResponse<Product>> {
  const res = await client.post<ApiResponse<Product>>(`/products/${id}/${action}`)
  return res.data
}
