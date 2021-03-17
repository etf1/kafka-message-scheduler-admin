import React, { Suspense } from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';


jest.mock("react-i18next", () => ({
  // this mock makes sure any components using the translate hook can use it without a warning being shown
  useTranslation: () => {
    return {
      t: (str: string) => str,
      i18n: {
        changeLanguage: () => new Promise(() => {}),
      },
    };
  },
}));

test('App should render', () => {
  render(  <Suspense fallback={<div></div>}><App /></Suspense>);
});
