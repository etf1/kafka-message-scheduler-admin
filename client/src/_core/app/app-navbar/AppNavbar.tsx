import { changeLanguage, Lang } from '_core/i18n';
import { ROUTE_HOME } from '_core/router/routes';
import clsx from 'clsx';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import Style from './AppNavbar.module.css';

const AppNavbar = () => {
  const { t } = useTranslation();
  const [isOpen, setIsOpen] = useState(false);

  const handleBurgerClick = () => setIsOpen((isOpen) => !isOpen);

  const setLang = (lang: Lang) => {
    changeLanguage(lang);
  };

  const getLangLabel = (lang: Lang) => {
    switch (lang) {
      case 'en-US': {
        return 'Menu-Display-In-English';
      }
      case 'fr-FR': {
        return 'Menu-Display-In-French';
      }
      default:
        return 'Menu-Display-In-English';
    }
  };

  return (
    <nav className={clsx('navbar', Style.Nav)}>
      <div className='container'>
        <div className='navbar-brand'>
          <button
            type='button'
            className={clsx('navbar-burger burger white', isOpen ? 'is-active' : null, Style.NavbarMenu)}
            aria-label='menu'
            aria-expanded={isOpen}
            data-target='navbarMenu'
            onClick={handleBurgerClick}
          >
            <span aria-hidden='true'></span>
            <span aria-hidden='true'></span>
            <span aria-hidden='true'></span>
          </button>
        </div>
        <div
          id='navbarMenu'
          className={clsx('navbar-menu', Style.NavbarMenu, isOpen ? 'is-active' : null)}
        >
          <div className='navbar-start'>
            <span className={clsx('navbar-item', Style.Brand)}>
              <a
                className={clsx('button is-white', Style.NavbarLink)}
                href={ROUTE_HOME}
              >
                <span className='icon'>
                  <i className='fas fa-calendar-alt'></i>
                </span>
                <span className={Style.BrandTitle}>{t('App-title')}</span>
              </a>
            </span>
          </div>
          <div className='navbar-end'>
            <div className='navbar-item'>
              <a
                className={clsx('button is-white', Style.NavbarLink)}
                href='https://github.com/Yoann-Abbes/kafka-message-scheduler-admin'
                target='_blank'
                rel='noopener noreferrer'
              >
                <span className='icon'>
                  <i className='fab fa-github'></i>
                </span>
                <span>Source</span>
              </a>
            </div>
            <div className={clsx('navbar-item has-dropdown is-hoverable', Style.NavbarDropdown)}>
              <span className={clsx('navbar-link', Style.NavbarLink)}>
                <span className='icon'>
                  <i className='fa fa-flag'></i>
                </span>
              </span>

              <div className='navbar-dropdown'>
                <button
                  type='button'
                  onClick={() => setLang('en-US')}
                  className={clsx('navbar-item', 'has-tooltip-left', Style.LangButton)}
                  data-tooltip={t(getLangLabel('en-US'))}
                >
                  <img
                    src='/asset/english_flag.svg'
                    width='32'
                    alt={t(getLangLabel('en-US'))}
                  />
                </button>
                <button
                  type='button'
                  onClick={() => setLang('fr-FR')}
                  className={clsx('navbar-item', 'has-tooltip-left', Style.LangButton)}
                  data-tooltip={t(getLangLabel('fr-FR'))}
                >
                  <img
                    src='/asset/french_flag.svg'
                    width='32'
                    alt={t(getLangLabel('fr-FR'))}
                  />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
};

export default AppNavbar;
