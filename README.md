# vim-for-babies
Sorta like Telescope+Keymaps in Neovim, but for non-Neovim users. Built with Bubbletea, a TUI library

The motivation for this project is for VSCode users using Vim keybindings to have a way to see what keybindings are available to them. I also wanted to remember the keybindings I often Googled, so I made this. I also wanted to take this chance to explore Bubbletea, a TUI library for Golang. Oh, and did I mention that Copilot taught me a whole bunch of Vim Motions I didn't even know about? That was pretty cool. 

## Features
- Shows some basic list of keybindings
- Provides fuzzy finding for keybindings
- **Add your own bindings!**
- Delete an existing binding

## Installation
This relies on Go & `go install`. I also recommend aliasing the command to something shorter, like `vfb`. Here's how you can do that:
```
go install github.com/tanzy96/vim-for-babies
alias vfb=vim-for-babies
```

If you wish to run the program without installing it, you can do so cloning this repo and then with `go run .`

## Future features
- [ ] Add a way to search for keybindings by action
- [ ] Differentiate modes (normal, insert, visual, etc.) for keybindings
- [ ] Add more keybindings

Anyway this is just for fun. After this I got a lot better at Vim, so I don't really use this anymore. But I hope it helps someone else!
