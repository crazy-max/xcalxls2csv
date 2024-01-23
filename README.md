[![xcalxls2csv](.github/xcalxls2csv.png)](https://github.com/crazy-max/xcalxls2csv)

## About

Converts [Xcalibur](https://www.thermofisher.com/nl/en/home/industrial/mass-spectrometry/liquid-chromatography-mass-spectrometry-lc-ms/lc-ms-software/lc-ms-data-acquisition-software/xcalibur-data-acquisition-interpretation-software.html)
XLS data frames to CSV format.

___

* [Download](#download)
* [Usage](#usage)
* [Build](#build)

## Download

xcalxls2csv binaries are available on the [release](https://github.com/crazy-max/xcalxls2csv/releases/latest)
page. Choose the binary matching your platform and test it with `./xcalxls2csv --help`
or move it to a permanent location:

```
$ ./xcalxls2csv --help
Usage: xcalxls2csv <xcalxls>

Xcalibur XLS data frames to CSV

Arguments:
  <xcalxls>    Xcalibur XLS file to convert.

Flags:
  -h, --help             Show context-sensitive help.
      --version
      --output=STRING    Custom output filename.
```

## Usage

`xcalxls2csv [<flags>] <xcalxls>`

`<xcalxls>` is the path to the Xcalibur XLS file to convert. Some samples are
available in [./pkg/xcal/fixtures](./pkg/xcal/fixtures):

```
$ xcalxls2csv ./pkg/xcal/fixtures/MCF001120_0_Short.xls
11:13AM INF Converting Xcalibur XLS file ./pkg/xcal/fixtures/MCF001120_0_Short.xls to ./pkg/xcal/fixtures/MCF001120_0_Short.csv
11:13AM INF Xcalibur XLS file converted successfully to ./pkg/xcal/fixtures/MCF001120_0_Short.csv
```

CSV file should be available in `./pkg/xcal/fixtures/MCF001120_0_Short.csv`.

## Build

```bash
# build and output to ./bin/build
$ docker buildx bake
# then run on linux
$ ./bin/build/xcalxls2csv --help
# on windows
$ .\bin\build\xcalxls2csv.exe --help
```

## License

MIT. See `LICENSE` for more details.
