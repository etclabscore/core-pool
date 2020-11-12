#!/bin/bash

# This is dirty af but will do for now.
cp fix/intl-format-cache/src/* node_modules/intl-format-cache/src/
cp fix/intl-format-cache/lib/* node_modules/intl-format-cache/lib/

./node_modules/.bin/ember build --environment production
