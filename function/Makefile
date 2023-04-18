function:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
	docker compose -f ../docker-compose.yml build function
	docker compose -f ../docker-compose.yml up function

remote:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
	docker build . -t test_function:latest
	docker tag test_function:latest 087297168378.dkr.ecr.eu-west-3.amazonaws.com/serverless-thesis:latest
	docker push 087297168378.dkr.ecr.eu-west-3.amazonaws.com/serverless-thesis:latest
	aws lambda update-function-code --region eu-west-3 --function-name test-function --image-uri 087297168378.dkr.ecr.eu-west-3.amazonaws.com/serverless-thesis:latest

login:
	aws ecr get-login-password --region eu-west-3 | docker login --username AWS --password-stdin 087297168378.dkr.ecr.eu-west-3.amazonaws.com
 
