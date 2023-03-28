# center-png
A command line tool to center png images based on their transparancy.

## Installation
- Clone the repo
- Run `go install`

## Usage
```
center-png [OPTIONS] FILE [FILEs...] .

Center images by cropping out transparent pixels.

OPTIONS:
  -f    Force overwrite of existing files
  -o string
        Output folder. If not specified, the output files will be written to 
                        a \"centered\" folder in the current working directory. (default "./centered")
  -p string
        Padding to add to all 4 sides of the image. 
                        Either an amount of pixels or a percentage relative to the output image
  -pb string
        Padding to add to the bottom sides of the image.
                        Either an amount of pixels or a percentage relative to the output image
  -pl string
        Padding to add to the left sides of the image.
                        Either an amount of pixels or a percentage relative to the output image
  -pr string
        Padding to add to the right sides of the image.
                        Either an amount of pixels or a percentage relative to the output image
  -pt string
        Padding to add to the top sides of the image.
                        Either an amount of pixels or a percentage relative to the output image
  -px string
        Padding to add to the left and right sides of the image. 
                        Either an amount of pixels or a percentage relative to the output image
  -py string
        Padding to add to the top and bottom sides of the image.
                        Either an amount of pixels or a percentage relative to the output image
  -s    Stop at the first error encountered. If not specified, the program will
                        continue processing the rest of the files.
  -t int
        Tolerance for detecting transparent pixels. 
                        0-255, 0 being exact and 255 being anything
```

## TODO
- Implementing parallel processing
- Adding support for resizing images
- Creating release binaries 
