// Три простые хеш-функции для строк. Метод деления (хеш % размер таблицы) применяется снаружи.
#pragma once
#include <string>

namespace ds {


inline unsigned int hashDjb2(const std::string& s) {
    unsigned int h = 5381;
    for (size_t i = 0; i < s.size(); ++i) {
        h = h * 33 + (unsigned char)s[i];
    }
    return h;
}

inline unsigned int hashFnv1a(const std::string& s) {
    unsigned int h = 2166136261u;
    for (size_t i = 0; i < s.size(); ++i) {
        h = h ^ (unsigned char)s[i];
        h = h * 16777619u;
    }
    return h;
}

inline unsigned int hashSdbm(const std::string& s) {
    unsigned int h = 0;
    for (size_t i = 0; i < s.size(); ++i) {
        h = (unsigned char)s[i] + (h << 6) + (h << 16) - h;
    }
    return h;
}

}  // namespace ds
