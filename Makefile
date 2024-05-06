.PHONY: create_measurements build create_measurements_1000

create_measurements: build
	./bin/1brc create_measurements 1000000000

create_measurements_1000: build
	./bin/1brc create_measurements 1000

calculate_average: build
	./bin/1brc calculate_average

build:
	@echo "building 1brc..."
	go build -o bin/1brc .
