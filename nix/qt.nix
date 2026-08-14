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
    home.sessionVariables = {
      QT_PLUGIN_PATH = "${qt5ct}/lib/qt-5/plugins:${qt6ct}/lib/qt-6/plugins:${breeze}/lib/qt-6/plugins";
      QT_QPA_PLATFORMTHEME = "qt5ct:qt6ct";
    };
    xdg.configFile = {
      "qt5ct/qt5ct.conf".text = qtConfig "qt5ct";
      "qt6ct/qt6ct.conf".text = qtConfig "qt6ct";
    };
  };
}
