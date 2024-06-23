FROM python:3.8-buster

WORKDIR /app

RUN apt-get update

RUN apt-get install -y libpq-dev
RUN pip3 install RPi.GPIO
RUN pip3 install pytz 
RUN pip3 install psycopg2
RUN pip3 install weatherhat
RUN pip3 install python-dotenv
# RUN pip3 install adafruit-circuitpython-scd30

COPY . .
CMD [ "python3", "main.py"]