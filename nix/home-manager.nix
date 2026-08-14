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
      };
    }
    // lib.optionalAttrs cfg.qbittorrent.enable {
      qbittorrent = {
        command = cfg.qbittorrent.command;
        rcc_command = cfg.qbittorrent.rccCommand;
        theme_dir = cfg.qbittorrent.themeDir;
        theme_file = cfg.qbittorrent.themeFile;
        config_file = cfg.qbittorrent.configFile;
      };
    }
    // lib.optionalAttrs cfg.portal.enable {
      portal = {
        theme_dir = cfg.portal.themeDir;
        config_home = cfg.portal.configHome;
        data_home = cfg.portal.dataHome;
        theme_name = cfg.portal.themeName;
        apply_global_gtk = cfg.portal.applyGlobalGtk;
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
  command = "${package}/bin/caelestia-extras --config ${config.xdg.configHome}/caelestia-extras/config.toml";
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
        description = "Directory containing generated Qt palette and stylesheet files.";
      };
      configHome = lib.mkOption {
        type = lib.types.str;
        default = config.xdg.configHome;
        description = "XDG configuration directory used by qt5ct and qt6ct.";
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
      enable = lib.mkEnableOption "Caelestia qBittorrent theme";
      command = lib.mkOption { type = lib.types.str; default = "qbittorrent"; };
      rccCommand = lib.mkOption { type = lib.types.str; default = "rcc"; };
      themeDir = lib.mkOption { type = lib.types.str; default = "${config.xdg.stateHome}/caelestia/theme"; };
      themeFile = lib.mkOption { type = lib.types.str; default = "${config.xdg.dataHome}/caelestia-extras/qbittorrent/caelestia.qbtheme"; };
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
      applyGlobalGtk = lib.mkOption {
        type = lib.types.bool;
        default = true;
        description = "Write Caelestia's generated GTK stylesheet to GTK 3 and GTK 4.";
      };
    };
  };

  config = lib.mkIf cfg.enable {
    home.packages = [package];
    xdg.configFile = {
      "caelestia-extras/config.toml".source = configFile;
    } // lib.optionalAttrs cfg.qt.enable {
      "caelestia/templates/qt-caelestia.conf".source = ../assets/manual/templates/qt-caelestia.conf;
      "caelestia/templates/qt-caelestia.qss".source = ../assets/manual/templates/qt-caelestia.qss;
    } // lib.optionalAttrs cfg.prismlauncher.enable {
      "caelestia/templates/prismlauncher.json".text = ''
        {
          "name": "Caelestia",
          "widgets": "Fusion",
          "qssFilePath": "themeStyle.css",
          "colors": {
            "Window": "#{{ surface.hex }}", "WindowText": "#{{ onSurface.hex }}",
            "Base": "#{{ surfaceContainerLowest.hex }}", "AlternateBase": "#{{ surfaceContainerLow.hex }}",
            "ToolTipBase": "#{{ inverseSurface.hex }}", "ToolTipText": "#{{ inverseOnSurface.hex }}",
            "Text": "#{{ onSurface.hex }}", "Button": "#{{ surfaceContainerHighest.hex }}", "ButtonText": "#{{ onSurface.hex }}",
            "BrightText": "#{{ error.hex }}", "Link": "#{{ primary.hex }}", "Highlight": "#{{ primary.hex }}", "HighlightedText": "#{{ onPrimary.hex }}",
            "fadeAmount": 0.42, "fadeColor": "#{{ surface.hex }}"
          },
          "logColors": {
            "Message": "#{{ onSurface.hex }}", "Launcher": "#{{ primary.hex }}", "Debug": "#{{ onSurfaceVariant.hex }}",
            "Warning": "#{{ tertiary.hex }}", "Error": "#{{ error.hex }}", "Fatal": "#{{ onErrorContainer.hex }}",
            "MessageHighlight": "#{{ surfaceContainerLow.hex }}", "LauncherHighlight": "#{{ primaryContainer.hex }}",
            "DebugHighlight": "#{{ surfaceContainer.hex }}", "WarningHighlight": "#{{ tertiaryContainer.hex }}",
            "ErrorHighlight": "#{{ errorContainer.hex }}", "FatalHighlight": "#{{ errorContainer.hex }}"
          }
        }
      '';
      "caelestia/templates/prismlauncher.qss".text = ''
        QWidget {
          color: #{{ onSurface.hex }};
          background-color: #{{ surface.hex }};
        }

        QMainWindow, QDialog, QStackedWidget, QScrollArea, QFrame {
          background-color: #{{ surface.hex }};
        }

        QAbstractItemView {
          alternate-background-color: #{{ surfaceContainerLow.hex }};
          selection-background-color: #{{ primaryContainer.hex }};
          selection-color: #{{ onPrimaryContainer.hex }};
          border: 1px solid #{{ outlineVariant.hex }};
          outline: 0;
        }

        QAbstractItemView::item {
          min-height: 28px;
          padding: 5px 8px;
          border: 0;
        }

        QAbstractItemView::item:hover {
          background-color: #{{ secondaryContainer.hex }};
          color: #{{ onSecondaryContainer.hex }};
        }

        QAbstractItemView::item:selected {
          background-color: #{{ primaryContainer.hex }};
          color: #{{ onPrimaryContainer.hex }};
        }

        QPushButton, QToolButton, QComboBox, QSpinBox, QLineEdit, QTextEdit, QPlainTextEdit {
          min-height: 30px;
          padding: 5px 10px;
          color: #{{ onSurface.hex }};
          background-color: #{{ surfaceContainer.hex }};
          border: 1px solid #{{ outlineVariant.hex }};
          border-radius: 7px;
        }

        QPushButton:hover, QToolButton:hover, QComboBox:hover, QSpinBox:hover,
        QLineEdit:hover, QTextEdit:hover, QPlainTextEdit:hover {
          background-color: #{{ surfaceContainerHigh.hex }};
          border-color: #{{ primary.hex }};
        }

        QPushButton:pressed, QToolButton:pressed, QPushButton:checked, QToolButton:checked {
          color: #{{ onPrimary.hex }};
          background-color: #{{ primary.hex }};
          border-color: #{{ primary.hex }};
        }

        QPushButton:disabled, QToolButton:disabled, QComboBox:disabled, QLineEdit:disabled {
          color: #{{ onSurfaceVariant.hex }};
          background-color: #{{ surfaceContainerLow.hex }};
          border-color: #{{ outlineVariant.hex }};
        }

        QLineEdit:focus, QTextEdit:focus, QPlainTextEdit:focus, QComboBox:focus, QSpinBox:focus {
          border: 2px solid #{{ primary.hex }};
          padding: 4px 9px;
        }

        QComboBox::drop-down { width: 28px; border: 0; }
        QComboBox QAbstractItemView { padding: 5px; background-color: #{{ surfaceContainerHigh.hex }}; }

        QTabWidget::pane {
          top: -1px;
          border: 1px solid #{{ outlineVariant.hex }};
          background-color: #{{ surface.hex }};
        }

        QTabBar::tab {
          min-height: 30px;
          padding: 6px 13px;
          color: #{{ onSurfaceVariant.hex }};
          background-color: #{{ surfaceContainerLow.hex }};
          border: 1px solid #{{ outlineVariant.hex }};
          border-bottom: 0;
          border-top-left-radius: 7px;
          border-top-right-radius: 7px;
        }

        QTabBar::tab:hover { color: #{{ onSurface.hex }}; background-color: #{{ surfaceContainer.hex }}; }
        QTabBar::tab:selected { color: #{{ onPrimaryContainer.hex }}; background-color: #{{ primaryContainer.hex }}; border-color: #{{ primary.hex }}; }

        QHeaderView::section {
          min-height: 30px;
          padding: 5px 9px;
          color: #{{ onSurfaceVariant.hex }};
          background-color: #{{ surfaceContainerHigh.hex }};
          border: 0;
          border-right: 1px solid #{{ outlineVariant.hex }};
          border-bottom: 1px solid #{{ outlineVariant.hex }};
        }

        QCheckBox, QRadioButton { spacing: 8px; }
        QCheckBox::indicator, QRadioButton::indicator {
          width: 17px;
          height: 17px;
          background-color: #{{ surfaceContainer.hex }};
          border: 1px solid #{{ outline.hex }};
          border-radius: 4px;
        }
        QRadioButton::indicator { border-radius: 9px; }
        QCheckBox::indicator:hover, QRadioButton::indicator:hover { border-color: #{{ primary.hex }}; }
        QCheckBox::indicator:checked, QRadioButton::indicator:checked { background-color: #{{ primary.hex }}; border-color: #{{ primary.hex }}; }

        QProgressBar {
          min-height: 10px;
          color: transparent;
          background-color: #{{ surfaceContainerHighest.hex }};
          border: 0;
          border-radius: 5px;
        }
        QProgressBar::chunk { background-color: #{{ primary.hex }}; border-radius: 5px; }

        QScrollBar:vertical { width: 12px; margin: 3px; background: transparent; }
        QScrollBar:horizontal { height: 12px; margin: 3px; background: transparent; }
        QScrollBar::handle:vertical, QScrollBar::handle:horizontal {
          min-height: 30px;
          min-width: 30px;
          background-color: #{{ outlineVariant.hex }};
          border-radius: 6px;
        }
        QScrollBar::handle:vertical:hover, QScrollBar::handle:horizontal:hover { background-color: #{{ primary.hex }}; }
        QScrollBar::add-line, QScrollBar::sub-line { width: 0; height: 0; }

        QMenu, QToolTip {
          padding: 6px;
          color: #{{ onSurface.hex }};
          background-color: #{{ surfaceContainerHigh.hex }};
          border: 1px solid #{{ outlineVariant.hex }};
          border-radius: 8px;
        }
        QMenu::item { min-height: 26px; padding: 5px 24px 5px 10px; border-radius: 5px; }
        QMenu::item:selected { color: #{{ onSecondaryContainer.hex }}; background-color: #{{ secondaryContainer.hex }}; }
      '';
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
      };
    home.activation.caelestiaExtrasPortal = lib.mkIf cfg.portal.enable (
      lib.hm.dag.entryAfter ["writeBoundary"] "${command} portal sync"
    );
    home.activation.caelestiaExtrasQt = lib.mkIf cfg.qt.enable (
      lib.hm.dag.entryAfter ["writeBoundary"] "${command} qt sync"
    );
    home.activation.caelestiaExtrasHyprtoolkit = lib.mkIf cfg.hyprtoolkit.enable (
      lib.hm.dag.entryAfter ["writeBoundary"] "${command} hyprtoolkit sync"
    );
    home.activation.caelestiaExtrasQBittorrent = lib.mkIf cfg.qbittorrent.enable (
      lib.hm.dag.entryAfter ["writeBoundary"] "${command} qbittorrent sync"
    );
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
      (lib.mkIf cfg.hyprtoolkit.enable {
        caelestia-extras-hyprtoolkit = {
          Unit = { Description = "Sync Hyprtoolkit with Caelestia"; After = ["graphical-session.target"]; };
          Service = {
            Type = "oneshot";
            ExecStart = "${command} hyprtoolkit sync";
          };
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.portal.enable {
        caelestia-extras-portal = {
          Unit = { Description = "Sync XDG portal theme with Caelestia"; After = ["graphical-session.target"]; };
          Service = {
            Type = "oneshot";
            ExecStart = "${command} portal sync";
          };
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.qt.enable {
        caelestia-extras-qt = {
          Unit = { Description = "Sync shared Qt theme with Caelestia"; After = ["graphical-session.target"]; };
          Service = {
            Type = "oneshot";
            ExecStart = "${command} qt sync";
          };
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.qbittorrent.enable {
        caelestia-extras-qbittorrent = {
          Unit = { Description = "Sync qBittorrent with Caelestia"; After = ["graphical-session.target"]; };
          Service = {
            Type = "oneshot";
            ExecStart = "${command} qbittorrent sync";
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
      (lib.mkIf cfg.hyprtoolkit.enable {
        caelestia-extras-hyprtoolkit = {
          Path.PathChanged = "${cfg.hyprtoolkit.themeDir}/hyprtoolkit.conf";
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.portal.enable {
        caelestia-extras-portal = {
          Path.PathChanged = [
            "${cfg.portal.themeDir}/gtk-portal.css"
            "${cfg.portal.themeDir}/gtk-global.css"
            "${cfg.portal.themeDir}/qt6ct-caelestia.conf"
            "${cfg.portal.themeDir}/qt6ct-portal.qss"
          ];
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.qt.enable {
        caelestia-extras-qt = {
          Path.PathChanged = [
            "${cfg.qt.themeDir}/qt-caelestia.conf"
            "${cfg.qt.themeDir}/qt-caelestia.qss"
          ];
          Install.WantedBy = ["graphical-session.target"];
        };
      })
      (lib.mkIf cfg.qbittorrent.enable {
        caelestia-extras-qbittorrent = {
          Path.PathChanged = [
            "${cfg.qbittorrent.themeDir}/qbittorrent.qss"
            "${cfg.qbittorrent.themeDir}/qbittorrent.json"
          ];
          Install.WantedBy = ["graphical-session.target"];
        };
      })
    ];
  };
}
