FROM python

COPY requirements.txt /

RUN pip install -r /requirements.txt

COPY . /app
WORKDIR /app

CMD ["python","-u","run.py"]