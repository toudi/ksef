/* https://github.com/nashwaan/xml-js/issues/53 */
import { xml2js } from "xml-js"
import utf8 from "utf8"

const nativeType = function (value) {
    let nValue = Number(value);
    if (!isNaN(nValue)) {
      return nValue;
    }
    let bValue = value.toLowerCase();
    if (bValue === 'true') {
      return true;
    } else if (bValue === 'false') {
      return false;
    }
    return utf8.decode(value);
  }
  
  const removeJsonTextAttribute = function (value, parentElement) {
    try {
      const parentOfParent = parentElement._parent;
      const pOpKeys = Object.keys(parentElement._parent);
      const keyNo = pOpKeys.length;
      const keyName = pOpKeys[keyNo - 1];
      const arrOfKey = parentElement._parent[keyName];
      const arrOfKeyLen = arrOfKey.length;
      if (arrOfKeyLen > 0) {
        const arr = arrOfKey;
        const arrIndex = arrOfKey.length - 1;
        arr[arrIndex] = value;
      } else {
        parentElement._parent[keyName] = nativeType(value);
      }
    } catch (e) { }
  };
  
export const xml2object = (content) => {
    return xml2js(content, {
        compact: true,
        ignoreAttributes: true,
        nativeType: false,
        textFn: removeJsonTextAttribute,
    });
}