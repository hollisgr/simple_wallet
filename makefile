SRC := cmd/app/main.go
EXEC := wallet_service

UUID := github.com/google/uuid
GIN := github.com/gin-gonic/gin
PGX := github.com/jackc/pgx github.com/jackc/pgx/v5/pgxpool
CLEANENV := github.com/ilyakaznacheev/cleanenv
VALIDATOR := github.com/go-playground/validator/v10

all: clean build run

build:
	go build -o $(EXEC) $(SRC)

run: 
	./$(EXEC)

clean:
	rm -f ./$(EXEC)

mod:
	go mod init $(SRC)

get:
	go get \
		$(UUID) \
		$(GIN) \
		$(PGX) \
		$(CLEANENV) \
		$(VALIDATOR)