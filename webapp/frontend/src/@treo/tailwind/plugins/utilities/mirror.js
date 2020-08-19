const plugin = require('tailwindcss/plugin');

/**
 * Adds utility classes for mirroring
 */
module.exports = plugin(({addUtilities, variants}) => {

        const utilities = {
            [`.mirror`]         : {
                transform: `scale(-1, 1)`
            },
            [`.mirror-vertical`]: {
                transform: `scale(1, -1)`
            }
        };

        addUtilities(utilities, variants('mirror'));
    },
    {
        variants: {
            mirror: []
        }
    }
);
