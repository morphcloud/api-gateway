build:
	docker build -t api-gateway-image:1.0.0 .

run:
	docker run -d -p 8000:8000 --name api-gateway-container api-gateway-image:1.0.0
