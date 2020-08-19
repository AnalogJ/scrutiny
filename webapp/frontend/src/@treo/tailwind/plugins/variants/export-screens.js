const plugin = require('tailwindcss/plugin');
const buildMediaQuery = require('tailwindcss/lib/util/buildMediaQuery').default;
const postcss = require('postcss');
const _ = require('lodash');

/**
 * Exports 'screens' configuration as an SCSS map
 */
module.exports = plugin(({addVariant, theme}) => {

    const variant = ({container}) => {

        let map = '';

        _.forEach(theme('screens'), (value, key) => {
            map = `${map} ${key}: '${buildMediaQuery(value)}',\n`;
        });

        container.append(
            postcss.decl({
                prop : '$treo-breakpoints',
                value: `(\n ${map} ) !default`
            })
        );
    };

    addVariant('export-screens', variant);
});
