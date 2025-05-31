# Stage 1: Builder
FROM python:3.13.2 as builder
WORKDIR /app

COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt


# Stage 2: Runtime
FROM python:3.13.2
WORKDIR /app

COPY --from=builder /root/.local /root/.local
COPY --from=builder /app/requirements.txt .

COPY cmd/recommend_service/ cmd/recommend_service/
COPY internal/recommend_service/ internal/recommend_service/
COPY pkg/proto/ pkg/proto/
COPY reluma.cbm reluma.cbm

ENV PATH=/root/.local/bin:$PATH
ENV PYTHONPATH=/app
ENV LOGURU_LEVEL=DEBUG

RUN pip check

CMD ["python", "cmd/recommend_service/main.py"]