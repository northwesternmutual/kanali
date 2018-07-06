#!/bin/bash

echo "Deleting old publication"
rm -rf public
mkdir public
git worktree prune
rm -rf .git/worktrees/public/

echo "Checking out gh-pages branch into public"
git worktree add -B gh-pages-static-gen public origin/gh-pages

echo "Removing existing files"
rm -rf public/*

echo "Generating site"
git submodule update
cd hugo
hugo
cd ..

echo "Updating gh-pages branch"
git add -f --all && git commit -m "Publishing to gh-pages-static-gen (publish.sh)"

echo "Ready to push"
cd ..