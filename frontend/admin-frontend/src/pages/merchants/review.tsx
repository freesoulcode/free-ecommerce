import { useTranslation } from 'react-i18next'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Info } from 'lucide-react'

export default function MerchantsPage() {
  const { t } = useTranslation()

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <CardHeader>
          <CardTitle>{t('merchants.list')}</CardTitle>
        </CardHeader>
        <CardContent>
          <Alert>
            <Info className="size-4" />
            <AlertDescription>
              {t('merchants.notImplemented')}
            </AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    </div>
  )
}
