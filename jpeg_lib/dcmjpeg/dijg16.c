#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <setjmp.h>
#include "jpeglib16.h"

// private error handler struct
struct DJDIJG16ErrorStruct{
  struct jpeg_error_mgr pub;
	jmp_buf setjmp_buffer;
	};

// private source manager struct
struct DJDIJG16SourceManagerStruct{
	struct jpeg_source_mgr pub;
	long skip_bytes;
	unsigned char *next_buffer;
	unsigned int next_buffer_size;
	};

void DJDIJG16ErrorExit(j_common_ptr cinfo){
	 struct DJDIJG16ErrorStruct  *jerr = (struct DJDIJG16ErrorStruct *) cinfo->err;
	 longjmp(jerr->setjmp_buffer, 1);
	 }
   
void DJDIJG16initSource(j_decompress_ptr cinfo)
{
}
   
void DJDIJG16termSource(j_decompress_ptr cinfo)
{
}

boolean DJDIJG16fillInputBuffer(j_decompress_ptr cinfo) {
  struct DJDIJG16SourceManagerStruct *src = (struct DJDIJG16SourceManagerStruct*) cinfo->src;

  // if we already have the next buffer, switch buffers
  if (src->next_buffer) {
    src->pub.next_input_byte    = src->next_buffer;
    src->pub.bytes_in_buffer    = (unsigned int) src->next_buffer_size;
    src->next_buffer            = NULL;
    src->next_buffer_size       = 0;

    // The suspension was caused by DJDIJG16skipInputData iff src->skip_bytes > 0.
    // In this case we must skip the remaining number of bytes here.
    if (src->skip_bytes > 0)
    {
      if (src->pub.bytes_in_buffer < (unsigned long) src->skip_bytes)
      {
        src->skip_bytes            -= (unsigned int) src->pub.bytes_in_buffer;
        src->pub.next_input_byte   += src->pub.bytes_in_buffer;
        src->pub.bytes_in_buffer    = 0;
        // cause a suspension return
        return FALSE;
      }
      else
      {
        src->pub.bytes_in_buffer   -= (unsigned int) src->skip_bytes;
        src->pub.next_input_byte   += src->skip_bytes;
        src->skip_bytes             = 0;
      }
    }
    return TRUE;
    }
return FALSE;
}

void DJDIJG16skipInputData(j_decompress_ptr cinfo, long num_bytes) {
  struct DJDIJG16SourceManagerStruct *src = (struct DJDIJG16SourceManagerStruct*) cinfo->src;

  if (src->pub.bytes_in_buffer < (size_t) num_bytes)
  {
    src->skip_bytes             = num_bytes - (unsigned int) src->pub.bytes_in_buffer;
    src->pub.next_input_byte   += src->pub.bytes_in_buffer;
    src->pub.bytes_in_buffer    = 0; // causes a suspension return
  }
  else
  {
    src->pub.bytes_in_buffer   -= (unsigned int) num_bytes;
    src->pub.next_input_byte   += num_bytes;
    src->skip_bytes             = 0;
  }
}

boolean decode16(unsigned char *jpeg_data, int jpeg_size, unsigned char *output_data, int output_size) {
  struct jpeg_decompress_struct cinfo;
  struct DJDIJG16ErrorStruct jerr;
  struct DJDIJG16SourceManagerStruct src;

  src.pub.init_source = DJDIJG16initSource;
  src.pub.fill_input_buffer = DJDIJG16fillInputBuffer;
  src.pub.skip_input_data   = DJDIJG16skipInputData;
  src.pub.resync_to_restart = jpeg_resync_to_restart;
  src.pub.term_source = DJDIJG16termSource;
  src.pub.bytes_in_buffer   = 0;
  src.pub.next_input_byte   = NULL;
  src.skip_bytes             = 0;
  src.next_buffer            = NULL;
  src.next_buffer_size       = 0;

  cinfo.err = jpeg_std_error(&jerr.pub);
  jerr.pub.error_exit = DJDIJG16ErrorExit;
  if(setjmp(jerr.setjmp_buffer)){
    char buffer[JMSG_LENGTH_MAX];
    cinfo.err->format_message((j_common_ptr)&cinfo, buffer);
    printf("ERROR, Exception, decode16, %s\r\n", buffer);
    jpeg_destroy_decompress(&cinfo);
    return FALSE;
    }

  jpeg_create_decompress(&cinfo);
  cinfo.src = &src.pub;
  src.next_buffer = jpeg_data;
  src.next_buffer_size = jpeg_size;
  jpeg_read_header(&cinfo, TRUE);
  
  printf("INFO, %d, %d\r\n", cinfo.image_width, cinfo.image_height);

  JSAMPARRAY buffer = NULL;
  int bufsize = 0;
  size_t rowsize = 0;
  void *jsampBuffer;

  if (jpeg_start_decompress(&cinfo)==TRUE) {
    bufsize = cinfo.output_width * cinfo.output_components; // number of JSAMPLEs per row
    rowsize = bufsize * sizeof(JSAMPLE); // number of bytes per row
    buffer = (*cinfo.mem->alloc_sarray)((j_common_ptr) &cinfo, JPOOL_IMAGE, bufsize, 1);
    if (buffer == NULL){
      puts("ERROR, decode, buffer==NULL");       
      return FALSE;
      }
    jsampBuffer = buffer;
  } else {
    bufsize = cinfo.output_width * cinfo.output_components;
    rowsize = bufsize * sizeof(JSAMPLE);
    buffer = (JSAMPARRAY) jsampBuffer;
  }

  if (output_size < rowsize * cinfo.output_height) {
    puts("ERROR, decode, output_size < rowsize*cinfo.output_height");    
    return FALSE;
    }

  while (cinfo.output_scanline < cinfo.output_height) {
    if (0 == jpeg_read_scanlines(&cinfo, buffer, 1)){
      puts("ERROR, decode, jpeg_read_scanlines");    
      return FALSE;
    }
  memcpy(output_data + (cinfo.output_scanline-1) * rowsize, *buffer, rowsize);
  }

  if (FALSE == jpeg_finish_decompress(&cinfo)) {
    puts("ERROR, decode, jpeg_finish_decompress");  
    return FALSE;
  }
    
  jpeg_destroy_decompress(&cinfo);  
  return TRUE;
}
