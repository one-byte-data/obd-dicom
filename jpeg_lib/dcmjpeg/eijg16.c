#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include "oftypes.h"
#include "jpeglib16.h"

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
void init_destination16 (j_compress_ptr cinfo) {
     mem_dest_ptr dest = (mem_dest_ptr) cinfo->dest;
     dest->pub.next_output_byte = dest->buffer;
     dest->pub.free_in_buffer = BUFFER_SIZE; /* input buffer size */
}

boolean empty_output_buffer16 (j_compress_ptr cinfo) {
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

void term_destination16 (j_compress_ptr cinfo) {
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

boolean encode16(Uint16 *image_buffer, Uint16 width, Uint16 height, Uint16 samplesPerPixel, Uint8 **jpegBuf, int *jpegSize, int mode) {
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
     dest->pub.init_destination = init_destination16;
     dest->pub.empty_output_buffer = empty_output_buffer16;
     dest->pub.term_destination = term_destination16;
     
     cinfo.image_width = width;
     cinfo.image_height = height;
     cinfo.input_components = samplesPerPixel;
	if(samplesPerPixel==3)
		cinfo.in_color_space = JCS_RGB;
	if(samplesPerPixel==1)
		cinfo.in_color_space = JCS_GRAYSCALE;

     jpeg_set_defaults(&cinfo);
  	jpeg_simple_lossless(&cinfo, 1, 0);

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

/*
int main() {
unsigned char *jpeg_data;
unsigned char *img;
int jpeg_size, size;
FILE *fp;

puts("INFO, Starting test");
if((fp=fopen("dicom.raw", "rb"))==NULL) {
  puts("ERROR, can't open test.raw");
	return -1;
  }
size = 1992*1936*2; 
img = malloc(size);
size=fread(img, 1, size, fp);
fclose(fp);

if(encode16((Uint16 *)img, 1992, 1936, 1, &jpeg_data, &jpeg_size, 0)==TRUE){
  puts("INFO, finished encoding");
  if(jpeg_size>0){
    if((fp=fopen("dicom.jpl", "wb"))==NULL) {
      puts("ERROR, can't write test.jpg");
      return -1;
      }
    fwrite(jpeg_data, 1, jpeg_size, fp);
    fclose(fp);
    }
  else {
    puts("ERROR, jpeg_size is zero");
    }
  }
else {
  puts("ERROR, encoding JPEG data");
}
free(jpeg_data);
free(img);
puts("INFO, Finish test");
return 0;
}
*/