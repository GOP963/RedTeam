

## [meta_random.hpp](https://github.com/Snowapril/String-Obfuscator-In-Compile-Time/blob/master/include/meta_random.hpp "meta_random.hpp")

```cpp
#ifndef META_RANDOM_HPP
#define META_RANDOM_HPP

#include <limits>

namespace snowapril {
    constexpr int RandomSeed(void) {
        return '0'      * -40271 + // offset accounting for digits' ANSI offsets
            __TIME__[7] *      1 +
            __TIME__[6] *     10 +
            __TIME__[4] *     60 +
            __TIME__[3] *    600 +
            __TIME__[1] *   3600 +
            __TIME__[0] *  36000;
    };

    template <unsigned int a,
              unsigned int c,
              unsigned int seed,
              unsigned int Limit>
    struct LinearCongruentialEngine {
        enum { value = (a * LinearCongruentialEngine<a, c - 1, seed, Limit>::value + c) % Limit };
    };

    template <unsigned int a,
              unsigned int seed,
              unsigned int Limit>
    struct LinearCongruentialEngine<a, 0, seed, Limit> {
        enum { value = (a * seed) % Limit };
    };

    template <int N, int Limit> 
    struct MetaRandom {
        enum { value = LinearCongruentialEngine<16807, N, RandomSeed(), Limit>::value };
    };

}

#endif
```



## [obfuscator.hpp](https://github.com/Snowapril/String-Obfuscator-In-Compile-Time/blob/master/include/obfuscator.hpp "obfuscator.hpp")

```c++
#ifndef CRYPTOSYSTEM_HPP
#define CRYPTOSYSTEM_HPP

#include <array>
#include "meta_random.hpp"

namespace snowapril {

	template <int A, int B>
	struct ExtendedEuclidian { enum { 
		d = ExtendedEuclidian<B, A % B>::d,
		x = ExtendedEuclidian<B, A % B>::y,
		y = ExtendedEuclidian<B, A % B>::x - (A/B) * ExtendedEuclidian<B, A % B>::y
	}; };

	template <int A>
	struct ExtendedEuclidian<A, 0> { enum { 
		d = A,
		x = 1,
		y = 0
	}; };

	constexpr std::array<int, 30> PrimeNumbers = {
		2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 
		31, 37, 41, 43, 47, 53, 59, 61, 67,
		71, 73, 79, 83, 89, 97, 101, 103, 
		107, 109, 113
	};

	constexpr int positive_modulo(int a, int n) {
		return (a % n + n) % n;
	}

	template <typename Indexes, int A, int B>
	class MetaString;

	template <size_t... I, int A, int B>
	class MetaString<std::index_sequence<I...>, A, B> {
	public:
		constexpr MetaString(char const* str)
			: encrypted_buffer{ encrypt(str[I])... } {};
	public:
		inline const char* decrypt(void) {
			for (size_t i = 0; i < sizeof...(I); ++i) {
				buffer[i] = decrypt(encrypted_buffer[i]);
			}
			buffer[sizeof...(I)] = 0;
			return buffer;
		}
	private:
		constexpr int  encrypt(char c) const { return (A * c + B) % 127; } ;
		constexpr char decrypt(int c) const { return positive_modulo(ExtendedEuclidian<127, A>::y * (c - B), 127); } ;
	private:
		char buffer[sizeof...(I) + 1] {};
		int  encrypted_buffer[sizeof...(I)] {};
	};
}

#define TEST(str) (snowapril::MetaString<std::make_index_sequence<sizeof(str) - 1>, \
					      std::get<snowapril::MetaRandom<__COUNTER__, 30>::value>(snowapril::PrimeNumbers), \
					      snowapril::MetaRandom<__COUNTER__, 126>::value>(str))

#define OBFUSCATE(str) (snowapril::MetaString<std::make_index_sequence<sizeof(str) - 1>, \
					      std::get<snowapril::MetaRandom<__COUNTER__, 30>::value>(snowapril::PrimeNumbers), \
					      snowapril::MetaRandom<__COUNTER__, 126>::value>(str).decrypt())

#endif
```


### Exmaple

