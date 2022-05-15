<div align="center">
     <img src=".github/media/avatar_round.png" width="200"/>
     <h1>~ YURI 69 ~</h1>
     <strong>Even faster, harder, louder!<br/>The best Discord Soundboard you can get - Successor of <a href="https://github.com/zekroTJA/yuri2">yuri2</a></strong><br><br>
     <img src="https://forthebadge.com/images/badges/made-with-go.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/uses-html.svg" height="30" />&nbsp;
     <img src="https://forthebadge.com/images/badges/uses-css.svg" height="30" />&nbsp;
     <a href="https://zekro.de/discord"><img src="https://img.shields.io/discord/307084334198816769.svg?logo=discord&style=for-the-badge" height="30"></a>
</div>

---

Yuri69 is the next level of Discord Sound Boards. Featuring a fancy web interface where you can manage all your sounds, hotkey control, an open REST API to connect to devices like a Stream Deck, for example, and much more!

## Features

You can log in to the web interface using your Discord Account (via [OAuth2](https://oauth.net/2/)). There you can list, search, filter, upload, edit, favorize and remove sounds as well as managing the player state (join/leave chanel, stop playback, manage volume). Also, you can set personal settings like the hotkey trigger or guild filters for random sound playback.

<img src=".github/media/ss/0.png" width="100%" />

You can directly upload sounds via the web interface. There you can specify a unique ID, displayname and tags. These tags can be used for filtering and searching sounds. All uploaded files are converted using [FFMPEG](https://ffmpeg.org/) to be stored in an optimized file state for playback. 

<img src=".github/media/ss/1.png" width="100%" />

When going to the setting page, you can scan the displayed QR code using your mobile device to control the sound board from another device! You can also obtain an API key there to, for exmaple, play sounds from a stream deck or batch-upload sounds you have laying around (see [scripts/upload.sh](scripts/upload.sh)).

<img src=".github/media/ss/2.png" width="100%" />

There are also some general playback stats like a playback log or a list of most played sounds.

<img src=".github/media/ss/3.png" width="100%" />

Things like deletings or editing sounds uploaded by other users require admin permissions. These can be set in the admin panel. You need to specify yourself as `owner` in the config to get initial access to this panel and to add other admins.

<img src=".github/media/ss/4.png" width="100%" />

## Setup

*Coming soon™️*

## API Documentation

*Coming soon™️*