FROM debian:stable-slim
COPY main.py main.py
COPY books/ books/

RUN apt-get update && apt-get install -y python3 python3-pip

CMD ["python3", "main.py"]
