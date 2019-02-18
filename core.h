#include <stdlib.h>
#include <dlfcn.h>

#define C_MAX_OBJECTS 1000
typedef struct bbox_t {
    unsigned int x, y, w, h;
    float prob;
    unsigned int obj_id;
    unsigned int track_id;
    unsigned int frames_counter;
} bbox;

typedef struct bbox_t_container {
    bbox candidates[C_MAX_OBJECTS];
} bbox_container;

int call_init(void*, const char*, const char*, int);
int call_detect_image(void*, const char*, bbox_container*);
int call_dispose(void*, void*);