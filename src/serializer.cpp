#include "ds/serializer.hpp"

#include <fstream>
#include <istream>
#include <ostream>

namespace ds {

namespace {

const char kMagic[4] = {'D', 'S', 'B', '1'};
const unsigned int kMaxItemSize = 64U * 1024U * 1024U;   // защита от повреждённых файлов

void putU32(std::ostream& out, unsigned int v) {          // пишем 4 байта
    unsigned char bytes[4];
    bytes[0] = (unsigned char)(v & 0xFFU);
    bytes[1] = (unsigned char)((v >> 8) & 0xFFU);
    bytes[2] = (unsigned char)((v >> 16) & 0xFFU);
    bytes[3] = (unsigned char)((v >> 24) & 0xFFU);
    out.write((const char*)bytes, 4);
}

unsigned int getU32(std::istream& in) {
    unsigned char b[4];
    in.read((char*)b, 4);
    if (!in) fail("binary data is truncated");
    unsigned int v = 0;
    v = v | (unsigned int)b[0];
    v = v | ((unsigned int)b[1] << 8);
    v = v | ((unsigned int)b[2] << 16);
    v = v | ((unsigned int)b[3] << 24);
    return v;
}

void requireEmpty(const IStorable& obj) {
    if (obj.size() != 0) fail("object must be empty before loading");
}

}  // namespace

void Serializer::saveText(const IStorable& obj, std::ostream& out) {
    out << obj.type();
    Strs items = obj.dump();
    for (size_t i = 0; i < items.size(); ++i) {
        out << ' ' << quote(items[i]);
    }
    out << '\n';
}

void Serializer::loadText(IStorable& obj, std::istream& in) {
    requireEmpty(obj);
    std::string line;
    if (!std::getline(in, line)) fail("text data is empty");
    Strs tok = tokenize(line);
    if (tok.empty() || tok[0] != std::string(1, obj.type())) fail("type mismatch in text data");
    Strs items;
    for (size_t i = 1; i < tok.size(); ++i) items.push_back(tok[i]);
    obj.load(items);
}

void Serializer::saveBinary(const IStorable& obj, std::ostream& out) {
    Strs items = obj.dump();
    out.write(kMagic, 4);
    out.put(obj.type());
    putU32(out, (unsigned int)items.size());
    for (size_t i = 0; i < items.size(); ++i) {
        putU32(out, (unsigned int)items[i].size());
        out.write(items[i].data(), (std::streamsize)items[i].size());
    }
}

void Serializer::loadBinary(IStorable& obj, std::istream& in) {
    requireEmpty(obj);
    char magic[4];
    in.read(magic, 4);
    bool ok = true;
    if (!in) ok = false;
    if (ok) {
        for (int i = 0; i < 4; ++i) if (magic[i] != kMagic[i]) ok = false;
    }
    if (!ok) fail("not a binary data file");
    char type = 0;
    in.get(type);
    if (!in || type != obj.type()) fail("type mismatch in binary data");
    unsigned int count = getU32(in);
    Strs items;
    for (unsigned int i = 0; i < count; ++i) {
        unsigned int len = getU32(in);
        if (len > kMaxItemSize) fail("binary data is corrupted");
        std::string s;
        s.resize(len);
        if (len > 0) in.read(&s[0], (std::streamsize)len);
        if (!in) fail("binary data is truncated");
        items.push_back(s);
    }
    obj.load(items);
}

void Serializer::saveTextFile(const IStorable& obj, const std::string& path) {
    std::ofstream out(path.c_str());
    if (!out) fail("cannot write file: " + path);
    saveText(obj, out);
}

void Serializer::loadTextFile(IStorable& obj, const std::string& path) {
    std::ifstream in(path.c_str());
    if (!in) fail("cannot open file: " + path);
    loadText(obj, in);
}

void Serializer::saveBinaryFile(const IStorable& obj, const std::string& path) {
    std::ofstream out(path.c_str(), std::ios::binary);
    if (!out) fail("cannot write file: " + path);
    saveBinary(obj, out);
}

void Serializer::loadBinaryFile(IStorable& obj, const std::string& path) {
    std::ifstream in(path.c_str(), std::ios::binary);
    if (!in) fail("cannot open file: " + path);
    loadBinary(obj, in);
}

}  // namespace ds
