FROM python:3.10
COPY requirements.txt .
# RUN apt-get update -y && apt-get upgrade -y
# RUN apt-get install tcpdump -y
RUN pip install -r requirements.txt
COPY . ./

ENV PYTHONUNBUFFERED=1

CMD ["python3","-u", "./app/src/flclient.py"]