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
      compositor.backend = cfg.compositor.backend;
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
    // lib.optionalAttrs cfg.hyprtoolkit.enable {
      hyprtoolkit = {
        theme_dir = cfg.hyprtoolkit.themeDir;
        config_file = cfg.hyprtoolkit.configFile;
      };
    }
    // lib.optionalAttrs cfg.pavucontrol.enable {
      pavucontrol.command = cfg.pavucontrol.command;
    }
    // lib.optionalAttrs cfg.qt.enable {
      qt = {
        theme_dir = cfg.qt.themeDir;
        config_home = cfg.qt.configHome;
        data_home = cfg.qt.dataHome;
      };
    }
    // lib.optionalAttrs cfg.qbittorrent.enable {
      qbittorrent = {
        command = cfg.qbittorrent.command;
        config_file = cfg.qbittorrent.configFile;
      };
    }
    // lib.optionalAttrs cfg.portal.enable {
      portal = {
        theme_dir = cfg.portal.themeDir;
        config_home = cfg.portal.configHome;
        data_home = cfg.portal.dataHome;
        theme_name = cfg.portal.themeName;
      };
    }
  );
  compositorRuntime = {
    hyprland = [pkgs.hyprland];
  };
  cursorRuntime = with pkgs; [
    coreutils
    dconf
    hyprcursor
    systemd
    util-linux
  ] ++ compositorRuntime.${cfg.compositor.backend} ++ lib.optionals cfg.cursor.xcursorFallback [cbmp clickgen];
  gtkRuntime = [pkgs.dconf];
  watcherRuntime = lib.unique (cursorRuntime ++ gtkRuntime ++ [pkgs.systemd]);
  runtimePackage = pkgs.symlinkJoin {
    name = "caelestia-extras-with-runtime";
    paths = [package];
    nativeBuildInputs = [pkgs.makeWrapper];
    postBuild = ''
      wrapProgram "$out/bin/caelestia-extras" \
        --prefix PATH : ${lib.makeBinPath watcherRuntime}
    '';
  };
  command = "${runtimePackage}/bin/caelestia-extras --config ${config.xdg.configHome}/caelestia-extras/config.toml";
  qt6PluginPath = lib.makeSearchPathOutput "lib" "lib/qt-6/plugins" [
    pkgs.qt6Packages.qt6ct
    pkgs.kdePackages.breeze
  ];
  qbittorrentLauncher = pkgs.writeShellScript "caelestia-qbittorrent" ''
    export QT_QPA_PLATFORMTHEME=qt6ct
    export QT_STYLE_OVERRIDE=Breeze
    export QT_PLUGIN_PATH=${lib.escapeShellArg qt6PluginPath}''${QT_PLUGIN_PATH:+:$QT_PLUGIN_PATH}
    ${command} qbittorrent sync || exit $?
    exec ${lib.escapeShellArg cfg.qbittorrent.command} \
      ${lib.escapeShellArg "-stylesheet=${../assets/manual/templates/qbittorrent.qss}"} \
      "$@"
  '';
