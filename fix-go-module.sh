#!/bin/bash

set -e

MODULE_PATH="github.com/thegostev/go-kubernetes-controllers"

echo "🔧 Ensuring go.mod module path is correct..."
MOD_PATH=$(grep '^module ' go.mod | awk '{print $2}')
if [[ "$MOD_PATH" != "$MODULE_PATH" ]]; then
  echo "⚠️  Updating go.mod module path from '$MOD_PATH' to '$MODULE_PATH'"
  sed -i.bak "s|^module .*|module $MODULE_PATH|" go.mod
  rm go.mod.bak
else
  echo "✅ go.mod module path is correct."
fi

echo "🔧 Rewriting incorrect import paths..."
BAD_IMPORTS=$(grep -r --include='*.go' -E 'import.*"(github.com/[^\"]+)"' . | grep -v "$MODULE_PATH" | grep -v '^\./vendor/' | awk -F'"' '{print $2}' | sort | uniq)
for bad in $BAD_IMPORTS; do
  echo "⚠️  Rewriting import: $bad -> $MODULE_PATH"
  find . -type f -name '*.go' -exec sed -i.bak "s|$bad|$MODULE_PATH|g" {} +
  find . -name '*.bak' -delete
done

echo "🔧 Ensuring all internal package directories exist..."
for pkg in $(grep -r --include='*.go' -oE "$MODULE_PATH/[^\"/]+" . | sort | uniq); do
  PKG_PATH="./${pkg#"$MODULE_PATH/"}"
  if [[ ! -d "$PKG_PATH" ]]; then
    echo "⚠️  Creating missing package directory: $PKG_PATH"
    mkdir -p "$PKG_PATH"
    touch "$PKG_PATH/.keep"
  fi
done

echo "🔧 Running go mod tidy and auto-staging changes..."
go mod tidy
if [[ -n $(git status --porcelain go.mod go.sum) ]]; then
  echo "⚠️  Staging go.mod and go.sum changes."
  git add go.mod go.sum
else
  echo "✅ go.mod and go.sum are tidy."
fi

echo "🎉 All Go module and import issues have been auto-fixed!" 