import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { getOrders } from '@/api/orders'
import { PaginationBar } from '@/components/pagination-bar'
import type { OrderStatus } from '@/types'

export default function OrdersPage() {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<OrderStatus | undefined>()

  const statusLabels: Record<OrderStatus, string> = {
    pending_payment: t('orders.status_pending_payment'),
    paid: t('orders.status_paid'),
    merchant_processing: t('orders.status_merchant_processing'),
    shipped: t('orders.status_shipped'),
    completed: t('orders.status_completed'),
    cancelled: t('orders.status_cancelled'),
    closed: t('orders.status_closed'),
  }

  const statusVariants: Record<OrderStatus, 'secondary' | 'default' | 'destructive' | 'outline'> = {
    pending_payment: 'secondary',
    paid: 'default',
    merchant_processing: 'default',
    shipped: 'default',
    completed: 'default',
    cancelled: 'destructive',
    closed: 'outline',
  }

  const { data, isLoading } = useQuery({
    queryKey: ['orders', page, statusFilter],
    queryFn: () => getOrders({ page, page_size: 20, status: statusFilter }),
  })

  const orders = data?.data?.order_groups || []
  const total = data?.data?.total || 0
  const totalPages = Math.ceil(total / 20)

  const formatTime = (ts: number) => {
    if (!ts || ts < 1000000000) return '-'
    return new Date(ts * 1000).toLocaleString('zh-CN')
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="mb-2">
        <Select value={statusFilter || 'all'} onValueChange={(v) => { setStatusFilter(v === 'all' ? undefined : v as OrderStatus); setPage(1) }}>
          <SelectTrigger className="w-36">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t('orders.all')}</SelectItem>
            {Object.entries(statusLabels).map(([key, label]) => (
              <SelectItem key={key} value={key}>{label}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
      {isLoading ? (
        <div className="flex flex-col gap-2">
          {Array.from({ length: 5 }).map((_, i) => (
            <Skeleton key={i} className="h-12 w-full" />
          ))}
        </div>
      ) : (
        <div className="rounded-lg border overflow-hidden">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>{t('orders.orderId')}</TableHead>
              <TableHead>{t('orders.buyerId')}</TableHead>
              <TableHead>{t('orders.shop')}</TableHead>
              <TableHead>{t('orders.amount')}</TableHead>
              <TableHead>{t('orders.status')}</TableHead>
              <TableHead>{t('orders.createdAt')}</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {orders.length === 0 ? (
              <TableRow>
                <TableCell colSpan={6} className="text-center text-muted-foreground">
                  {t('orders.noData')}
                </TableCell>
              </TableRow>
            ) : (
              orders.map((o) => (
                <TableRow key={o.id}>
                  <TableCell className="font-mono text-sm">{o.id}</TableCell>
                  <TableCell className="font-mono text-sm">{o.user_id}</TableCell>
                  <TableCell>
                    {o.shop_orders?.map((s) => s.shop_name).join(', ') || '-'}
                  </TableCell>
                  <TableCell>¥{(o.total_pay_amount / 100).toFixed(2)}</TableCell>
                  <TableCell>
                    <Badge variant={statusVariants[o.status]}>{statusLabels[o.status]}</Badge>
                  </TableCell>
                  <TableCell>{formatTime(o.created_at)}</TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
        </div>
      )}

      <PaginationBar page={page} totalPages={totalPages} onPageChange={setPage} />
    </div>
  )
}
