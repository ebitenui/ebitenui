[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/ebitenui/ebitenui?include_prereleases&label=Release)](https://github.com/ebitenui/ebitenui/releases)
[![GitHub](https://img.shields.io/github/license/ebitenui/ebitenui?color=blue&label=License)](https://opensource.org/licenses/MIT)
[![GoDoc](https://img.shields.io/badge/Go-Reference-blue?style=flat)](https://pkg.go.dev/github.com/ebitenui/ebitenui)
[![GoDoc](https://img.shields.io/badge/Go-Documentation-blue?style=flat)](https://ebitenui.github.io/)

[![Discord](https://img.shields.io/discord/958140778931175424?color=green&label=Discord&logo=discord&logoColor=white)](https://discord.gg/ujEeeHgptU)
[![Subreddit subscribers](https://img.shields.io/reddit/subreddit-subscribers/birdmtndev?color=green&label=r%2FBirdMtnDev&logo=reddit&logoColor=white)](https://www.reddit.com/r/birdmtndev/)

Ebiten UI
=========

**A user interface engine and widget library for [Ebitengine](https://ebitengine.org/)**

Ebiten UI is an extension to Ebitengine that provides the ability to render a complete user interface,
with widgets such as buttons, lists, combo boxes, and so on. It uses the [retained mode] model.
All graphics used by Ebiten UI can be fully customized, so you can really make your UI your own.

Documentation on how to use and extend Ebiten UI is available at [ebitenui.github.io](https://ebitenui.github.io).

![Screenshots](ebiten-ui.gif)

Quick Start
------
Ebiten UI is written in Go 1.19 which is available at [https://go.dev/](https://go.dev/).

There are Ebiten UI examples that can be found in the `_examples/` folder. 

They can be run from the root directory of the project with the following commands:
* Ebiten UI complete demo: `go run github.com/ebitenui/ebitenui/_examples/demo`
* Ebiten UI widget: `go run github.com/ebitenui/ebitenui/_examples/widget_demos/<folder_name>`

The examples can also be tested as WASM by running the following commands and opening your browser to [http://localhost:6262](http://localhost:6262):
* Ebiten UI complete demo: `go run github.com/hajimehoshi/wasmserve@latest -http=:6262 ./_examples/demo`
* Ebiten UI widget: `go run github.com/hajimehoshi/wasmserve@latest -http=:6262 ./_examples/widget_demos/<folder_name>`

Used By
------
* [Roboden by quasilyte](https://quasilyte.itch.io/roboden)
* [BANKWAVE: Neon Networth by Frabjous Studios](https://bankwave.frabjousstudios.com/)
* [Networked Game by Nmorenor](https://nmorenor.com/)
* [Sinecord by quasilyte](https://quasilyte.itch.io/sinecord)
* [Cavebots by quasilyte](https://quasilyte.itch.io/cavebots)


Social Media
-------
* [Discord](https://discord.gg/ujEeeHgptU)

* [Reddit](https://www.reddit.com/r/birdmtndev/)


License
-------

Ebiten UI is licensed under the [MIT license](https://opensource.org/licenses/MIT).

Maintainers
-------
* Currently maintained by Mark Carpenter <mark@bird-mtn.dev>
* Originally written by Maik Schreiber <blizzy@blizzy.de>


Contributing
-------
Want to help develop Ebiten UI? Check out our [current issues](https://github.com/ebitenui/ebitenui/issues). Want to know the steps on how to start contributing, take a look at the [open source guide](https://opensource.guide/how-to-contribute/).
