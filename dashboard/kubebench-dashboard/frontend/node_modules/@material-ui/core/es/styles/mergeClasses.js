import _objectSpread from "@babel/runtime/helpers/objectSpread";
import warning from 'warning';
import getDisplayName from '../utils/getDisplayName';

function mergeClasses(options = {}) {
  const {
    baseClasses,
    newClasses,
    Component
  } = options;

  if (!newClasses) {
    return baseClasses;
  }

  const nextClasses = _objectSpread({}, baseClasses);

  Object.keys(newClasses).forEach(key => {
    process.env.NODE_ENV !== "production" ? warning(baseClasses[key] || !newClasses[key], [`Material-UI: the key \`${key}\` ` + `provided to the classes property is not implemented in ${getDisplayName(Component)}.`, `You can only override one of the following: ${Object.keys(baseClasses).join(',')}.`].join('\n')) : void 0;
    process.env.NODE_ENV !== "production" ? warning(!newClasses[key] || typeof newClasses[key] === 'string', [`Material-UI: the key \`${key}\` ` + `provided to the classes property is not valid for ${getDisplayName(Component)}.`, `You need to provide a non empty string instead of: ${newClasses[key]}.`].join('\n')) : void 0;

    if (newClasses[key]) {
      nextClasses[key] = `${baseClasses[key]} ${newClasses[key]}`;
    }
  });
  return nextClasses;
}

export default mergeClasses;