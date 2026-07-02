import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Info, Search } from 'lucide-react'
import { getUser } from '@/api/users'
import type { UserRole } from '@/types'

export default function UsersPage() {
  const { t } = useTranslation()
  const [userId, setUserId] = useState('')
  const [searchId, setSearchId] = useState('')

  const roleLabels: Record<UserRole, string> = {
    admin: t('users.role_admin'),
    customer_service: t('users.role_customer_service'),
    operations: t('users.role_operations'),
  }

  const roleVariants: Record<UserRole, 'default' | 'secondary' | 'outline'> = {
    admin: 'default',
    customer_service: 'secondary',
    operations: 'outline',
  }

  const { data, isLoading, isError } = useQuery({
    queryKey: ['user', searchId],
    queryFn: () => getUser(searchId),
    enabled: !!searchId,
  })

  const user = data?.data

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader>
          <CardTitle>{t('users.query')}</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex gap-2">
            <div className="flex flex-1 flex-col gap-2">
              <Label htmlFor="userId">{t('users.userId')}</Label>
              <Input
                id="userId"
                placeholder={t('users.userIdPlaceholder')}
                value={userId}
                onChange={(e) => setUserId(e.target.value)}
                onKeyDown={(e) => { if (e.key === 'Enter') setSearchId(userId) }}
              />
            </div>
            <Button
              className="mt-auto"
              onClick={() => setSearchId(userId)}
              disabled={!userId}
            >
              <Search className="size-4" />
              {t('users.search')}
            </Button>
          </div>

          <Separator />

          {isLoading && (
            <div className="flex flex-col gap-2">
              <Skeleton className="h-8 w-48" />
              <Skeleton className="h-6 w-32" />
              <Skeleton className="h-6 w-64" />
            </div>
          )}

          {isError && (
            <Alert variant="destructive">
              <AlertDescription>{t('users.notFound')}</AlertDescription>
            </Alert>
          )}

          {user && (
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2">
                <span className="font-semibold">{t('users.id')}:</span>
                <code className="text-sm">{user.id}</code>
              </div>
              <div className="flex items-center gap-2">
                <span className="font-semibold">{t('users.nickname')}:</span>
                <span>{user.nickname}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="font-semibold">{t('users.email')}:</span>
                <span>{user.email}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="font-semibold">{t('users.role')}:</span>
                <Badge variant={roleVariants[user.role]}>{roleLabels[user.role]}</Badge>
              </div>
            </div>
          )}

          {!searchId && !isLoading && (
            <Alert>
              <Info className="size-4" />
              <AlertDescription>
                {t('users.noListApi')}
              </AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
