ARG PORT=5000

FROM ubuntu:22.04

WORKDIR /app

RUN apt update -y && \
    apt install -y wget locales

RUN locale-gen pt_BR.UTF-8
ENV LC_ALL pt_BR.UTF-8
ENV LANG pt_BR.UTF-8

RUN wget --no-verbose -O /tmp/chrome.deb https://dl.google.com/linux/chrome/deb/pool/main/g/google-chrome-stable/google-chrome-stable_114.0.5735.198-1_amd64.deb \
    && apt install -y /tmp/chrome.deb \
    && rm /tmp/chrome.deb

RUN apt install python3 -y

COPY requirements.txt .

RUN apt install -y python3-pip && \
    pip install -r requirements.txt

COPY . .

CMD uvicorn --host 0.0.0.0 --port $PORT main:asgi_app --timeout-keep-alive 120
