# disgoboard

**Discord Soundboard CLI in Go.**

Building this to solve a very niche limitation in Discord. There doesn't seem to be a way to keybind specific soundboard sounds, so I figured I would make my own workaround. Not because it's needed, or even that I should, but because I wanted to.

There are a million ways to use your own soundboard and hotkey sounds, but not with Discord's native soundboard - so where's the fun in that?

### Status

Very early development. Not sure how far I want to take this.

### Current Limitations

- When playing a sound using Discord's API, the sound is only audible for others in the channel. Planning to add local playback to fix this.

- Hotkeys are handled by your system's keyboard shortcuts, not disgoboard itself.
