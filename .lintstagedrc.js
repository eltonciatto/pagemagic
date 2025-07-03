module.exports = {
  '*.{js,jsx,ts,tsx}': [
    'eslint --fix',
    'prettier --write',
  ],
  '*.{json,md,yml,yaml}': [
    'prettier --write',
  ],
  '*.{go}': [
    'gofmt -w',
    'go vet',
  ],
  '*.{rs}': [
    'cargo fmt',
    'cargo clippy --fix --allow-dirty',
  ],
  '*.{py}': [
    'black',
    'isort',
    'flake8',
  ],
};
