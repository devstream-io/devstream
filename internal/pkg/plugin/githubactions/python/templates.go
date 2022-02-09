package python

var mainPipeline = `name: Tests

on: [push, pull_request]

jobs:
  build:

    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: [2.7, 3.5, 3.6, 3.7]

    steps:
    - uses: actions/checkout@v2
    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v1
      with:
        python-version: ${{ matrix.python-version }}
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install .[test]
    - name: Lint with flake8
      run: |
        pip install -Iv enum34==1.1.6 # https://bitbucket.org/stoneleaf/enum34/issues/27/enum34-118-broken
        pip install flake8
        flake8 . --count --show-source --statistics
    - name: Run unit tests
      run: |
        python -m unittest discover -v tests/unit
`
