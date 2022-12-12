FROM python:3.10-alpine

WORKDIR /code

COPY ./requirements.txt /code/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

COPY ./wsgi.py /code/
COPY ./app /code/app

USER 1000

CMD ["gunicorn", "--conf", "app/gunicorn.conf.py", "--bind", "0.0.0.0:8080", "wsgi:app"]
