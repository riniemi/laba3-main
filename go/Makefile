# ЛР3, Go. Сборка, тесты (testing + testify), покрытие с HTML-отчетом, бенчмарк, линтеры.
#
#   make tidy     скачать зависимости (testify), создать go.sum
#   make fmt      привести код к стилю gofmt
#   make test     тесты (кратко, с процентом покрытия)
#   make test-v   тесты подробно
#   make cover    покрытие + HTML-отчет coverage.html
#   make bench    бенчмарки
#   make lint     gofmt, go vet, staticcheck, golangci-lint
#   make demo     демонстрация фабрики, полиморфизма и сериализации
#   make clean

PKG := ./ds

tidy:
	go mod tidy

fmt:
	gofmt -w .

test:
	go test -count=1 -cover ./...

test-v:
	go test -v -count=1 ./...

cover:
	go test -count=1 -covermode=atomic -coverprofile=coverage.out $(PKG)
	go tool cover -func=coverage.out > coverage_func.txt
	go tool cover -html=coverage.out -o coverage.html
	@echo
	@echo "== Покрытие по функциям (полный список: coverage_func.txt)"
	@cat coverage_func.txt
	@echo
	@awk '/^total:/ {print "== ИТОГО ПОКРЫТИЕ: " $$3}' coverage_func.txt
	@echo "== HTML-отчет: coverage.html  (xdg-open coverage.html)"

bench:
	go test $(PKG) -run '^$$' -bench . -benchmem -benchtime=200ms | tee bench_raw.txt | python3 tools/benchtable.py | tee bench_result.txt

lint:
	@echo "== gofmt (пустой вывод = все файлы отформатированы)"
	@gofmt -l .
	@echo "== go vet"
	go vet ./...
	@echo "== staticcheck"
	staticcheck ./...
	@echo "== golangci-lint"
	golangci-lint run ./...

demo:
	go run ./cmd/demo

clean:
	rm -f coverage.out coverage.html coverage_func.txt bench_raw.txt bench_result.txt demo_*.txt demo_*.bin

.PHONY: tidy fmt test test-v cover bench lint demo clean
