import { useTranslation } from 'react-i18next'
const Loader = () => {
  const { t } = useTranslation()
  return (
    <strong className='animate-opacity gray italic mtop-20'>
      {t('Loading')}
    </strong>
  )
}

export default Loader
