<div align="center">
  <a href="https://ebitenui.github.io/">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="ebitenui-logo-dark.svg">
      <img alt="ebitenui-logo" src="ebitenui-logo-light.svg" height="128">
    </picture>
  </a>
  <h1>Ebiten UI</h1>

[![Release](https://img.shields.io/github/v/release/ebitenui/ebitenui?style=for-the-badge&logo=github&labelColor=%23202e3bff&color=%235a7d93ff%20)](https://github.com/ebitenui/ebitenui/releases)
[![License](https://img.shields.io/github/license/ebitenui/ebitenui?style=for-the-badge&logo=github&labelColor=%23202e3bff&color=%235a7d93ff%20)](https://opensource.org/licenses/MIT)
[![Github](https://img.shields.io/badge/code-5a7d93ff?style=for-the-badge&logo=github&label=github&labelColor=%23202e3bff&color=%235a7d93ff)](https://github.com/ebitenui/ebitenui)
[![Docs](https://img.shields.io/badge/ebitenui.github.io-5a7d93ff?style=for-the-badge&logo=go&logoColor=white&label=docs&labelColor=%23202e3bff&color=%235a7d93ff)](https://ebitenui.github.io)
[![Discord](https://img.shields.io/discord/958140778931175424?style=for-the-badge&labelColor=%23202e3bff&color=%235a7d93ff%20&label=Discord&logo=discord&logoColor=white)](https://discord.gg/ujEeeHgptU)
[![Subreddit](https://img.shields.io/reddit/subreddit-subscribers/birdmtndev?style=for-the-badge&logo=reddit&logoColor=white&label=r%2Fbirdmtndev&labelColor=%23202e3bff&color=%235a7d93ff&cacheSeconds=120)](https://www.reddit.com/r/birdmtndev)
</div>
<br>

**A user interface engine and widget library for [Ebitengine](https://ebitengine.org/)**

**>> Note: This library is separate from [Ebitengine](https://ebitengine.org/). Please reach out to the linked Discord server for questions or suggestions**

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
