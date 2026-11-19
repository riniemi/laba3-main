// Общие мелочи для тестов всех трех фреймворков.
#pragma once
#include <initializer_list>
#include <string>

#include "ds/common.hpp"

// Список строк одной записью: S({"a", "b"}). Скобки вокруг нужны, чтобы запятые не ломали макросы проверок.
inline ds::Strs S(std::initializer_list<std::string> items) { return ds::Strs(items); }
