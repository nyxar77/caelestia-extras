{
  config,
  lib,
  pkgs,
  ...
}: let
  cfg = config.programs.caelestia-extras.qt;
  qt5ct = pkgs.libsForQt5.qt5ct;
  qt6ct = pkgs.qt6Packages.qt6ct;
  breeze = pkgs.kdePackages.breeze;
  qtPluginPath = "${qt5ct}/lib/qt-5/plugins:${qt6ct}/lib/qt-6/plugins:${breeze}/lib/qt-6/plugins";
  qtEnvironment = {
    QT_PLUGIN_PATH = qtPluginPath;
    QT_QPA_PLATFORMTHEME = "qt5ct:qt6ct";
  };
  qtConfig = version: ''
    [Appearance]
    color_scheme_path=${cfg.configHome}/${version}/colors/caelestia.conf
    custom_palette=true
    icon_theme=${cfg.iconTheme}
    style=${if version == "qt6ct" then "Breeze" else "Fusion"}
  '';
in {
  config = lib.mkIf cfg.enable {
    home.packages = [qt5ct qt6ct breeze];
    home.sessionVariables = qtEnvironment;
    home.activation.caelestiaExtrasQtEnvironment = lib.hm.dag.entryAfter ["writeBoundary"] ''
      if [ -n "''${DBUS_SESSION_BUS_ADDRESS:-}" ]; then
        ${pkgs.systemd}/bin/systemctl --user set-environment \
          QT_PLUGIN_PATH=${lib.escapeShellArg qtPluginPath} \
          QT_QPA_PLATFORMTHEME=qt5ct:qt6ct
        ${pkgs.dbus}/bin/dbus-update-activation-environment --systemd \
          QT_PLUGIN_PATH=${lib.escapeShellArg qtPluginPath} \
          QT_QPA_PLATFORMTHEME=qt5ct:qt6ct
      fi
    '';
    xdg.configFile = {
      "environment.d/10-caelestia-qt.conf".text = ''
        QT_PLUGIN_PATH=${qtPluginPath}
        QT_QPA_PLATFORMTHEME=qt5ct:qt6ct
      '';
      "qt5ct/qt5ct.conf".text = qtConfig "qt5ct";
      "qt6ct/qt6ct.conf".text = qtConfig "qt6ct";
    };
  };
}
