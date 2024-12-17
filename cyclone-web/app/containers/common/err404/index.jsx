/**
 * NotFoundPage
 *
 * This is the page we show when the user visits a url that doesn't have a route
 */

import React from 'react';
import Exception from 'components/Exception';
import Layout from 'components/layout/page-layout';

export function NotFound(props) {
  return (
    <Layout>
      <Exception type='404'/>
    </Layout>
  );
}

export default NotFound;
