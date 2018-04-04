#!/bin/bash

echo "Deleting old publication"
rm -rf public
mkdir public
git worktree prune
rm -rf .git/worktrees/public/

echo "Checking out gh-pages branch into public"
git worktree add -B gh-pages public origin/gh-pages

echo "Removing existing files"
rm -rf public/*

echo "Generating site"
cd hugo
hugo
cd ..

echo "Adding CNAME entry"
cd public
echo "kanali.io" >> CNAME

echo "Updating gh-pages branch"
git add -f --all && git commit -m "Publishing to gh-pages (publish.sh)"

echo "Ready to push"
cd ..