import { useSanitize } from '@/scripts/util.js';

function useAttr(json) {
  let attrs = '';

  for (const attr in json) {
    if (json[attr] != null) {
      attrs += ' ' + attr;
      attrs += '="';
      attrs += ('' + json[attr]).replace(/"/g, '&quote;');
      attrs += '"';
    }
  }

  return attrs;
}

function useElems(json) {
  let elems = '';

  json.forEach(elem => {
    switch (typeof elem) {
      case 'object':
        elems += '<' + elem.name + useAttr(elem.attr);

        if (elem?.child?.length > 0) {
          elems += '>';
          elems += useElems(elem.child);
          elems += '</' + elem.name + '>';
        } else elems += '/>';
        break;
      case 'string':
        elems += useSanitize(elem);
        break;
    }
  });

  return elems;
}

export function useXML(json) {
  json = JSON.parse(JSON.stringify(json));

  let xml = '<?xml version="1.0" encoding="utf-8"?>';
  xml += useElems(json);

  return xml;
}
