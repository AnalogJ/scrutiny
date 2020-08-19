const plugin = require('tailwindcss/plugin');
const postcss = require('postcss');
const _ = require('lodash');

/**
 * Exports 'boxShadow' configuration as an SCSS map
 */
module.exports = plugin(({addVariant, theme}) => {

    const variant = ({container}) => {

        let map = '';

        _.forEach(theme('boxShadow'), (value, key) => {
            map = `${map} '${key}': '${theme('boxShadow.' + key)}',\n`;
        });

        container.append(
            postcss.decl({
                prop : '$treo-elevations',
                value: `(\n ${map} ) !default`
            })
        );
    };

    addVariant('export-boxShadow', variant);
});
