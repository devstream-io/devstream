# [[.AppName]]

This is a repo for app [[.AppName]]; bootstrapped by DevStream.

By default, the automatically generated scaffolding contains:

- a piece of sample Python [Flask](https://flask.palletsprojects.com/en/2.2.0/) web app
- sample unittest
- .gitignore
- requirements.txt
- wsgi.py
- Dockerfile
- Helm chart

## Dev

### Run Locally

```shell
python3 -m venv .venv
. .venv/bin/activate
pip install -r requirements.txt
flask --app app/[[.AppName]].py run
```

### Run in Docker

```shell
docker build . -t [[.imageRepo]]:latest
docker run -it -d -p 8080:8080 [[.imageRepo]]:latest
```

### Run in K8s

```shell
# install
helm install [[.AppName]] helm/[[.AppName]]/

# uninstall
helm delete [[.AppName]]
```

## Test

```shell
python3 -m unittest
```
