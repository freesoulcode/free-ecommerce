// API 统一响应格式
export interface ApiResponse<T = unknown> {
  code: string
  message: string
  data: T
  request_id: string
  trace_id: string
}

// 分页参数
export interface PaginationParams {
  page: number
  page_size: number
}

// 分页响应
export interface PaginatedData<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

// 用户角色
export type UserRole = 'admin' | 'customer_service' | 'operations'

// 用户信息
export interface UserInfo {
  id: string
  email: string
  role: UserRole
  nickname: string
  avatar?: string
}

// 登录请求
export interface LoginRequest {
  email: string
  password: string
}

// 登录响应（buyer-api 格式）
export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
  access_token_expires_at: number
  refresh_token_expires_at: number
  refresh_session_id: number
  buyer: {
    id: number
    email: string
    nickname: string
    phone: string
    email_verified: boolean
  }
}

// 商家状态
export type MerchantStatus = 'pending_review' | 'active' | 'rejected' | 'disabled'

// 商家信息
export interface Merchant {
  id: string
  user_id: string
  shop_name: string
  contact_name: string
  contact_phone: string
  status: MerchantStatus
  created_at: string
  updated_at: string
}

// 商品状态
export type ProductStatus = 'draft' | 'pending_review' | 'review_rejected' | 'approved' | 'on_sale' | 'off_sale'

// 商品信息
export interface Product {
  id: number
  shop_id: number
  shop_name: string
  title: string
  sub_title: string
  main_image_url: string
  review_status: ProductStatus
  sale_status: string
  min_price_amount: number
  max_price_amount: number
  total_stock: number
  currency: string
  created_at: number
  updated_at: number
}

// 订单状态
export type OrderStatus = 'pending_payment' | 'paid' | 'merchant_processing' | 'shipped' | 'completed' | 'cancelled' | 'closed'

// 订单信息
export interface Order {
  id: number
  user_id: number
  status: OrderStatus
  total_pay_amount: number
  total_item_amount: number
  total_shipping_amount: number
  currency: string
  item_count: number
  shop_order_count: number
  shop_orders: ShopOrder[]
  created_at: number
  updated_at: number
}

export interface ShopOrder {
  id: number
  shop_id: number
  shop_name: string
  status: string
  pay_amount: number
  item_count: number
}
