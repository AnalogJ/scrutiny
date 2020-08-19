#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const resolveConfig = require('tailwindcss/resolveConfig');
const buildMediaQuery = require('tailwindcss/lib/util/buildMediaQuery').default;

if ( !process.argv[3] || !process.argv[5] )
{
    console.error('Usage: -c [Relative path to Tailwind config file] -o [Relative path to Output file]');
    process.exit(1);
}

const tailwindConfig = require(path.join(process.cwd(), process.argv[3]));
const output = process.argv[5];
let outputFileContents = '';

// Read screens and build media queries
const screens = resolveConfig(tailwindConfig).theme.screens;
let queries = {};
Object.keys(screens).forEach((key) => {
    queries[key] = buildMediaQuery(screens[key])
});
queries = JSON.stringify(queries);
queries = queries.replace(/"/g, '\'').replace(/,/g, ', ').replace(/:/g, ': ');
outputFileContents = `${outputFileContents}export const treoBreakpoints = ${queries};\n`;

// Write the output file
fs.writeFile(output, outputFileContents, (err) => {
    if ( err )
    {
        return console.log(err);
    }
});
