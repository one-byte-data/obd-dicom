#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include "oftypes.h"
#include "jpeglib12.h"

#define BUFFER_SIZE 16384

typedef struct {
     struct jpeg_destination_mgr pub; /* base class */
     unsigned char *buffer; /* buffer start address */
     int bufsize; /* size of buffer */
     unsigned char *jpeg_image;
     int jpeg_size;
} memory_destination_mgr;

typedef memory_destination_mgr* mem_dest_ptr;

/* This function is called by the library before any data gets written */
void init_destination12 (j_compress_ptr cinfo) {
     mem_dest_ptr dest = (mem_dest_ptr) cinfo->dest;
     dest->pub.next_output_byte = dest->buffer;
     dest->pub.free_in_buffer = BUFFER_SIZE; /* input buffer size */
}

boolean empty_output_buffer12 (j_compress_ptr cinfo) {
     mem_dest_ptr dest = (mem_dest_ptr) cinfo->dest;
     if(dest->jpeg_size>0){
          unsigned char *temp = malloc(dest->jpeg_size);
          memcpy(temp, dest->jpeg_image, dest->jpeg_size);
          free(dest->jpeg_image);
          dest->jpeg_image = malloc(dest->jpeg_size+BUFFER_SIZE);
          memcpy(dest->jpeg_image, temp, dest->jpeg_size);
          free(temp);
     }              
     memcpy(dest->jpeg_image+dest->jpeg_size, dest->buffer, BUFFER_SIZE);
     dest->jpeg_size=dest->jpeg_size+BUFFER_SIZE;
     dest->pub.next_output_byte = dest->buffer;
     dest->pub.free_in_buffer = BUFFER_SIZE;
     return TRUE;
}

void term_destination12 (j_compress_ptr cinfo) {
     int count;
     mem_dest_ptr dest = (mem_dest_ptr) cinfo->dest;
     count = BUFFER_SIZE - dest->pub.free_in_buffer;
     if (count) {                            
          unsigned char *temp = malloc(dest->jpeg_size);
          memcpy(temp, dest->jpeg_image, dest->jpeg_size);
          free(dest->jpeg_image);
          dest->jpeg_image = malloc(dest->jpeg_size+count);
          memcpy(dest->jpeg_image, temp, dest->jpeg_size);
          free(temp);                  
          memcpy(dest->jpeg_image+dest->jpeg_size, dest->buffer, count);
          dest->jpeg_size=dest->jpeg_size+count;
    }
}

boolean encode12(Uint16 *image_buffer, Uint16 width, Uint16 height, Uint16 samplesPerPixel, Uint8 **jpegBuf, int *jpegSize, int mode) {
     int quality=90;
     struct jpeg_compress_struct cinfo;
     struct jpeg_error_mgr jerr;
  	mem_dest_ptr dest;
     JSAMPROW row_pointer[1];
     int row_stride;
     cinfo.err = jpeg_std_error(&jerr);

     jpeg_create_compress(&cinfo);
     /* set method callbacks */
     /* first call for this instance - need to setup */
     if (cinfo.dest == 0) {
          cinfo.dest = (struct jpeg_destination_mgr *)
          (*cinfo.mem->alloc_small) ((j_common_ptr) &cinfo, JPOOL_PERMANENT, sizeof (memory_destination_mgr));
          }

     dest = (mem_dest_ptr) cinfo.dest;
     dest->buffer = malloc(BUFFER_SIZE);
     dest->jpeg_image = malloc(BUFFER_SIZE);
     dest->jpeg_size = 0;
     dest->pub.init_destination = init_destination12;
     dest->pub.empty_output_buffer = empty_output_buffer12;
     dest->pub.term_destination = term_destination12;
     
     cinfo.image_width = width;
     cinfo.image_height = height;
     cinfo.input_components = samplesPerPixel;

     if(samplesPerPixel==3)
	     cinfo.in_color_space = JCS_RGB;
	if(samplesPerPixel==1)
	     cinfo.in_color_space = JCS_GRAYSCALE;

     jpeg_set_defaults(&cinfo);

//    case EJM_baseline:
	jpeg_set_quality(&cinfo, quality, 1);
//    case EJM_lossless:
     // always disables any kind of color space conversion
//     jpeg_simple_lossless(&cinfo, psv, pt);

	if(cinfo.jpeg_color_space == JCS_YCbCr){
          cinfo.comp_info[0].h_samp_factor=1;
          cinfo.comp_info[0].v_samp_factor=1;
          }
     for(int sfi=1; sfi< MAX_COMPONENTS; sfi++){
          cinfo.comp_info[sfi].h_samp_factor=1;
          cinfo.comp_info[sfi].v_samp_factor=1;
          }

     jpeg_start_compress(&cinfo,TRUE);
     row_stride = width * samplesPerPixel;
     while (cinfo.next_scanline < cinfo.image_height){
          row_pointer[0] = &image_buffer[cinfo.next_scanline * row_stride];
          jpeg_write_scanlines(&cinfo, row_pointer, 1);
     }

     jpeg_finish_compress(&cinfo);
     *jpegBuf = malloc(dest->jpeg_size);
     memcpy(*jpegBuf, dest->jpeg_image, dest->jpeg_size);
     *jpegSize = dest->jpeg_size;
     free(dest->buffer);
     free(dest->jpeg_image);
     jpeg_destroy_compress(&cinfo);
     return TRUE;
}
