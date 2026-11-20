#include <iostream>
#include <windows.h>
#include <cstdint>
#include <cstdlib>
#include <array>
#include <cerrno>
void alpha()
{
    MessageBoxA(
        NULL,
        "Hello From Alpha",
        "Obfuscate",
        MB_OK
    );
}

void brayan()
{
    MessageBoxA(
        NULL,
        "Hello From Brayan",
        "Obfuscate",
        MB_OK
    );
}
void brayan1()
{
    MessageBoxA(
        NULL,
        "Hello From Brayan1",
        "Obfuscate",
        MB_OK
    );
}

class Obfuscate
{
private:

    uint32_t input;
    uint32_t key;
    uint32_t index;

    using FunctionPointer = void(*)();

    std::array<FunctionPointer, 3> functions =
    {
        alpha,
        brayan,
        brayan1
    };
// The input is decoded by XORing it with the key (90 / 0x5A).
// For 10 functions, valid decoded indices are 0 through 9:
//
//   Function index:  0  1  2  3  4  5  6  7  8  9
//   Required input: 90 91 88 89 94 95 92 93 82 83
//
// Example: input 83 decodes to index 9, so the tenth function is called.

    bool Decode()
    {
        uint32_t decoded = input ^ key;

        if (decoded >= functions.size())
        {
            std::cout << "Invalid input!" << std::endl;
            return false;
        }

        index = decoded;

        return true;
    }

    void Dispatch()
    {
        FunctionPointer fp = functions[index];

        fp();
    }

public:

    Obfuscate(uint32_t userInput)
    {
        input = userInput;
        key = 0x5A;
        index = 0;
    }

    void Execute()
    {
        if (!Decode())
            return;

        Dispatch();
    }
};

int main(int argc, char** argv)
{
    if (argc < 2)
    {
        std::cout << "Usage: program.exe <input>" << std::endl;
        return 1;
    }

    char* end = nullptr;
    errno = 0;

    unsigned long parsed = std::strtoul(argv[1], &end, 10);

    if (errno != 0 || end == argv[1] || *end != '\0' ||
        parsed > UINT32_MAX)
    {
        std::cout << "Invalid numeric input!" << std::endl;
        return 1;
    }

    Obfuscate file(static_cast<uint32_t>(parsed));

    file.Execute();

    return 0;
}
