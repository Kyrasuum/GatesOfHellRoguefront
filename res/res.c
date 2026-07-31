#include "res.h"
#include "stddef.h"
#include "string.h"
#include "stdio.h"
#include "stdlib.h"
#include "raymath.h"

#define TINYOBJ_MALLOC RL_MALLOC
#define TINYOBJ_CALLOC RL_CALLOC
#define TINYOBJ_REALLOC RL_REALLOC
#define TINYOBJ_FREE RL_FREE

char* ReadData(const char* file, const char* dir) {
    size_t dir_len = (unsigned int)strlen(dir);
    size_t fnm_len = (unsigned int)strlen(file);
    char* cfile = (char*)RL_CALLOC(fnm_len, sizeof(char));
    if (cfile == NULL) {
        return NULL;
    }
    char* cdir = (char*)RL_CALLOC(dir_len, sizeof(char));
    if (cdir == NULL) {
        free(cfile);
        return NULL;
    }
    strncpy(cfile, file, fnm_len);
    strncpy(cdir, dir, dir_len);

    return GetData(cfile, cdir);
}

unsigned char* ReadFileDataOverride(const char* fileName, int *dataSize){return ReadData(fileName, "");};
char* ReadFileTextOverride(const char* fileName){return ReadData(fileName, "");};

bool SaveFileDataOverride(const char *fileName, void *data, int dataSize){return true;};
bool SaveFileTextOverride(const char *fileName, const char *text){return true;};

void Init() {
	SetLoadFileDataCallback(ReadFileDataOverride);
	SetLoadFileTextCallback(ReadFileTextOverride);
	SetSaveFileDataCallback(SaveFileDataOverride);
	SetSaveFileTextCallback(SaveFileTextOverride);
}
