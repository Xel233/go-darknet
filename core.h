#include <stdlib.h>
#define C_SHARP_MAX_OBJECTS 1000
typedef struct bbox_t {
    unsigned int x, y, w, h;
    float prob;
    unsigned int obj_id;
    unsigned int track_id;
    unsigned int frames_counter;
} bbox;

typedef struct bbox_t_container {
    bbox candidates[C_SHARP_MAX_OBJECTS];
} bbox_container;

int init(const char*, const char*, int);

int detect_image(const char*, bbox_container*);

int dispose();