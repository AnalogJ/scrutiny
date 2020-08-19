const plugin = require('tailwindcss/plugin');
const _ = require('lodash');

/**
 * Adds utility classes for contrasting colors such as
 * 'text-red-200-contrast' and 'bg-blue-contrast'
 */
module.exports = plugin(({addUtilities, variants, theme, e}) => {

        const utilities = _.map(theme('colorContrasts'), (value, colorName) => {

            if ( _.isObject(value) )
            {
                return _.map(value, (color, hueName) => {

                    hueName = hueName === 'default' ? '' : `-${hueName}`;

                    return {
                        [`.${e(`text-${colorName}${hueName}-contrast`)}`]: {
                            color: color
                        },
                        [`.${e(`bg-${colorName}${hueName}-contrast`)}`]  : {
                            backgroundColor: color
                        }
                    }
                });
            }
            else
            {
                return {
                    [`.${e(`text-${colorName}-contrast`)}`]: {
                        color: value
                    },
                    [`.${e(`bg-${colorName}-contrast`)}`]  : {
                        backgroundColor: value
                    }
                }
            }
        });

        addUtilities(utilities, variants('colorContrasts'));
    },
    {
        theme   : {
            colorContrasts: theme => ({
                black      : theme('colors.white'),
                white      : theme('colors.gray.800'),
                gray       : {
                    50     : theme('colors.gray.900'),
                    100    : theme('colors.gray.900'),
                    200    : theme('colors.gray.900'),
                    300    : theme('colors.gray.900'),
                    400    : theme('colors.gray.900'),
                    500    : theme('colors.gray.900'),
                    600    : theme('colors.gray.50'),
                    700    : theme('colors.gray.50'),
                    800    : theme('colors.gray.50'),
                    900    : theme('colors.gray.50'),
                    default: theme('colors.gray.900')
                },
                'cool-gray': {
                    50     : theme('colors.cool-gray.900'),
                    100    : theme('colors.cool-gray.900'),
                    200    : theme('colors.cool-gray.900'),
                    300    : theme('colors.cool-gray.900'),
                    400    : theme('colors.cool-gray.900'),
                    500    : theme('colors.cool-gray.900'),
                    600    : theme('colors.cool-gray.50'),
                    700    : theme('colors.cool-gray.50'),
                    800    : theme('colors.cool-gray.50'),
                    900    : theme('colors.cool-gray.50'),
                    default: theme('colors.cool-gray.900')
                },
                red        : {
                    50     : theme('colors.red.900'),
                    100    : theme('colors.red.900'),
                    200    : theme('colors.red.900'),
                    300    : theme('colors.red.900'),
                    400    : theme('colors.red.900'),
                    500    : theme('colors.red.900'),
                    600    : theme('colors.red.50'),
                    700    : theme('colors.red.50'),
                    800    : theme('colors.red.50'),
                    900    : theme('colors.red.50'),
                    default: theme('colors.red.900')
                },
                orange     : {
                    50     : theme('colors.orange.900'),
                    100    : theme('colors.orange.900'),
                    200    : theme('colors.orange.900'),
                    300    : theme('colors.orange.900'),
                    400    : theme('colors.orange.900'),
                    500    : theme('colors.orange.900'),
                    600    : theme('colors.orange.50'),
                    700    : theme('colors.orange.50'),
                    800    : theme('colors.orange.50'),
                    900    : theme('colors.orange.50'),
                    default: theme('colors.orange.900')
                },
                yellow     : {
                    50     : theme('colors.yellow.900'),
                    100    : theme('colors.yellow.900'),
                    200    : theme('colors.yellow.900'),
                    300    : theme('colors.yellow.900'),
                    400    : theme('colors.yellow.900'),
                    500    : theme('colors.yellow.900'),
                    600    : theme('colors.yellow.50'),
                    700    : theme('colors.yellow.50'),
                    800    : theme('colors.yellow.50'),
                    900    : theme('colors.yellow.50'),
                    default: theme('colors.yellow.900')
                },
                green      : {
                    50     : theme('colors.green.900'),
                    100    : theme('colors.green.900'),
                    200    : theme('colors.green.900'),
                    300    : theme('colors.green.900'),
                    400    : theme('colors.green.900'),
                    500    : theme('colors.green.50'),
                    600    : theme('colors.green.50'),
                    700    : theme('colors.green.50'),
                    800    : theme('colors.green.50'),
                    900    : theme('colors.green.50'),
                    default: theme('colors.green.50')
                },
                teal       : {
                    50     : theme('colors.teal.900'),
                    100    : theme('colors.teal.900'),
                    200    : theme('colors.teal.900'),
                    300    : theme('colors.teal.900'),
                    400    : theme('colors.teal.900'),
                    500    : theme('colors.teal.50'),
                    600    : theme('colors.teal.50'),
                    700    : theme('colors.teal.50'),
                    800    : theme('colors.teal.50'),
                    900    : theme('colors.teal.50'),
                    default: theme('colors.teal.50')
                },
                blue       : {
                    50     : theme('colors.blue.900'),
                    100    : theme('colors.blue.900'),
                    200    : theme('colors.blue.900'),
                    300    : theme('colors.blue.900'),
                    400    : theme('colors.blue.900'),
                    500    : theme('colors.blue.50'),
                    600    : theme('colors.blue.50'),
                    700    : theme('colors.blue.50'),
                    800    : theme('colors.blue.50'),
                    900    : theme('colors.blue.50'),
                    default: theme('colors.blue.50')
                },
                indigo     : {
                    50     : theme('colors.indigo.900'),
                    100    : theme('colors.indigo.900'),
                    200    : theme('colors.indigo.900'),
                    300    : theme('colors.indigo.900'),
                    400    : theme('colors.indigo.900'),
                    500    : theme('colors.indigo.50'),
                    600    : theme('colors.indigo.50'),
                    700    : theme('colors.indigo.50'),
                    800    : theme('colors.indigo.50'),
                    900    : theme('colors.indigo.50'),
                    default: theme('colors.indigo.50')
                },
                purple     : {
                    50     : theme('colors.purple.900'),
                    100    : theme('colors.purple.900'),
                    200    : theme('colors.purple.900'),
                    300    : theme('colors.purple.900'),
                    400    : theme('colors.purple.900'),
                    500    : theme('colors.purple.50'),
                    600    : theme('colors.purple.50'),
                    700    : theme('colors.purple.50'),
                    800    : theme('colors.purple.50'),
                    900    : theme('colors.purple.50'),
                    default: theme('colors.purple.50')
                },
                pink       : {
                    50     : theme('colors.pink.900'),
                    100    : theme('colors.pink.900'),
                    200    : theme('colors.pink.900'),
                    300    : theme('colors.pink.900'),
                    400    : theme('colors.pink.900'),
                    500    : theme('colors.pink.50'),
                    600    : theme('colors.pink.50'),
                    700    : theme('colors.pink.50'),
                    800    : theme('colors.pink.50'),
                    900    : theme('colors.pink.50'),
                    default: theme('colors.pink.50')
                }
            })
        },
        variants: {
            colorContrasts: []
        }
    }
);
