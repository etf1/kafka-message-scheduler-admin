import { useTranslation } from 'react-i18next'
const NoResult = () => {
  const { t } = useTranslation()
  return (
    <strong
    className='animate-opacity gray italic'
    >
      { t('NoResults')}
    </strong>
  )
}

export default NoResult;
