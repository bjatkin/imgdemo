# ImgDeom

This repo is a simple demonstration of some of the image processing one can acomplish with gos standard library.
It has no dependencies and relies only on Go's standard library.
In particular it relies on the [image package](https://pkg.go.dev/image@go1.23.0)

### Root

ImgDemo is a cli tool.
The root command for the cli is `imgdemo`.
To learn more about what you can do with the cli run `imgdemo help`

### Hide

The `hide` command can be used to hide secret data in a PNG image.
```sh
$ imgdemo hide help
hide: hide data inside an image using steganography
USAGE:  hide [INPUT IMAGE PATH] [DATA FILE PATH] [OUTPUT IMAGE PATH]
EXAMPLES:
hide data from 'secret.dat' in the in 'img.png'
$ hide src.jpeg secret.dat img.png

only png files are supported as output files
$ hide src.jpeg secret.dat img.jpeg
        command failed: png is the only supported output image format
```

### Find

The `find` comman searches for hidden data in a PNG image.
It will only be able to find secrete data hidden by this tool.
```sh
$ imgdemo find help
find: find data hidden inside an image
USAGE:  find [IMAGE PATH]
EXAMPLES:
find hidden data from inside 'img.png'
$ find img.png
        Heres the hidden data

searching for hidden data in 'img.png' fails
$ find img.png
        command failed: magic number does not match
```

### Ishihara

[Ishihara test plates](https://en.wikipedia.org/wiki/Ishihara_test) are used to asses color blindness.
This command can be used to generate these images using a custom color palette and image mask.
```sh
$ imgdemo ishihara help
ishihara: create an ishihara image using the given color pallets and mask image
USAGE:  ishihara [PRIMARY COLORS] [SECONDARY COLORS] [MASK IMAGE PATH] [OUTPUT IMAGE PATH]
EXAMPLES:
create a red green colorblind test image
$ ishihara 3a6a2f,76cd63 a32222,db5f5f mask.png red_green.png

```
