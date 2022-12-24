---
title: Media
layout: default
grand_parent: Documentation
parent: Actions
nav_order: 5
---

# Media Actions
{: .no_toc }

## Table of contents
{: .no_toc .text-delta }

1. TOC
{:toc}

## Camera

### `takePhoto(showPreview)`

Takes a photo. `showPreview` is a optional `boolean` that defaults to `true`.

---

### `takePhotos(number)`

Takes `number` photo(s).

---

### `takeVideo(camera,quality,startImmediately)`

Takes a video using `camera` in `quality`.

- `camera` is a string value of `Front` or `Back`.
- `quality` is a string value of `Low`, `Medium`, `High`. Default is `Medium`.
- `startImmediately` is an optional boolean value set to `false` by default.

## Music

### `getCurrentSong()`

Gets the current song.

---

### `clearUpNext()`

Clears the up next songs.

## Photos

### `latestPhotoImport()`

Gets the latest photo import.

## Playback

### `setVolume(number)`

Set device volume to `number`.

## Video

### `trimVideo(video)`

Prompts the user to trim `video`. Returns the trimmed video.
