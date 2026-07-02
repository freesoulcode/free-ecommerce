import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
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
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { getProducts, reviewProduct } from '@/api/products'
import { PaginationBar } from '@/components/pagination-bar'
import type { Product, ProductStatus } from '@/types'

export default function ProductsPage() {
  const { t } = useTranslation()
  const [page, setPage] = useState(1)
  const [statusFilter, setStatusFilter] = useState<ProductStatus | undefined>('pending_review')
  const [reviewDialog, setReviewDialog] = useState<{ product: Product; action: 'approve' | 'reject' } | null>(null)
  const queryClient = useQueryClient()

  const statusLabels: Record<ProductStatus, string> = {
    draft: t('products.status_draft'),
    pending_review: t('products.status_pending_review'),
    review_rejected: t('products.status_review_rejected'),
    approved: t('products.status_approved'),
    on_sale: t('products.status_on_sale'),
    off_sale: t('products.status_off_sale'),
  }

  const statusVariants: Record<ProductStatus, 'secondary' | 'default' | 'destructive' | 'outline'> = {
    draft: 'outline',
    pending_review: 'secondary',
    review_rejected: 'destructive',
    approved: 'default',
    on_sale: 'default',
    off_sale: 'outline',
  }

  const { data, isLoading } = useQuery({
    queryKey: ['products', page, statusFilter],
    queryFn: () => getProducts({ page, page_size: 20, status: statusFilter }),
  })

  const reviewMutation = useMutation({
    mutationFn: ({ id, action }: { id: number; action: 'approve' | 'reject' }) =>
      reviewProduct(String(id), action),
    onSuccess: () => {
      toast.success(t('products.operateSuccess'))
      queryClient.invalidateQueries({ queryKey: ['products'] })
      setReviewDialog(null)
    },
    onError: () => {
      toast.error(t('products.operateFailed'))
    },
  })

  const products = data?.data?.products || []
  const total = data?.data?.total || 0
  const totalPages = Math.ceil(total / 20)

  const formatTime = (ts: number) => {
    if (!ts || ts < 1000000000) return '-'
    return new Date(ts * 1000).toLocaleString('zh-CN')
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="mb-2">
        <Select value={statusFilter || 'all'} onValueChange={(v) => { setStatusFilter(v === 'all' ? undefined : v as ProductStatus); setPage(1) }}>
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
                  <TableHead>{t('products.name')}</TableHead>
                  <TableHead>{t('products.shop')}</TableHead>
                  <TableHead>{t('products.price')}</TableHead>
                  <TableHead>{t('products.status')}</TableHead>
                  <TableHead>{t('products.createdAt')}</TableHead>
                  <TableHead className="text-right">{t('products.actions')}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {products.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={6} className="text-center text-muted-foreground">
                      {t('products.noData')}
                    </TableCell>
                  </TableRow>
                ) : (
                  products.map((p) => (
                    <TableRow key={p.id}>
                      <TableCell className="font-medium">{p.title}</TableCell>
                      <TableCell>{p.shop_name}</TableCell>
                      <TableCell>¥{(p.min_price_amount / 100).toFixed(2)}</TableCell>
                      <TableCell>
                        <Badge variant={statusVariants[p.review_status]}>{statusLabels[p.review_status]}</Badge>
                      </TableCell>
                      <TableCell>{formatTime(p.created_at)}</TableCell>
                      <TableCell className="text-right">
                        {p.review_status === 'pending_review' && (
                          <div className="flex justify-end gap-2">
                            <Button
                              size="sm"
                              onClick={() => setReviewDialog({ product: p, action: 'approve' })}
                            >
                              {t('products.approve')}
                            </Button>
                            <Button
                              size="sm"
                              variant="destructive"
                              onClick={() => setReviewDialog({ product: p, action: 'reject' })}
                            >
                              {t('products.reject')}
                            </Button>
                          </div>
                        )}
                      </TableCell>
                    </TableRow>
                  ))
                )}
              </TableBody>
            </Table>
            </div>
          )}

      <PaginationBar page={page} totalPages={totalPages} onPageChange={setPage} />

      <Dialog open={!!reviewDialog} onOpenChange={() => setReviewDialog(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>
              {reviewDialog?.action === 'approve' ? t('products.confirmApprove') : t('products.confirmReject')}
            </DialogTitle>
            <DialogDescription>
              {reviewDialog?.action === 'approve'
                ? t('products.confirmApproveMsg', { name: reviewDialog?.product.title })
                : t('products.confirmRejectMsg', { name: reviewDialog?.product.title })}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setReviewDialog(null)}>
              {t('products.cancel')}
            </Button>
            <Button
              variant={reviewDialog?.action === 'approve' ? 'default' : 'destructive'}
              onClick={() => {
                if (reviewDialog) {
                  reviewMutation.mutate({ id: reviewDialog.product.id, action: reviewDialog.action })
                }
              }}
              disabled={reviewMutation.isPending}
            >
              {t('products.confirm')}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
