FROM python:3.9-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY app.py .

# Data directory
ENV DATA_DIR=/data
RUN mkdir -p /data/received

EXPOSE 5000

CMD ["python", "app.py"]
