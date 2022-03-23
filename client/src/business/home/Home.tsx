import { useTranslation } from 'react-i18next'

import {
  AppStat,
  getAppStats
} from 'business/scheduler/service/SchedulerService'
import { useEffect, useState } from 'react'
import { clear } from '_common/service/SessionStorageService'
import AppStatCard from './AppStartCard'
import Icon from '_common/component/element/icon/Icon'
import Loader from '_common/component/element/Loader'

const Home = () => {
  const { t } = useTranslation()
  const [error, setError] = useState<Error>()
  const [stats, setStats] = useState<AppStat[]>()
  clear(() => false)
  useEffect(() => {
    ;(async () => {
      try {
        const stats = await getAppStats()

        setStats(stats)
        setError(undefined)
      } catch (err) {
        console.error(err)
        setError(err)
      }
    })()
  }, [])

  return (
    <div className='columns' style={{ margin: '3rem', marginTop: '6rem' }}>
      {error && (
        <div
          className='animate-opacity'
          style={{ fontWeight: 800, color: 'red' }}
        >
          <Icon name='exclamation-triangle' /> {t('LoadingError')}
        </div>
      )}
      {stats === undefined && (
        <Loader />
      )}
      {!error &&
        stats && stats.map(st => {
          return <AppStatCard key={st.scheduler} stat={st} />
        })}


    </div>
  )
}

export default Home
