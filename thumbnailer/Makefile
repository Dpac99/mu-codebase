function:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
	docker compose build thumbnailer
	docker compose up thumbnailer

remote:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
	docker build . -t baseline_function:latest
	docker tag baseline_function:latest 087297168378.dkr.ecr.eu-west-3.amazonaws.com/baseline_function:latest
	docker push 087297168378.dkr.ecr.eu-west-3.amazonaws.com/baseline_function:latest
	aws lambda update-function-code --region eu-west-3 --function-name baseline_function --image-uri 087297168378.dkr.ecr.eu-west-3.amazonaws.com/baseline_function:latest

login:
	aws ecr get-login-password --region eu-west-3 | docker login --username AWS --password-stdin 087297168378.dkr.ecr.eu-west-3.amazonaws.com
 
