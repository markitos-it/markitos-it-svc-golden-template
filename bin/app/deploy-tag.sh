#!/bin/bash
#:[.'.']:>- ===================================================================================
#:[.'.']:>- Marco Antonio - markitos devsecops kulture
#:[.'.']:>- The Way of the Artisan
#:[.'.']:>- markitos.es.info@gmail.com
#:[.'.']:>- 🌍 https://github.com/orgs/markitos-it/repositories
#:[.'.']:>- 🌍 https://github.com/orgs/markitos-public/repositories
#:[.'.']:>- 📺 https://www.youtube.com/@markitos_devsecops
#:[.'.']:>- ===================================================================================
set -e

if [ -z "$1" ]; then
    echo "❌ Error: Version is required"
    echo "Usage: make deploy-tag 1.2.3"
    exit 1
fi

VERSION=$1
TAG="${VERSION}"

echo "🏷️  Creating and pushing tag: ${TAG}"
git tag -a "${TAG}" -m "Release ${TAG}"
git push origin "${TAG}"
echo "✅ Tag ${TAG} created and pushed"
echo "🚀 GitHub Actions will build and deploy automatically"
