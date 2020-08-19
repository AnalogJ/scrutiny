const plugin = require('tailwindcss/plugin');
const flattenColorPalette = require('tailwindcss/lib/util/flattenColorPalette').default;
const _ = require('lodash');

module.exports = plugin(({addUtilities, e, theme, variants}) => {

        const utilities = _.fromPairs(
            _.map(flattenColorPalette(theme('iconColor')), (value, modifier) => {
                return [
                    `.${e(`icon-${modifier}`)}`,
                    {
                        [`.mat-icon`]: {
                            color: value
                        }
                    }
                ]
            })
        );

        addUtilities(utilities, variants('iconColor'))
    },
    {
        theme   : {
            iconColor: theme => theme('colors')
        },
        variants: {
            iconColor: []
        }
    });
