/**
 * i18n.js
 *
 * This will setup the i18n language files and locale data for your app.
 *
 */
import { addLocaleData } from 'react-intl';
import enLocaleData from 'react-intl/locale-data/en';
import zhLocaleData from 'react-intl/locale-data/zh';
import objectPath from 'common/object-path';
import IntlMessageFormat from 'intl-messageformat';

const translationsReq = require.context('./translations', true, /\.js$/);

function parseTranslation() {
  const ret = {};
  translationsReq.keys().forEach((key) => {
    const trans = translationsReq(key).default;
    const matches = key.match(/\.\/(.*)\.js/);
    if (matches) {
      ret[matches[1]] = trans;
    }
  });

  const flattened = objectPath.flatten(ret, '');
  const translations = {};

  Object.keys(flattened).forEach((key) => {
    const parts = key.split('.');
    const language = parts.pop();
    if (!translations[language]) {
      translations[language] = {};
    }
    translations[language][parts.join('.')] = flattened[key];
  });

  return translations;
}

addLocaleData(enLocaleData);
addLocaleData(zhLocaleData);

export const appLocales = [
  'en',
  'zh'
];

export const translationMessages = parseTranslation();

export const formatMessage = (config) => {
  const locale = config.locale || 'zh';
  const localedMessage = translationMessages[locale];
  const formatter = new IntlMessageFormat(localedMessage[config.id], locale);
  return formatter.format(config.data);
};
