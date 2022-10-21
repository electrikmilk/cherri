[Back](../actions.md)

# Media Actions

### `clearUpNext()`

Clears the up next songs.

---

### `getCurrentSong()`

Gets the current song.

---

### `latestPhotoImport()`

Gets the latest photo import.

---

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

---

### `trimVideo(video)`

Prompts the user to trim `video`. Returns the trimmed video.

---

### `setVolume(number)`

Set device volume to `number`.
