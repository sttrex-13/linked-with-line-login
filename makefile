run:
	docker compose up --build
hosting:
	ngrok http --domain=polished-pheasant-explicitly.ngrok-free.app 8080