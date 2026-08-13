{
  config,
  lib,
  pkgs,
  ...
}: let
  cfg = config.programs.caelestia-extras;
  toml = pkgs.formats.toml {};
  package = pkgs.callPackage ./package.nix {};
  bibataSource = pkgs.runCommand "bibata-cursor-source" {} ''
    cp -r ${pkgs.bibata-cursors.src} "$out"
  '';
  configFile = toml.generate "caelestia-extras.toml" (
    {
      scheme.file = cfg.schemeFile;
    }
    // lib.optionalAttrs cfg.cursor.enable {
      cursor = {
        source = cfg.cursor.source;
        build_config = cfg.cursor.buildConfig;
        icon_dir = cfg.cursor.iconDir;
        theme = cfg.cursor.theme;
        size = cfg.cursor.size;
        xcursor_sizes = cfg.cursor.xcursorSizes;
        xcursor_fallback = cfg.cursor.xcursorFallback;
        update_gtk = cfg.cursor.updateGtk;
      };
    }
    // lib.optionalAttrs cfg.gtk.enable {
      gtk = {
        dark_theme = cfg.gtk.darkTheme;
        light_theme = cfg.gtk.lightTheme;
      };
    }
    // lib.optionalAttrs cfg.pavucontrol.enable {
      pavucontrol.command = cfg.pavucontrol.command;
    }
  );
  cursorRuntime = with pkgs; [
    coreutils
    dconf
    hyprcursor
    hyprland
    systemd
    util-linux
  ] ++ lib.optionals cfg.cursor.xcursorFallback [cbmp clickgen];
  gtkRuntime = [pkgs.dconf];
  command = "${package}/bin/caelestia-extras --config ${config.xdg.configHome}/caelestia-extras/config.toml";
in {
  options.programs.caelestia-extras = {
    enable = lib.mkEnableOption "optional Caelestia integrations";
    schemeFile = lib.mkOption {
      type = lib.types.str;
      default = "${config.xdg.stateHome}/caelestia/scheme.json";
    };
    cursor = {
      enable = lib.mkEnableOption "dynamic Bibata cursor";
      source = lib.mkOption { type = lib.types.str; default = "${bibataSource}/svg/modern"; };
      buildConfig = lib.mkOption { type = lib.types.str; default = "${bibataSource}/configs/normal/x.build.toml"; };
      iconDir = lib.mkOption { type = lib.types.str; default = "${config.xdg.dataHome}/icons"; };
      theme = lib.mkOption { type = lib.types.str; default = "Bibata-Caelestia"; };
      size = lib.mkOption { type = lib.types.positive; default = 20; };
      xcursorSizes = lib.mkOption { type = lib.types.listOf lib.types.positive; default = [20 24 32]; };
      xcursorFallback = lib.mkOption { type = lib.types.bool; default = true; };
      updateGtk = lib.mkOption { type = lib.types.bool; default = true; };
    };
    gtk = {
      enable = lib.mkEnableOption "Caelestia GTK preference sync";
      darkTheme = lib.mkOption { type = lib.types.str; default = "adw-gtk3-dark"; };
      lightTheme = lib.mkOption { type = lib.types.str; default = "adw-gtk3"; };
    };
    pavucontrol = {
      enable = lib.mkEnableOption "Caelestia-themed pavucontrol-qt launcher";
      command = lib.mkOption { type = lib.types.str; default = "pavucontrol-qt"; };
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [package] ++ lib.optionals cfg.pavucontrol.enable [pkgs.lxqt.pavucontrol-qt];
    xdg.configFile."caelestia-extras/config.toml".source = configFile;
    xdg.desktopEntries.pavucontrol-qt = lib.mkIf cfg.pavucontrol.enable {
      name = "PulseAudio Volume Control";
      genericName = "Volume Control";
      exec = "${command} pavucontrol";
      icon = "multimedia-volume-control";
      categories = ["AudioVideo" "Audio" "Mixer" "Qt"];
    };
    systemd.user.services = lib.mkMerge [
      (lib.mkIf cfg.cursor.enable {
        caelestia-extras-cursor = {
          Unit = { Description = "Sync cursor with the active Caelestia scheme"; After = ["graphical-session.target"]; };
          Service = {
            Type = "oneshot";
            Environment = "PATH=${lib.makeBinPath cursorRuntime}";
            ExecStart = "${command} cursor sync";
          };
          Install.WantedBy = ["graphical-session.target"];
        };
        caelestia-extras-xcursor = {
          Unit.Description = "Refresh the XCursor fallback for Caelestia";
          Service = {
            Type = "oneshot";
            Environment = "PATH=${lib.makeBinPath cursorRuntime}";
            ExecStart = "${command} cursor sync-xcursor";
          };
        };
      })
      (lib.mkIf cfg.gtk.enable {
        caelestia-extras-gtk = {
          Unit = { Description = "Sync GTK preferences with the active Caelestia scheme"; After = ["graphical-session.target"]; };
          Service = {
            Type = "oneshot";
            Environment = "PATH=${lib.makeBinPath gtkRuntime}";
            ExecStart = "${command} gtk sync";
          };
          Install.WantedBy = ["graphical-session.target"];
        };
      })
    ];
    systemd.user.paths = lib.mkMerge [
      (lib.mkIf cfg.cursor.enable {
        caelestia-extras-cursor = {
          Path.PathChanged = cfg.schemeFile;
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.gtk.enable {
        caelestia-extras-gtk = {
          Path.PathChanged = cfg.schemeFile;
          Install.WantedBy = ["graphical-session.target"];
        };
      })
    ];
  };
}
