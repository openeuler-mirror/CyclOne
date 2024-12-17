/**
 * [description]
 * @return {[type]} [description]
 */
export default function getHashParams() {
  const params = location.hash.split('?')[1] || '';
  return JSON.parse(
    `{"${decodeURI(params)
      .replace(/"/g, '\\"')
      .replace(/&/g, '","')
      .replace(/=/g, '":"')}"}`
  );
}

export function parseQuery(param) {
  const lookup = {};

  const parts = param.split('&');
  for (let i = 0; i < parts.length; i++) {
    const t = parts[i].split('=');

    //Fix booleans
    if (t[1] === 'true') {
      lookup[decodeURIComponent(t[0])] = true;
    } else if (t[1] === 'false') {
      lookup[decodeURIComponent(t[0])] = false;
    } else if (t[1] == parseFloat(t[1])) {
      //Fix numbers
      lookup[decodeURIComponent(t[0])] = parseFloat(t[1]);
    } else if (t[1].indexOf('|') !== -1) {
      //Fix arrays (has a pipe | in the string).
      lookup[decodeURIComponent(t[0])] = decodeURIComponent(t[1]).split('|');
    } else {
      //Everything else (Assume it's a string).
      lookup[decodeURIComponent(t[0])] = decodeURIComponent(t[1]); //Decode that ugly URI stuff.
    }
  }

  return lookup;
}

export function encodeURI(queryData) {
  const query = [];
  Object.keys(queryData).forEach(key => {
    const t = [];
    t.push(encodeURIComponent(key));
    t.push(encodeURIComponent(queryData[key]));
    query.push(t.join('='));
  });

  return query.join('&');
}

export function combineURI(router, queryData) {
  const uri = [];
  uri.push(router);
  uri.push(encodeURI(queryData));
  return uri.join('?');
}
