# Chromatika

Chromatika is a lightweight image color extraction tool and palette generator written in Go.

---

## Features

- Extract the dominant colors of an image via Fuzzy C-Means clustering over a quantized Lab color space.
- Generate harmonious palettes from an image's dominant colors.
- Export generated palettes with different formats.

---

## Installation

Make sure you have **Go 1.25.5+** installed

```bash
go install github.com/mermonia/chromatika@v0.1.0
```

The binary will be installed in $GOPATH/bin (usually ~/go/bin).

---

## Usage

### Generate a palette from an image

You can use the **chromatika palette** command to generate a palette from an image:

```bash
chromatika palette path-to-your-image.png
```

The command takes the path to the image file you want as a base for the palette as one of its arguments,
and outputs the generated palette (in the specified format) to the standard output. This makes it easier
to pipe this command to other tools, or to a file for future reference.

The **chromatika palette** command takes several flags, which you can just to adjust the algorithm's
parameters, as well as some optimizations and properties of the resulting palette. For example:

```bash
chromatika palette path-to-your-image.png -w 512 --darkmode --format toml -m 1.5 >myPalette.toml
```

The previous command will generate a darkmode palette from the specified image, encode it in a toml format
and output it to the **myPalette.toml** file. Before performing any calculations on the image, chromatika will
**downscale it to a width of 512px**. The performed FCM algorithm will use a **fuzziness value** of 1.5.

### Extract the dominant colors of an image

When generating a palette, chromatika might generate additional colors from the original dominant colors of
an image. Besides, generated palettes always have the same amount of colors, and must have neutral light and
dark colors.

If you ever want to simply *extract* the actual dominant colors of an image, you can use the **chromatika extract** command:

```bash
chromatika extract path-to-your-image.png
```

Just like the **palette** command, extract also allows you to configure several parameters, including the
amount of generated clusters/dominant colors:

```bash
chromatika extract path-to-your-image.png -k 5
```

## Algorithm

The algorithm used to extract dominant colors from images is FCM (Fuzzy C-Means), a well known clustering
algorithm that assigns, to each sample, a degree of membership to each one of the resulting clusters.

Running the algorithm with every single pixel of the image as a sample would be completely unfeasible though,
which is why chromatika uses color quantization to generate an admissible amount of samples:

- First, the source image's pixels are all converted to the (relatively) perceptually uniform color Lab space 
(since performing both quantization and clustering on a non-uniform space would produce poor results).
- Then, every pixel of the image is "assigned" to one single section of the quantized Lab space, and sections
which have no pixels assigned to them are discarded.
- After that, FCM is performed on the remaining color sections, each of them having a different weight (linearly 
depending on the amount of pixels assigned to them).

> Note: Conversion is performed from sRGB space to Lab space, so images in a color space different from sRGB will
not be properly analyzed. If this is a concern, consider converting them to sRGB space before analysis.

The centroids of each resulting cluster are sections of the Lab color space. The representative color
for each bin (section) is its center.

## Output formats

For now, there are only 2 available output formats for the **palette** command (ascii and toml),
while the **extract** command always outputs a list of hex colors. In the near future, both commands will
include ascii, toml, hex and json formats.

## Licence

[MIT](https://choosealicense.com/licenses/mit/)



