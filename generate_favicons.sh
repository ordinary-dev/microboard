#!/bin/sh

cd assets/img/favicon

convert favicon.svg -resize 16x16 favicon-16.png
convert favicon.svg -resize 32x32 favicon-32.png

convert favicon.svg \
    \( -clone 0 -resize 16x16 \) \
    \( -clone 0 -resize 32x32 \) \
    \( -clone 0 -resize 48x48 \) \
    -delete 0 -alpha remove -colors 256 favicon.ico

convert favicon.svg -resize 180x180 favicon-apple.png

cd -
