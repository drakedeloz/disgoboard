# disgoboard

**Discord Soundboard CLI in Go.**

Building this to solve a very niche limitation in Discord. There doesn't seem to be a way to keybind specific soundboard sounds, so I figured I would make my own workaround. Not because it's needed, or even that I should, but because I wanted to.

There are a million ways to use your own soundboard and hotkey sounds, but not with Discord's native soundboard - so where's the fun in that?

### Status

Current build works, but could be refined.

- Might add an auth flow for initial setup
- Might add an easier way to import sounds

### Usage

Before using, ensure you've properly setup your config file as outlined in ```example_config.json```. The config should be placed in ```/home/<user>/.config/disgoboard```. I haven't yet added a way to set this up automatically. Feel free to reach out if you need help.

As of right now, disgoboard will only work in a single discord, but it will work in any channel within that discord. Some of the API calls require a bot token, so create a bot, invite it to your discord, and paste the bot token in your config file.

```disgoboard add <sound id> <guild source id>``` will fetch the sound from discord's cdn and cache a local mp3 file. 

```disgoboard list``` will print a list of all imported sounds.

```disgoboard play <sound>``` will play the sound in the channel you are connected to and on your local machine.

### Current Limitations

- Hotkeys are handled by your system's keyboard shortcuts, not disgoboard itself.
