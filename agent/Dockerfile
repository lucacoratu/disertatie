# Build the agent
FROM golang:1.23-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /agent .

# The agent runs the python classifier scripts, so the runtime image contains python and scikit-learn
# (the models were serialized with scikit-learn 1.6.1)
FROM python:3.13-slim
RUN pip install --no-cache-dir scikit-learn==1.6.1 numpy==2.5.5 pandas==3.1.1 scipy==1.18.1 joblib==1.6.0
WORKDIR /app
COPY --from=build /agent /usr/local/bin/agent
# The agent loads these files relative to its working directory
COPY rules ./rules
COPY detection/ai/models ./detection/ai/models
COPY detection/ai/scripts ./detection/ai/scripts
COPY static ./static
COPY lists ./lists
COPY templates ./templates
RUN mkdir -p datasets
# The configuration is mounted in /config (the agent writes the UUID received from the API in it)
ENTRYPOINT ["agent", "-config", "/config/agent.conf"]
