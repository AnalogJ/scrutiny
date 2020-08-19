const plugin = require('tailwindcss/plugin');
const _ = require('lodash');

/**
 * Adds utility classes for .mat-icon size
 */
module.exports = plugin(({addUtilities, variants, theme, e}) => {

        const utilities = _.map(theme('iconSize'), (value, key) => {

            return {
                [`.${e(`icon-size-${key}`)}`]: {
                    width     : value,
                    height    : value,
                    minWidth  : value,
                    minHeight : value,
                    fontSize  : value,
                    lineHeight: value,
                    [`svg`]   : {
                        width     : value,
                        height    : value
                    }
                }
            }
        });

        addUtilities(utilities, variants('iconSize'));
    },
    {
        theme   : {
            iconSize: {
                12: '12px',
                14: '14px',
                16: '16px',
                18: '18px',
                20: '20px',
                24: '24px',
                32: '32px',
                40: '40px',
                48: '48px',
                56: '56px',
                64: '64px',
                72: '72px',
                80: '80px',
                88: '88px',
                96: '96px'
            }
        },
        variants: {
            iconSize: []
        }
    }
);
