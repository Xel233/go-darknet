#include "core.h"

//description: ((return_data_type (*)(input_data_type))bridge_input_function_pointer) (bridge_input_value)
int call_init(void* fn, const char* config, const char* weight, int gpu) {
	return ((int (*)(const char*, const char*, int))fn)(config, weight, gpu);
}

int call_detect_image(void* fn, const char* img_path, bbox_container* container) {
	return ((int (*)(const char*, bbox_container*))fn)(img_path, container);
}

int call_dispose(void* fn, void* handle) {
	return ((int (*)(void*))fn)(handle);
}