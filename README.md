# Hyperpipe

<div align="center">

![](https://codeberg.org/Hyperpipe/Hyperpipe/raw/branch/main/public/android-chrome-192x192.png)

A Privacy Respecting Frontend for YouTube Music inspired and built with the help of [Piped][piped] and YouTube's InnerTube API.

![](https://img.shields.io/badge/youtube-music-red?style=for-the-badge&logo=youtube-music)
![Offical instance](https://img.shields.io/website?down_color=red&down_message=offline&label=status&style=for-the-badge&up_color=cornflowerblue&up_message=online&url=https%3A%2F%2Fhyperpipe.surge.sh%2f1x1.png)

</div>

## Disclaimer

**_HYPERPIPE IS VERY EXPIRIMENTAL, EXPECT_**

- Bugs
- Messy Code
- Stuff that used to work to break in the future.

_But if you see any of the following, Please open an issue_

PS: Please don't forget to support your favorite artists :)

## Usage

```sh
npm install
```

### Compile and Hot-Reload for Development

#### Local

```sh
npm run dev
```

#### Local + Network

```sh
npm run host
```

### Compile and Minify for Production

```sh
npm run build
```

## Instances

Please see [hyperpipe.codeberg.page](https://hyperpipe.codeberg.page)

## Automatic Redirection

### LibRedirect

Use [LibRedirect](https://github.com/libredirect/libredirect) to automatically redirect YouTube Music links to Hyperpipe.

### Redirector

You can use the [Redirector](https://github.com/einaregilsson/Redirector) extension to redirect **_MOST_** YouTube Music links to Hyperpipe with the configuration below (_Note:_ replace `hyperpipe.surge.sh` with any instance URL):

- Description: YouTube Music to Hyperpipe
- Example URL: https://music.youtube.com/channel/channel-id
- Include pattern: `https://music.youtube.com/*`
- Redirect to: `https://hyperpipe.surge.sh/$1`
- Pattern type: Wildcard Expression
- Pattern Description: Redirects all YouTube Music URLs to Hyperpipe
- Advanced options:
  - Apply to:
    - [x] Main window (address bar)

## Translation

### Weblate

[![Visit Hyperpipe on Weblate](https://hosted.weblate.org/widgets/hyperpipe/-/open-graph.png)](https://hosted.weblate.org/engage/hyperpipe/)

![](https://hosted.weblate.org/widgets/hyperpipe/-/multi-auto.svg)

### Codeberg

[See History](https://codeberg.org/Hyperpipe/Hyperpipe/pulls?q=&type=all&labels=69451)

## LICENSE

### AGPL-3.0-or-later

Hyperpipe

Copyright (C) 2022-24 Shiny Nematoda <snematoda.751k2@aleeas.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

Please refer to [LICENSE][license] for more details.

## Help

Feel free to join in on [Matrix!](https://matrix.to/#/#hyperpipe:nitro.chat)

You can reach out to me personally on:

- snematoda [dot] 751k2 [at] aleeas [dot] com
- [Matrix - @snematoda:nitro.chat](https://matrix.to/#/@snematoda:nitro.chat)

### Dependencies & Acknowledgments

- [Denevola](https://codeberg.org/Hyperpipe/Hyperpipe/pulls/22) -> for Hyperpipe's logo
- VueJS -> [MIT][vue]
- ViteJS -> [MIT][vite]
- shaka-player -> [Apache-2.0][shaka]
- PeerJS -> [MIT][peer]
- SortableJS -> [MIT][sortable]
- Bootstrap Icons -> [MIT][bi]
- VueJS theme -> [MIT][vuetheme]
- Dracula theme -> [MIT][dracula]
- Nord theme -> [MIT][nord]


### Similar Projects

*Hyperpipe is not affiliated with any of these projects*

- [Beatbump](https://github.com/snuffyDev/Beatbump) -> Alternative YouTube Music frontend built with Svelte/SvelteKit

- [VibeYou](https://github.com/you-apps/VibeYou) -> Privacy focused music player built with MD3

---

*All content on Hyperpipe is from YouTube Music. YouTube and YouTube Music are trademarks of Google LLC. Hyperpipe is not affiliated with YouTube or YouTube Music.*

[hypipe]: https://hyperpipe.surge.sh
[piped]: https://github.com/TeamPiped/Piped
[license]: https://codeberg.org/Hyperpipe/Hyperpipe/src/branch/main/LICENSE.md
[vue]: https://github.com/vuejs/core/blob/main/LICENSE
[vite]: https://github.com/vitejs/vite/blob/main/LICENSE
[bi]: https://github.com/twbs/icons/blob/main/LICENSE.md
[peer]: https://github.com/peers/peerjs/blob/master/LICENSE
[shaka]: https://github.com/shaka-project/shaka-player/blob/main/LICENSE
[sortable]: https://github.com/SortableJS/Sortable/blob/master/LICENSE
[nord]: https://github.com/arcticicestudio/nord/blob/develop/LICENSE.md
[vuetheme]: https://github.com/vuejs/theme/blob/main/LICENSE
[dracula]: https://github.com/dracula/dracula-theme/blob/master/LICENSE