in {
  imports = [
    ./hyprtoolkit.nix
    ./portal.nix
    ./qt.nix
  ];

  options.programs.caelestia-extras = {
    enable = lib.mkEnableOption "optional Caelestia integrations";
    compositor.backend = lib.mkOption {
      type = lib.types.enum ["hyprland"];
      default = "hyprland";
      description = "Compositor backend used for compositor-specific integrations.";
    };
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
      size = lib.mkOption { type = lib.types.ints.positive; default = 20; };
      xcursorSizes = lib.mkOption { type = lib.types.listOf lib.types.ints.positive; default = [20 24 32]; };
      xcursorFallback = lib.mkOption { type = lib.types.bool; default = true; };
      updateGtk = lib.mkOption { type = lib.types.bool; default = true; };
    };
    gtk = {
      enable = lib.mkEnableOption "Caelestia GTK preference sync";
      darkTheme = lib.mkOption { type = lib.types.str; default = "adw-gtk3-dark"; };
      lightTheme = lib.mkOption { type = lib.types.str; default = "adw-gtk3"; };
      themeDir = lib.mkOption {
        type = lib.types.str;
        default = "${config.xdg.stateHome}/caelestia/theme";
        description = "Directory containing Caelestia's generated GTK colour stylesheet.";
      };
      directLaunch = lib.mkOption {
        type = lib.types.attrsOf (lib.types.submodule {
          options = {
            name = lib.mkOption { type = lib.types.str; };
            exec = lib.mkOption { type = lib.types.str; };
            icon = lib.mkOption { type = lib.types.nullOr (lib.types.either lib.types.str lib.types.path); default = null; };
            comment = lib.mkOption { type = lib.types.nullOr lib.types.str; default = null; };
            genericName = lib.mkOption { type = lib.types.nullOr lib.types.str; default = null; };
            categories = lib.mkOption { type = lib.types.nullOr (lib.types.listOf lib.types.str); default = null; };
            mimeType = lib.mkOption { type = lib.types.nullOr (lib.types.listOf lib.types.str); default = null; };
            startupNotify = lib.mkOption { type = lib.types.nullOr lib.types.bool; default = null; };
            settings = lib.mkOption { type = lib.types.attrsOf lib.types.str; default = {}; };
          };
        });
        default = {};
        description = ''
          Desktop-entry overrides that launch GTK applications directly instead
          of through D-Bus activation. Use the upstream desktop ID as the
          attribute name. This is useful for applications that only load
          dynamic GTK styling in a newly started process.
        '';
      };
    };
    hyprtoolkit = {
      enable = lib.mkEnableOption "Caelestia-generated Hyprtoolkit configuration";
      themeDir = lib.mkOption {
        type = lib.types.str;
        default = "${config.xdg.stateHome}/caelestia/theme";
        description = "Directory containing the generated Hyprtoolkit configuration.";
      };
      configFile = lib.mkOption {
        type = lib.types.str;
        default = "${config.xdg.configHome}/hypr/hyprtoolkit.conf";
        description = "Destination for Hyprtoolkit's active configuration.";
      };
    };
    pavucontrol = {
      enable = lib.mkEnableOption "Caelestia-themed pavucontrol-qt launcher";
      command = lib.mkOption { type = lib.types.str; default = "pavucontrol-qt"; };
    };
    qt = {
      enable = lib.mkEnableOption "shared Caelestia Qt theming";
      themeDir = lib.mkOption {
        type = lib.types.str;
        default = "${config.xdg.stateHome}/caelestia/theme";
        description = "Directory containing generated Qt palette and Breeze colour-scheme files.";
      };
      configHome = lib.mkOption {
        type = lib.types.str;
        default = config.xdg.configHome;
        description = "XDG configuration directory used by qt5ct and qt6ct.";
      };
      dataHome = lib.mkOption {
        type = lib.types.str;
        default = config.xdg.dataHome;
        description = "XDG data directory used for the generated Breeze colour scheme.";
      };
      iconTheme = lib.mkOption {
        type = lib.types.str;
        default = "Papirus-Dark";
        description = "Icon theme used by qt5ct and qt6ct.";
      };
    };
    prismlauncher = {
      enable = lib.mkEnableOption "Caelestia PrismLauncher theme";
      themeDir = lib.mkOption {
        type = lib.types.str;
        default = "${config.xdg.stateHome}/caelestia/theme";
        description = "Directory containing the generated PrismLauncher theme.";
      };
      themeName = lib.mkOption { type = lib.types.str; default = "caelestia-breeze"; };
    };
    qbittorrent = {
      enable = lib.mkEnableOption "native Breeze integration for qBittorrent";
      command = lib.mkOption { type = lib.types.str; default = "qbittorrent"; };
      configFile = lib.mkOption { type = lib.types.str; default = "${config.xdg.configHome}/qBittorrent/qBittorrent.conf"; };
    };
    portal = {
      enable = lib.mkEnableOption "Caelestia-themed XDG desktop portals";
      themeDir = lib.mkOption {
        type = lib.types.str;
        default = "${config.xdg.stateHome}/caelestia/theme";
        description = "Directory containing the Caelestia-generated portal theme files.";
      };
      configHome = lib.mkOption {
        type = lib.types.str;
        default = config.xdg.configHome;
        description = "XDG configuration directory used by portal-specific qt6ct.";
      };
      dataHome = lib.mkOption {
        type = lib.types.str;
        default = config.xdg.dataHome;
        description = "XDG data directory used for the generated GTK portal theme.";
      };
      themeName = lib.mkOption {
        type = lib.types.str;
        default = "Caelestia-Portal";
        description = "GTK theme name exposed exclusively to xdg-desktop-portal-gtk.";
      };
      iconTheme = lib.mkOption {
        type = lib.types.str;
        default = "Papirus-Dark";
        description = "Icon theme used by GTK and Qt portal file choosers.";
      };
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [runtimePackage];
    xdg.configFile = {
      "caelestia-extras/config.toml".source = configFile;
    } // lib.optionalAttrs cfg.gtk.enable {
      "caelestia/templates/gtk.css".source = ../assets/manual/templates/gtk.css;
      "caelestia/templates/gtk4.css".source = ../assets/manual/templates/gtk4.css;
      "gtk-3.0/gtk.css".source = config.lib.file.mkOutOfStoreSymlink "${cfg.gtk.themeDir}/gtk.css";
      "gtk-4.0/gtk.css".source = config.lib.file.mkOutOfStoreSymlink "${cfg.gtk.themeDir}/gtk4.css";
    } // lib.optionalAttrs (cfg.qt.enable || cfg.portal.enable) {
      "caelestia/templates/qt-caelestia.conf".source = ../assets/manual/templates/qt-caelestia.conf;
    } // lib.optionalAttrs cfg.qt.enable {
      "caelestia/templates/breeze-caelestia.colors".source = ../assets/manual/templates/breeze-caelestia.colors;
    } // lib.optionalAttrs cfg.pavucontrol.enable {
      "caelestia/templates/pavucontrol-qt.qss".source = ../assets/manual/templates/pavucontrol-qt.qss;
    } // lib.optionalAttrs cfg.prismlauncher.enable {
      "caelestia/templates/prismlauncher.json".source = ../assets/manual/templates/prismlauncher.json;
      "caelestia/templates/prismlauncher.qss".source = ../assets/manual/templates/qt6ct-caelestia.qss;
    };
    xdg.dataFile = lib.mkIf cfg.prismlauncher.enable {
      "PrismLauncher/themes/${cfg.prismlauncher.themeName}/theme.json".source =
        config.lib.file.mkOutOfStoreSymlink "${cfg.prismlauncher.themeDir}/prismlauncher.json";
      "PrismLauncher/themes/${cfg.prismlauncher.themeName}/themeStyle.css".source =
        config.lib.file.mkOutOfStoreSymlink "${cfg.prismlauncher.themeDir}/prismlauncher.qss";
    };
    xdg.desktopEntries =
      lib.optionalAttrs cfg.gtk.enable (lib.mapAttrs (_: entry: {
        inherit (entry) name exec icon comment genericName categories mimeType startupNotify;
        terminal = false;
        settings = entry.settings // {DBusActivatable = "false";};
      }) cfg.gtk.directLaunch)
      // lib.optionalAttrs cfg.pavucontrol.enable {
        pavucontrol-qt = {
          name = "PulseAudio Volume Control";
          genericName = "Volume Control";
          exec = "${command} pavucontrol";
          icon = "multimedia-volume-control";
          categories = ["AudioVideo" "Audio" "Mixer" "Qt"];
        };
      }
      // lib.optionalAttrs cfg.qbittorrent.enable {
        "org.qbittorrent.qBittorrent" = {
          name = "qBittorrent";
          genericName = "BitTorrent client";
          comment = "Download and share files over BitTorrent";
          exec = "${qbittorrentLauncher} %U";
          icon = "qbittorrent";
          categories = ["Network" "FileTransfer" "P2P" "Qt"];
          mimeType = ["application/x-bittorrent" "x-scheme-handler/magnet"];
          terminal = false;
          startupNotify = false;
          settings = {
            SingleMainWindow = "true";
            StartupWMClass = "qbittorrent";
          };
        };
      };
    home.activation.caelestiaExtrasRetirePathUnits = lib.hm.dag.entryAfter ["writeBoundary"] ''
      if [ -n "''${DBUS_SESSION_BUS_ADDRESS:-}" ]; then
        ${pkgs.systemd}/bin/systemctl --user disable --now \
          caelestia-extras-cursor.path \
          caelestia-extras-gtk.path \
          caelestia-extras-hyprtoolkit.path \
          caelestia-extras-portal.path \
          caelestia-extras-qt.path \
          caelestia-extras-qbittorrent.path >/dev/null 2>&1 || true
        ${pkgs.systemd}/bin/systemctl --user stop caelestia-extras-xcursor.service >/dev/null 2>&1 || true
        ${pkgs.systemd}/bin/systemctl --user reset-failed \
          caelestia-extras-cursor.path caelestia-extras-cursor.service \
          caelestia-extras-gtk.path caelestia-extras-gtk.service \
          caelestia-extras-hyprtoolkit.path caelestia-extras-hyprtoolkit.service \
          caelestia-extras-portal.path caelestia-extras-portal.service \
          caelestia-extras-qt.path caelestia-extras-qt.service \
          caelestia-extras-qbittorrent.path caelestia-extras-qbittorrent.service \
          caelestia-extras-xcursor.service >/dev/null 2>&1 || true
      fi
    '';
    home.activation.caelestiaExtrasSync = lib.hm.dag.entryAfter ["caelestiaExtrasRetirePathUnits"] ''
      if [ -n "''${DBUS_SESSION_BUS_ADDRESS:-}" ]; then
        PATH=${lib.makeBinPath watcherRuntime} ${command} sync
      fi
    '';
    systemd.user.services.caelestia-extras-watch = {
      Unit = {
        Description = "Keep desktop integrations synced with Caelestia";
        After = ["graphical-session.target"];
        PartOf = ["graphical-session.target"];
        StartLimitIntervalSec = 0;
      };
      Service = {
        Type = "simple";
        Environment = "PATH=${lib.makeBinPath watcherRuntime}";
        ExecStart = "${command} watch";
        Restart = "on-failure";
        RestartSec = 1;
      };
      Install.WantedBy = ["graphical-session.target"];
    };
  };
}
