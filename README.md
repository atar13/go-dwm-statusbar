# go-dwm-statusbar

Lightweight status bar for the dwm window manager. 

### Modules:
- Date
- Time
- Battery
    - Requires acpi package
- Brightness
    - Requires xbacklight (package: xorg-xbacklight)
- CPU
- RAM
- MPRIS (Media playback status)
- Pulseaudo Volume Status
    - Requires pulseaudio package

### Download
```
git clone https://github.com/atar13/go-dwm-statusbar.git
```

### Install 

```
cd go-dwm-statusbar
sudo make install
```

### Configuration

Edit or create a config file at ```$HOME/.config/go-dwm-statusbar/config.yaml```

Full description of configuration options can be found in the [```config-sample.yaml```](./config-sample.yaml) file