FROM ani5msr/microservices-python-flask-grpc:0.1
COPY . /api_services
WORKDIR /api_services
EXPOSE 5000
ENTRYPOINT FLASK_DEBUG=1 python run.py
