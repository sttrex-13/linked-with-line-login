run:
	docker compose up --build
hosting:
	ngrok http --domain=positive-muskrat-emerging.ngrok-free.app 8080