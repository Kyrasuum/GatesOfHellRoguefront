#include "raylib.h"

extern char* GetData(char* file, char* dir);

char* ReadData(const char* file, const char* dir);

unsigned char* ReadFileDataOverride(const char* fileName, int *dataSize);
char* ReadFileTextOverride(const char* fileName);

bool SaveFileDataOverride(const char *fileName, void *data, int dataSize);
bool SaveFileTextOverride(const char *fileName, const char *text);

void Init();
