helm-repo-html
==============

A quick little plugin that will generate a index.html file based on your index.yaml

## Quick Install

```bash
helm plugin install https://github.com/halkeye/helm-repo-html
```

## Quick Usage

```bash
# Generate your yaml and html
helm repo index ./
helm repo-html
# Now save the two file
git add index.yaml index.html
git commit -m "Update release"
git push
```

## Example

https://halkeye.github.io/helm-charts/
