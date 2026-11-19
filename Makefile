# Корневой Makefile ЛР3: C++ (include, src, tests, bench) и Go (go/).
#
#   make cpp          собрать библиотеку C++ (build/libds.a)
#   make cpp-test     прогнать тесты C++ (gtest, Catch2, Boost.Test)
#   make cpp-bench    прогнать бенчмарк C++
#   make cpp-cover    покрытие C++ + HTML-отчет coverage_cpp.html (нужен gcovr)
#   make cpp-lint     строгая проверка компиляции C++ (предупреждения как ошибки)
#   make go-test      тесты Go
#   make go-cover     покрытие Go + HTML-отчет
#   make go-bench     бенчмарки Go
#   make go-lint      линтеры Go
#   make tidy         скачать Go-зависимости (testify), нужен интернет  (то же: go-tidy)
#   make fmt          форматирование Go (gofmt)                         (то же: go-fmt)
#   make clean        удалить артефакты сборки

CXX      := g++
CXXFLAGS := -std=c++17 -Wall -Wextra -O2
INCLUDES := -Iinclude -Itests

SRC   := $(wildcard src/*.cpp)
BUILD := build
OBJ   := $(patsubst src/%.cpp,$(BUILD)/%.o,$(SRC))

# C++: объектные файлы и статическая библиотека
cpp: $(BUILD)/libds.a

$(BUILD)/libds.a: $(OBJ)
	ar rcs $@ $^

$(BUILD)/%.o: src/%.cpp
	@mkdir -p $(BUILD)
	$(CXX) $(CXXFLAGS) $(INCLUDES) -c $< -o $@

# C++: тесты на трех фреймворках (нужны установленные библиотеки)
GTEST_LIBS  := -lgtest -lgtest_main -pthread
CATCH2_LIBS := -lCatch2Main -lCatch2

cpp-test: cpp-test-gtest cpp-test-catch2 cpp-test-boost

cpp-test-gtest: $(BUILD)/test_gtest
	$(BUILD)/test_gtest

cpp-test-catch2: $(BUILD)/test_catch2
	$(BUILD)/test_catch2

cpp-test-boost: $(BUILD)/test_boost
	$(BUILD)/test_boost

$(BUILD)/test_gtest: tests/test_gtest.cpp $(BUILD)/libds.a
	@mkdir -p $(BUILD)
	$(CXX) $(CXXFLAGS) $(INCLUDES) $< $(BUILD)/libds.a $(GTEST_LIBS) -o $@

$(BUILD)/test_catch2: tests/test_catch2.cpp $(BUILD)/libds.a
	@mkdir -p $(BUILD)
	$(CXX) $(CXXFLAGS) $(INCLUDES) $< $(BUILD)/libds.a $(CATCH2_LIBS) -o $@

$(BUILD)/test_boost: tests/test_boost.cpp $(BUILD)/libds.a
	@mkdir -p $(BUILD)
	$(CXX) $(CXXFLAGS) $(INCLUDES) $< $(BUILD)/libds.a -o $@

# C++: бенчмарк
cpp-bench: $(BUILD)/benchmark
	$(BUILD)/benchmark

$(BUILD)/benchmark: bench/benchmark.cpp $(BUILD)/libds.a
	@mkdir -p $(BUILD)
	$(CXX) $(CXXFLAGS) $(INCLUDES) $< $(BUILD)/libds.a -o $@

# C++: строгая проверка компиляции (любое предупреждение - ошибка)
cpp-lint:
	$(CXX) -std=c++17 -Wall -Wextra -Werror $(INCLUDES) -fsyntax-only $(SRC)

# C++: покрытие с HTML-отчетом (gcov + gcovr).
# Весь код компилируется заново с флагом --coverage (он же включает
# -fprofile-arcs и -ftest-coverage), чтобы после запуска тестов появились файлы
# .gcno/.gcda. Для покрытия запускаются все три тестовых набора.
COVER_DIR := build/coverage
COV_OBJ   := $(patsubst src/%.cpp,$(COVER_DIR)/%.o,$(SRC))

$(COVER_DIR)/%.o: src/%.cpp
	@mkdir -p $(COVER_DIR)
	$(CXX) $(CXXFLAGS) $(INCLUDES) --coverage -c $< -o $@

$(COVER_DIR)/test_gtest.o: tests/test_gtest.cpp
	@mkdir -p $(COVER_DIR)
	$(CXX) $(CXXFLAGS) $(INCLUDES) --coverage -c $< -o $@

$(COVER_DIR)/test_catch2.o: tests/test_catch2.cpp
	@mkdir -p $(COVER_DIR)
	$(CXX) $(CXXFLAGS) $(INCLUDES) --coverage -c $< -o $@

$(COVER_DIR)/test_boost.o: tests/test_boost.cpp
	@mkdir -p $(COVER_DIR)
	$(CXX) $(CXXFLAGS) $(INCLUDES) --coverage -c $< -o $@

$(COVER_DIR)/test_gtest_cov: $(COV_OBJ) $(COVER_DIR)/test_gtest.o
	$(CXX) $(CXXFLAGS) --coverage $^ $(GTEST_LIBS) -o $@

$(COVER_DIR)/test_catch2_cov: $(COV_OBJ) $(COVER_DIR)/test_catch2.o
	$(CXX) $(CXXFLAGS) --coverage $^ $(CATCH2_LIBS) -o $@

$(COVER_DIR)/test_boost_cov: $(COV_OBJ) $(COVER_DIR)/test_boost.o
	$(CXX) $(CXXFLAGS) --coverage $^ -o $@

cpp-cover: $(COVER_DIR)/test_gtest_cov $(COVER_DIR)/test_catch2_cov $(COVER_DIR)/test_boost_cov
	rm -f $(COVER_DIR)/*.gcda
	$(COVER_DIR)/test_gtest_cov
	$(COVER_DIR)/test_catch2_cov
	$(COVER_DIR)/test_boost_cov
	gcovr --root . --filter 'src/' --exclude-throw-branches --exclude-unreachable-branches --html-details coverage_cpp.html --fail-under-line 85 --fail-under-function 85 --fail-under-branch 85
	gcovr --root . --filter 'src/' --exclude-throw-branches --exclude-unreachable-branches --json-summary-pretty -o $(COVER_DIR)/summary.json
	python3 -c 'import json, sys; files=json.load(open("$(COVER_DIR)/summary.json"))["files"]; low=[f for f in files if f["line_percent"] < 85 or f["function_percent"] < 85]; print("\n".join("{}: lines {:.1f}%, functions {:.1f}%".format(f["filename"], f["line_percent"], f["function_percent"]) for f in low) if low else "Every C++ source file has at least 85% line and function coverage."); sys.exit(bool(low))'
	@echo
	@echo "== HTML-отчет покрытия C++: coverage_cpp.html  (xdg-open coverage_cpp.html)"

# Go: цели уходят в подпроект go/
go-test:  ; $(MAKE) -C go test
go-cover: ; $(MAKE) -C go cover
go-bench: ; $(MAKE) -C go bench
go-lint:  ; $(MAKE) -C go lint
go-tidy:  ; $(MAKE) -C go tidy
go-fmt:   ; $(MAKE) -C go fmt

# короткие псевдонимы (только Go, чтобы не печатать префикс go-)
tidy:     ; $(MAKE) -C go tidy
fmt:      ; $(MAKE) -C go fmt

clean:
	rm -rf $(BUILD)
	rm -f coverage_cpp.html coverage_cpp.*.html
	$(MAKE) -C go clean

.PHONY: cpp cpp-test cpp-test-gtest cpp-test-catch2 cpp-test-boost cpp-bench cpp-cover cpp-lint \
        go-test go-cover go-bench go-lint go-tidy go-fmt tidy fmt clean
