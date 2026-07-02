import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Store, Package, ShoppingCart, Users } from 'lucide-react'

export default function DashboardPage() {
  const { t } = useTranslation()

  const stats = [
    { title: t('dashboard.pendingMerchants'), value: '--', icon: Store, description: t('dashboard.pendingMerchantsDesc') },
    { title: t('dashboard.pendingProducts'), value: '--', icon: Package, description: t('dashboard.pendingProductsDesc') },
    { title: t('dashboard.totalOrders'), value: '--', icon: ShoppingCart, description: t('dashboard.totalOrdersDesc') },
    { title: t('dashboard.registeredUsers'), value: '--', icon: Users, description: t('dashboard.registeredUsersDesc') },
  ]

  return (
    <div className="flex flex-col gap-6">
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Card key={stat.title}>
            <CardHeader className="flex flex-row items-center justify-between pb-2">
              <CardTitle className="text-sm font-medium text-muted-foreground">
                {stat.title}
              </CardTitle>
              <stat.icon className="size-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stat.value}</div>
              <p className="text-xs text-muted-foreground">{stat.description}</p>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}