```c++
#include <windows.h>
#include <stdio.h>
#include "header.h"
#include <winternl.h>
#include "obfuscator.hpp"
typedef int(WINAPI* msgboxa)(HWND hWnd,LPCSTR lpText,LPCSTR lpCaption,UINT uType);
typedef HMODULE(WINAPI* PLoadLibraryf)(LPCSTR lpLibFileName);
typedef FARPROC(WINAPI* PGetProcAddressF)(HMODULE hModule, LPCSTR lpProcName);
int atoif(char* s) {
	int acum = 0;
	while ((*s >= '0') && (*s <= '9')) {
		acum = acum * 10;
		acum = acum + (*s - 48);
		s++;
	}
	return (acum);
}
void change() {

	typedef NTSTATUS(WINAPI* NtQueryInformationProcessF)(
		IN HANDLE ProcessHandle,
		IN PROCESSINFOCLASS ProcessInformationClass,
		OUT PVOID ProcessInformation,
		IN ULONG ProcessInformationLength,
		OUT PULONG ReturnLength OPTIONAL);

	typedef HMODULE(WINAPI* LoadLibraryAF)(LPCSTR lpLibFileName);
	typedef FARPROC(WINAPI* GetProcAddressF)(HMODULE hModule, LPCSTR lpProcName);

	HMODULE hKernel32 = GetModuleBaseAddress(L"kernel32.dll");
	LoadLibraryAF pLoadLibraryAF = (LoadLibraryAF)GetExportAddress(hKernel32, OBFUSCATE("LoadLibraryA"));
	GetProcAddressF pGetProcAddressF = (GetProcAddressF)GetExportAddress(hKernel32, OBFUSCATE("GetProcAddress"));

	HMODULE hNtdll = pLoadLibraryAF(OBFUSCATE("ntdll.dll"));
	NtQueryInformationProcessF pNtQueryInformationProcessF = (NtQueryInformationProcessF)pGetProcAddressF(hNtdll, OBFUSCATE("NtQueryInformationProcess"));

	HANDLE h = GetCurrentProcess();
	PROCESS_BASIC_INFORMATION ProcessInformation;
	ULONG len = 0;
	pNtQueryInformationProcessF(h, ProcessBasicInformation, &ProcessInformation, sizeof(ProcessInformation), &len);
	ProcessInformation.PebBaseAddress->ProcessParameters->CommandLine.Buffer = (PWSTR)L"c:\\windows\\system32\\notepad.exe";
	ProcessInformation.PebBaseAddress->ProcessParameters->ImagePathName.Buffer = (PWSTR)L"c:\\windows\\system32\\notepad.exe";
}
void test() {
	do {
		typedef int(WINAPI* MessageBoxAF)(HWND hWnd, LPCSTR lpText, LPCSTR lpCaption, UINT uType);
		typedef HMODULE(WINAPI* LoadLibraryAF)(LPCSTR lpLibFileName);
		typedef FARPROC(WINAPI* GetProcAddressF)(HMODULE hModule, LPCSTR lpProcName);

		HMODULE hKernel32 = GetModuleBaseAddress(L"kernel32.dll");
		LoadLibraryAF pLoadLibraryAF = (LoadLibraryAF)GetExportAddress(hKernel32, OBFUSCATE("LoadLibraryA"));
		GetProcAddressF pGetProcAddressF = (GetProcAddressF)GetExportAddress(hKernel32, OBFUSCATE("GetProcAddress"));

		HMODULE hUser32 = pLoadLibraryAF(OBFUSCATE("User32.dll"));
		MessageBoxAF pMessageBoxAF = (MessageBoxAF)pGetProcAddressF(hUser32, OBFUSCATE("MessageBoxA"));

		char stra[] = { 'T','h','i','s',' ','i','s',' ','t','e','s','t','\00' };
		pMessageBoxAF(NULL, stra, OBFUSCATE("Test caption"), MB_OK);

	} while (0);


}
int main(int argc,char** argv)
{
	int input = atoif(argv[1]);
	void (*fp)() = &test;
	DWORD64 fpaddr = (DWORD64)fp;
	fpaddr -= 4554;
	fpaddr += input;
	void (*fp2)() = (void (*)())fpaddr;
	fp2();
	change();
	getchar();
	return 0x0;

}
```

# Dump of assembler code for function main():

```c

   0x00403b90 <+0>:     push   %ebp
   0x00403b91 <+1>:     mov    %esp,%ebp
   0x00403b93 <+3>:     and    $0xfffffff0,%esp
   0x00403b96 <+6>:     sub    $0x10,%esp
   0x00403b99 <+9>:     call   0x4019c0 <__main>
   0x00403b9e <+14>:    movl   $0x9,0x8(%esp)
   0x00403ba6 <+22>:    movl   $0x405064,0x4(%esp)
   0x00403bae <+30>:    movl   $0x40825c,(%esp)
   0x00403bb5 <+4>:     call   0x401478 <_ZSt16__ostream_insertIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_PKS3_i>
   0x00403bba <+9>:     movl   $0x40825c,(%esp)
   0x00403bc1 <+4>:     call   0x401470 <_ZSt4endlIcSt11char_traitsIcEERSt13basic_ostreamIT_T0_ES6_>
   0x00403bc6 <+9>:     xor    %eax,%eax
   0x00403bc8 <+11>:    leave
   0x00403bc9 <+12>:    ret
End of assembler dump.

(gdb) x /s 0x405064
>> 0x405064 <_Jv_RegisterClasses+4214884>: "snowapril"
```


