{
  config,
  lib,
  pkgs,
  ...
}: let
  cfg = config.programs.caelestia-extras.qt;
  qt5ct = pkgs.libsForQt5.qt5ct;
  qt6ct = pkgs.qt6Packages.qt6ct;
  qtConfig = version: ''
    [Appearance]
    color_scheme_path=${cfg.configHome}/${version}/colors/caelestia.conf
    custom_palette=true
    icon_theme=${cfg.iconTheme}
    style=Fusion

    [Interface]
    stylesheets=${cfg.configHome}/${version}/qss/caelestia.qss
  '';
in {
  config = lib.mkIf cfg.enable {
    home.packages = [qt5ct qt6ct];
    home.sessionVariables = {
      QT_PLUGIN_PATH = "${qt5ct}/lib/qt-5/plugins:${qt6ct}/lib/qt-6/plugins";
      QT_QPA_PLATFORMTHEME = "qt5ct:qt6ct";
      QT_STYLE_OVERRIDE = "Fusion";
    };
    xdg.configFile = {
      "qt5ct/qt5ct.conf".text = qtConfig "qt5ct";
      "qt6ct/qt6ct.conf".text = qtConfig "qt6ct";
    };
  };
}
