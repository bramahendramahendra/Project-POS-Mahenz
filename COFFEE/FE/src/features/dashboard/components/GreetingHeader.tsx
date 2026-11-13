import { APP_TIMEZONE, getWIBNow } from '@/shared/utils/date'

import { getGreeting } from '../dashboard.utils'

export function GreetingHeader({ fullName }: { fullName: string }) {
  return (
    <div className="mb-6">
      <h1 className="text-2xl font-bold text-gray-800">
        {getGreeting()}, {fullName}
      </h1>
      <p className="mt-0.5 text-sm text-gray-500">
        {getWIBNow()
          .toDate()
          .toLocaleDateString('id-ID', {
            weekday: 'long',
            day: 'numeric',
            month: 'long',
            year: 'numeric',
            timeZone: APP_TIMEZONE,
          })}
      </p>
    </div>
  )
}
