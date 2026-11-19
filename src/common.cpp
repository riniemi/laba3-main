#include "ds/common.hpp"

#include <cctype>
#include <stdexcept>

using namespace std;

namespace ds {


void fail(const string& msg) { throw runtime_error(msg); }

string join(const Strs& v, const string& sep) {
    string r;
    for (size_t i = 0; i < v.size(); ++i) {
        if (i > 0) r += sep;
        r += v[i];
    }
    return r;
}

string quote(const string& s) {
    bool need = s.empty();
    for (char c : s)
        if (isspace((unsigned char)c) || c == '"' || c == '\\') need = true;
    if (!need) return s;
    string r = "\"";
    for (char c : s) {
        if (c == '"' || c == '\\') r += '\\';
        r += c;
    }
    return r + "\"";
}

string spaced(const Strs& v) {
    if (v.empty()) return "(empty)";
    Strs q;
    for (const auto& s : v) q.push_back(quote(s));
    return join(q, " ");
}

string upper(string s) {
    for (auto& c : s) c = toupper((unsigned char)c);
    return s;
}

bool endsWith(const string& s, const string& suf) {
    return s.size() >= suf.size() && s.compare(s.size() - suf.size(), suf.size(), suf) == 0;
}

Strs tokenize(const string& line) {
    Strs r;
    size_t i = 0, n = line.size();
    while (i < n) {
        while (i < n && isspace((unsigned char)line[i])) ++i;
        if (i >= n) break;
        string tok;
        if (line[i] == '"') {
            ++i;
            bool closed = false;
            while (i < n) {
                char c = line[i++];
                if (c == '\\' && i < n) tok += line[i++];
                else if (c == '"') { closed = true; break; }
                else tok += c;
            }
            if (!closed) fail("unterminated quote in: " + line);
        } else {
            while (i < n && !isspace((unsigned char)line[i])) tok += line[i++];
        }
        r.push_back(tok);
    }
    return r;
}

const string& arg(const Strs& a, size_t i) {
    if (i >= a.size()) fail("missing argument for " + a[0]);
    return a[i];
}

size_t toIndex(const string& s) {
    if (s.empty() || s.size() > 9) fail("bad index: " + s);
    size_t r = 0;
    for (size_t i = 0; i < s.size(); ++i) {
        if (s[i] < '0' || s[i] > '9') fail("bad index: " + s);
        r = r * 10 + (size_t)(s[i] - '0');
    }
    return r;
}

bool isInteger(const string& s) {
    if (s.empty()) return false;
    size_t i = 0;
    if (s[0] == '-') i = 1;
    if (i >= s.size()) return false;
    for (; i < s.size(); ++i) {
        if (s[i] < '0' || s[i] > '9') return false;
    }
    return true;
}

string normInt(const string& s) {
    if (!isInteger(s)) fail("tree stores integers only, got: " + s);
    bool minus = false;
    size_t i = 0;
    if (s[0] == '-') { minus = true; i = 1; }
    while (i + 1 < s.size() && s[i] == '0') i++; // убираем лишние нули
    string res = s.substr(i);
    if (res == "0") return "0";
    if (minus) return "-" + res;
    return res;
}

}  // namespace ds
