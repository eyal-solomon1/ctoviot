createdb:
	docker exec -it server-postgres-1 createdb --username=postgres --owner=postgres ctoviot

dropdb:
	docker exec -it server-postgres-1 dropdb ctoviot

# create a new migration: 'migrate create -ext sql -dir db/migrations -seq <NAME>' 
migrateup: 
	migrate -path ./db/migrations  -database "postgres://postgres:password@localhost:5432/ctoviot?sslmode=disable" -verbose up

migrateuplast: 
	migrate -path ./db/migrations  -database "postgres://postgres:password@localhost:5432/ctoviot?sslmode=disable" -verbose up 1

migratedown:
	migrate -path ./db/migrations  -database "postgres://postgres:password@localhost:5432/ctoviot?sslmode=disable" -verbose down -all

migratedownlast:
	migrate -path ./db/migrations  -database "postgres://postgres:password@localhost:5432/ctoviot?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate
	
test:
	go test -v -count=1 -cover ./... 

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/eyal-solomon1/ctoviot/db/sqlc Store

mock-aws:
	mockgen -package mockdb -destination db/mock/aws.go github.com/eyal-solomon1/ctoviot/internal/aws AWS

mock-ffmpeg:
	mockgen -package mockdb -destination db/mock/ffmpeg.go github.com/eyal-solomon1/ctoviot/internal/ffmpeg FFMPEG

mock-openai:
	mockgen -package mockdb -destination db/mock/openai.go github.com/eyal-solomon1/ctoviot/internal/openai OpenAI

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migratedownlast migrateuplast mock-aws mock-ffmpeg mock-openai