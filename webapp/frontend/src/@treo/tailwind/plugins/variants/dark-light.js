const plugin = require('tailwindcss/plugin');

/**
 * Adds 'dark-light' variants
 */
module.exports = plugin(({addVariant, e}) => {

    const variant = ({modifySelectors, separator}) => {
        modifySelectors(({className}) => {
            return `[class*="theme-dark"].${e(`dark${separator}${className}`)}, [class*="theme-dark"] .${e(`dark${separator}${className}`)}, [class*="theme-light"].${e(`light${separator}${className}`)}, [class*="theme-light"] .${e(`light${separator}${className}`)}`
        })
    };

    addVariant('dark-light', variant);
});
